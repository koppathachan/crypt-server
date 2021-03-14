package crypt

import (
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func must(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func okay(t bool, m string) {
	if !t {
		log.Panicln(errors.New(m))
	}
}

// Encrypter interface to implement for whomsoever wishes to encrypt request
type Encrypter interface {
	Key() string
	Encrypt(p string) string
	Hash(p string) string
	KeyEncrypt(k string) string
}

// Decrypter interface to implement for whomsoever wishes to decrypt response
type Decrypter interface {
	Decrypt(r io.Reader) string
}

type dec struct {
	el   openpgp.EntityList
	pass []byte
}

func (d dec) Decrypt(r io.Reader) string {
	block, err := armor.Decode(r)
	must(err)
	okay(block.Type == "PGP MESSAGE", "Not a PGP Message")
	md, err := openpgp.ReadMessage(block.Body, d.el, nil, nil)
	must(err)
	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	must(err)
	return string(bytes)
}

func NewDecrypter(filepath string, pb64 string) Decrypter {
	f, err := os.Open(filepath)
	must(err)
	defer f.Close()
	el, err := openpgp.ReadArmoredKeyRing(f)
	must(err)
	pass, err := base64.StdEncoding.DecodeString(pb64)
	must(err)
	e := el[0]
	for _, sk := range e.Subkeys {
		sk.PrivateKey.Decrypt(pass)
	}
	return dec{pass: pass, el: el}
}
