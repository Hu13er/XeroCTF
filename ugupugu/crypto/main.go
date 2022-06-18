package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := readenv("UGUPUGU_ADDR", "0.0.0.0:9000")
	pubkey := readenv("UGUPUGU_PUB_KEY", "")
	privkey := readenv("UGUPUGU_PRIV_KEY", "")

	s := &Server{}
	var err error
	s.PrivateKey, err = loadPrivKey(privkey)
	if err != nil {
		log.Fatalln("could not parse PrivateKey:", err)
	}
	s.PublicKey, err = loadPubKey(pubkey)
	if err != nil {
		log.Println("could not parse PublicKey:", err)
		if s.PrivateKey != nil {
			log.Println("extracting PublicKey from PrivateKey.")
			s.PublicKey = &s.PrivateKey.PublicKey
		} else {
			os.Exit(1)
		}
	}

	log.Println("Listening on", addr)
	if err := http.ListenAndServe(":9000", s.Handler()); err != nil {
		log.Fatalln(err)
	}
}

func readenv(key, def string) string {
	o := os.Getenv(key)
	if o == "" {
		if def == "" {
			log.Fatalln(key, "not provided.")
		}
		log.Println(key, "not provided. using default:", def)
		o = def
	}
	return o
}
