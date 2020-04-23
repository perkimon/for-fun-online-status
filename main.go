package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
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

type Config struct {
}

type statusRequest struct {
	UserID    int   `json:"user_id"`
	FriendIDs []int `json:"friends"`
}

type friendResponse struct {
	UserID int  `json:"user_id"`
	Online bool `json:"online"`
}

var tracker *status.Tracker
var resetCh map[int]chan bool

func main() {
	config := &Config{}
	LaunchInfo(config)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tracker = status.NewTracker()
	resetCh = make(map[int]chan bool)

	r := mux.NewRouter()
	r.HandleFunc("/status", HandleStatus)

	go func() {
		err := http.ListenAndServe(":2000", r)
		if err != nil {
			log.Println(err)
		}
	}()

	err := listenForUDP(":2000")
	if err != nil {
		log.Println(err)
	}

	WaitForCtrlC()
	return
}

func listenForUDP(addr string) error {
	ladd, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", ladd)
	if err != nil {
		return err
	}
	size := 1024 * 1024
	err = udpConn.SetReadBuffer(size)
	if err != nil {
		return err
	}
	b := make([]byte, 1024, 1024)
	oob := make([]byte, 1024, 1024)
	go func() {
		for {
			n, _, _, raddr, err := udpConn.ReadMsgUDP(b, oob)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if n > 0 {
				// Copy into new byte array in case multiple go routines are used to process array
				// Avoids overwriting data from byte array re-use
				UDPin := make([]byte, n)
				copy(UDPin, b[:n])
				processUDP(UDPin, raddr, udpConn)
			}
		}
		defer udpConn.Close()

	}()

	return nil
}

func LaunchInfo(c *Config) {
	fmt.Println("Launched...  Exit with CTRL-C")
	fmt.Println("Use curl to update status and listen for changes")
	fmt.Println(`curl -X POST -d '{"user_id": 1, "friends": [2, 3, 4]}' http://localhost:2000/status`)
	fmt.Println(`curl -X POST -d '{"user_id": 2, "friends": [1, 3, 4]}' http://localhost:2000/status`)
	fmt.Println(`curl -X POST -d '{"user_id": 3, "friends": [1, 2, 4]}' http://localhost:2000/status`)
	fmt.Println(`curl -X POST -d '{"user_id": 4, "friends": [1, 2, 3]}' http://localhost:2000/status`)
	fmt.Println(`Can also run netcat for UDP connections on nc -u localhost 2000 using the same json`)
}

func WaitForCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until CTRL-C is received.
	<-c
}
