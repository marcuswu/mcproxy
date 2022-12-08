package cheat

import (
	//"github.com/rs/zerolog/log"
	"github.com/marcuswu/mcproxy/world/chunk"
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
	case packet.IDPlayerAction:
		return proxy.HandlePlayerAction(pk.(*packet.PlayerAction))
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
	case packet.IDUpdateBlock:
		return proxy.HandleUpdateBlock(pk.(*packet.UpdateBlock))
	}
	return pk, true
}

func AbilitiesToString(abilities uint32) []string {
	have := make([]string, 0, 0)
	if abilities & protocol.AbilityBuild != 0 {
		have = append(have, "Build")
	}
	if abilities & protocol.AbilityMine != 0 {
		have = append(have, "Mine")
	}
	if abilities & protocol.AbilityDoorsAndSwitches != 0 {
		have = append(have, "DoorsAndSwitches")
	}
	if abilities & protocol.AbilityOpenContainers != 0 {
		have = append(have, "OpenContainers")
	}
	if abilities & protocol.AbilityAttackPlayers != 0 {
		have = append(have, "AttackPlayers")
	}
	if abilities & protocol.AbilityAttackMobs != 0 {
		have = append(have, "AttackMobs")
	}
	if abilities & protocol.AbilityOperatorCommands != 0 {
		have = append(have, "OperatorCommands")
	}
	if abilities & protocol.AbilityTeleport != 0 {
		have = append(have, "Teleport")
	}
	if abilities & protocol.AbilityInvulnerable != 0 {
		have = append(have, "Invulnerable")
	}
	if abilities & protocol.AbilityFlying != 0 {
		have = append(have, "Flying")
	}
	if abilities & protocol.AbilityMayFly != 0 {
		have = append(have, "MayFly")
	}
	if abilities & protocol.AbilityInstantBuild != 0 {
		have = append(have, "InstantBuild")
	}
	if abilities & protocol.AbilityLightning != 0 {
		have = append(have, "Lightning")
	}
	if abilities & protocol.AbilityFlySpeed != 0 {
		have = append(have, "FlySpeed")
	}
	if abilities & protocol.AbilityWalkSpeed != 0 {
		have = append(have, "WalkSpeed")
	}
	if abilities & protocol.AbilityMuted != 0 {
		have = append(have, "Muted")
	}
	if abilities & protocol.AbilityWorldBuilder != 0 {
		have = append(have, "WorldBuilder")
	}
	if abilities & protocol.AbilityNoClip != 0 {
		have = append(have, "NoClip")
	}
	return have
}
