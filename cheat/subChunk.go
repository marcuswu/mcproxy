package cheat

import (
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
	log.Debug().Int8("X", int8(chunkPos.X())).Int8("Z", int8(chunkPos.Z())).Msg("Subchunk entry offset")
	// ch, ok := proxy.Chunks[chunkPos]
	// if !ok {
	// 	return c, true
	// }
	// A chunk will already have all subchunks.
	// * Parse the Subchunk entries
	// * calculate the subchunk it belongs to
	// * iterate the the subchunk entry and add blocks to the Subchunk
	for _, e := range c.SubChunkEntries {
		if e.Result == protocol.SubChunkResultSuccess {
			log.Debug().Int8("X", e.Offset[0]).Int8("Y", e.Offset[1]).Int8("Z", e.Offset[2]).Msg("Subchunk entry offset")
			// var Y uint8
			// buf := bytes.NewBuffer(e.RawPayload)
			// sch, err := chunk.DecodeSubChunk(buf, ch, &Y, chunk.NetworkEncoding)
			// for i, layer := range sch.Layers() {
			// e.Offset
			// layer.
			// }
			/*idx := ch.SubIndex(int16(Y))
			if len(ch.Sub()) <= int(idx) {
				ch.SetSubChunk(s, int16(Y))
			} else {
				ch.SubChunk(idx).CombineSubChunk(s)
			}*/
			// if err != nil {
			// return c, true
			// }
		}
	}

	return c, true
}
