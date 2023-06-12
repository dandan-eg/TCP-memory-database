package db

import (
	"log"
	"net"
)

func (m *MemoryDB) dispatch(conn net.Conn, action, key, value string) {
	switch action {
	case "SET":
		if value == "" || key == "" {
			m.respond(conn, "bad request")
			break
		}

		m.set(key, value)
		m.respond(conn, "OK")
	case "GET":

		v, ok := m.get(key)
		if !ok {
			m.respond(conn, "not found")
		} else {
			m.respond(conn, v)
		}

	case "EXIT":
		m.respond(conn, "exited")
		m.closeConn(conn)

	case "CLOSE":
		m.respondAll("db closed")
		m.close()
	default:
		m.respond(conn, "bad request")
	}

}

func (m *MemoryDB) close() {
	for _, conn := range m.conns {
		m.mu.Lock()
		m.closeConn(conn)
		m.mu.Unlock()

	}

	m.Quit <- true
}

func (m *MemoryDB) closeConn(conn net.Conn) {
	m.deregister(conn)
	err := conn.Close()
	m.handleErr(err)

}

func (m *MemoryDB) set(k, v string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.records[k] = v

}

func (m *MemoryDB) get(k string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.records[k]
	return v, ok

}

func (m *MemoryDB) handleErr(err error) {
	if err != nil {
		log.Println(err)
		m.respondAll("internal error")
		m.close()
	}
}
