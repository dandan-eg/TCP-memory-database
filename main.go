package main

import (
	"TCP-memory-database/db"
	"TCP-memory-database/saver"
	"flag"
	"log"
	"net"
)

func main() {

	saverFormat := flag.String("save", "json", "")
	saverPath := flag.String("out", "./", "")
	flag.Parse()

	create, ok := saver.Factory[*saverFormat]
	if !ok {
		log.Fatalf("%s is not a supported format", *saverFormat)
	}

	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	defer li.Close()

	sv := create(*saverPath)
	memory := db.New(sv)

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
