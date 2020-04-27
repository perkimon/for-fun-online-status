package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/perkimon/for-fun-online-status/status"
)

/*
# Backend code challenge
Write a TCP service in a language of your choice. If you are comfortable with Go, use that, but other languages are acceptable as long as you can explain the pros/cons and scaling characteristics.The service should have the following end points:1.

1. Start listening for tcp/udp connections.
2. Be able to accept connections.
3. Read json payload ({"user_id": 1, "friends": [2, 3, 4]})
3. After establishing successful connection - "store" it in memory the way you like.
4. When another connection established with the user_id from the list of any other user's "friends" section, they should be notified about it with message {"online": true}
5. When the user goes offline, his "friends" (if it has any and any of them online) should receive a message {"online": false}

Questions:
1. What changes if we switch TCP to UDP?
2. How would you detect the user (connection) is still online?
3. What happens if the user has a lot of friends?
4. How design of your application will change if there will be more than one instance of the server.
*/

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	LaunchInfo()
	status.Do()
	WaitForCtrlC()
	return
}

func LaunchInfo() {
	fmt.Println("Launched...  Exit with CTRL-C")
	fmt.Println(`nc localhost 2000 '{"user_id": 1, "friends": [2, 3, 4]}'`)
	fmt.Println(`nc localhost 2000 '{"user_id": 2, "friends": [1, 3, 4]}'`)
	fmt.Println(`nc localhost 2000 '{"user_id": 3, "friends": [1, 2, 4]}'`)
	fmt.Println(`nc -u localhost 2000 '{"user_id": 4, "friends": [1, 2, 3]}'`)
	fmt.Println(`Can also run netcat for UDP connections on nc -u localhost 2000 using the same json`)
}

func WaitForCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until CTRL-C is received.
	<-c
}
