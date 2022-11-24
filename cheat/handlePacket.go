package cheat

import (
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Proxy struct {
	ClientConn *minecraft.Conn
	ServerConn *minecraft.Conn
}

func (proxy Proxy) HandleClientPacket(pk packet.Packet) (packet.Packet, bool) {
	log.Trace().Msgf("Client packet id %v", pk.ID())
	switch pk.ID() {
	case packet.IDCommandRequest:
		return proxy.HandleCommand(pk.(*packet.CommandRequest))
	}

	return pk, true
}

func (proxy Proxy) HandleServerPacket(pk packet.Packet) (packet.Packet, bool) {
	log.Trace().Msgf("Server packet id %v", pk.ID())
	switch pk.ID() {
	case packet.IDSubChunk:
		return proxy.HandleSubChunk(pk.(*packet.SubChunk))
	}
	return pk, true
}
