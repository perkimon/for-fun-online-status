package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func main() {
	config := &Config{}

	LaunchInfo(config)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	tracker := status.NewTracker()
	r := mux.NewRouter()
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Method Post is required"))
			return
		}
		sr := &statusRequest{}
		postdata, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte("Cannot read post data1"))
			return
		}
		err = json.Unmarshal(postdata, sr)
		if err != nil {
			w.Write([]byte("Bad post JSON"))
			return
		}
		w.WriteHeader(http.StatusOK)
		//Make sure that the writer supports flushing.
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("X-Accel-Buffering", "no")
		flusher.Flush()
		notify := r.Context().Done()
		fmt.Println(sr.UserID, "Joined")
		unregisterFn := tracker.Add(sr.UserID, sr.FriendIDs, func(friendID int, online bool) {

			response := friendResponse{
				UserID: friendID,
				Online: online,
			}
			data, err := json.Marshal(response)
			if err != nil {
				log.Println("Json error")
				return
			}
			out := fmt.Sprintf("%v\n", string(data))
			fmt.Fprintf(w, out)
			fmt.Println(string(data))
			flusher.Flush()
		})

		<-notify
		fmt.Println(sr.UserID, "Disconnected")
		unregisterFn()

	})

	go func() {
		err := http.ListenAndServe("localhost:2000", r)
		if err != nil {
			log.Println(err)
		}
	}()
	WaitForCtrlC()
	return
}

func LaunchInfo(c *Config) {
	fmt.Println("Launched...  Exit with CTRL-C")
	fmt.Println("Use curl to update status and listen for changes")
	fmt.Println(`curl -X POST -d '{"user_id": 1, "friends": [2, 3, 4]}' http://localhost:2000/status`)
	fmt.Println(`curl -X POST -d '{"user_id": 2, "friends": [1, 3, 4]}' http://localhost:2000/status`)
	fmt.Println(`curl -X POST -d '{"user_id": 3, "friends": [1, 2, 4]}' http://localhost:2000/status`)
	fmt.Println(`curl -X POST -d '{"user_id": 4, "friends": [1, 2, 3]}' http://localhost:2000/status`)
}

func WaitForCtrlC() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	// Block until CTRL-C is received.
	<-c
}
