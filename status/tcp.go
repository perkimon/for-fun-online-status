package status

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
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
	_, err = ur.conn.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (ur *TCPResponder) IsStateless() bool {
	return false
}

func tcpListener(incomingCh chan<- *requestContext) error {
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
				request := &statusRequest{}
				err = json.Unmarshal(bytes, request)
				if err != nil {
					log.Println("JSON error", err)
					return
				}

				incomingCh <- &requestContext{
					statusRequest: request,
					responder: &TCPResponder{
						conn: conn,
					},
					Action: Joining,
				}
				if !deferExitSet {
					deferExitSet = true
					defer func() {

						incomingCh <- &requestContext{
							statusRequest: &statusRequest{
								UserID:    request.UserID,
								FriendIDs: request.FriendIDs,
							},
							responder: &TCPResponder{
								conn: conn,
							},
							Action: Leaving,
						}
					}()
				}
			}

		}(conn)
	}

	return nil
}
