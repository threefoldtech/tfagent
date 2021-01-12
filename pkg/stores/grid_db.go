package stores

import (
	"bytes"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v2"
	"github.com/centrifuge/go-substrate-rpc-client/v2/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v2/types"
)

// Client is a struct that holds the api client
type Client struct {
	api *gsrpc.SubstrateAPI
}

type Twin struct {
	TwinID   types.U64
	Pubkey   types.AccountID
	PeerID   []types.U8
	Entities []entityProof
}

type entityProof struct {
	entityID  types.U64
	signature []types.U8
}

// NewGridDB creates a new substrate api client
func NewGridDB(url string) (*Client, error) {
	if url == "" {
		url = "ws://localhost:9944"
	}
	api, err := gsrpc.NewSubstrateAPI(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		api: api,
	}, nil
}

// GetTwin gets a twin by id from storage
func (c *Client) GetTwin(twinID uint64) (Twin, error) {
	meta, err := c.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return Twin{}, err
	}

	buf := bytes.NewBuffer(nil)
	enc := scale.NewEncoder(buf)
	if err := enc.Encode(twinID); err != nil {
		return Twin{}, err
	}
	key := buf.Bytes()

	key, err = types.CreateStorageKey(meta, "TemplateModule", "Twins", key, nil)
	if err != nil {
		return Twin{}, err
	}

	var twin Twin
	ok, err := c.api.RPC.State.GetStorageLatest(key, &twin)
	if err != nil || !ok {
		return Twin{}, err
	}

	return twin, nil
}

func byteSliceToString(bs []types.U8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}
