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
		m.respondTernary(conn, ok, v, "not found")
	case "DEL":
		ok := m.delete(key)
		m.respondTernary(conn, ok, "OK", "not found")

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
		go m.closeConn(conn)
	}

	m.Quit <- true
}

func (m *MemoryDB) closeConn(conn net.Conn) {
	err := conn.Close()
	m.handleErr(err)

	m.deregister(conn)
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

func (m *MemoryDB) delete(k string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.records[k]; !ok {
		return false
	}

	delete(m.records, k)

	return true
}

func (m *MemoryDB) handleErr(err error) {
	if err != nil {
		m.respondAll("internal error")
		log.Fatalf("[FATAL] %s", err)
	}
}
