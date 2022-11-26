package cheat

import (
	"bytes"

	"github.com/marcuswu/mcproxy/chunk"
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleSubChunk(c *packet.SubChunk) (*packet.SubChunk, bool) {
	log.Info().Msgf("got subchunk with %d entries", len(c.SubChunkEntries))

	// for _, e := range chunk.SubChunkEntries {
	// e.Offset
	// }
	// chunk

	// The raw data for a chunk is a varint followed by that many varints representing indexes to the global block palette
	// find the chunk for this subchunk
	// var index uint8
	// c.sub[index], err = chunk.DecodeSubChunk(bytes.NewBuffer(c.SubChunkEntries[0].RawPayload), c, &index, chunk.NetworkEncoding)

	chunkPos := protocol.ChunkPos{c.Position.X(), c.Position.Z()}
	ch, ok := proxy.Chunks[chunkPos]
	if !ok {
		return c, true
	}
	// A chunk will already have all subchunks.
	// * Parse the Subchunk entries
	// * calculate the subchunk it belongs to
	// * iterate the the subchunk entry and add blocks to the Subchunk
	for _, e := range c.SubChunkEntries {
		if e.Result == protocol.SubChunkResultSuccess {
			var Y uint8
			buf := bytes.NewBuffer(e.RawPayload)
			s, err := chunk.DecodeSubChunk(buf, ch, &Y, chunk.NetworkEncoding)
			ch.SubChunk(int16(Y)).CombineSubChunk(s)
			if err != nil {
				return c, true
			}
		}
	}

	return c, true
}
