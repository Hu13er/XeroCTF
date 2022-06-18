package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func loadPubKey(file string) (*rsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func loadPrivKey(file string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bytes)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
