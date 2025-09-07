package protocol

// types
const (
	TYPE_MESSAGE = "message"
)

const (
	CLIENT_VERSION = "0.1.5"
	SERVER_VERSION = "0.1.0"
)

const USER_NOT_FOUND = "user_not_found"

const (
	SERVER_PORT = ":3000"
	CLIENT_PORT = ":3001"
)

// format types for sprintf
const (
	LOOK               = "<b>%s</b> %s<br><img src=\"/img/%s\" class=\"img-fluid w-100 m-4 rounded-lg starting:opacity-0 opacity-100 transition-all duration-700 delay-200 ease-in-out \" />"
	LOOK_NO_IMAGE      = "<b>%s</b> %s"
	SAY                = "<span class=\"inline-flex items-baseline\"> <img src=\"%s\" class=\"mx-1 size-7 self-center rounded-full\" /> <span>%s</span></span> says: \"%s\""
	I_DONT_KNOW_HOW_TO = "I don't know how to %s"
	SHOUT              = "%s shouts: <b>\"%s\"</b>"
	LOL                = "%s laughs out loud! 😂"
	YOU_ARE_IN         = "You are in %s<br><img src=\"/img/%s\" class=\"img-fluid w-100 m-4 rounded-lg starting:opacity-0 opacity-100 transition-all duration-700 delay-200 ease-in-out \" />"
	IMAGE              = "<img src=\"/img/%s\" class=\"img-fluid w-100 m-4 rounded-lg starting:opacity-0 opacity-100 transition-all duration-700 delay-200 ease-in-out \" />"
	WHISPER            = "%s whispers to %s: \"<span class=\"italic\">%s</span>\""
	WHISPER_FAIL       = "There's no adventurer %s here."
	DO                 = "%s starts to %s"
)

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// joining messages
const (
	JOINING_MESSAGE           = "%s has arrived."
	JOINING_STUMBLES_IN       = "%s stumbles in."
	JOINING_CREEPS_IN         = "%s creeps in."
	JOINING_ENTERS_CAUTIOUSLY = "%s arrives cautiously."
)

// Leaving messages
const (
	LEAVING_MESSAGE           = "%s has left, going %s."
	LEAVING_STUMBLES_OUT      = "%s stumbles out, going %s."
	LEAVING_CREEPS_OUT        = "%s creeps out, going %s."
	LEAVING_ENTERS_CAUTIOUSLY = "%s leaves cautiously, going %s."
)

const (
	STATE_TYPE_USER  = "user"
	STATE_TYPE_ENEMY = "enemy"
)
