package main

//
// support for Raft and kvraft to save persistent
// Raft state (log &c) and k/v server snapshots.
//
// a “real” implementation would do this by writing Raft's persistent state
// to disk each time it changes, and reading the latest saved state from disk
// when restarting after a reboot.
// this implementation won't use the disk; instead, it will save and restore
// persistent state from a Persister object. Whoever calls Raft.Make()
// supplies a Persister that initially holds Raft's most recently persisted state (if any).
// Raft should initialize its state from that Persister, and should use it to
// save its persistent state each time the state changes.
//
// we will use the original persister.go to test your code for grading.
// so, while you can modify this code to help you debug, please
// test with the original before submitting.
//

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)
var lock 		sync.Mutex
var mu        	sync.Mutex

// Persister structure
type Persister struct {
	raftstate []byte
	snapshot  []byte
}

// Marshal is a function that marshals the object into an
// io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}


// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	//bytes_v, err := json.Marshal(v)
	//if err != nil {
	//	return err
	//}
	//return ioutil.WriteFile(path, bytes_v, 0777)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = os.Chmod(path, 0777)
	if err != nil {
		log.Fatal(err)
	}
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	//file, _ := ioutil.ReadFile(path)
	//return json.Unmarshal([]byte(file), v)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return Unmarshal(f, v)
}

// MakePersister create a Persister instance
func MakePersister() *Persister {
	return &Persister{}
}

// Copy a Persister
func (ps *Persister) Copy() *Persister {
	mu.Lock()
	defer mu.Unlock()
	np := MakePersister()
	np.raftstate = ps.raftstate
	np.snapshot = ps.snapshot
	return np
}

// SaveRaftState save data in a list of byte
func (ps *Persister) SaveRaftState(data []byte) {
	mu.Lock()
	defer mu.Unlock()
	ps.raftstate = data
}

// ReadRaftState return a list of byte
func (ps *Persister) ReadRaftState() []byte {
	mu.Lock()
	defer mu.Unlock()
	return ps.raftstate
}

// RaftStateSize return state size in int
func (ps *Persister) RaftStateSize() int {
	mu.Lock()
	defer mu.Unlock()
	return len(ps.raftstate)
}

// SaveSnapshot save a snapshot data in a list of byte
func (ps *Persister) SaveSnapshot(snapshot []byte) {
	mu.Lock()
	defer mu.Unlock()
	ps.snapshot = snapshot
}

// ReadSnapshot read data in list of byte
func (ps *Persister) ReadSnapshot() []byte {
	mu.Lock()
	defer mu.Unlock()
	return ps.snapshot
}

// SnapshotSize return the value in int
func (ps *Persister) SnapshotSize() int {
	mu.Lock()
	defer mu.Unlock()
	return len(ps.snapshot)
}

