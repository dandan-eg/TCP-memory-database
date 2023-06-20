package memory

import (
	"io"
	"log"
)

func (m *DB) dispatch(conn io.ReadWriteCloser, action, key, value string) {
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

	case "SAVE":
		m.save()
		m.respond(conn, "SAVED")

	case "EXIT":
		m.respond(conn, "exited")
		m.exit(conn)

	case "CLOSE":
		m.respondAll("db closed")
		m.close()

	default:
		m.respond(conn, "bad request")
	}

}

func (m *DB) close() {

	for _, conn := range m.conns {
		if err := conn.Close(); err != nil {
			m.internalError(err)
		}
	}

	m.conns = nil

	close(m.quit)
}

func (m *DB) exit(conn io.ReadWriteCloser) {
	if err := conn.Close(); err != nil {
		m.internalError(err)
	}
	m.deregister(conn)
}

func (m *DB) set(k, v string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[k] = v

}

func (m *DB) get(k string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	v, ok := m.data[k]
	return v, ok

}

func (m *DB) delete(k string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[k]; !ok {
		return false
	}

	delete(m.data, k)

	return true
}

func (m *DB) save() {
	err := m.Saver.Save(m.data)
	if err != nil {
		m.internalError(err)
	}

}

func (m *DB) internalError(err error) {

	m.respondAll("internal error")
	log.Fatalf("[FATAL] %s", err)
}
