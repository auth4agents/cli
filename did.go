package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const DIDMethod = "did:agent"

type DIDDocument struct {
	Context              []string             `json:"@context"`
	ID                   string               `json:"id"`
	Controller           string               `json:"controller"`
	CreatedAt            string               `json:"created_at"`
	VerificationMethod   []VerificationMethod `json:"verificationMethod"`
	Authentication       []string             `json:"authentication"`
	AssertionMethod      []string             `json:"assertionMethod"`
	Service              []ServiceEndpoint    `json:"service,omitempty"`
}

type VerificationMethod struct {
	ID                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type ServiceEndpoint struct {
	ID              string      `json:"id"`
	Type            string      `json:"type"`
	ServiceEndpoint interface{} `json:"serviceEndpoint"`
}

func GenerateDID(
	operatorDomain string,
	publicKey ed25519.PublicKey,
) string {

	hash := sha256.Sum256(publicKey)

	suffix := hex.EncodeToString(hash[:16])

	return fmt.Sprintf(
		"%s:%s:%s",
		DIDMethod,
		NormalizeDomain(operatorDomain),
		suffix,
	)
}

func GenerateDIDDocument(
	did string,
	publicKey ed25519.PublicKey,
) *DIDDocument {

	publicKeyB64 := base64.StdEncoding.EncodeToString(publicKey)

	keyID := did + "#key-1"

	return &DIDDocument{
		Context: []string{
			"https://www.w3.org/ns/did/v1",
		},

		ID:         did,
		Controller: did,

		CreatedAt: time.Now().UTC().Format(time.RFC3339),

		VerificationMethod: []VerificationMethod{
			{
				ID:                 keyID,
				Type:               "Ed25519VerificationKey2020",
				Controller:         did,
				PublicKeyMultibase: publicKeyB64,
			},
		},

		Authentication: []string{
			keyID,
		},

		AssertionMethod: []string{
			keyID,
		},
	}
}

func ValidateDID(did string) error {

	if did == "" {
		return errors.New("did is empty")
	}

	if !strings.HasPrefix(did, DIDMethod+":") {
		return fmt.Errorf(
			"invalid DID method: expected %s",
			DIDMethod,
		)
	}

	parts := strings.Split(did, ":")

	if len(parts) != 4 {
		return fmt.Errorf(
			"invalid DID format: expected 4 parts, got %d",
			len(parts),
		)
	}

	method := parts[1]
	if method != "agent" {
		return fmt.Errorf(
			"invalid DID method identifier: %s",
			method,
		)
	}

	domain := parts[2]
	if domain == "" {
		return errors.New("missing operator domain")
	}

	fingerprint := parts[3]

	if len(fingerprint) < 16 {
		return errors.New(
			"invalid DID fingerprint length",
		)
	}

	return nil
}

func ParseDID(did string) (
	method string,
	domain string,
	fingerprint string,
	err error,
) {

	err = ValidateDID(did)
	if err != nil {
		return "", "", "", err
	}

	parts := strings.Split(did, ":")

	return parts[1], parts[2], parts[3], nil
}

func DIDDocumentToJSON(
	doc *DIDDocument,
) ([]byte, error) {

	return json.MarshalIndent(
		doc,
		"",
		"  ",
	)
}

func NormalizeDomain(domain string) string {

	domain = strings.TrimSpace(domain)

	domain = strings.ToLower(domain)

	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	domain = strings.TrimSuffix(domain, "/")

	return domain
}

func PublicKeyFingerprint(
	publicKey ed25519.PublicKey,
) string {

	hash := sha256.Sum256(publicKey)

	return hex.EncodeToString(hash[:16])
}