package cheat

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/marcuswu/mcproxy/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	// "github.com/rs/zerolog/log"
)

func (proxy *Proxy) HandleLevelChunk(c *packet.LevelChunk) (*packet.LevelChunk, bool) {
	air, ok := chunk.StateToRuntimeID("minecraft:air", nil)
	if !ok {
		return c, true
	}
	levelChunk, err := chunk.NetworkDecode(air, c.RawPayload, int(c.SubChunkCount), cube.Range{-64, 319})
	if err != nil {
		return c, true
	}

	proxy.Chunks[c.Position] = levelChunk

	return c, true
}
