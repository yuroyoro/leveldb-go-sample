package main

import (
	"log"
)

type Backend struct {
	db     *LevelDB
	putch  chan putRequest
	delch  chan delRequest
	quitch chan bool
}

type putRequest struct {
	key string
	val string
}

type delRequest struct {
	key string
}

func NewBackend(dbname string) (backend *Backend, err error) {
	db, err := OpenLevelDB(dbname)

	if err != nil {
		return
	}

	log.Printf("LevelDB opened : name -> %s", dbname)

	backend = &Backend{
		db,
		make(chan putRequest),
		make(chan delRequest),
		make(chan bool),
	}

	return
}

func (backend *Backend) Start() {
	go func() {
		for {
			select {
			case putreq := <-backend.putch:
				backend.db.Put(putreq.key, putreq.val)
				log.Printf("Backend.Put(%s, %v)\n", putreq.key, putreq.val)
			case delreq := <-backend.delch:
				backend.db.Delete(delreq.key)
				log.Printf("Backend.Delete(%s)\n", delreq.key)
			case <-backend.quitch:
				log.Printf("Backend stoped")
				return
			}
		}
	}()

	log.Printf("Backend started")
}

func (backend *Backend) Shutdown() {
	backend.quitch <- true
}

func (backend *Backend) Get(key string) (val string, err error) {
	return backend.db.Get(key)
}

func (backend *Backend) Put(key, val string) {
	backend.putch <- putRequest{key, val}

}

func (backend *Backend) Delete(key string) {
	backend.delch <- delRequest{key}
}
