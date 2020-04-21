package status

import "sync"

type Tracker struct {
	FriendsList   map[int][]int
	userIDsOnline map[int]bool
	notifyMap     map[int]func(friendID int, online bool)
	mu            *sync.Mutex
}

func NewTracker() *Tracker {
	t := &Tracker{
		FriendsList:   make(map[int][]int),
		userIDsOnline: make(map[int]bool),
		notifyMap:     make(map[int]func(friendID int, online bool)),
		mu:            &sync.Mutex{},
	}

	return t
}

// Add a connection to the Tracker which will call the function to notify
// when a friend is online
// Returns a function to call to unsubscribe from the notifications
func (t *Tracker) Add(userID int, friendsIDs []int, notifyFn func(friendID int, online bool)) func() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.userIDsOnline[userID] = true
	t.FriendsList[userID] = friendsIDs
	//fmt.Println(t.FriendsList)
	t.notifyMap[userID] = notifyFn

	t.notifyUserOfFriendsStatus(userID)
	t.notifyFriendsOf(userID)

	return func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.userIDsOnline[userID] = false
		t.notifyFriendsOf(userID)
		delete(t.notifyMap, userID)
	}
}

func (t *Tracker) notifyUserOfFriendsStatus(userID int) {
	notifyFn, ok := t.notifyMap[userID]
	if !ok {
		return
	}
	friendIDs, ok := t.FriendsList[userID]
	//fmt.Println("Friends for ", userID, friendIDs)
	if !ok {
		//no friends found
		//fmt.Println("no friends list found")
		return
	}

	for _, friendID := range friendIDs {
		online, ok := t.userIDsOnline[friendID]
		if !ok {
			//fmt.Println("userIDsOnline user not in the map")
			online = false
		}
		notifyFn(friendID, online)
	}

}

func (t *Tracker) notifyFriendsOf(userID int) {
	online, ok := t.userIDsOnline[userID]
	if !ok {
		//fmt.Println("userIDsOnline user not in the map")
		online = false
	}

	friendIDs, ok := t.FriendsList[userID]
	//fmt.Println("Friends for ", userID, friendIDs)
	if !ok {
		//no friends found
		//fmt.Println("no friends list found")
		return
	}

	for _, friend := range friendIDs {
		notifyFn, ok := t.notifyMap[friend]
		if !ok {
			//fmt.Println("Notify FN not found")
			continue
		}
		notifyFn(userID, online)

	}
}
