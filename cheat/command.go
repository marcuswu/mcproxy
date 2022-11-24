package cheat

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func HandleCommand(command *packet.CommandRequest, proxy Proxy) (*packet.CommandRequest, bool) {
	log.Info().Msgf("Got command %v", command.CommandLine)
	if strings.HasPrefix(command.CommandLine, "/give") {
		serverMessage := packet.Text{TextType: packet.TextTypeSystem, Message: "This command is not implemented yet"}
		_ = proxy.ClientConn.WritePacket(&serverMessage)
		return command, false
	}

	return command, true
}
