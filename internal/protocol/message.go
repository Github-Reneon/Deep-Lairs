package protocol

// types
var (
	TYPE_MESSAGE = "message"
)

var (
	PORT = ":3000"
)

// format types for sprintf
var (
	LOOK               = "<b>%s</b> %s<br><img src=\"/img/%s\" class=\"img-fluid w-100 m-4 rounded-lg starting:opacity-0 opacity-100 transition-all duration-700 delay-200 ease-in-out \" />"
	LOOK_NO_IMAGE      = "<b>%s</b> %s"
	SAY                = "%s says: \"%s\""
	I_DONT_KNOW_HOW_TO = "I don't know how to %s"
	SHOUT              = "%s shouts: <b>\"%s\"</b>"
	LOL                = "%s laughs out loud! 😂"
	YOU_ARE_IN         = "You are in %s<br><img src=\"/img/%s\" class=\"img-fluid w-100 m-4 rounded-lg starting:opacity-0 opacity-100 transition-all duration-700 delay-200 ease-in-out \" />"
	WHISPER            = "%s whispers to %s: \"<span class=\"italic\">%s</span>\""
	WHISPER_FAIL       = "There's no adventurer %s here."
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// joining messages
var (
	JOINING_MESSAGE   = "%s has joined the %s."
	STUMBLES_IN       = "%s stumbles in."
	CREEPS_IN         = "%s creeps in."
	ENTERS_CAUTIOUSLY = "%s enters cautiously."
)
