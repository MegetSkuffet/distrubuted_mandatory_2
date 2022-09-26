package main

import (
	"bufio"
	"fmt"
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
	to      string
	from    client
	message string
	sync    int
	ack     int
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
		server.receive <- packet
	}

	time.Sleep(20 * time.Second)
}

func startServer(s server) {
	for {
		select {
		case p := <-s.receive:
			{
				if p.message == "" {
					//Packet is handshake
					fmt.Println("server received sync", p.sync)
					p.sync++
					p.ack = p.sync
					fmt.Println("server is sending sync", p.sync, "and acknowlegdement", p.ack)
					p.from.receive <- p
				} else {
					//Packet has message after handshake
					fmt.Println("server received sync", p.sync, "and acknowlegdement", p.ack)
					fmt.Println("Message received: ", p.message)
				}
			}

		default:
			continue
		}
	}
}

func birthClient(c client, s server) {
	for {
		select {
		case p := <-c.receive:
			{
				fmt.Println("client received ackknowlegdement", p.ack, "and sync", p.sync)
				if p.ack == 2 {
					p.sync++
					p.ack = p.sync
					fmt.Println("What message do you wish to send?")
					sc := bufio.NewScanner(os.Stdin)
					sc.Scan()
					p.message = sc.Text()

					fmt.Println("client is sending sync", p.sync, "and acknowlegdement", p.ack, "and message \"", p.message, "\"")
					s.receive <- p
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
