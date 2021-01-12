package pkg

import (
	"bufio"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
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
)

const protocolID = "/echo/proto/1.0.0"

//ConnectionManager handles streams amd connections
type ConnectionManager struct {
	Ctx      context.Context
	Host     host.Host
	Routing  routing.PeerRouting
	Messages chan Msg
}

// Msg chan
type Msg struct {
	Remote peer.ID
	M      string
}

// Send sends a message to a list of addresses
func (c *ConnectionManager) Send(message []byte, peerID peer.ID) error {
	if c.Ctx.Err() != nil {
		return errors.Wrap(c.Ctx.Err(), "failed to send message")
	}
	ctx, cancel := context.WithCancel(c.Ctx)
	defer cancel()

	s, err := c.Host.NewStream(ctx, peerID, protocolID)
	if err != nil {
		return errors.Wrap(err, "could not open new stream to remote")
	}

	writer := bufio.NewWriter(s)
	_, err = writer.WriteString(string(message) + "\n")

	err = writer.Flush()
	if err != nil {
		log.Println("Failed to send the transaction to sign to", peerID, ":", err)
	}
	log.Println("[DEBUG] Sent message to", peerID, ":'", string(message), "'")
	cancel()

	return err
}

//Start creates a libp2p host and starts handling connections
func (c *ConnectionManager) Start(ctx context.Context, privateKey crypto.PrivKey) (err error) {
	c.Ctx = ctx
	libp2pCtx, unused := context.WithCancel(ctx)
	_ = unused // pacify vet lostcancel check: libp2pCtx is always canceled through its parent
	c.Host, c.Routing, err = createLibp2pHost(libp2pCtx, privateKey)
	if err != nil {
		return
	}

	fmt.Println("Started dht peer", c.Host.ID().Pretty())
	if c.Messages == nil {
		c.Messages = make(chan Msg)
	}

	c.Host.SetStreamHandler(protocolID, func(s p2pnetwork.Stream) {
		connection := s.Conn()
		keybytes, err := connection.RemotePublicKey().Bytes()
		if err != nil {
			log.Println("Failed to read message")
			return
		}
		log.Println("Got a new stream from", connection.RemotePeer())
		log.Println("remote pubkey", hex.EncodeToString(keybytes))
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		message, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Failed to read message")
			return
		}
		message = strings.TrimSuffix(message, "\n")
		c.Messages <- Msg{
			M:      message,
			Remote: connection.RemotePeer(),
		}
		s.Close() //TODO: don't close it immediately but reuse when possible.
	})

	return nil
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
