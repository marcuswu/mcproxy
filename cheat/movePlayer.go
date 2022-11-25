package cheat

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func (proxy *Proxy) HandleMovePlayer(mp *packet.MovePlayer) (*packet.MovePlayer, bool) {
	if mp.EntityRuntimeID == proxy.PlayerID {
		proxy.PlayerPos = mp
	}

	return mp, true
}
