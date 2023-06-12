package db

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func command(t *testing.T, db *MemoryDB, c string) string {
	c1, c2 := net.Pipe()
	db.register(c1)

	errCh := make(chan error)

	go func() {
		action, key, value := inputs(c)
		db.dispatch(c1, action, key, value)

		errCh <- c1.Close()
	}()

	bs, errRead := io.ReadAll(c2)

	if errRead != nil {
		t.Fatal(errRead)
	}

	if err := <-errCh; err != nil {
		t.Fatal(errRead)
	}

	return string(bs)
}

func TestMemoryDB_GET(t *testing.T) {

	db := New()
	cases := []struct {
		key    string
		value  string
		insert bool
	}{
		{key: "name1", value: "john", insert: true},
		{key: "name2", value: "daniel", insert: false},
		{key: "name3", value: "john", insert: false},
	}

	for _, s := range cases {
		if s.insert {
			db.records[s.key] = s.value

		}

		returned, ok := db.get(s.key)
		if !ok && s.insert {
			t.Errorf("key \"%s\" does not exists, expected ok=true got=false", s.key)
		}

		if returned != s.value && s.insert {
			t.Errorf("the value returned is not equal to the expected value, expected=\"%s\" got=\"%s\"", s.value, returned)
		}

	}
}

func TestMemoryDB_SET_Integration(t *testing.T) {
	db := New()

	tests := []struct {
		cmd      string
		expected string
	}{
		{cmd: "SET name john", expected: "OK"},
		{cmd: "SET name maria", expected: "OK"},
		{cmd: "SET ", expected: "bad request"},
		{cmd: "notfound ", expected: "bad request"},
	}

	for _, test := range tests {

		msg := command(t, db, test.cmd)

		if len(db.records) == 0 {
			t.Errorf("no records, expected=%d, got=0", len(db.records))
			return
		}

		if msg != test.expected {
			t.Errorf("msg expected=OK, got=%s", msg)
			return
		}

		_, ok := db.records["name"]

		if !ok {
			t.Error("key \"names\" should exists, got=false")
			return
		}
	}
}

func TestMemoryDB_GET_Integration(t *testing.T) {
	db := New()
	db.set("name1", "john")
	db.set("name2", "maria")

	tests := []struct {
		key      string
		expected string
	}{
		{"name1", "john"},
		{"name2", "maria"},
		{"404", "not found"},
	}

	for _, test := range tests {

		msg := command(t, db, "GET "+test.key)

		if msg != test.expected {
			t.Errorf("msg expected=\"john\", got=\"%s\"", msg)
		}

	}

}

func TestMemoryDB_EXIT_Integration(t *testing.T) {
	db := New()
	msg := command(t, db, "EXIT")

	if msg != "exited" {
		t.Errorf("")
	}

	if len(db.conns) > 0 {
		t.Errorf("")
	}

}

func TestMemoryDB_CLOSE_Integration(t *testing.T) {
	db := New()
	msg := command(t, db, "CLOSE")

	<-db.Quit
	fmt.Println(msg)

}
