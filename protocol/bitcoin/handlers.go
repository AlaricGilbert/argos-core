package bitcoin

import (
	"fmt"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/cloudwego/netpoll"
)

// Ctx stores the context for a message handler
type Ctx struct {
	peer        *Peer
	header      MessageHeader
	command     string
	payloadhash [32]byte
	checksum    uint32
	payload     *netpoll.LinkBuffer
	err         error
}

type CommandHandler func(ctx *Ctx)

var commandHandlers = map[string]CommandHandler{
	CommandReject:      handleNop,
	CommandVerack:      handleNop,
	CommandPong:        handleNop,
	CommandSendHeaders: handleSendHeaders,
	CommandVersion:     handleVersion,
	CommandInv:         handleInv,
	CommandNotFound:    handleNotFound,
	CommandGetData:     handleGetData,
	CommandTx:          handleTx,
	CommandPing:        handlePing,
	CommandAddr:        handleAddr,
	CommandFilterAdd:   handleFilterAdd,
	CommandFilterClear: handleFilterClear,
	CommandFilterLoad:  handleFilterLoad,
	CommandGetHeaders:  handleGetHeaders,
	CommandHeaders:     handleHeaders,
	CommandSendCmpct:   handleSendCmpct,
	CommandFeeFilter:   handleFeeFilter,
}

func deserializePayload[T any](ctx *Ctx) *T {
	var t T
	if _, ctx.err = serialization.Deserialize(ctx.payload, &t); ctx.err != nil {
		ctx.peer.logger().WithError(ctx.err).Info("bitcoin peer deserialize payload failed")
		return nil
	}

	ctx.peer.logger().WithField("payload", t).WithField("type", fmt.Sprintf("%T", t)).Info("bitcoin peer received payload")
	return &t
}

func handleNop(ctx *Ctx) {}

func handleSendHeaders(ctx *Ctx) {
	ctx.peer.sendheaders = true
}

func handleVersion(ctx *Ctx) {
	if _ = deserializePayload[Version](ctx); ctx.err == nil {
		ctx.err = ctx.peer.sendVerack()
	}
}

func handleInv(ctx *Ctx) {
	if inv := deserializePayload[Inv](ctx); ctx.err == nil {
		revTime := time.Now()
		for _, ii := range inv.Inventory {
			// we only support transactions here
			if ii.Type.Tx() {
				ctx.peer.s.NotifyTransaction(argos.TransactionNotify{
					Source:    ctx.peer.addr.TCPAddr,
					Timestamp: revTime,
					TxID:      ii.Hash,
				})
			}
		}
	}
}

func handleNotFound(ctx *Ctx) {
	if nf := deserializePayload[NotFound](ctx); ctx.err == nil {
		for _, ii := range nf.Inventory {
			if ii.Type.Tx() {
				ctx.peer.logger().WithField("inv", ii).Warn("bitcoin peer transaction notfound")
			}
		}
	}
}

func handleTx(ctx *Ctx) {
	if tx := deserializePayload[Transaction](ctx); ctx.err == nil {
		ctx.peer.txs[ctx.payloadhash] = *tx
	}
}

func handlePing(ctx *Ctx) {
	if ping := deserializePayload[Ping](ctx); ctx.err == nil {
		ctx.err = ctx.peer.sendPong(ping.Nonce)
	}
}

func handleAddr(ctx *Ctx) {
	if addr := deserializePayload[Addr](ctx); ctx.err == nil {
		var addrlist []net.TCPAddr
		for _, address := range addr.AddrList {
			addrlist = append(addrlist, *address.TCPAddr())
		}
		ctx.peer.s.NodeConn(ctx.peer.addr.TCPAddr, addrlist)
	}
}

func handleFilterAdd(ctx *Ctx) {
	if add := deserializePayload[FilterAdd](ctx); ctx.err == nil {
		if ctx.peer.filterLoad != nil {
			ctx.peer.filterLoad.Filter = append(ctx.peer.filterLoad.Filter, add.Data...)
		}
	}
}

func handleFilterLoad(ctx *Ctx) {
	if load := deserializePayload[FilterLoad](ctx); ctx.err == nil {
		ctx.peer.filterLoad = load
	}
}

func handleFilterClear(ctx *Ctx) {
	ctx.peer.filterLoad = nil
}

func handleFeeFilter(ctx *Ctx) {
	if filter := deserializePayload[FeeFilter](ctx); ctx.err == nil {
		ctx.peer.feeFilter = int64(*filter)
	}
}

func handleGetHeaders(ctx *Ctx) {
	if getheaders := deserializePayload[GetHeaders](ctx); ctx.err == nil {
		var invs []Inventory
		for _, bl := range getheaders.BlockLocatorHashes {
			invs = append(invs, Inventory{
				Type: MSG_BLOCK,
				Hash: bl,
			})
		}
		ctx.err = ctx.peer.sendNotFound(invs...)
	}
}

func handleGetData(ctx *Ctx) {
	if getdata := deserializePayload[GetData](ctx); ctx.err == nil {
		ctx.err = ctx.peer.sendNotFound(getdata.Inventory...)
	}
}

func handleHeaders(ctx *Ctx) {
	if _ = deserializePayload[Headers](ctx); ctx.err != nil {
	}
}

func handleSendCmpct(ctx *Ctx) {
	if cmpct := deserializePayload[SendCmpct](ctx); ctx.err != nil {
		ctx.peer.announce = cmpct.Announce
	}
}
