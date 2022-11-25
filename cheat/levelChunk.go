package cheat

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/marcuswu/mcproxy/latestmappings"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleLevelChunk(c *packet.LevelChunk) (*packet.LevelChunk, bool) {
	// Decode raw payload to generate palette
	latestmappings.StateToRuntimeID("minecraft:air", nil)
	levelChunk, err := chunk.NetworkDecode(0, c.RawPayload, int(c.SubChunkCount), cube.Range{-64, int(c.HighestSubChunk)})
	if err != nil {
		return c, true
	}

	proxy.Chunks[c.Position] = levelChunk
	// log.Info().
	// 	Int32("Xmin", c.Position.X()<<4).
	// 	Int32("Xmax", (c.Position.X()<<4)+15).
	// 	Int32("Zmin", c.Position.Z()<<4).
	// 	Int32("Zmax", (c.Position.Z()<<4)+15).
	// 	Msgf("Loaded LevelChunk")

	return c, true
}
