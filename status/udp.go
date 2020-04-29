package status

import (
	"encoding/json"
	"log"
	"net"
)

type UDPResponder struct {
	raddr *net.UDPAddr
	conn  *net.UDPConn
}

func (ur *UDPResponder) Reply(fr *friendResponse) error {
	data, err := json.Marshal(fr)
	if err != nil {
		log.Println(err)
		return nil
	}
	data = append(data, []byte("\n")...)
	_, err = ur.conn.WriteToUDP(data, ur.raddr)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UDPResponder) IsStateless() bool {
	return true
}

func udpListener(incomingCh chan<- *requestContext) error {
	ladd, err := net.ResolveUDPAddr("udp", ":2000")
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
	go func() {
		b := make([]byte, 1024, 1024)
		oob := make([]byte, 1024, 1024)
		defer udpConn.Close()
		for {
			n, _, _, raddr, err := udpConn.ReadMsgUDP(b, oob)
			if err != nil {
				log.Println(err)
				return
			}
			if n > 0 {
				if n < 2 {
					continue
				}
				request := &statusRequest{}
				err := json.Unmarshal(b[:n], request)
				if err != nil {
					log.Println("JSON error", err)
					continue
				}

				if request.Action == Empty {
					request.Action = Joining
				}

				action := allowedUserActions(request.Action)

				incomingCh <- &requestContext{
					statusRequest: request,
					responder: &UDPResponder{
						raddr: raddr,
						conn:  udpConn,
					},
					action: action,
				}

			}
		}

	}()
	return nil
}
