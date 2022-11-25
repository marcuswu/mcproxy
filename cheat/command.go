package cheat

import (
	"fmt"
	"strings"

	"github.com/marcuswu/mcproxy/latestmappings"
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleCommand(command *packet.CommandRequest) (*packet.CommandRequest, bool) {
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
		log.Info().Msgf("sending diamond block inventory transaction")
		err := proxy.ServerConn.WritePacket(&transaction)
		if err != nil {
			log.Info().Msgf("error sending inventory transaction %v", err)
		}

		return command, false
	}
	if strings.HasPrefix(command.CommandLine, "/find") {
		commandParts := strings.Split(command.CommandLine, " ")
		blockName := commandParts[1]
		// Loop through chunks looking for diamond ore nearest the player
		translatePosition := func(chunkX int32, chunkZ int32, x uint8, z uint8, Y int16, dist int16) (int32, int32, int32) {
			return (chunkX << 4) + int32(x), int32(Y + dist), (chunkZ << 4) + int32(z)
		}
		logFoundLocation := func(chunkX int32, chunkZ int32, x uint8, z uint8, Y int16, dist int16) {
			foundX, foundY, foundZ := translatePosition(chunkX, chunkZ, x, z, Y, dist)
			serverMessage := packet.Text{TextType: packet.TextTypeSystem, Message: fmt.Sprintf("Found %s at %d, %d, %d", blockName, foundX, foundY, foundZ)}
			_ = proxy.ClientConn.WritePacket(&serverMessage)
		}
		chunkX := int32(proxy.PlayerPos.Position.X()) >> 4
		chunkZ := int32(proxy.PlayerPos.Position.Z()) >> 4
		Y := int16(proxy.PlayerPos.Position.Y())
		lookingFor, found := latestmappings.StateToRuntimeID(fmt.Sprintf("minecraft:%s", blockName), nil)
		if !found {
			return command, false
		}
		c, ok := proxy.Chunks[protocol.ChunkPos{chunkX, chunkZ}]
		if !ok {
			return command, false
		}
		found = false
		for dist := int16(0); int(Y+dist) < c.Range().Max() && int(Y-dist) > c.Range().Min(); dist++ {
			for x := uint8(0); x < 16; x++ {
				for z := uint8(0); z < 16; z++ {
					if lookingFor == c.Block(x, Y+dist, z, 0) {
						found = true
						logFoundLocation(chunkX, chunkZ, x, z, Y, dist)
					}
					if lookingFor == c.Block(x, Y-dist, z, 0) {
						found = true
						logFoundLocation(chunkX, chunkZ, x, z, Y, dist)
					}
				}
			}
			if found {
				return command, false
			}
		}
	}

	return command, true
}
