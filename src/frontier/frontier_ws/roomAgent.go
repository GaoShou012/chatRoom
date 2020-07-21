package frontier_ws

import "sync"

type RoomCode string

type RoomAgent struct {
	mutex       sync.Mutex
	rooms       map[RoomCode]*Room
	connections map[*Conn]RoomCode
}

func (ra *RoomAgent) Init() {
	ra.rooms = make(map[RoomCode]*Room)
	ra.connections = make(map[*Conn]RoomCode)
}

func (ra *RoomAgent) Join(roomCode RoomCode, conn *Conn) error {
	ra.mutex.Lock()
	defer ra.mutex.Unlock()

	r, ok := ra.rooms[roomCode]
	if ok {
		return r.Join(conn)
	}
	r = &Room{}
	r.Init()
	ra.rooms[roomCode] = r
	return r.Join(conn)
}
func (ra *RoomAgent) Leave(conn *Conn) {
	ra.mutex.Lock()
	defer ra.mutex.Unlock()

	roomCode,ok := ra.connections[conn]
	if !ok {
		return
	}
	r,ok := ra.rooms[roomCode]
	if !ok {
		return
	}

}

func (ra *RoomAgent) Broadcast(roomCode string, message Message) {

}
