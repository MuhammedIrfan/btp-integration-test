package icon

import (
	"github.com/gorilla/websocket"
)

type Chain struct {
	Height int
}

func (c Chain) Handle(p []byte, con *websocket.Conn, m int) {
	if string(p) == "hi" {
		con.WriteMessage(m, []byte("hello"))
	}
}

type Result struct {
	S string
}

type Res struct {}

func (r Res) Replay(st *Result, reply *string) error {

	s := st.S
	s = "hi " + s
	*reply = s
	return nil
}