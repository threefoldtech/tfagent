package pkg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
	secio "github.com/libp2p/go-libp2p-secio"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const protocolID = "/tfagent/message/1.0.0"

// P2PNode handles streams amd connections
type P2PNode struct {
	ctx     context.Context
	host    host.Host
	routing routing.PeerRouting
	msgChan chan<- Message
}

func NewP2PNode(msgChan chan<- Message) *P2PNode {
	return &P2PNode{
		msgChan: msgChan,
	}
}

// Send sends a message to a list of addresses
func (c *P2PNode) Send(message Message, peerID peer.ID, timeout time.Duration) error {
	if c.ctx.Err() != nil {
		return errors.Wrap(c.ctx.Err(), "failed to send message")
	}
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	defer cancel()

	s, err := c.host.NewStream(ctx, peerID, protocolID)
	if err != nil {
		return errors.Wrap(err, "could not open new stream to remote")
	}

	if err = json.NewEncoder(s).Encode(message); err != nil {
		log.Error().Err(err).Str("peerID", string(peerID)).Msg("could not send message to peer")
		return err
	}

	log.Debug().Str("peerID", string(peerID)).Msg("sent message")

	return err
}

// Start creates a libp2p host and starts handling connections
func (c *P2PNode) Start(ctx context.Context, privateKey crypto.PrivKey) error {
	c.ctx = ctx
	var err error
	c.host, c.routing, err = createLibp2pHost(ctx, privateKey)
	if err != nil {
		return err
	}

	log.Info().Str("ID", c.host.ID().Pretty()).Msg("started dht peer")

	c.host.SetStreamHandler(protocolID, func(s p2pnetwork.Stream) {
		connection := s.Conn()

		log.Debug().Str("peerID", connection.RemotePeer().Pretty()).Msg("got a new stream from remote")

		msg := Message{}
		if err = json.NewDecoder(s).Decode(&msg); err != nil {
			log.Debug().Err(err).Msg("could not decode message from peer")
			return
		}

		c.msgChan <- msg
		s.Close() // TODO: don't close it immediately but reuse when possible.
	})

	return nil
}

func (c *P2PNode) PeerID() string {
	return c.host.ID().Pretty()
}

func createLibp2pHost(ctx context.Context, privateKey crypto.PrivKey) (host.Host, routing.PeerRouting, error) {
	var idht *dht.IpfsDHT
	var err error
	libp2phost, err := libp2p.New(ctx,
		// Use the keypair we generated
		libp2p.Identity(privateKey),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",      // regular tcp connections
			"/ip4/0.0.0.0/udp/0/quic", // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support secio connections
		libp2p.Security(secio.ID, secio.New),
		// support QUIC
		libp2p.Transport(libp2pquic.NewTransport),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr.NewConnManager(
			100,         // Lowwater
			400,         // HighWater,
			time.Minute, // GracePeriod
		)),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),
		// Let this host use relays and advertise itself on relays if
		// it finds it is behind NAT. Use libp2p.Relay(options...) to
		// enable active relays and more.
		libp2p.EnableAutoRelay(),
	)
	// This connects to public bootstrappers
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, _ := peer.AddrInfoFromP2pAddr(addr)
		// We ignore errors as some bootstrap peers may be down
		// and that is fine.
		libp2phost.Connect(ctx, *pi)
	}
	return libp2phost, idht, err
}
