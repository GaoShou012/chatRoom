package frontier_ws

type Message struct {
	RuledOut map[*Conn]bool
	Body []byte
}
