package main

import (
	"crypto/rand"
	"encoding/base64"
)

var (
	plain []byte
)

func HalfTimePath() []byte {
	l := randInt(20, 100)
	key := make([]byte, l)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	cycle := make([]byte, len(plain))
	for i := range cycle {
		cycle[i] = key[i%l]
	}

	cipher := make([]byte, len(plain))
	for i := range cipher {
		cipher[i] = cycle[i] ^ plain[i]
	}

	b64 := base64.StdEncoding.EncodeToString(cipher)
	return []byte(b64)
}

func randInt(a, b int) int {
	buf := make([]byte, 1)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return int(buf[0])%(b-a+1) + a
}
