package protocol

// types
var (
	TYPE_MESSAGE = "message"
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
