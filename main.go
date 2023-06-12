package main

import (
	"TCP-memory-database/db"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer li.Close()
	memory := db.New()

loop:
	for {
		select {
		default:
			conn, err := li.Accept()
			if err != nil {
				log.Fatal(err)
			}

			go memory.Handle(conn)

		case <-memory.Quit:
			break loop
		}

	}

}
