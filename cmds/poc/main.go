package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/threefoldtech/tfagent/pkg"
)

func main() {
	shouldReply := true
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(os.Args) != 2 && len(os.Args) != 4 {
		fmt.Println("Usage: ./poc <dht user count> (<remote> <message>)")
		return
	}

	amount, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Second arg should be the amount of dhts")
		return
	}

	nodes := make(map[peer.ID]*pkg.ConnectionManager)
	for i := 0; i < amount; i++ {
		conmgr := pkg.ConnectionManager{}
		priv, pub, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			fmt.Println("could not generate key", err.Error())
			return
		}
		if err = conmgr.Start(ctx, priv); err != nil {
			fmt.Println("could not start host:", err.Error())
			return
		}
		pid, err := peer.IDFromPublicKey(pub)
		if err != nil {
			fmt.Println("could not generate peer ID", err.Error())
			return
		}
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-conmgr.Messages:
					fmt.Println("Received message from", msg.Remote, "-", msg.M)
					if shouldReply {
						data := []byte(msg.M)
						reverse(data)
						conmgr.Send(data, msg.Remote)
					}
				}
			}
		}()
		nodes[pid] = &conmgr
	}

	// send message, wait untill our nodes are initialized a bit
	// time.Sleep(time.Second * 5)
	if len(os.Args) == 4 {
		shouldReply = false
		pid, err := peer.Decode(os.Args[2])
		if err != nil {
			fmt.Println("could not interpret peer ID")
			return
		}
		msg := os.Args[3]

		// select some peer to send from
		sendPid := getSender(nodes)
		if sendPid == nil {
			fmt.Println("could not find sender")
			return
		}

		if err = sendPid.Send([]byte(msg), pid); err != nil {
			fmt.Println("failed to send message:", err.Error())
		}

	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan

	return
}

func getSender(nodes map[peer.ID]*pkg.ConnectionManager) *pkg.ConnectionManager {
	for k := range nodes {
		return nodes[k]
	}

	return nil
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
