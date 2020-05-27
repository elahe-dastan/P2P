package server

import (
	"fmt"
	"github.com/elahe-dstn/p2p/request"
	"net"
)

type Server struct {
	IP string
}

func New() Server {
	return Server{
		IP: "127.0.0.1",
	}
}

func (s *Server) Up() {
	addr := net.UDPAddr{
		IP:   net.ParseIP(s.IP),
		Port: 1378,
	}

	add, err := net.ResolveUDPAddr("udp", addr.String())
	print(add)

	ser, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Println(err)
		return
	}

	message := make([]byte, 2048)

	_, remoteAddr, err := ser.ReadFromUDP(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	req := request.Unmarshal(string(message))

	protocol(req, ser, remoteAddr)
}

func protocol(req request.Request, ser *net.UDPConn, remoteAddr *net.UDPAddr) {
	switch req.(type) {
	case request.Discover:
		go transfer(ser, remoteAddr, "this should be the list")
		//case request.File:
		//	f := req.(request.File)
		//	if n.Search(f.Name) {
		//		go transfer(ser, remoteAddr, response.File{Answer: true, TcpPort: n.TcpPort}.Marshal())
		//	}
	}
	// if t == "list" {
	//	n.merge(protocol[1:])
	//}
	//} else if t == "ans" {
	//	if n.waiting {
	//		n.check(protocol[1:])
	//	}
	//}
}

func transfer(conn *net.UDPConn, addr *net.UDPAddr, message string) {
	_, err := conn.WriteToUDP([]byte(message), addr)
	if err != nil {
		fmt.Println(err)
		return
	}
}