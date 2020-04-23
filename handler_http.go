package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleStatus(w http.ResponseWriter, r *http.Request) {

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

}
