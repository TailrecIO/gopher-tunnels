package commons

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/nacl/box"
	"io"
)

type PublicKey = [32]byte
type PrivateKey = [32]byte

type KeyPair struct {
	PublicKey *PublicKey
	PrivateKey *PrivateKey
}

func (kp *KeyPair) GetHexEncodedPublicKey() *string {
	encodedKey := hex.EncodeToString(kp.PublicKey[:])
	return &encodedKey
}

func (kp *KeyPair) GetHexEncodedPrivateKey() *string {
	encodedKey := hex.EncodeToString(kp.PrivateKey[:])
	return &encodedKey
}

func NewKeyPair() *KeyPair {

	publicKey, privateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return &KeyPair {PublicKey: publicKey, PrivateKey: privateKey}
}

func generateNonce() [24]byte {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func Encrypt(publicKey *PublicKey, privateKey *PrivateKey, message *string) *string {
	nonce := generateNonce()
	msgBytes := []byte(*message)
	// This encrypts msg and appends the result to the nonce.
	encrypted := box.Seal(nonce[:], msgBytes, &nonce, publicKey, privateKey)
	cipher := base64.StdEncoding.EncodeToString(encrypted)
	return &cipher
}

func Decrypt(publicKey *PublicKey, privateKey *PrivateKey, cipher *string) (*string, error) {
	cipherBytes, err := base64.StdEncoding.DecodeString(*cipher)
	if err != nil {
		return nil, errors.New("invalid cipher")
	}
	var decryptNonce [24]byte
	copy(decryptNonce[:], cipherBytes[:24])
	decrypted, ok := box.Open(nil, cipherBytes[24:], &decryptNonce, publicKey, privateKey)
	if !ok {
		return nil, errors.New("decrypt error")
	}
	plainText := string(decrypted)
	return &plainText, nil
}

