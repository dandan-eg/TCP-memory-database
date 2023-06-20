package main

import (
	"TCP-memory-database/memory"
	"TCP-memory-database/saver"
	"TCP-memory-database/src"
	"flag"
	"log"
	"net"
)

func main() {
	// Command line flags
	saverPath := flag.String("out", "./", "Save path")
	saverType := flag.String("save", "json", "Save format (json, csv)")
	srcPath := flag.String("src", "", "Source file path")

	flag.Parse()

	save, err := saver.New(*saverPath, *saverType)
	if err != nil {
		log.Fatal(err)
	}

	db := memory.NewDB(save)

	var source src.Sourcer

	if *srcPath != "" {
		source, err = src.New(*srcPath)

		if err != nil {
			log.Fatal(err)
		}
	}

	if err := db.Load(source); err != nil {
		log.Fatal(err)
	}

	// Start the server
	go serve(db, ":8080")
	db.WaitForClose()
}

func serve(db *memory.DB, port string) {
	li, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Server listening on %s\n", port)

	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go db.Handle(conn)
	}

}
