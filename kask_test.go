package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func SetupDB(t *testing.T) (http.Handler, Store) {
	hostname     := Getenv("CASSANDRA_HOST", "localhost")
	port         := Getenv("CASSANDRA_PORT", "9042")
	keyspace     := Getenv("CASSANDRA_KEYSPACE", "kask_test_keyspace")
	table        := Getenv("CASSANDRA_TABLE", "test_table")
	portNum, err := strconv.Atoi(port)

	logger := NewLogger("kask_testing")

	store, err := NewCassandraStore(hostname, portNum, keyspace, table)
	if err != nil {
		logger.Error("Error connecting to Cassandra: %s", err)
		log.Fatal("Error connecting to Cassandra: ", err)
	}

	handler := HttpHandler{store, &logger}
	handle  := ParseKeyMiddleware("/sessions/v1/", http.HandlerFunc(handler.Dispatch))

	return handle, store
}

func TestGETSuccess (t *testing.T) {
	handle, store   := SetupDB(t)
	ts              := httptest.NewServer(handle)
	defer ts.Close()

	req, err        := http.NewRequest("GET", ts.URL+"/foo", nil)
	expected        := "bar"

	store.Set("foo", []byte(expected))

	if err != nil {
		t.Error(err)
	}

	resp, er := http.DefaultClient.Do(req)

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code got %v expected %v", resp.StatusCode, http.StatusOK)
	}

	if er != nil {
		t.Error(er)
	}

	if string(body) != expected {
		t.Errorf("handler returned unexpected body: got %v expected %v", string(body), expected)
	}
}

func TestGETNotFound (t *testing.T) {
	handle, _   := SetupDB(t)
	ts          := httptest.NewServer(handle)
	defer ts.Close()

	req, err    := http.NewRequest("GET", ts.URL+"/cat", nil)

	if err != nil {
		t.Error(err)
	}

	resp, er := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("handler returned wrong status code got %v expected %v", resp.StatusCode, http.StatusNotFound)
	}

	if er != nil {
		t.Error(er)
	}
}

func TestPOST (t *testing.T) {
	handle, store   := SetupDB(t)
	ts              := httptest.NewServer(handle)
	defer ts.Close()

	expected        := []byte("woof")
	req, err        := http.NewRequest("POST", ts.URL+"/dog",  bytes.NewReader(expected))

	if err != nil {
		t.Error(err)
	}

	resp, er := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("handler returned wrong status code got %v expected %v", resp.StatusCode, http.StatusCreated)
	}

	if er != nil {
		t.Error(er)
	}

	if value, _ := store.Get("dog"); !bytes.Equal(value, expected) {
		t.Errorf("store returned value of %s expected %s", value, expected)
	}
}

func TestDELETE (t *testing.T) {
	handle, store   := SetupDB(t)
	ts              := httptest.NewServer(handle)
	defer ts.Close()

	req, err        := http.NewRequest("DELETE", ts.URL+"/duck", nil)
	store.Set("duck", []byte("quack"))

	if err != nil {
		t.Error(err)
	}

	resp, er := http.DefaultClient.Do(req)

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("handler returned wrong status code got %v expected %v", resp.StatusCode, http.StatusNoContent)
	}

	if er != nil {
		t.Error(er)
	}

	if value, _ := store.Get("duck"); len(value) > 0 {
		t.Errorf("DELETE did not remove key: duck and value: %s", value)
	}
}