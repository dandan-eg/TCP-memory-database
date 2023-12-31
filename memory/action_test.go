package memory

import (
	"bytes"
	"io"
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

// newConn creates a new TestReadWriteCloser and uses it to simulate a connection to the DB.
// It writes the given command to the TestReadWriteCloser and passes it to the DB's handle method.
// It returns the TestReadWriteCloser for further inspection.
func execute(db *DB, cmd string) string {
	rwc := &TestReadWriteCloser{
		Buffer: &bytes.Buffer{},
	}

	action, k, v := inputs(cmd)
	db.dispatch(rwc, action, k, v)

	return rwc.String()
}

func TestMemoryDB_SET_Unit(t *testing.T) {
	//setup
	db := NewDB(nil)

	var before, after int

	//add new value

	before = len(db.data)
	db.set("key", "value")
	after = len(db.data)

	if before+1 != after {
		t.Fatal("length after should be one more than before")
	}

	//replace the value
	after = len(db.data)
	db.set("key", "another value")
	before = len(db.data)

	if before != after {
		t.Fatal("replace should not change the length")
	}

}

func TestMemoryDB_SET_Integration(t *testing.T) {
	//setup
	db := NewDB(nil)

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

		msg := execute(db, test.cmd)

		if msg != test.expected {
			t.Fatalf("msg expected=%s, got=%s", test.expected, msg)
		}

		if msg == "bad request" {
			//if it is a bad request as expected, no need to go further
			return
		}

		_, ok := db.data["name"]

		if !ok {
			t.Error("key \"names\" should exists, got=false")
			return
		}

	}
}

func TestMemoryDB_GET_Unit(t *testing.T) {

	//setup
	db := NewDB(nil)
	db.data = map[string]string{
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
	db := NewDB(nil)
	db.data = map[string]string{
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

		msg := execute(db, test.cmd)

		if msg != test.expected {
			t.Errorf("msg expected=\"%s\", got=\"%s\"", test.expected, msg)
		}

	}

}

func TestMemoryDB_Delete_Unit(t *testing.T) {
	//setup
	db := NewDB(nil)
	db.data = map[string]string{
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
	db := NewDB(nil)
	db.data = map[string]string{
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

		msg := execute(db, test.cmd)

		if msg != test.expected {
			t.Fatalf("expected message=\"%s\", got=\"%s\"", test.expected, msg)
		}

		key := strings.Replace(test.cmd, "DEL ", "", 1)
		if _, ok := db.data[key]; ok {
			t.Fatalf("key=%s should be deleted, got ok=true", key)
		}

	}

}

func TestMemoryDB_Exit_Unit(t *testing.T) {
	db := NewDB(nil)
	conn := &TestReadWriteCloser{}

	db.conns = []io.ReadWriteCloser{
		conn,
	}

	db.exit(conn)

	if len(db.conns) > 0 {
		t.Errorf("conn pool should be equals 0, got=\"%d\"", len(db.conns))
	}
}

func TestMemoryDB_Exit_Integration(t *testing.T) {
	db := NewDB(nil)
	msg := execute(db, "EXIT")

	if msg != "exited" {
		t.Errorf("msg does not correspond, expected=\"exited\" got=%s", msg)
	}

	if len(db.conns) > 0 {
		t.Errorf("conn pool should be equals 0, got=%d", len(db.conns))
	}
}

func TestMemoryDB_Close_Unit(t *testing.T) {
	db := NewDB(nil)
	db.conns = []io.ReadWriteCloser{
		&TestReadWriteCloser{},
		&TestReadWriteCloser{},
		&TestReadWriteCloser{},
	}

	go db.close()

	<-db.quit

	if len(db.conns) > 0 {
		t.Errorf("conn pool should be equals 0, got=%d", len(db.conns))
	}
}

type TestReadWriteCloser struct {
	*bytes.Buffer
}
