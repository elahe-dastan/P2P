package node

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elahe-dstn/p2p/cluster"
	"github.com/elahe-dstn/p2p/config"
	"github.com/elahe-dstn/p2p/tcp/client"
	tcp "github.com/elahe-dstn/p2p/tcp/server"
	udp "github.com/elahe-dstn/p2p/udp/server"
)

type Node struct {
	UDPServer udp.Server
	TCPServer tcp.Server
	TCPClient client.Client
	TCPPort   chan int
	Addr      chan string
	fName     chan string
}

func New(folder string, c []string) Node {
	clu := cluster.New(c)

	cfg := config.Read()

	ip := cfg.Host
	port := cfg.Port
	d := cfg.DiscoveryPeriod
	waitingDuration := cfg.WaitingTime

	return Node{
		UDPServer: udp.New(ip, port, &clu, time.NewTicker(time.Duration(d)*time.Second), waitingDuration, folder),
		TCPServer: tcp.New(folder),
		TCPClient: client.New(folder),
		TCPPort:   make(chan int),
		Addr:      make(chan string, 1),
		fName:     make(chan string),
	}
}

func (n *Node) Run() {
	reader := bufio.NewReader(os.Stdin)

	go n.TCPServer.Up(n.TCPPort)

	go n.UDPServer.Up(n.TCPPort, n.Addr, n.fName)

	go n.UDPServer.Discover()

	go n.TCPClient.Connect(n.Addr, n.fName)

	for {
		fmt.Println("Enter a file you want to download or list to see the cluster")

		text, err := reader.ReadString('\n')


		if err != nil {
			fmt.Println(err)
			return
		}

		text = strings.TrimSuffix(text, "\n")

		fmt.Println(text)

		req := strings.Split(text, " ")

		if req[0] == "list" {
			fmt.Println(n.UDPServer.Cluster.List)
		}else if req[0] == "get" {
			n.UDPServer.Req = req[1]
			n.UDPServer.File()
			//n.fName <- req[1]
		}
	}
}
