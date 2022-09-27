package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
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
	testSlice    []int
}

func main() {

	sc := bufio.NewScanner(os.Stdin)
	var server = server{receive: make(chan packet), write: make(chan packet)}
	var client = client{receive: make(chan packet), write: make(chan packet)}

	go startServer(server)
	go birthClient(client, server)
	fmt.Println("### Please input command:")
	fmt.Println("### Enter 1 to simulate sending a single message to the server")
	fmt.Println("### Enter 2 to simulate sending several packets to the server and handling message reordering")

	var commandAsString string
	sc.Scan()
	commandAsString = sc.Text()

	switch commandAsString {
	case "1":
		sendSingleMessage(client, server)
	case "2":
		sendTestSlice(client, server)
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
				} else if p.testSlice != nil {
					fmt.Println("Server recieved simulated unsorted array and is sorting it: ")
					sort.Ints(p.testSlice)
					for _, v := range p.testSlice {
						fmt.Println(v)
					}
					p.message, _ = json.Marshal("Confirm")
					p.from.receive <- p

				} else {
					//Packet has message after handshake
					var s string
					json.Unmarshal(p.message, &s)
					localHash := hash(s)

					time.Sleep(1 * time.Second)

					fmt.Println("server received sync", p.sync, "and acknowlegdement", p.ack)

					time.Sleep(1 * time.Second)

					if p.hash != localHash {
						fmt.Println("Shits fucked")
					} else {

						fmt.Println("## Hash is correct ##")
						var msg string
						json.Unmarshal(p.message, &msg)
						fmt.Println("Message received: ", msg)
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
			var msg string
			json.Unmarshal(p.message, msg)
			{
				if p.ack == p.originalSync+1 && msg != "Confirm" && p.testSlice == nil {
					time.Sleep(1 * time.Second)

					fmt.Println("client received ackknowlegdement", p.ack, "and sync", p.sync)
					time.Sleep(1 * time.Second)

					p.ack = p.sync

					p.sync++
					fmt.Println("What message do you wish to send?")
					sc := bufio.NewScanner(os.Stdin)
					sc.Scan()
					p.message, _ = json.Marshal(sc.Text())
					p.hash = hash(sc.Text())

					time.Sleep(1 * time.Second)
					fmt.Println("client is sending sync", p.sync, "and acknowlegdement", p.ack, "and message \"", p.message, "\"")
					s.receive <- p

				} else if p.testSlice != nil {
					time.Sleep(1 * time.Second)

					fmt.Println("client received ackknowlegdement", p.ack, "and sync", p.sync)
					time.Sleep(1 * time.Second)

					p.ack = p.sync

					p.sync++
					noMsg := "noMsg"
					p.message, _ = json.Marshal(noMsg)
					p.hash = hash(noMsg)

					time.Sleep(1 * time.Second)
					fmt.Println("client is sending sync", p.sync, ", acknowlegdement", p.ack, " and sending simulated unsorted array")
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

func sendSingleMessage(c client, s server) {
	var packet = packet{from: c, sync: 1, ack: 0}
	packet.originalSync = packet.sync
	s.receive <- packet
}

func sendTestSlice(c client, s server) {
	slice := []int{1, 2, 3, 4, 5}
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	var packet = packet{from: c, sync: 1, ack: 0, testSlice: slice}
	packet.originalSync = packet.sync
	s.receive <- packet
}
