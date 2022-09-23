package main

import (
	"fmt"
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

	var server = server{receive: make(chan packet), write: make(chan packet)}
	var client = client{receive: make(chan packet), write: make(chan packet)}

	go startServer(server)
	go birthClient(client, server)
	var packet = packet{from: client, sync: 1, ack: 0}
	server.receive <- packet
	/* for {
			fmt.Println("pls input command")

			var command,err = fmt.Scanln()
			_ = err

			switch command {
	        case "":
	            fmt.Println("1")
			}


		} */
	time.Sleep(20 * time.Second)
}

func startServer(s server) {
	for {
		select {
		case p := <-s.receive:
			{
				if p.message == "" {
					fmt.Println("server received sync", p.sync)
					p.sync++
					p.ack = p.sync
					fmt.Println("server is sending sync", p.sync, "and acknowlegdement", p.ack)
					p.from.receive <- p
				} else {
					fmt.Println("server received sync", p.sync, "and acknowlegdement", p.ack)
					fmt.Println(p.message)
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
					p.message = "hej med dig johan"
					fmt.Println("client is sending sync", p.sync, "and acknowlegdement", p.ack, "and message \"", p.message, "\"")
					s.receive <- p
				}
			}

		default:
			continue
		}
	}
}
