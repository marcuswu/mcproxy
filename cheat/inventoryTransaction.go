package cheat

import (
	"github.com/rs/zerolog/log"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

func (proxy *Proxy) HandleInventoryTransaction(transaction *packet.InventoryTransaction) (*packet.InventoryTransaction, bool) {
	log.Info().Msgf("inventory transaction")
	log.Info().Msgf("transaction id %d", transaction.LegacyRequestID)

	for _, slot := range transaction.LegacySetItemSlots {
		log.Info().Msgf("slot %d, %v", slot.ContainerID, slot.Slots)
	}

	for i, action := range transaction.Actions {
		log.Info().
			Uint32("source type", action.SourceType).
			Int32("window id", action.WindowID).
			Uint32("source flags", action.SourceFlags).
			Uint32("slot", action.InventorySlot).
			Int32("old item stack network id", action.OldItem.StackNetworkID).
			Int32("old item stack block id", action.OldItem.Stack.BlockRuntimeID).
			Uint32("old item stack count", uint32(action.OldItem.Stack.Count)).
			Bool("new item has network id", action.NewItem.Stack.HasNetworkID).
			Int32("new item stack network id", action.NewItem.StackNetworkID).
			Int32("new item stack block id", action.NewItem.Stack.BlockRuntimeID).
			Uint32("new item stack count", uint32(action.NewItem.Stack.Count)).
			Bool("new item has network id", action.NewItem.Stack.HasNetworkID).
			Msgf("action %d", i)
	}
	switch transaction.TransactionData.(type) {
	case nil, *protocol.NormalTransactionData:
		log.Info().Msgf("Normal transaction")
	case *protocol.MismatchTransactionData:
		log.Info().Msgf("Mismatch transaction")
	case *protocol.UseItemTransactionData:
		log.Info().Msgf("Use item transaction")
	case *protocol.UseItemOnEntityTransactionData:
		log.Info().Msgf("Use item on entity transaction")
	case *protocol.ReleaseItemTransactionData:
		log.Info().Msgf("Release item transaction")
		released := transaction.TransactionData.(*protocol.ReleaseItemTransactionData)
		log.Info().
			Int32("hot bar slot", released.HotBarSlot).
			Int32("item stack network id", released.HeldItem.StackNetworkID).
			Int32("item block id", released.HeldItem.Stack.BlockRuntimeID).
			Uint16("item count", released.HeldItem.Stack.Count).
			Bool("has network id", released.HeldItem.Stack.HasNetworkID).
			Msgf("ReleaseItemTransactionData")
	}

	return transaction, true
}
