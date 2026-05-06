package main

import (
	"encoding/json"
	"errors"
	"time"
)

const (
	ProofTypeEd25519 = "Ed25519Signature2020"
	ProofPurposeAuth = "authentication"
)

type Proof struct {
	Type               string `json:"type"`
	Created            string `json:"created"`
	ProofPurpose       string `json:"proofPurpose"`
	VerificationMethod string `json:"verificationMethod"`
	Challenge          string `json:"challenge"`
	Nonce              string `json:"nonce"`
	Signature          string `json:"signature"`
}

type TokenExchangeRequest struct {
	DID       string `json:"did"`
	Scope     string `json:"scope"`
	Audience  string `json:"audience"`
	TTL       string `json:"ttl"`
	Proof     Proof  `json:"proof"`
}

type ChallengePayload struct {
	DID string `json:"did"`
}

type ChallengeResult struct {
	Challenge string `json:"challenge"`
	Nonce     string `json:"nonce"`
	ExpiresAt string `json:"expires_at"`
}

func NewProof(
	did string,
	challenge string,
	nonce string,
	signature string,
) *Proof {

	return &Proof{
		Type:         ProofTypeEd25519,

		Created: time.Now().
			UTC().
			Format(time.RFC3339),

		ProofPurpose: ProofPurposeAuth,

		VerificationMethod: did + "#key-1",

		Challenge: challenge,
		Nonce:     nonce,
		Signature: signature,
	}
}

func (p *Proof) Validate() error {

	if p.Type == "" {
		return errors.New(
			"proof type is required",
		)
	}

	if p.ProofPurpose == "" {
		return errors.New(
			"proof purpose is required",
		)
	}

	if p.VerificationMethod == "" {
		return errors.New(
			"verification method is required",
		)
	}

	if p.Challenge == "" {
		return errors.New(
			"challenge is required",
		)
	}

	if p.Nonce == "" {
		return errors.New(
			"nonce is required",
		)
	}

	if p.Signature == "" {
		return errors.New(
			"signature is required",
		)
	}

	return nil
}

func (p *Proof) ToJSON() ([]byte, error) {

	return json.MarshalIndent(
		p,
		"",
		"  ",
	)
}

func ParseProof(
	data []byte,
) (*Proof, error) {

	var proof Proof

	err := json.Unmarshal(data, &proof)
	if err != nil {
		return nil, err
	}

	err = proof.Validate()
	if err != nil {
		return nil, err
	}

	return &proof, nil
}

func BuildTokenExchangeRequest(
	did string,
	scope string,
	audience string,
	ttl string,
	proof *Proof,
) *TokenExchangeRequest {

	return &TokenExchangeRequest{
		DID:      did,
		Scope:    scope,
		Audience: audience,
		TTL:      ttl,
		Proof:    *proof,
	}
}