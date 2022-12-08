package cheat

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/marcuswu/mcproxy/world/chunk"
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleCommand(command *packet.CommandRequest) (*packet.CommandRequest, bool) {
	log.Info().Msgf("Got command %v", command.CommandLine)
	if strings.HasPrefix(command.CommandLine, "/update") {
		serverMessage := packet.SubChunkRequest{
			Dimension: 0, 
			Position:protocol.SubChunkPos{0, 0, 0}, 
			Offsets: []protocol.SubChunkOffset{protocol.SubChunkOffset{0, 0, 0}},
		}
		_ = proxy.ClientConn.WritePacket(&serverMessage)
		return command, false
	}
	if strings.HasPrefix(command.CommandLine, "/identify") {
		commandParts := strings.Split(command.CommandLine, " ")
		if len(commandParts) < 4 {
			return command, false
		}
		X, err := strconv.Atoi(commandParts[1])
		if err != nil {
			log.Debug().Str("coord", commandParts[1]).Msg("Could not parse X coordinate")
			return command, false
		}
		Y, err := strconv.Atoi(commandParts[2])
		if err != nil {
			log.Debug().Str("coord", commandParts[2]).Msg("Could not parse Y coordinate")
			return command, false
		}
		Z, err := strconv.Atoi(commandParts[3])
		if err != nil {
			log.Debug().Str("coord", commandParts[3]).Msg("Could not parse Z coordinate")
			return command, false
		}
		chunkX := int32(X) >> 4
		chunkY := int16(Y)
		chunkZ := int32(Z) >> 4
		xOffset := uint8(int32(X) - (chunkX << 4))
		zOffset := uint8(int32(Z) - (chunkZ << 4))
		log.Debug().Int32("X", chunkX).Int32("Z", chunkZ).Msg("chunk position")
		log.Debug().Uint8("X", xOffset).Uint8("Z", zOffset).Msg("chunk offset")

		c, ok := proxy.Chunks[protocol.ChunkPos{chunkX, chunkZ}]
		if !ok {
			log.Debug().Msg("Could not find chunk")
			return command, false
		}
		yIdx := c.SubIndex(chunkY)
		yOffset := chunkY - c.SubY(yIdx)
		log.Debug().Int16("Y", chunkY).Int16("index", c.SubIndex(chunkY)).Int16("yOffset", yOffset).Msg("y sub chunk index")
		id := c.Block(xOffset, chunkY, zOffset, 0)
		name, _, found:= chunk.RuntimeIDToState(id)
		if !found {
			log.Debug().Msgf("Could not find id %s", id)
		}
		serverMessage := packet.Text{TextType: packet.TextTypeSystem, Message: fmt.Sprintf("Block at %d, %d, %d has id %d (%s)", X, chunkY, Z, id, name)}
		_ = proxy.ClientConn.WritePacket(&serverMessage)
		return command, false
	}
	if strings.HasPrefix(command.CommandLine, "/find") {
		if proxy.PlayerPos == nil {
			log.Debug().Msg("PlayerPos is unknown")
			return command, false
		}
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
		lookingFor, found := chunk.StateToRuntimeID(fmt.Sprintf("minecraft:%s", blockName), nil)
		if !found {
			log.Debug().Msg("Could not find block runtime id")
			return command, false
		}
		c, ok := proxy.Chunks[protocol.ChunkPos{chunkX, chunkZ}]
		if !ok {
			log.Debug().Int32("chunkX", chunkX).Int32("chunkZ", chunkZ).Msg("Could not find chunk")
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
						logFoundLocation(chunkX, chunkZ, x, z, Y, -dist)
					}
				}
			}
			if found {
				log.Debug().Msg("Found matching block, exiting find")
				return command, false
			}
		}
		log.Debug().Msg("Could not find matching block")
		serverMessage := packet.Text{TextType: packet.TextTypeSystem, Message: "Could not find matching block"}
		_ = proxy.ClientConn.WritePacket(&serverMessage)
		return command, false
	}

	return command, true
}
