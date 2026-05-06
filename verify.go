package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
)

var verifyToken string
var verifyOnline bool
var verifyJSON bool

type VerifyClaims struct {
	jwt.RegisteredClaims

	Scope string `json:"scope"`

	OperatorID string `json:"operator_id"`
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify JWT token",

	Run: func(cmd *cobra.Command, args []string) {

		cfg, err := LoadAgentConfig()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to load config: %v\n",
				err,
			)

			os.Exit(1)
		}

		client := NewClient(
			cfg.ServerURL,
		)

		ctx, cancel := context.WithTimeout(
			context.Background(),
			15*time.Second,
		)

		defer cancel()

		jwks, err := FetchJWKS(
			ctx,
			client,
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"failed to fetch jwks: %v\n",
				err,
			)

			os.Exit(1)
		}

		token, err := jwt.ParseWithClaims(
			verifyToken,
			&VerifyClaims{},
			func(token *jwt.Token) (interface{}, error) {

				if token.Method.Alg() != jwt.SigningMethodEdDSA.Alg() {
					return nil, fmt.Errorf(
						"invalid signing algorithm",
					)
				}

				kidRaw, ok := token.Header["kid"]

				if !ok {
					return nil, fmt.Errorf(
						"missing kid",
					)
				}

				kid, ok := kidRaw.(string)

				if !ok {
					return nil, fmt.Errorf(
						"invalid kid",
					)
				}

				for _, jwk := range jwks.Keys {

					if jwk.KID != kid {
						continue
					}

					return JWKToPublicKey(jwk)
				}

				return nil, fmt.Errorf(
					"matching jwk not found",
				)
			},
		)

		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"verification failed: %v\n",
				err,
			)

			os.Exit(1)
		}

		if !token.Valid {
			fmt.Fprintf(
				os.Stderr,
				"token invalid\n",
			)

			os.Exit(1)
		}

		claims, ok := token.Claims.(*VerifyClaims)

		if !ok {
			fmt.Fprintf(
				os.Stderr,
				"invalid claims\n",
			)

			os.Exit(1)
		}

		if claims.ExpiresAt == nil {
			fmt.Fprintf(
				os.Stderr,
				"missing expiration\n",
			)

			os.Exit(1)
		}

		expired := time.Now().UTC().After(
			claims.ExpiresAt.Time,
		)

		result := map[string]interface{}{
			"valid":    true,
			"expired":  expired,
			"issuer":   claims.Issuer,
			"subject":  claims.Subject,
			"audience": claims.Audience,
			"scope":    claims.Scope,
			"issued": claims.IssuedAt.
				Time.
				UTC(),
			"expires": claims.ExpiresAt.
				Time.
				UTC(),
			"operator_id": claims.OperatorID,
		}

		if verifyOnline {

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

			result["online"] = verifyResp
		}

		if verifyJSON {

			_ = json.NewEncoder(os.Stdout).
				Encode(result)

			return
		}

		fmt.Printf(
			"token verification\n\n",
		)

		fmt.Printf(
			"valid: %v\n",
			true,
		)

		fmt.Printf(
			"expired: %v\n",
			expired,
		)

		fmt.Printf(
			"issuer: %s\n",
			claims.Issuer,
		)

		fmt.Printf(
			"subject: %s\n",
			claims.Subject,
		)

		fmt.Printf(
			"audience: %v\n",
			claims.Audience,
		)

		fmt.Printf(
			"scope: %s\n",
			claims.Scope,
		)

		fmt.Printf(
			"operator_id: %s\n",
			claims.OperatorID,
		)

		fmt.Printf(
			"issued_at: %s\n",
			claims.IssuedAt.Time.UTC(),
		)

		fmt.Printf(
			"expires_at: %s\n",
			claims.ExpiresAt.Time.UTC(),
		)

		if verifyOnline {

			fmt.Printf(
				"\nonline verification: success\n",
			)
		}
	},
}

func JWKToPublicKey(
	jwk JWK,
) (ed25519.PublicKey, error) {

	if jwk.KTY != "OKP" {
		return nil, fmt.Errorf(
			"unsupported key type",
		)
	}

	if jwk.CRV != "Ed25519" {
		return nil, fmt.Errorf(
			"unsupported curve",
		)
	}

	keyBytes, err := base64.RawURLEncoding.
		DecodeString(jwk.X)

	if err != nil {
		return nil, err
	}

	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf(
			"invalid public key size",
		)
	}

	return ed25519.PublicKey(
		keyBytes,
	), nil
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