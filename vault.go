package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault-client-go"
	"time"
)

func readVaultSecret(secretName, secretPath string) (map[string]string, error) {
	var result = map[string]string{}
	ctx := context.Background()
	client, err := vault.New(
		vault.WithAddress(configuration.Vault.Server),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// Аутентифкация (переделать)
	if err := client.SetToken(configuration.Vault.Token); err != nil {
		return nil, errors.New(err.Error())
	}

	//Читаем секрет
	s, err := client.Secrets.KvV2Read(ctx, secretName, vault.WithMountPath(secretPath))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	for a, b := range s.Data.Data {
		result[a] = fmt.Sprintln(b)
	}
	fmt.Println(result)
	return result, nil
}
