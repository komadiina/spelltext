package server

import (
	"sync"
	"time"

	pb "github.com/komadiina/spelltext/proto/chat"
)

type Hub struct {
	mutex      sync.RWMutex
	subs       map[int]chan *pb.ChatMessage
	next       int
	bufferSize int
}

func NewHub(buffer int) *Hub {
	return &Hub{subs: make(map[int]chan *pb.ChatMessage), bufferSize: buffer}
}

func (hub *Hub) Add() (id int, ch <-chan *pb.ChatMessage) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	id = hub.next
	c := make(chan *pb.ChatMessage, hub.bufferSize)
	hub.subs[id] = c
	hub.next++
	return id, c
}

// removes subscriber by id
func (hub *Hub) Remove(id int) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()
	if c, ok := hub.subs[id]; ok {
		delete(hub.subs, id)
		close(c)
	}
}

func (hub *Hub) Broadcast(sender, text string) *pb.ChatMessage {
	msg := &pb.ChatMessage{
		Sender:    sender,
		Message:   text,
		Timestamp: uint64(time.Now().Unix()),
	}

	hub.mutex.RLock()
	for _, ch := range hub.subs {
		// ??
		select {
		case ch <- msg:
		default:
			continue
		}
	}
	hub.mutex.RUnlock()

	return msg
}
