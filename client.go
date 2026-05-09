package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type ChallengeRequest struct {
	DID string `json:"did"`
}

type ChallengeResponse struct {
	Challenge string `json:"challenge"`
	Nonce     string `json:"nonce"`
	ExpiresAt string `json:"expires_at"`
}

type TokenRequest struct {
	DID       string `json:"did"`
	Challenge string `json:"challenge"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Scope     string `json:"scope"`
	Audience  string `json:"audience"`
	TTL       string `json:"ttl"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresAt string `json:"expires_at"`
}

type VerifyRequest struct {
	Token string `json:"token"`
}

type VerifyResponse struct {
	Valid      bool                   `json:"valid"`
	Revoked    bool                   `json:"revoked"`
	Claims     map[string]interface{} `json:"claims"`
	VerifiedAt string                 `json:"verified_at"`
}

type RegisterOperatorRequest struct {
	Domain        string `json:"domain"`
	RootPublicKey string `json:"root_public_key"`
}

type RegisterOperatorResponse struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`
	Status string `json:"status"`
}

type RegisterAgentRequest struct {
	DID            string      `json:"did"`
	DIDDocument    interface{} `json:"did_document"`
	AgentPublicKey string      `json:"agent_public_key"`
	OperatorID     string      `json:"operator_id"`
}

type RegisterAgentResponse struct {
	ID     string `json:"id"`
	DID    string `json:"did"`
	Status string `json:"status"`
}

func (c *Client) RegisterOperator(
	ctx context.Context,
	req RegisterOperatorRequest,
) (*RegisterOperatorResponse, error) {

	var out RegisterOperatorResponse

	err := c.Post(
		ctx,
		"/v1/operators",
		req,
		&out,
	)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) RegisterAgent(
	ctx context.Context,
	req RegisterAgentRequest,
) (*RegisterAgentResponse, error) {

	var out RegisterAgentResponse

	path := fmt.Sprintf(
		"/v1/operators/%s/agents",
		req.OperatorID,
	)

	err := c.Post(
		ctx,
		path,
		req,
		&out,
	)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) GetChallenge(
	ctx context.Context,
	did string,
) (*ChallengeResponse, error) {

	reqBody := ChallengeRequest{
		DID: did,
	}

	var out ChallengeResponse

	err := c.Post(
		ctx,
		"/v1/challenge",
		reqBody,
		&out,
	)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) ExchangeToken(
	ctx context.Context,
	req TokenRequest,
) (*TokenResponse, error) {

	var out TokenResponse

	err := c.Post(
		ctx,
		"/v1/token",
		req,
		&out,
	)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) VerifyToken(
	ctx context.Context,
	token string,
) (*VerifyResponse, error) {

	reqBody := VerifyRequest{
		Token: token,
	}

	var out VerifyResponse

	err := c.Post(
		ctx,
		"/v1/verify",
		reqBody,
		&out,
	)

	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) Get(
	ctx context.Context,
	path string,
	out interface{},
) error {

	return c.do(
		ctx,
		http.MethodGet,
		path,
		nil,
		out,
	)
}


func (c *Client) Post(
	ctx context.Context,
	path string,
	payload interface{},
	out interface{},
) error {

	return c.do(
		ctx,
		http.MethodPost,
		path,
		payload,
		out,
	)
}

func (c *Client) do(
	ctx context.Context,
	method string,
	path string,
	payload interface{},
	out interface{},
) error {
	var body io.Reader
	
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(data)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		
		// Include status code in error for better handling
		return fmt.Errorf("%s %s failed (status %d): %s",
			method, path, resp.StatusCode,
			strings.TrimSpace(string(respBody)),
		)
	}
	
	if out == nil {
		return nil
	}
	
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) Patch(
	ctx context.Context,
	path string,
	payload interface{},
	out interface{},
) error {

	return c.do(
		ctx,
		http.MethodPatch,
		path,
		payload,
		out,
	)
}
