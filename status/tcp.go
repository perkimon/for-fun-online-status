package status

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"
)

type TCPResponder struct {
	conn net.Conn
}

func (ur *TCPResponder) Reply(fr *friendResponse) error {
	data, err := json.Marshal(fr)
	if err != nil {
		log.Println(err)
		return nil
	}
	data = append(data, []byte("\n")...)

	//set timeout
	ur.conn.SetWriteDeadline(time.Now().Add(time.Duration(time.Millisecond) * 500))
	_, err = ur.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (ur *TCPResponder) IsStateless() bool {
	return false
}

func tcpListener(incomingCh chan<- workIn) error {
	ln, err := net.Listen("tcp", ":2000")
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			log.Println("Connection")
			defer conn.Close()

			bufReader := bufio.NewReader(conn)
			deferExitSet := false
			for {
				bytes, err := bufReader.ReadBytes('\n')
				if err != nil {
					if err != io.EOF {
						log.Println("Buff error", err)
					}
					return
				}
				if len(bytes) < 2 {
					continue
				}
				request := &statusRequest{}
				err = json.Unmarshal(bytes, request)
				if err != nil {
					log.Println("JSON error", err)
					continue
				}

				incomingCh <- workIn{
					action: allowedUserActions(Joining),
					payload: &userContext{
						Responder: &TCPResponder{
							conn: conn,
						},
						ID:      request.UserID,
						Friends: request.FriendIDs,
					},
				}

				if !deferExitSet {
					//When the TCP connection exits send a Leaving message
					deferExitSet = true
					defer func() {
						incomingCh <- workIn{
							action: allowedUserActions(Leaving),
							payload: &userContext{
								Responder: &TCPResponder{
									conn: conn,
								},
								ID:      request.UserID,
								Friends: request.FriendIDs,
							},
						}
					}()
				}
			}

		}(conn)
	}

	return nil
}
