package pkg

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	p2pnetwork "github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/pkg/errors"
	"github.com/stellar/go/strkey"
	"github.com/threefoldtech/tfagent/pkg/workloads"

	"github.com/libp2p/go-libp2p-core/host"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"
	secio "github.com/libp2p/go-libp2p-secio"
	libp2ptls "github.com/libp2p/go-libp2p-tls"

	"github.com/libp2p/go-libp2p-core/crypto"
)

const (
	protocolIDWorkloadDeploy = "/tfagent/workload_deploy/1.0.0"
	protocolIDWorkloadDelete = "/tfagent/workload_delete/1.0.0"
	// TODO
)

//ConnectionManager handles streams amd connections
type ConnectionManager struct {
	Ctx      context.Context
	Host     host.Host
	Routing  routing.PeerRouting
	Messages chan string
	Log      Log
}

// IncommingMessage from a remote
type IncommingMessage struct {
	PeerID  peer.ID
	PeerKey crypto.PubKey
}

//NewConnectionManager creates a new ConnectionManager
func NewConnectionManager(coSigners []string) *ConnectionManager {
	c := &ConnectionManager{
		//All received messages are sent through this channel
		Messages: make(chan string),
	}
	return c
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

	c.Host.SetStreamHandler(protocolIDWorkloadDeploy, func(s p2pnetwork.Stream) {
		connection := s.Conn()
		remoteAddress, err := stellarAddressFromP2PPublicKey(connection.RemotePublicKey())
		if err != nil {
			log.Println("Failed to get the Stellar address from remote connection", connection.RemotePeer())
			return
		}
		log.Println("Got a new stream from", remoteAddress)
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		message, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Failed to read message")
			return
		}
		message = strings.TrimSuffix(message, "\n")
		c.Messages <- message
		s.Close() //TODO: don't close it immediately but reuse when possible.
	})

	return nil
}

// RegisterWorkloadInterest notifies the network handler that we are intereseted
// in receiving workload deployment and delete requests.
func (c *ConnectionManager) RegisterWorkloadInterest(deployChan chan<- workloads.WorkloadInfo, deleteChan chan<- workloads.WorkloadInfo) {
	c.Host.SetStreamHandler(protocolIDWorkloadDeploy, func(s p2pnetwork.Stream) {
		connection := s.Conn()
		remoteAddress, err := stellarAddressFromP2PPublicKey(connection.RemotePublicKey())
		if err != nil {
			log.Println("Failed to get the Stellar address from remote connection", connection.RemotePeer())
			return
		}
		log.Println("Got a new stream from", remoteAddress)
		dec := json.NewDecoder(s)
		wenv := workloads.WorkloadEnveloppe{}
		if err = dec.Decode(&wenv); err != nil {
			s.Reset()
			return
		}

		var workload workloads.Workload
		switch wenv.Type {
		case workloads.WorkloadTypeZDB:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeContainer:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeVolume:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeNetwork:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeKubernetes:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeProxy:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeReverseProxy:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeSubDomain:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeDomainDelegate:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeGateway4To6:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypeNetworkResource:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		case workloads.WorkloadTypePublicIP:
			var w workloads.ZDB
			if err = json.Unmarshal(wenv.WorkloadData, &w); err != nil {
				s.Reset()
				return
			}
			workload = w
		}
		deployChan <- workloads.WorkloadInfo{Name: wenv.Name, Peer: remoteAddress, Workload: workload}
		s.Close()
	})

	c.Host.SetStreamHandler(protocolIDWorkloadDelete, func(s p2pnetwork.Stream) {
		connection := s.Conn()
		remoteAddress, err := stellarAddressFromP2PPublicKey(connection.RemotePublicKey())
		if err != nil {
			log.Println("Failed to get the Stellar address from remote connection", connection.RemotePeer())
			return
		}
		log.Println("Got a new stream from", remoteAddress)
		dec := json.NewDecoder(s)
		var workloadName string
		if err = dec.Decode(&workloadName); err != nil {
			s.Reset()
			return
		}

		deleteChan <- workloads.WorkloadInfo{Name: workloadName, Peer: remoteAddress}
		s.Close()
	})
}

func (c *ConnectionManager) RequestProvision(ctx context.Context, workload workloads.Workload) error {

}

//send sends a message to a list of addresses
func (c *ConnectionManager) send(message []byte, address string, protocolID protocol.ID) error {
	if c.Ctx.Err() != nil {
		return errors.Wrap(c.Ctx.Err(), "failed to send message")
	}
	peerID, err := getPeerIDFromStellarAddress(address)
	if err != nil {
		return errors.Wrap(err, "could not get peerID from address")
	}
	ctx, cancel := context.WithCancel(c.Ctx)
	defer cancel()

	s, err := c.Host.NewStream(ctx, peerID, protocolID)
	if err != nil {
		return errors.Wrap(err, "could not open new stream to remote")
	}

	io.Copy(
	_, err = writer.WriteString(message + "\n")

	err = writer.Flush()
	if err != nil {
		log.Println("Failed to send the transaction to sign to", cosignerAddress, ":", err)
	}
	log.Println("[DEBUG] Sent message to", cosignerAddress, ":'", message, "'")
	cancel()
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

// func (c *ConnectionManager) connectToPeer(ctx context.Context, peerID peer.ID) (err error) {
// 	findPeerCtx, cancel := context.WithCancel(ctx)
// 	defer cancel()
// 	peeraddrInfo, err := c.Routing.FindPeer(findPeerCtx, peerID)
// 	if err != nil {
// 		return
// 	}
// 	ConnectPeerCtx, cancel := context.WithCancel(ctx)
// 	defer cancel()
// 	err = c.Host.Connect(ConnectPeerCtx, peeraddrInfo)
// 	return
// }

// func getLibp2pPrivateKeyFromStellarSeed(seed string) crypto.PrivKey {
// 	versionbyte, rawSecret, err := strkey.DecodeAny(seed)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	if versionbyte != strkey.VersionByteSeed {
// 		log.Fatalf("%s is not a valid Stellar seed", seed)
// 	}
//
// 	secretreader := bytes.NewReader(rawSecret)
// 	libp2pPrivKey, _, err := crypto.GenerateEd25519Key(secretreader)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return libp2pPrivKey
// }

func getPeerIDFromStellarAddress(address string) (peerID peer.ID, err error) {
	versionbyte, pubkeydata, err := strkey.DecodeAny(address)
	if err != nil {
		return
	}
	if versionbyte != strkey.VersionByteAccountID {
		err = fmt.Errorf("%s is not a valid Stellar address", address)
		return
	}
	libp2pPubKey, err := crypto.UnmarshalEd25519PublicKey(pubkeydata)
	if err != nil {
		return
	}

	peerID, err = peer.IDFromPublicKey(libp2pPubKey)
	return peerID, err
}

func stellarAddressFromP2PPublicKey(pubKey crypto.PubKey) (address string, err error) {
	rawPubKey, err := pubKey.Raw()
	if err != nil {
		return
	}
	address, err = strkey.Encode(strkey.VersionByteAccountID, rawPubKey)
	return
}
