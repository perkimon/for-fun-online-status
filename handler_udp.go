package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func processUDP(b []byte, raddr *net.UDPAddr, conn *net.UDPConn) {
	sr := &statusRequest{}
	err := json.Unmarshal(b, sr)
	if err != nil {
		fmt.Println("Bad post JSON")
		return
	}
	tracker.Add(sr.UserID, sr.FriendIDs, func(friendID int, online bool) {
		response := friendResponse{
			UserID: friendID,
			Online: online,
		}
		data, err := json.Marshal(response)
		if err != nil {
			log.Println("Json error")
			return
		}
		data = append(data, []byte("\n")...)
		_, err = conn.WriteToUDP(data, raddr)
		if err != nil {
			fmt.Println("Error responding")
			return
		}
	})
	tracker.Stateless(sr.UserID)
	return
}
