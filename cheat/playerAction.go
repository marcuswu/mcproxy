package cheat

import (
	//"github.com/rs/zerolog/log"
	//"time"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	//"github.com/sandertv/gophertunnel/minecraft/protocol"
)

func (proxy *Proxy) HandlePlayerAction(pa *packet.PlayerAction) (*packet.PlayerAction, bool) {

	return pa, true
}
