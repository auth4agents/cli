package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AgentConfig - for agent identities
type AgentConfig struct {
	DID        string `json:"did"`
	PrivateKey string `json:"private_key"`
	ServerURL  string `json:"server_url"`
	OperatorID string `json:"operator_id"`
}

// OperatorConfig - for operator identities  
type OperatorConfig struct {
	ID               string `json:"id"`
	Domain           string `json:"domain"`
	RootPublicKey    string `json:"root_public_key"`
	RootPrivateKey   string `json:"root_private_key"`
	ServerURL        string `json:"server_url"`
	DomainVerifiedAt string `json:"domain_verified_at,omitempty"`
}

func GetAgentConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agentauth", "agent.json"), nil
}

func GetOperatorConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agentauth", "operator.json"), nil
}

func SaveAgentConfig(cfg *AgentConfig) error {
	path, err := GetAgentConfigPath()
	if err != nil {
		return err
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0600)
}

func LoadAgentConfig() (*AgentConfig, error) {
	path, err := GetAgentConfigPath()
	if err != nil {
		return nil, err
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var cfg AgentConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}

func SaveOperatorConfig(cfg *OperatorConfig) error {
	path, err := GetOperatorConfigPath()
	if err != nil {
		return err
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0600)
}

func LoadOperatorConfig() (*OperatorConfig, error) {
	path, err := GetOperatorConfigPath()
	if err != nil {
		return nil, err
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var cfg OperatorConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	
	return &cfg, nil
}