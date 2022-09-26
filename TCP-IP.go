package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"time"
)

type server struct {
	receive chan packet
	write   chan packet
}

type client struct {
	receive chan packet
	write   chan packet
}

type packet struct {
	to           string
	from         client
	message      []byte
	sync         int
	originalSync int
	ack          int
	hash         uint32
}

func main() {

	sc := bufio.NewScanner(os.Stdin)
	var server = server{receive: make(chan packet), write: make(chan packet)}
	var client = client{receive: make(chan packet), write: make(chan packet)}

	go startServer(server)
	go birthClient(client, server)
	fmt.Println("## pls input command ##")

	var commandAsString string
	sc.Scan()
	commandAsString = sc.Text()

	switch commandAsString {
	case "send message":
		var packet = packet{from: client, sync: 1, ack: 0}
		packet.originalSync = packet.sync
		server.receive <- packet
	}

	time.Sleep(20 * time.Second)
}

func startServer(s server) {
	for {
		select {
		case p := <-s.receive:
			{
				var s string
				json.Unmarshal(p.message, &s)
				if s == "" {
					//Packet is handshake
					fmt.Println("server received sync", p.sync)
					time.Sleep(1 * time.Second)
					p.sync++
					p.ack = p.sync
					fmt.Println("server is sending sync", p.sync, "and acknowlegdement", p.ack)
					p.from.receive <- p
				} else {
					//Packet has message after handshake
					var s string
					json.Unmarshal(p.message, &s)
					localHash := hash(s)

					fmt.Println("server received sync", p.sync, "and acknowlegdement", p.ack)

					if p.hash != localHash {
						fmt.Println("Shits fucked")
					} else {

						fmt.Println("## Hash is correct ##")
						fmt.Println("Message received: ", p.message)
						p.message, _ = json.Marshal("Confirm")

						p.from.receive <- p
					}

				}
			}

		default:
			continue
		}
	}
}

func birthClient(c client, s server) {
	var timeout int
	for {
		select {
		case p := <-c.receive:
			{
				if p.ack == p.originalSync+1 && p.message != "Confirm" {
					fmt.Println("client received ackknowlegdement", p.ack, "and sync", p.sync)

					p.ack = p.sync

					p.sync++
					fmt.Println("What message do you wish to send?")
					sc := bufio.NewScanner(os.Stdin)
					sc.Scan()
					p.message = sc.Text()
					p.hash = hash(p.message)
					a, _ := json.Marshal(p.message)

					fmt.Println("client is sending sync", p.sync, "and acknowlegdement", p.ack, "and message \"", p.message, "\"")
					s.receive <- p

				}
				for i := 0; i < 6; i++ {

					select {
					case <-c.receive:
						fmt.Println("Client received answer correctly")
						break
					default:
						break
					}

					if i == 5 {
						s.receive <- p
						i = 0
						timeout++
					} else if timeout == 3 {
						fmt.Println("could not send message")
						os.Exit(0)
					}

				}
			}

		default:
			continue
		}
	}
}

func clientSendMessage(c client, s server, msg string) {
	var packet = packet{from: c, sync: 1, ack: 0}
	s.receive <- packet

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
