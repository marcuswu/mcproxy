package cheat

import (
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleSubChunk(c *packet.SubChunk) (*packet.SubChunk, bool) {
	log.Info().Msgf("got subchunk with position %d, %d, %d", c.Position.X(), c.Position.Y(), c.Position.Z())

	// for _, e := range chunk.SubChunkEntries {
	// e.Offset
	// }
	// chunk

	// The raw data for a chunk is a varint followed by that many varints representing indexes to the global block palette

	return c, true
}
