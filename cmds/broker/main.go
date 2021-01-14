package main

import (
	"context"
	"crypto/rand"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/rs/zerolog/log"
	"github.com/threefoldtech/tfagent/pkg"
	"github.com/threefoldtech/tfagent/pkg/stores"
)

func main() {
	ctx := context.Background()

	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		log.Fatal().Err(err).Msg("could not generate key")
		return
	}

	store := stores.MockStore{}
	node := pkg.NewBufferedNode(store)
	node.Start(ctx, priv)

	server, err := pkg.NewServer(ctx, 8888, store, node)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get server")
	}
	defer server.Close()

	server.Run()
}
