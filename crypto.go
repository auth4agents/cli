package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

func GenerateEd25519Keypair() (
	ed25519.PublicKey,
	ed25519.PrivateKey,
	error,
) {

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	return pub, priv, nil
}

func EncodePublicKey(
	publicKey ed25519.PublicKey,
) string {

	return base64.StdEncoding.EncodeToString(
		publicKey,
	)
}

func EncodePrivateKey(
	privateKey ed25519.PrivateKey,
) string {

	return base64.StdEncoding.EncodeToString(
		privateKey,
	)
}

func DecodePublicKey(
	value string,
) (ed25519.PublicKey, error) {

	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	if len(raw) != ed25519.PublicKeySize {
		return nil, fmt.Errorf(
			"invalid public key length: %d",
			len(raw),
		)
	}

	return ed25519.PublicKey(raw), nil
}

func DecodePrivateKey(
	value string,
) (ed25519.PrivateKey, error) {

	raw, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}

	if len(raw) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf(
			"invalid private key length: %d",
			len(raw),
		)
	}

	return ed25519.PrivateKey(raw), nil
}

func SignChallenge(
	privateKey ed25519.PrivateKey,
	challenge string,
	nonce string,
) (string, error) {

	if challenge == "" {
		return "", errors.New(
			"challenge is empty",
		)
	}

	if nonce == "" {
		return "", errors.New(
			"nonce is empty",
		)
	}

	message := BuildProofMessage(
		challenge,
		nonce,
	)

	signature := ed25519.Sign(
		privateKey,
		[]byte(message),
	)

	return base64.StdEncoding.EncodeToString(
		signature,
	), nil
}

func VerifyChallengeSignature(
	publicKey ed25519.PublicKey,
	challenge string,
	nonce string,
	signatureB64 string,
) error {

	signature, err := base64.StdEncoding.DecodeString(
		signatureB64,
	)

	if err != nil {
		return err
	}

	message := BuildProofMessage(
		challenge,
		nonce,
	)

	valid := ed25519.Verify(
		publicKey,
		[]byte(message),
		signature,
	)

	if !valid {
		return errors.New(
			"invalid signature",
		)
	}

	return nil
}

func BuildProofMessage(
	challenge string,
	nonce string,
) string {

	return fmt.Sprintf(
		"auth4agents-challenge:%s:%s",
		challenge,
		nonce,
	)
}

func SHA256Hex(
	value string,
) string {

	hash := sha256.Sum256(
		[]byte(value),
	)

	return hex.EncodeToString(hash[:])
}

func GenerateNonce() (string, error) {

	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
