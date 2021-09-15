package api

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
)

type VaultClient interface {
	QuerySecret(path string, field string) (string, error)
	QuerySecretMap(path string) (map[string]interface{}, error)
}

type vaultClient struct {
	apiClient *api.Client
	secretVersion string
}

func NewVaultClient(vaultEndpoint string, vaultToken string, secretVersion string) (VaultClient, error) {
	apiClient, err := api.NewClient(&api.Config{
		Address: vaultEndpoint,
	})

	if err != nil {
		return nil, err
	}

	apiClient.SetToken(strings.TrimSpace(vaultToken))

	vaultClient := &vaultClient{
		apiClient: apiClient,
		secretVersion: "?version=" + secretVersion 
	}

	return vaultClient, nil
}

func (c *vaultClient) QuerySecretMap(path string) (map[string]interface{}, error) {
	secret, err := c.apiClient.Logical().Read(path + c.secretVersion)

	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("path '%s' is not found in version '%s'", path, c.secretVersion)
	}

	return secret.Data, nil
}

func (c *vaultClient) QuerySecret(path string, field string) (string, error) {
	secret, err := c.apiClient.Logical().Read(path + c.secretVersion)

	if err != nil {
		return "", err
	}

	secretValue, ok := secret.Data[field]
	if !ok {
		m, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			fmt.Printf("%T %#v\n", secret.Data["data"], secret.Data["data"])
			return "", fmt.Errorf("error reading path '%s'", path)
		}

		secretValue, ok := m[field]

		if !ok {
			return "", fmt.Errorf("secret at path '%s' in version '%s' has no field '%s'", path, c.secretVersion, field)
		}
		return secretValue.(string), nil
	}

	return secretValue.(string), nil
}
