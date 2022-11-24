package cheat

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func HandleCommand(command *packet.CommandRequest, proxy Proxy) (*packet.CommandRequest, bool) {
	log.Info().Msgf("Got command %v", command.CommandLine)
	if strings.HasPrefix(command.CommandLine, "/gdbs") {
		// serverMessage := packet.Text{TextType: packet.TextTypeSystem, Message: "This command is not implemented yet"}
		// _ = proxy.ClientConn.WritePacket(&serverMessage)
		transaction := packet.InventoryTransaction{
			Actions: []protocol.InventoryAction{
				{
					SourceType:    protocol.InventoryActionSourceWorld,
					InventorySlot: 0,
					NewItem: protocol.ItemInstance{
						StackNetworkID: 1,
						Stack: protocol.ItemStack{
							BlockRuntimeID: 57,
							Count:          64,
							HasNetworkID:   false,
						},
					},
				},
			},
			TransactionData: &protocol.NormalTransactionData{},
		}
		_ = proxy.ClientConn.WritePacket(&transaction)
		return command, false
	}

	return command, true
}
