package send_message

import "github.com/gorilla/websocket"

type Chhain struct {
	Height int
}

func (c Chhain) Handle(p []byte, con *websocket.Conn, m int) {
	if string(p) == "hi" {
		con.WriteMessage(m, []byte("hello"))
	}
}

type Result struct {
	S string
}

type Res string

func (r *Res) Replay(st *Result, reply *string) error {

	s := st.S
	s = "hi " + s
	*reply = s
	return nil
}
