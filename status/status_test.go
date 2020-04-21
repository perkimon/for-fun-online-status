package status

import (
	"testing"
)

func TestAdd(t *testing.T) {

	friend1results := make(map[int]bool)
	friend2results := make(map[int]bool)
	friend3results := make(map[int]bool)
	friend4results := make(map[int]bool)

	tracker := NewTracker()
	user1LogsoutFn := tracker.Add(1, []int{2, 3, 4}, func(friendID int, online bool) {
		//fmt.Println("userID 1: FriendID", friendID, "online:", online)
		friend1results[friendID] = online
	})

	user2LogsoutFn := tracker.Add(2, []int{1, 3, 4}, func(friendID int, online bool) {
		//fmt.Println("userID 2: FriendID", friendID, "online:", online)
		friend2results[friendID] = online
	})

	user3LogsoutFn := tracker.Add(3, []int{1, 2, 4}, func(friendID int, online bool) {
		//fmt.Println("userID 3: FriendID", friendID, "online:", online)
		friend3results[friendID] = online
	})

	user4LogsoutFn := tracker.Add(4, []int{1, 2, 3}, func(friendID int, online bool) {
		//fmt.Println("userID 3: FriendID", friendID, "online:", online)
		friend4results[friendID] = online
	})

	for _, friendID := range []int{2, 3, 4} {
		if friend1results[friendID] != true {
			t.Fatalf("Friend 1 should have friend %v online", friendID)
			return
		}
	}

	for _, friendID := range []int{1, 3, 4} {
		if friend2results[friendID] != true {
			t.Fatalf("Friend 2 should have friend %v online", friendID)
			return
		}
	}

	for _, friendID := range []int{1, 2, 4} {
		if friend3results[friendID] != true {
			t.Fatalf("Friend 3 should have friend %v online", friendID)
			return
		}
	}

	for _, friendID := range []int{1, 2, 3} {
		if friend4results[friendID] != true {
			t.Fatalf("Friend 4 should have friend %v online", friendID)
			return
		}
	}

	// fmt.Println(1, friend1results)
	// fmt.Println(2, friend2results)
	// fmt.Println(3, friend3results)
	// fmt.Println(4, friend4results)

	user1LogsoutFn()
	for _, friendID := range []int{1} {
		if friend2results[friendID] != false {
			t.Fatalf("Friend 2 should have friend %v offline", friendID)
			return
		}
	}

	for _, friendID := range []int{1} {
		if friend3results[friendID] != false {
			t.Fatalf("Friend 3 should have friend %v offline", friendID)
			return
		}
	}

	for _, friendID := range []int{1} {
		if friend4results[friendID] != false {
			t.Fatalf("Friend 4 should have friend %v offline", friendID)
			return
		}
	}

	user2LogsoutFn()

	for _, friendID := range []int{1, 2} {
		if friend3results[friendID] != false {
			t.Fatalf("Friend 3 should have friend %v offline", friendID)
			return
		}
	}

	for _, friendID := range []int{1, 2} {
		if friend4results[friendID] != false {
			t.Fatalf("Friend 4 should have friend %v offline", friendID)
			return
		}
	}
	user3LogsoutFn()

	for _, friendID := range []int{1, 2, 3} {
		if friend4results[friendID] != false {
			t.Fatalf("Friend 4 should have friend %v offline", friendID)
			return
		}
	}

	user4LogsoutFn()

}
