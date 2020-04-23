package status

import (
	"fmt"
	"sync"
	"time"
)

const StatelessTimeoutInSecs = 5

type Tracker struct {
	users     map[int]*User //user state
	mu        *sync.Mutex
	stateless map[int]bool //users connections which are stateless (UDP)
}

type User struct {
	id           int
	friends      []int
	online       bool
	lastSeen     time.Time
	notifyFn     func(friendID int, online bool)
	disconnectFn func()
}

func NewUser(id int, friends []int, online bool, lastseen time.Time, notifyFn func(friendID int, online bool)) *User {
	u := &User{
		id:       id,
		friends:  friends,
		online:   online,
		lastSeen: lastseen,
		notifyFn: notifyFn,
	}
	return u
}

func NewTracker() *Tracker {
	t := &Tracker{
		users:     make(map[int]*User),
		stateless: make(map[int]bool),
		mu:        &sync.Mutex{},
	}
	t.startTimeoutWorker(1)
	return t
}

// Add a connection to the Tracker which will call the function to notify
// when a friend is online
// Returns a function to call to unsubscribe from the notifications
func (t *Tracker) Add(userID int, friendsIDs []int, notifyFn func(friendID int, online bool)) func() {
	t.mu.Lock()
	defer t.mu.Unlock()
	//clear out memory for existing user
	delete(t.users, userID)

	user := NewUser(userID, friendsIDs, true, time.Now(), notifyFn)
	t.users[userID] = user

	t.notifyUserOfFriendsStatus(userID)
	t.notifyFriendsOf(userID)

	disconnectFn := func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.disconnectUser(userID)
	}
	return disconnectFn
}

func (t *Tracker) Stateless(userID int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stateless[userID] = true
}

func (t *Tracker) disconnectUser(userID int) {
	t.users[userID].online = false
	t.notifyFriendsOf(userID)
	delete(t.users, userID)
	fmt.Println(userID, "Disconnected FN")
}

func (t *Tracker) notifyUserOfFriendsStatus(userID int) {

	user, ok := t.users[userID]
	if !ok {
		return
	}
	friendIDs := user.friends

	for _, friendID := range friendIDs {
		online := false
		friend, ok := t.users[friendID]
		if !ok {
			//fmt.Println("userIDsOnline user not in the map")
			online = false
		} else {
			online = friend.online
		}

		user.notifyFn(friendID, online)
	}

}

func (t *Tracker) notifyFriendsOf(userID int) {
	user, ok := t.users[userID]
	if !ok {
		return
	}
	online := user.online
	friendIDs := user.friends

	for _, friend := range friendIDs {
		friend, ok := t.users[friend]
		if !ok {
			continue
		}
		friend.notifyFn(userID, online)
	}
}

func (t *Tracker) startTimeoutWorker(intervalInSeconds int) {
	go func() {
		for {
			time.Sleep(time.Duration(intervalInSeconds) * time.Second)
			t.disconnectOldStatelessConnections()
		}
	}()
}

func (t *Tracker) disconnectOldStatelessConnections() {
	t.mu.Lock()
	defer t.mu.Unlock()
	for userID, _ := range t.stateless {
		user, ok := t.users[userID]
		if !ok {
			delete(t.stateless, userID)
			continue
		}
		if time.Now().After(user.lastSeen.Add(time.Duration(1) * StatelessTimeoutInSecs * time.Second)) {
			t.disconnectUser(user.id)
			delete(t.stateless, user.id)
		}
	}
}
