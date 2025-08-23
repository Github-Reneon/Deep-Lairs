package protocol

// types
var (
	TYPE_MESSAGE = "message"
)

// format types for sprintf
var (
	LOOK               = "<b>%s</b><br><img src=\"/img/%s\" class=\"img-fluid w-70 rounded-lg\" />"
	LOOK_NO_IMAGE      = "<b>%s</b>"
	SAY                = "%s says: \"%s\""
	I_DONT_KNOW_HOW_TO = "I don't know how to %s"
	SHOUT              = "%s shouts: <b>\"%s\"</b>"
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}
