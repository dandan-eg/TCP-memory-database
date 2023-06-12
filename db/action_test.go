package db

import (
	"bytes"
	"strings"
	"testing"
)

// Close is an implementation of the Close method of the io.Closer interface.
func (t *TestReadWriteCloser) Close() error {
	// This implementation is empty since we don't need to perform any close operation in this case.
	return nil
}

// String is a custom method added to TestReadWriteCloser to provide a string representation of its contents.
// It replaces the first occurrence of "\r\n" with an empty string to remove any line breaks.
func (t *TestReadWriteCloser) String() string {
	s := t.Buffer.String()

	return strings.Replace(s, "\r\n", "", 1)
}

// newConn creates a new TestReadWriteCloser and uses it to simulate a connection to the MemoryDB.
// It writes the given command to the TestReadWriteCloser and passes it to the MemoryDB's Handle method.
// It returns the TestReadWriteCloser for further inspection.
func newConn(db *MemoryDB, cmd string) *TestReadWriteCloser {
	rwc := &TestReadWriteCloser{
		Buffer: bytes.NewBuffer(nil),
	}

	rwc.WriteString(cmd)
	db.Handle(rwc)

	return rwc

}

func TestMemoryDB_SET_Unit(t *testing.T) {
	//setup
	db := New()

	var before, after int

	//add new value

	before = len(db.records)
	db.set("key", "value")
	after = len(db.records)

	if before+1 != after {
		t.Fatal("length after should be one more than before")
	}

	//replace the value
	after = len(db.records)
	db.set("key", "another value")
	before = len(db.records)

	if before != after {
		t.Fatal("replace should not change the length")
	}

}

func TestMemoryDB_SET_Integration(t *testing.T) {
	//setup
	db := New()

	//scenarios
	tests := []struct {
		cmd      string
		expected string
	}{
		{cmd: "SET name john", expected: "OK"},
		{cmd: "SET name maria", expected: "OK"},
		{cmd: "SET ", expected: "bad request"},
		{cmd: "notfound ", expected: "bad request"},
	}

	for i, test := range tests {
		//run each scenario
		t.Logf("Test %d : %s", i+1, test.cmd)

		conn := newConn(db, test.cmd)

		msg := conn.String()

		if msg != test.expected {
			t.Fatalf("msg expected=%s, got=%s", test.expected, msg)
		}

		if msg == "bad request" {
			//if it is a bad request as expected, no need to go further
			return
		}

		_, ok := db.records["name"]

		if !ok {
			t.Error("key \"names\" should exists, got=false")
			return
		}

	}
}

func TestMemoryDB_GET_Unit(t *testing.T) {

	//setup
	db := New()
	db.records = map[string]string{
		"1": "john",
		"2": "daniel",
	}

	//scenario
	tests := []struct {
		key    string
		value  string
		insert bool
	}{
		{key: "1", value: "john", insert: true},
		{key: "2", value: "daniel", insert: true},
		{key: "not insert", value: "daniel", insert: false},
		{key: "not insert", value: "john", insert: false},
	}

	for i, test := range tests {
		//run each scenario
		t.Logf("Test %d : key=\"%s\" insert=%v", i+1, test.key, test.insert)

		returned, ok := db.get(test.key)

		if ok != test.insert {
			t.Fatal("Insertion state does not match the ok value")
		}

		if test.insert == false {
			//if no inserted no need to go further
			return
		}

		if returned != test.value {
			t.Fatalf("the value returned is not equal to the expected value, expected=\"%s\" got=\"%s\"", test.value, returned)
		}

	}

}

func TestMemoryDB_GET_Integration(t *testing.T) {
	//setup
	db := New()
	db.records = map[string]string{
		"name1": "john",
		"name2": "maria",
	}

	tests := []struct {
		cmd      string
		expected string
	}{
		{"GET name1", "john"},
		{"GET name2", "maria"},
		{"GET 404", "not found"},
	}

	for i, test := range tests {
		t.Logf("Test %d : \"%s\"", i+1, test.cmd)

		conn := newConn(db, test.cmd)

		msg := conn.String()

		if msg != test.expected {
			t.Errorf("msg expected=\"%s\", got=\"%s\"", test.expected, msg)
		}

	}

}

func TestMemoryDB_Delete_Unit(t *testing.T) {
	//setup
	db := New()
	db.records = map[string]string{
		"key": "value",
	}

	var deleted bool

	//delete existing key
	deleted = db.delete("key")
	if !deleted {
		t.Error("existing key should be deleted, got delete=false")
	}

	//delete non-existing key
	deleted = db.delete("not existing")
	if deleted {
		t.Error("existing key should not be deleted, got delete=true")
	}
}

func TestMemoryDB_Delete_Integration(t *testing.T) {
	//setup
	db := New()
	db.records = map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	//scenarios
	tests := []struct {
		cmd, expected string
	}{
		{cmd: "DEL key1", expected: "OK"},
		{cmd: "DEL key2", expected: "OK"},
		{cmd: "DEL key404", expected: "not found"},
	}

	for i, test := range tests {
		//test each scenario
		t.Logf("Test %d : \"%s\" \"%s\"", i+1, test.cmd, test.expected)

		conn := newConn(db, test.cmd)

		msg := conn.String()

		if msg != test.expected {
			t.Fatalf("expected message=\"%s\", got=\"%s\"", test.expected, msg)
		}

		key := strings.Replace(test.cmd, "DEL ", "", 1)
		if _, ok := db.records[key]; ok {
			t.Fatalf("key=%s should be deleted, got ok=true", key)
		}

	}

}

func TestMemoryDB_Save_Unit(t *testing.T) {

}

type TestReadWriteCloser struct {
	*bytes.Buffer
}
