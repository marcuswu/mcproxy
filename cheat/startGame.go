package cheat

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

func (proxy *Proxy) HandleStartGame(sg *packet.StartGame) (*packet.StartGame, bool) {
	proxy.PlayerID = sg.EntityRuntimeID

	return sg, true
}
