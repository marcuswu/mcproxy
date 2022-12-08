package cheat

import (
	//"github.com/rs/zerolog/log"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/marcuswu/mcproxy/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
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
	//log.Debug().Int32("X", c.Position.X()).Int32("Z", c.Position.Z()).Msg("got chunk")

	proxy.Chunks[c.Position] = levelChunk

	return c, true
}
