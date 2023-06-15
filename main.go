package main

import (
	"TCP-memory-database/memory"
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

	sv := create(*saverPath)
	db := memory.NewDB(sv)

	go serve(db)
	db.WaitForClose()
}

func serve(db *memory.DB) {
	li, err := net.ListenTCP("tcp", localhost())
	if err != nil {
		log.Fatal(err)
	}

	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go db.Handle(conn)
	}

}

func localhost() *net.TCPAddr {
	return &net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8080}
}
