package db

import (
	"fmt"
	"io"
	"strings"
)

func (m *MemoryDB) respondTernary(conn io.ReadWriteCloser, ok bool, strue, sfalse string) {
	if ok {
		m.respond(conn, strue)
	} else {
		m.respond(conn, sfalse)
	}
}

func (m *MemoryDB) respond(conn io.ReadWriteCloser, msg string) {
	_, err := fmt.Fprintf(conn, fmt.Sprintf("%s\r\n", msg))
	if err != nil {
		m.internalError(err)

	}
}

func (m *MemoryDB) respondAll(msg string) {

	for _, conn := range m.conns {
		m.respond(conn, msg)
	}

}

func (m *MemoryDB) String() string {
	sb := strings.Builder{}

	sb.WriteString("Memory Database methods :\r\n\n")
	sb.WriteString("GET <key> : Retrieve the value associated with a key from the database.\r\n")
	sb.WriteString("SET <key> <value> : Set a key-value pair in the database.\r\n")
	sb.WriteString("EXIT : Exit the database application.\r\n")
	sb.WriteString("CLOSE : Close the database connection.\r\n")

	return sb.String()
}
