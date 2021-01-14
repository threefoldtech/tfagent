package pkg

import (
	"context"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/pkg/errors"
)

type BufferedNode struct {
	node    *P2PNode
	msgChan <-chan Message

	peerStore PeerStore

	// receiving queue, messages are kept in the order they are received
	recvQ     []Message
	recvQLock sync.Mutex
	// sending queue, message are kept in the order they are submitted
	sendQ     []Message
	sendQLock sync.Mutex

	ctx context.Context
}

const singleMessageSendTTL = time.Second * 20 // 20 seconds by default to send a message

// NewBufferedNode creates a new buffered node embedding a regular P2PNode
func NewBufferedNode(store PeerStore) *BufferedNode {
	msgChan := make(chan Message)
	return &BufferedNode{
		node:      NewP2PNode(msgChan),
		msgChan:   msgChan,
		peerStore: store,
		recvQ:     []Message{},
		sendQ:     []Message{},
	}
}

func (bn *BufferedNode) Send(message Message) error {
	if err := bn.ctx.Err(); err != nil {
		return errors.Wrap(err, "could not send message")
	}

	peerIDStr, err := bn.peerStore.PeerID(message.Receiver)
	if err != nil {
		return errors.Wrap(err, "could not load receiver peerID")
	}

	peerID, err := peer.IDFromString(peerIDStr)
	if err != nil {
		return errors.Wrap(err, "invalid receiver peerID")
	}

	err = bn.node.Send(message, peerID, singleMessageSendTTL)
	if errors.Is(err, context.DeadlineExceeded) {
		bn.sendQLock.Lock()
		defer bn.sendQLock.Unlock()
		bn.sendQ = append(bn.sendQ, message)
		return nil // TODO: return ErrQueued?
	} else if err != nil {
		return errors.Wrap(err, "could not send message")
	}

	return nil
}

func (bn *BufferedNode) Start(ctx context.Context, privateKey crypto.PrivKey) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-bn.msgChan:
				bn.recvQLock.Lock()
				bn.recvQ = append(bn.recvQ, msg)
				bn.recvQLock.Unlock()
			}
		}
	}()
	// TODO: Start maintenance routine to kick expired Messages
	bn.ctx = ctx
	return bn.node.Start(ctx, privateKey)
}

// PeerID returns the underlying nodes PeerID
func (bn *BufferedNode) PeerID() string {
	return bn.node.PeerID()
}
