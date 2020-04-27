package status

import (
	"log"
	"time"
)

// Example status requests with optional action
//{"user_id": 2, "friends": [1, 3, 4],"Action":2}
//{"user_id": 2, "friends": [1, 3, 4],"Action":1}

type Responder interface {
	Reply(fr *friendResponse) error
	IsStateless() bool
}

type requestContext struct {
	statusRequest *statusRequest
	responder     Responder
	Action        int
}

type userContext struct {
	ID        int
	Friends   []int
	Responder Responder
	Online    bool
	LastSeen  time.Time
}

const (
	Empty = iota
	Joining
	Leaving
	Waving //for connectionless transports
	Check
)

const OnlineTTL = 30

type statusRequest struct {
	UserID    int   `json:"user_id"`
	FriendIDs []int `json:"friends"`
	Action    int   `json:",omitempty"`
}

type friendResponse struct {
	UserID int  `json:"user_id"`
	Online bool `json:"online"`
}

func Do() {
	incomingCh := make(chan *requestContext, 0)
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

func startConsumer(incomingCh chan *requestContext) error {
	go func() {
		users := make(map[int]*userContext)
		for rc := range incomingCh {
			//fmt.Println(rc.statusRequest.UserID, rc.statusRequest.FriendIDs)
			switch rc.Action {

			case Joining:
				//save state
				uc := &userContext{
					ID:        rc.statusRequest.UserID,
					Friends:   rc.statusRequest.FriendIDs,
					Responder: rc.responder,
					Online:    true,
					LastSeen:  time.Now(),
				}
				users[rc.statusRequest.UserID] = uc

				asyncCheck(uc, rc, incomingCh)

				//send user friend status'
				for _, friendID := range uc.Friends {
					friendContext, ok := users[friendID]
					online := false
					if !ok {
						online = false
					} else {
						online = users[friendID].Online
						if friendContext.Responder.IsStateless() {
							if time.Now().After(friendContext.LastSeen.Add(time.Duration(1) * OnlineTTL * time.Second)) {
								log.Println("UserId not seen in a while:", friendContext.ID, friendContext.LastSeen)
								online = false
							}
						}
					}

					fr := &friendResponse{
						UserID: friendID,
						Online: online,
					}
					uc.Responder.Reply(fr)
				}

				//notify friends of user coming online
				for _, friendID := range uc.Friends {
					friendContext, ok := users[friendID]
					if !ok {
						//friend doesn't exist
						continue
					}

					fr := &friendResponse{
						UserID: uc.ID,
						Online: true,
					}

					err := friendContext.Responder.Reply(fr)
					if err != nil {
						log.Println("error responding, deleting user")
						delete(users, uc.ID)
					}
				}
				log.Println("UserID", uc.ID, "Joined")
			case Leaving:
				uc, ok := users[rc.statusRequest.UserID]

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
					fr := &friendResponse{
						UserID: uc.ID,
						Online: false,
					}
					err := friendContext.Responder.Reply(fr)
					if err != nil {
						log.Println("error responding, deleting user", fr.UserID)
						delete(users, uc.ID)
					}
				}
				log.Println("UserID", uc.ID, "Left")
				//delete user
				delete(users, uc.ID)

			case Waving:
				//Update last seen
				uc, ok := users[rc.statusRequest.UserID]
				if ok {
					log.Println("UserID", rc.statusRequest.UserID, "Waved")
					users[rc.statusRequest.UserID].LastSeen = time.Now()
					asyncCheck(uc, rc, incomingCh)
				}
			case Check:
				//Update last seen
				uc, ok := users[rc.statusRequest.UserID]
				if ok {
					log.Println("UserID", uc.ID, "Checking")
					if time.Now().After(uc.LastSeen.Add(time.Duration(1) * OnlineTTL * time.Second)) {
						log.Println("UserID", uc.ID, "Timedout - should send a wave to avoid this")
						go func() {
							incomingCh <- &requestContext{
								statusRequest: rc.statusRequest,
								responder:     rc.responder,
								Action:        Leaving,
							}
						}()
					}
				}
			default:

			}
		}
	}()
	return nil
}

func asyncCheck(uc *userContext, rc *requestContext, incomingCh chan *requestContext) {
	if uc.Responder.IsStateless() {
		go func() {
			time.Sleep(time.Duration(time.Second) * OnlineTTL)
			incomingCh <- &requestContext{
				statusRequest: rc.statusRequest,
				responder:     rc.responder,
				Action:        Check,
			}
		}()
	}
}
