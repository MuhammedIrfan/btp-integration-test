package setup

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type HexInt string
type HexBytes string

type EventNotification struct {
	Hash   HexBytes `json:"hash"`
	Height HexInt   `json:"height"`
	Index  HexInt   `json:"index"`
	Events []HexInt `json:"events,omitempty"`
}

type BlockNotification struct {
	Hash    HexBytes     `json:"hash"`
	Height  HexInt       `json:"height"`
	Indexes [][]HexInt   `json:"indexes,omitempty"`
	Events  [][][]HexInt `json:"events,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type WSClient struct {
	Conn  *websocket.Conn
	Chain Chain
	Mu    sync.Mutex
}

type outbound struct {
	Body string `json:"body"`
}

type Chain interface {
	Handle([]byte, *websocket.Conn, int)
}

func (c *WSClient) send(v interface{}) error {

	c.Mu.Lock()
	defer c.Mu.Unlock()

	return c.Conn.WriteJSON(v)
}

func (c *WSClient) BlockHandler() {

	data := &BlockNotification{

		Hash:   "hxb51a65420ce5199e538f21fc614eacf4234454fe",
		Height: "0x1",
		Indexes: [][]HexInt{
			[]HexInt{"0x0"},
		},
		Events: [][][]HexInt{
			[][]HexInt{
				[]HexInt{"0x0"},
			},
		},
	}
	for {
		// _, p, err := c.Conn.ReadMessage()

		// fmt.Println(string(p))
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }

		c.Conn.WriteJSON(data)

	}

}

func (c *WSClient) EventHandler() {

	data := &EventNotification{

		Hash:   "hxb51a65420ce5199e538f21fc614eacf4234454fe",
		Height: "0x10",
		Index:  HexInt(0),
		Events: []HexInt{
			HexInt(0),
		},
	}
	for {
		// messageType, _, err := c.conn.ReadMessage()
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }

		c.Conn.WriteJSON(data)

	}
}

func (c *WSClient) SendMessage(message string) {
	err := c.Conn.WriteMessage(1, []byte(message))
	if err != nil {
		handleError(c.Conn, err)
	}
}

func (c *WSClient) Block(ce echo.Context) error {

	ws, err := upgrader.Upgrade(ce.Response(), ce.Request(), nil)
	if err != nil {
		log.Println(err)
	}
	c.Conn = ws

	defer c.Conn.Close()

	c.BlockHandler()

	return nil
}
func (c *WSClient) Event(ce echo.Context) error {

	ws, err := upgrader.Upgrade(ce.Response(), ce.Request(), nil)
	if err != nil {
		handleError(c.Conn, err)
	}
	c.Conn = ws
	defer c.Conn.Close()
	c.EventHandler()

	return nil
}

func handleError(ws *websocket.Conn, err error) {
	log.Println("Error:", err)

	b, err := json.Marshal(&outbound{Body: err.Error()})
	if err != nil {
		log.Println("Error:", err)
	}

	err = ws.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Println("Error:", err)
	}
}

// data := &BlockNotification{

// 	Hash:   "hxb51a65420ce5199e538f21fc614eacf4234454fe",
// 	Height: "0x10",
// 	Indexes: [][]HexInt{
// 		[]HexInt{"0x0"},
// 	},
// 	Events: [][][]HexInt{
// 		[][]HexInt{
// 			[]HexInt{"0x0"},
// 		},
// 	},
// }
// data := &EventNotification{
// 	Hash:   "hxb51a65420ce5199e538f21fc614eacf4234454fe",
// 	Height: "0x10",
// 	Index:  "0x0",
// 	Events: []HexInt{
// 		"0x0",
// 	},
// }
