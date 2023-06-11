package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	li, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handle(conn)

	}

}

func handle(conn net.Conn) {
	defer conn.Close()
	sc := bufio.NewScanner(conn)

	for sc.Scan() {
		fmt.Println(sc.Text())
	}
}
