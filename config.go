package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	ConfigDirName = ".auth4agents"

	AgentConfigFile    = "agent.json"
	OperatorConfigFile = "operator.json"

	DefaultKeyProvider = "file"
)

type KeyConfig struct {
	Provider string `json:"provider"`
	Ref      string `json:"ref"`
}

type AgentConfig struct {
	DID        string `json:"did"`
	OperatorID string `json:"operator_id"`
	ServerURL  string `json:"server_url"`

	PublicKey string `json:"public_key"`

	KeyConfig KeyConfig `json:"key"`

	DIDDocument *DIDDocument `json:"did_document"`
}

type OperatorConfig struct {
	ID     string `json:"id"`
	Domain string `json:"domain"`

	ServerURL string `json:"server_url"`

	RootPublicKey string `json:"root_public_key"`

	KeyConfig KeyConfig `json:"key"`

	DomainVerifiedAt string `json:"domain_verified_at,omitempty"`
}

func GetConfigDir() (string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		home,
		ConfigDirName,
	), nil
}

func EnsureConfigDir() error {

	dir, err := GetConfigDir()
	if err != nil {
		return err
	}

	return os.MkdirAll(
		dir,
		0700,
	)
}

func GetAgentConfigPath() (string, error) {

	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		dir,
		AgentConfigFile,
	), nil
}

func GetOperatorConfigPath() (string, error) {

	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		dir,
		OperatorConfigFile,
	), nil
}

func GetKeysDir() (string, error) {

	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(
		dir,
		"keys",
	), nil
}

func EnsureKeysDir() error {

	keysDir, err := GetKeysDir()
	if err != nil {
		return err
	}

	return os.MkdirAll(
		keysDir,
		0700,
	)
}

func SaveAgentConfig(
	cfg *AgentConfig,
) error {

	path, err := GetAgentConfigPath()
	if err != nil {
		return err
	}

	return SaveJSONFile(
		path,
		cfg,
	)
}

func LoadAgentConfig() (
	*AgentConfig,
	error,
) {

	path, err := GetAgentConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg AgentConfig

	err = LoadJSONFile(
		path,
		&cfg,
	)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveOperatorConfig(
	cfg *OperatorConfig,
) error {

	path, err := GetOperatorConfigPath()
	if err != nil {
		return err
	}

	return SaveJSONFile(
		path,
		cfg,
	)
}

func LoadOperatorConfig() (
	*OperatorConfig,
	error,
) {

	path, err := GetOperatorConfigPath()
	if err != nil {
		return nil, err
	}

	var cfg OperatorConfig

	err = LoadJSONFile(
		path,
		&cfg,
	)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveJSONFile(
	path string,
	value interface{},
) error {

	err := EnsureConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(
		value,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0600,
	)
}

func LoadJSONFile(
	path string,
	out interface{},
) error {

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(
		data,
		out,
	)
}

func SavePrivateKeyFile(
	filename string,
	privateKey string,
) (string, error) {

	err := EnsureKeysDir()
	if err != nil {
		return "", err
	}

	keysDir, err := GetKeysDir()
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(
		keysDir,
		filename,
	)

	err = os.WriteFile(
		fullPath,
		[]byte(privateKey),
		0600,
	)

	if err != nil {
		return "", err
	}

	return fullPath, nil
}

func LoadPrivateKeyFile(
	path string,
) (string, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
