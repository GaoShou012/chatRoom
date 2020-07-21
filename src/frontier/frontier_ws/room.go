package frontier_ws

import (
	"container/list"
	"github.com/src/frontier"
	"sync"
)

type Room struct {
	mutex sync.Mutex
	connections *list.List
	connectionsAnchor map[*Conn]*list.Element
	bucket chan Message
}

func (r *Room) Init(){
	r.connections = list.New()
	r.connectionsAnchor = make(map[*Conn]*list.Element)
	r.bucket = make(chan Message,10000)
	go func() {
		for{
			select{
			case message := <- r.bucket:
				body := message.Body
				ruledOut := message.RuledOut
				for ele := r.connections.Front(); ele != nil; ele = ele.Next() {
					conn := ele.Value.(*Conn)
					if _,ok := ruledOut[conn]; ok {continue}
					frontier.Sender.Push(conn,body)
				}
			}
		}
	}()
}

func (r *Room) Join(conn *Conn) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	anchor := r.connections.PushFront(conn)
	r.connectionsAnchor[conn] = anchor
	return nil
}

func (r *Room) Leave(conn *Conn) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	anchor,ok := r.connectionsAnchor[conn]
	if !ok {
		return
	}
	r.connections.Remove(anchor)
	delete(r.connectionsAnchor,conn)
}

func (r *Room) Broadcast(message Message){
	r.bucket <- message
}