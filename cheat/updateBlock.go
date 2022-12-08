package cheat

import (
	"github.com/rs/zerolog/log"
	"github.com/marcuswu/mcproxy/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func (proxy *Proxy) HandleUpdateBlock(ub *packet.UpdateBlock) (*packet.UpdateBlock, bool) {
	chunkX := ub.Position.X() >> 4
	chunkY := int16(ub.Position.Y())
	chunkZ := ub.Position.Z() >> 4
	xOffset := uint8(int32(ub.Position.X()) - (chunkX << 4))
	zOffset := uint8(int32(ub.Position.Z()) - (chunkZ << 4))

	c, ok := proxy.Chunks[protocol.ChunkPos{chunkX, chunkZ}]
	if !ok {
		log.Debug().Int32("chunkX", chunkX).Int32("chunkZ", chunkZ).Msg("UpdateBlock could not find chunk")
		return ub, true
	}

	name, _, _:= chunk.RuntimeIDToState(ub.NewBlockRuntimeID)
	oldID := c.Block(xOffset, chunkY, zOffset, 0)
	oldName, _, _:= chunk.RuntimeIDToState(oldID)

	if name == "minecraft:wool" || oldName == "minecraft:wool" {
	log.Debug().
	Int32("X", ub.Position.X()).
	Int32("Y", ub.Position.Y()).
	Int32("Z", ub.Position.Z()).
	Uint32("old runtime id", oldID).
	Str("old name", oldName).
	Uint32("new runtime id", ub.NewBlockRuntimeID).
	Str("new name", name).
	Msg("Update Block")
	}
	// yIdx := c.SubIndex(chunkY)
	// yOffset := chunkY - c.SubY(yIdx)
	c.SetBlock(xOffset, chunkY, zOffset, uint8(ub.Layer), ub.NewBlockRuntimeID)
	return ub, true
}
