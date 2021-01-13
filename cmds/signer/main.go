package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	msg := "A"

	fmt.Println("hex sig", hex.EncodeToString(ed25519.Sign(priv, []byte(msg))))
	fmt.Println("hex key", hex.EncodeToString(pub))
}
