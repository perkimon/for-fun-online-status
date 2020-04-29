//go:generate mockgen -package status -source=channelsolution.go -destination channelsolution_mock.go
package status

import (
	"log"
	"time"
)

const (
	Empty = iota
	Joining
	Leaving
	Waving //for connectionless transports
	Check
	Remove
)

var OnlineTTL = time.Duration(30) * time.Second

type Responder interface {
	Reply(fr *friendResponse) error
	IsStateless() bool
}

type workIn struct {
	action  int
	payload interface{}
}

type requestContext struct {
	statusRequest *statusRequest
	responder     Responder
	action        int
	payload       interface{}
}

type userContext struct {
	ID        int
	Friends   []int
	Responder Responder
	Online    bool
	LastSeen  time.Time
}

type statusRequest struct {
	UserID    int   `json:"user_id"`
	FriendIDs []int `json:"friends"`
	Action    int   `json:"action,omitempty"`
}

type friendResponse struct {
	UserID int  `json:"user_id"`
	Online bool `json:"online"`
}

func Do() {
	incomingCh := make(chan workIn, 0)
	err := startConsumer(incomingCh)
	if err != nil {
		log.Println(err)
		return
	}
	err = udpListener(incomingCh)
	if err != nil {
		log.Println(err)
		return
	}

	err = tcpListener(incomingCh)
	if err != nil {
		log.Println(err)
		return
	}

}

func startConsumer(incomingCh chan workIn) error {

	go func() {
		users := make(map[int]*userContext)
		for work := range incomingCh {

			switch work.action {

			case Joining:
				//save state
				uc, ok := work.payload.(*userContext)
				if !ok {
					log.Println("Not a userContext")
					continue
				}
				users[uc.ID] = uc
				uc.Online = true
				uc.LastSeen = time.Now()

				asyncCheck(uc, incomingCh)

				//send user friend status'
				for _, friendID := range uc.Friends {
					friendContext, ok := users[friendID]
					online := false
					if !ok {
						online = false
					} else {
						online = users[friendID].Online
						if friendContext.Responder.IsStateless() {
							if time.Now().After(friendContext.LastSeen.Add(OnlineTTL)) {
								log.Println("UserId not seen in a while:", friendContext.ID, friendContext.LastSeen)
								online = false
							}
						}
					}

					Reply(incomingCh, uc, friendID, online)
				}

				//notify friends of user coming online
				for _, friendID := range uc.Friends {
					friendContext, ok := users[friendID]
					if !ok {
						//friend doesn't exist
						continue
					}

					Reply(incomingCh, friendContext, uc.ID, true)

				}
				log.Println("UserID", uc.ID, "Joined")
			case Leaving:
				uc, ok := work.payload.(*userContext)
				if !ok {
					log.Println("Not a userContext")
					continue
				}
				uc, ok = users[uc.ID]
				if !ok {
					log.Println("no user found")
					continue
				}

				//send notifications to all friends
				for _, friendID := range uc.Friends {

					friendContext, ok := users[friendID]
					if !ok {
						//friend doesn't exist
						continue
					}

					Reply(incomingCh, friendContext, uc.ID, false)
				}
				log.Println("UserID", uc.ID, "Left")
				//delete user
				delete(users, uc.ID)

			case Waving:
				//Update last seen
				uc, ok := work.payload.(*userContext)
				if !ok {
					log.Println("Not a userContext")
					continue
				}
				uc, ok = users[uc.ID]
				if ok {
					log.Println("UserID", uc.ID, "Waved")
					users[uc.ID].LastSeen = time.Now()
					asyncCheck(uc, incomingCh)
				}
			case Check:
				uc, ok := work.payload.(*userContext)
				if !ok {
					log.Println("Not a userContext")
					continue
				}

				uc, ok = users[uc.ID]
				if ok {
					log.Println("UserID", uc.ID, "Checking")
					if time.Now().After(uc.LastSeen.Add(OnlineTTL)) {
						log.Println("UserID", uc.ID, "Timedout - should send a wave to avoid this")
						go func() {
							incomingCh <- workIn{
								action:  Leaving,
								payload: uc,
							}
						}()
					}
				}
			case Remove:
				userID, ok := work.payload.(int)
				if !ok {
					log.Println("Payload was not an integer, user not deleted")
					continue
				}
				log.Println("Removing", userID)
				delete(users, userID)
			case Empty:
				log.Println("Empty action overridden")
			default:
				log.Println("Default action executed")
			}
		}
	}()
	return nil
}

func asyncCheck(uc *userContext, incomingCh chan workIn) {
	if uc.Responder.IsStateless() {
		go func() {
			time.Sleep(OnlineTTL)
			incomingCh <- workIn{
				action:  Check,
				payload: uc,
			}
		}()
	}
}

func allowedUserActions(action int) int {
	switch action {
	case Joining:
		return Joining
	case Leaving:
		return Leaving
	case Waving:
		return Waving
	default:
		return Empty
	}

	return Empty
}

// Reply to a connected user without blocking
// errors are handled by sending another request to delete the recipient on failure
func Reply(incomingCh chan workIn, notify *userContext, userID int, online bool) {
	go func() {
		err := notify.Responder.Reply(&friendResponse{
			UserID: userID,
			Online: online,
		})
		if err != nil {
			log.Println("error replying:", err, "deleting user", notify.ID)
			incomingCh <- workIn{
				action:  Remove,
				payload: notify.ID,
			}
		}
	}()
}
