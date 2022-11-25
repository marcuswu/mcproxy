package cheat

import (
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Proxy struct {
	ClientConn *minecraft.Conn
	ServerConn *minecraft.Conn
	Chunks     map[protocol.ChunkPos]*chunk.Chunk
	PlayerID   uint64
	PlayerPos  *packet.MovePlayer
}

func (proxy *Proxy) HandleClientPacket(pk packet.Packet) (packet.Packet, bool) {
	// log.Trace().Msgf("Client packet id %v", pk.ID())
	switch pk.ID() {
	case packet.IDCommandRequest:
		return proxy.HandleCommand(pk.(*packet.CommandRequest))
	case packet.IDInventoryTransaction:
		return proxy.HandleInventoryTransaction(pk.(*packet.InventoryTransaction))
	case packet.IDMovePlayer:
		return proxy.HandleMovePlayer(pk.(*packet.MovePlayer), true)
	}

	return pk, true
}

func (proxy *Proxy) HandleServerPacket(pk packet.Packet) (packet.Packet, bool) {
	// log.Trace().Msgf("Server packet id %v", pk.ID())
	switch pk.ID() {
	case packet.IDLevelChunk:
		return proxy.HandleLevelChunk(pk.(*packet.LevelChunk))
	case packet.IDSubChunk:
		return proxy.HandleSubChunk(pk.(*packet.SubChunk))
	case packet.IDMovePlayer:
		return proxy.HandleMovePlayer(pk.(*packet.MovePlayer), false)
	}
	return pk, true
}
