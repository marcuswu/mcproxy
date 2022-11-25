package cheat

import (
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleMovePlayer(mp *packet.MovePlayer, fromClient bool) (*packet.MovePlayer, bool) {
	if fromClient {
		proxy.PlayerID = mp.EntityRuntimeID
	}
	if mp.EntityRuntimeID == proxy.PlayerID {
		proxy.PlayerPos = mp
		log.Info().Msgf("got player %d position %d, %d, %d", mp.EntityRuntimeID, mp.Position.X(), mp.Position.Y(), mp.Position.Z())
	}

	return mp, true
}
