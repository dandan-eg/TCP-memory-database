package db

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"testing"
)

// command is a utility function used to simulate command execution on a MemoryDB in the context of unit tests.
// It takes a testing.T object, a MemoryDB instance, and a command string as input.
// It returns the response message from the command execution.
func command(t *testing.T, db *MemoryDB, cmd string) string {
	c1, c2 := net.Pipe()

	// Register the first connection (c1) with the MemoryDB for receiving responses
	// and writing them in the second connection(c2)
	db.register(c1)

	go func() {
		action, key, value := inputs(cmd)
		db.dispatch(c1, action, key, value)

		err := c1.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Read the response from the second connection (c2) wrote in by the first connection (c1)
	bs, errRead := io.ReadAll(c2)

	if errRead != nil {
		t.Fatal(errRead)
	}

	if err := c2.Close(); err != nil {
		t.Fatal(err)
	}

	if panicErr := recover(); panicErr != nil {
		t.Fatal(panicErr)
	}

	msg := string(bs)

	if strings.HasSuffix(msg, "\r\n") {
		msg = strings.Replace(msg, "\r\n", "", -1)
	}

	return msg
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
		log.Printf("Test %d : %s", i+1, test.cmd)

		msg := command(t, db, test.cmd)

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
		log.Printf("Test %d : key=\"%s\" insert=%v", i+1, test.key, test.insert)

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
		key      string
		expected string
	}{
		{"name1", "john"},
		{"name2", "maria"},
		{"404", "not found"},
	}

	for i, test := range tests {
		log.Printf("Test %d : \"%s\"", i+1, test.key)

		msg := command(t, db, "GET "+test.key)

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
		key, expected string
	}{
		{key: "key1", expected: "OK"},
		{key: "key2", expected: "OK"},
		{key: "key404", expected: "not found"},
	}

	for i, test := range tests {
		//test each scenario
		t.Logf("Test %d : \"%s\" \"%s\"", i+1, test.key, test.expected)

		msg := command(t, db, "DEL "+test.key)

		if msg != test.expected {
			t.Fatalf("expected message=\"%s\", got=\"%s\"", test.expected, msg)
		}

		if _, ok := db.records[test.key]; ok {
			t.Fatalf("key=%s should be deleted, got ok=true", test.key)
		}

	}

}

func TestMemoryDB_Save_Unit(t *testing.T) {

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
