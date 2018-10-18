module github.com/justanotherorganization/justanotherbotkit/transport

require (
	github.com/bwmarrin/discordgo v0.18.0
	github.com/gogo/protobuf v1.1.1
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/justanotherorganization/justanotherbotkit/proto v0.0.0
	github.com/nlopes/slack v0.4.0
	github.com/pkg/errors v0.8.0 // indirect
	golang.org/x/crypto v0.0.0-20181001203147-e3636079e1a4 // indirect
)

replace github.com/justanotherorganization/justanotherbotkit/proto => ../proto
