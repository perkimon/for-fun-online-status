package status

import "time"

type User struct {
	id           int
	friends      map[int]bool
	online       bool
	lastSeen     time.Time
	notifyFn     func(friendID int, online bool)
	disconnectFn func()
}

func NewUser(id int, friends []int, online bool, lastseen time.Time, notifyFn func(friendID int, online bool)) *User {

	//create map of friends to make verification easier
	friendsMap := make(map[int]bool)
	for _, friendID := range friends {
		friendsMap[friendID] = true
	}
	u := &User{
		id:       id,
		friends:  friendsMap,
		online:   online,
		lastSeen: lastseen,
		notifyFn: notifyFn,
	}
	return u
}

func (u *User) IsFriend(friendID int) bool {

	_, ok := u.friends[friendID]
	if !ok {
		return false
	}
	return true
}
