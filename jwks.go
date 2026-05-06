package main

import (
	"context"
	"encoding/json"
	"fmt"
)

type JWK struct {
	KTY string `json:"kty"`
	CRV string `json:"crv"`
	ALG string `json:"alg"`
	USE string `json:"use"`
	KID string `json:"kid"`
	X   string `json:"x"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func FetchJWKS(
	ctx context.Context,
	client *Client,
) (*JWKS, error) {

	var out JWKS

	err := client.Get(
		ctx,
		"/.well-known/jwks.json",
		&out,
	)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func PrintJWKS(
	jwks *JWKS,
) {

	data, _ := json.MarshalIndent(
		jwks,
		"",
		"  ",
	)

	fmt.Println(string(data))
}