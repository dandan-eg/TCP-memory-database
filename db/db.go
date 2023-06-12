package db

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"
)

type MemoryDB struct {
	Quit    chan bool
	records map[string]string
	conns   []io.ReadWriteCloser
	mu      *sync.RWMutex
	writer  io.Writer
}

func New() *MemoryDB {

	return &MemoryDB{
		Quit:    make(chan bool),
		records: make(map[string]string),
		conns:   make([]io.ReadWriteCloser, 0, 0),
		mu:      &sync.RWMutex{},
	}
}

func (m *MemoryDB) register(conn io.ReadWriteCloser) {

	m.conns = append(m.conns, conn)
	//m.respond(conn, m.String())
}

func (m *MemoryDB) deregister(conn io.ReadWriteCloser) {
	last := len(m.conns) - 1
	m.mu.Lock()

	for i, c := range m.conns {
		if c == conn {

			// Swap the current connection with the last connection
			m.conns[last], m.conns[i] = m.conns[i], m.conns[last]

			// Remove the last connection from the slice
			m.conns = m.conns[:last]
			return
		}
	}

	m.mu.Unlock()
}

func (m *MemoryDB) Handle(conn io.ReadWriteCloser) {

	m.register(conn)

	sc := bufio.NewScanner(conn)
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		request := sc.Text()

		log.Printf("incomming request : %q\n", sc.Bytes())

		action, key, value := inputs(request)

		m.dispatch(conn, action, key, value)
	}
}

func inputs(txt string) (string, string, string) {
	trim := strings.TrimSpace(txt)

	if trim == "" {
		return "", "", ""
	}

	if trim == "EXIT" {
		return "EXIT", "", ""
	}

	if trim == "CLOSE" {
		return "CLOSE", "", ""
	}

	fs := strings.Fields(trim)
	var (
		action string
		key    string
		value  string
	)

	if len(fs) < 2 {
		return "", "", ""
	} else {
		action = fs[0]
		key = fs[1]
	}

	if len(fs) == 3 {
		value = fs[2]
	}

	return action, key, value
}
