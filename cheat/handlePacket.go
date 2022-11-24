package cheat

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type Proxy struct {
	ClientConn *minecraft.Conn
	ServerConn *minecraft.Conn
}

func HandleClientPacket(pk packet.Packet, proxy Proxy) (packet.Packet, bool) {
	switch pk.ID() {
	case packet.IDCommandRequest:
		return HandleCommand(pk.(*packet.CommandRequest), proxy)
	}

	return pk, true
}

func HandleServerPacket(pk packet.Packet, proxy Proxy) (packet.Packet, bool) {
	return pk, true
}
