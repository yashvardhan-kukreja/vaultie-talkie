package main

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	"github.com/yashvardhan-kukreja/vaultie-talkie/internal/target"
)

type VaultSettings struct {
	Host string
	Port int64

	PathToWatch string
	AccessToken string
}

func (v VaultSettings) InitClient() (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = fmt.Sprintf("http://%s:%d", v.Host, v.Port)

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting the new client from vault: %w", err)
	}
	client.SetToken(v.AccessToken)
	return client, nil
}

func renderKeyStore(client *vault.Client, path string) (target.KeyStore, error) {
	secret, err := client.KVv2("secret").Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("error occurred while listing the secret contents at the path '%s': %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, nil
	}
	return target.KeyStore(secret.Data), nil
}
