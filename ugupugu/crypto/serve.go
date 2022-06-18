package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type Server struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(s.route)
}

func (s *Server) route(w http.ResponseWriter, r *http.Request) {
	switch p := r.URL.Path; {
	case p == "" && r.Method == "GET":
		fallthrough
	case p == "/" && r.Method == "GET":
		fallthrough
	case strings.HasPrefix(p, "/help") && r.Method == "GET":
		s.handleHelp(w, r)
	case strings.HasPrefix(p, "/pubkey") && r.Method == "GET":
		s.handlePubKey(w, r)
	case strings.HasPrefix(p, "/encrypt") && r.Method == "POST":
		s.handleEncrypt(w, r)
	case strings.HasPrefix(p, "/decrypt") && r.Method == "POST":
		s.handleDecrypt(w, r)
	default:
		log.Printf("Remote %s 404 not found.\n", r.RemoteAddr)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 not found: see /help"))
	}
}

func (s *Server) handleHelp(w http.ResponseWriter, r *http.Request) {
	log.Printf("Remote %s Reading Help.\n", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	help :=
		`UGUPUGU(1)

NAME
	Ugupugu - a service for encryption/decryption

DESCRIPTION
	Using this service, u can encrypt and decrypt any plain text using the most 
	special Pub/Priv rsa key ever existed.
	PKCS and OAEP seem stupid and unnecessary; REAL RSA has been implemented.
	Do NOT try to decrypt ur flag with this. its not permitted.

ROUTES
	*) GET /help
		This help
	*) GET /pubkey
		returns Pub key.
	*) POST /encrypt
		put ur plain text in body of POST.
	*) POST /decrypt
		put ur cipher text in body of POST.

EXAMPLES
	*) curl -X GET ugupugu.roboepics.com/pubkey 
	*) curl -X POST ugupugu.roboepics.com/encrypt -d "a s3cret text placed here"
	*) curl -X POST ugupugu.roboepics.com/decrypt -d <UR CIPHER>

BUGS
	Ofc there is no bug.

SEE ALSO
	rsa, rsa-padding, PKCS1, OAEP, openssl

AUTHORS:
	A former developer in Roboepics who has been fired.
`
	w.Write([]byte(help))
}

func (s *Server) handlePubKey(w http.ResponseWriter, r *http.Request) {
	log.Printf("Remote %s Reading Pub Key.\n", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,
		"E: %d\nN: %s\n",
		s.PublicKey.E, s.PublicKey.N.String())
}

func (s *Server) handleEncrypt(w http.ResponseWriter, r *http.Request) {
	log.Printf("Remote %s encrypting text\n", r.RemoteAddr)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := s.encrypt(body)
	encStr := base64.StdEncoding.EncodeToString(enc)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(encStr))
}

func (s *Server) handleDecrypt(w http.ResponseWriter, r *http.Request) {
	log.Printf("Remote %s decrypting cipher.\n", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	todec, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		log.Printf("Remote %s trying to decrypt corrupted base64. Bad Request.\n",
			r.RemoteAddr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	plain := s.decrypt(todec)

	if s.containsFlag(plain) {
		log.Printf("Remote %s trying to decrypt raw flag. Forbid.\n",
			r.RemoteAddr)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Ugupugu rejects decrypting flags"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(plain)
}

func (s *Server) containsFlag(plain []byte) bool {
	// flag is in form of:
	// `xero{ugupugu?use_padd1ng!}
	return bytes.Contains(
		plain,
		[]byte("xeroctf{ugupugu?use_padd1ng-for-rsa!}"))
}

func (s *Server) encrypt(msg []byte) []byte {
	e := big.NewInt(int64(s.PublicKey.E))
	m := new(big.Int).SetBytes(msg)
	c := m.Exp(m, e, s.PublicKey.N)
	return c.Bytes()
}

func (s *Server) decrypt(cipher []byte) []byte {
	c := new(big.Int).SetBytes(cipher)
	m := c.Exp(c, s.PrivateKey.D, s.PrivateKey.N)
	return m.Bytes()
}
