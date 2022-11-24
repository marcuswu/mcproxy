package cheat

import (
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy Proxy) HandleSubChunk(chunk *packet.SubChunk) (*packet.SubChunk, bool) {
	log.Info().Msgf("got subchunk with position %d, %d, %d", chunk.Position.X(), chunk.Position.Y(), chunk.Position.Z())
	// for _, e := range chunk.SubChunkEntries {
	// e.Offset
	// }
	return chunk, true
}
