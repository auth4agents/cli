package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var verifyToken string
var verifyServer string
var verifyOnline bool
var verifyJSON bool

type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid,omitempty"`
}

type JWTClaims struct {
	ISS string   `json:"iss"`
	SUB string   `json:"sub"`
	AUD []string `json:"aud"`
	IAT int64    `json:"iat"`
	EXP int64    `json:"exp"`

	JTI string `json:"jti,omitempty"`

	Scope string `json:"scope,omitempty"`

	OperatorID string `json:"operator_id,omitempty"`

	Revoked bool `json:"revoked,omitempty"`
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify and inspect JWT token",

	Run: func(cmd *cobra.Command, args []string) {

		header, claims, err := DecodeJWT(
			verifyToken,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"invalid token: %v\n",
				err,
			)

			os.Exit(1)
		}

		expired := IsTokenExpired(claims)

		result := map[string]interface{}{
			"valid_structure": true,
			"expired":         expired,
			"header":          header,
			"claims":          claims,
		}

		if verifyOnline {

			server := verifyServer

			if server == "" {

				cfg, err := LoadAgentConfig()
				if err == nil {
					server = cfg.ServerURL
				}
			}

			if server == "" {
				fmt.Fprintf(
					os.Stderr,
					"--server is required for online verification\n",
				)

				os.Exit(1)
			}

			client := NewClient(server)

			ctx, cancel := context.WithTimeout(
				context.Background(),
				10*time.Second,
			)

			defer cancel()

			verifyResp, err := client.VerifyToken(
				ctx,
				verifyToken,
			)

			if err != nil {
				fmt.Fprintf(
					os.Stderr,
					"online verification failed: %v\n",
					err,
				)

				os.Exit(1)
			}

			result["online_verification"] = verifyResp
		}

		if verifyJSON {

			_ = json.NewEncoder(os.Stdout).
				Encode(result)

			return
		}

		fmt.Printf(
			"token inspection\n\n",
		)

		fmt.Printf(
			"issuer: %s\n",
			claims.ISS,
		)

		fmt.Printf(
			"subject: %s\n",
			claims.SUB,
		)

		fmt.Printf(
			"audience: %s\n",
			claims.AUD,
		)

		fmt.Printf(
			"scope: %s\n",
			claims.Scope,
		)

		fmt.Printf(
			"issued_at: %d\n",
			claims.IAT,
		)

		fmt.Printf(
			"expires_at: %d\n",
			claims.EXP,
		)

		fmt.Printf(
			"expired: %v\n",
			expired,
		)

		if verifyOnline {

			online := result["online_verification"].(*VerifyResponse)

			fmt.Printf(
				"revoked: %v\n",
				online.Revoked,
			)

			fmt.Printf(
				"verified_online: %v\n",
				online.Valid,
			)
		}
	},
}

func DecodeJWT(
	token string,
) (*JWTHeader, *JWTClaims, error) {

	parts := strings.Split(token, ".")

	if len(parts) != 3 {
		return nil, nil, fmt.Errorf(
			"invalid JWT structure",
		)
	}

	headerBytes, err := base64.RawURLEncoding.
		DecodeString(parts[0])

	if err != nil {
		return nil, nil, err
	}

	claimsBytes, err := base64.RawURLEncoding.
		DecodeString(parts[1])

	if err != nil {
		return nil, nil, err
	}

	var header JWTHeader

	err = json.Unmarshal(
		headerBytes,
		&header,
	)

	if err != nil {
		return nil, nil, err
	}

	var claims JWTClaims

	err = json.Unmarshal(
		claimsBytes,
		&claims,
	)

	if err != nil {
		return nil, nil, err
	}

	return &header, &claims, nil
}

func IsTokenExpired(
	claims *JWTClaims,
) bool {

	now := time.Now().Unix()

	return now >= claims.EXP
}

func init() {

	verifyCmd.Flags().
		StringVar(
			&verifyToken,
			"token",
			"",
			"JWT token",
		)

	verifyCmd.Flags().
		StringVar(
			&verifyServer,
			"server",
			"",
			"Auth4Agents server URL",
		)

	verifyCmd.Flags().
		BoolVar(
			&verifyOnline,
			"online",
			false,
			"perform online verification",
		)

	verifyCmd.Flags().
		BoolVar(
			&verifyJSON,
			"json",
			false,
			"output JSON",
		)

	_ = verifyCmd.MarkFlagRequired(
		"token",
	)

	rootCmd.AddCommand(
		verifyCmd,
	)
}
