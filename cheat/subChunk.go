package cheat

import (
	"bytes"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/marcuswu/mcproxy/world/chunk"
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleSubChunk(c *packet.SubChunk) (*packet.SubChunk, bool) {
	// A chunk will already have all its subchunks, but filled with air. This packet give us the subchunk block data.
	// * Parse the Subchunk entries
	// * Calculate the subchunk each belongs to
	// * Iterate the the subchunk entries and add its blocks to the Subchunk
	for _, e := range c.SubChunkEntries {
		if e.Result == protocol.SubChunkResultSuccess {
			var uIndex uint8
			chunkPos := protocol.ChunkPos{c.Position.X() + int32(e.Offset[0]), c.Position.Z() + int32(e.Offset[2])}
			//log.Debug().Int32("X", chunkPos.X()).Int32("Z", chunkPos.Z()).Msg("got subchunk")
			ch, ok := proxy.Chunks[chunkPos]
			if !ok {
				air, ok := chunk.StateToRuntimeID("minecraft:air", nil)
				if !ok {
					return c, true
				}
				ch = chunk.New(air, cube.Range{-64, 319})
				proxy.Chunks[chunkPos] = ch
			}
			buf := bytes.NewBuffer(e.RawPayload)
			sch, err := chunk.DecodeSubChunk(buf, ch, &uIndex, chunk.NetworkEncoding)
			if err != nil {
				continue
			}
			Y := ch.SubY(int16(int8(uIndex)))
			idx := ch.SubIndex(int16(Y))
			if len(ch.Sub()) <= int(idx) {
				log.Debug().Msg("Could not add subchunk data!")
			} else {
				// This function loops over every potential block in the subchunk (16*16*16)
				// There is probably a better way...
				// Maybe if the subchunk is not empty, just assign it into the slice
				ch.Sub()[idx].CombineSubChunk(sch)
			}
		}
	}

	return c, true
}
