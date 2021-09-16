package api

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"strings"
	"reflect"
	"encoding/json"
	"github.com/fathoniadi/vault-template/pkg/libraries"

)

type VaultClient interface {
	QuerySecret(path string, field string, parameters ...string) (string, error)
	QuerySecretMap(path string, parameters ...string) (map[string]interface{}, error)
}

type vaultClient struct {
	apiClient *api.Client
	pathHandler libraries.PathHandler
}


func NewVaultClient(vaultEndpoint string, vaultToken string, envPath string) (VaultClient, error) {

	apiClient, err := api.NewClient(&api.Config{
		Address: vaultEndpoint,
	})

	pathHandler := libraries.NewPathHandler(envPath)

	if err != nil {
		return nil, err
	}

	apiClient.SetToken(strings.TrimSpace(vaultToken))

	vaultClient := &vaultClient{
		apiClient: apiClient,
		pathHandler: pathHandler,
	}

	return vaultClient, nil
}


func (c *vaultClient) QuerySecretMap(path string, parameters ...string) (map[string]interface{}, error) {
	var versionError string

	data, err := libraries.ParamParsingHandler(parameters)

	if _, ok := data["version"]; ok {
		versionError = " in version " + data["version"][0]
	}

	if err != nil  {
		return nil, err
	}

	path, err = c.pathHandler.RenderPath(path)

	if err != nil {
		return nil, err
	}

	secret, err := c.apiClient.Logical().ReadWithData(string(c.pathHandler.PathV2(path)), data)

	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, fmt.Errorf("path '%s' is not found'%s'", path, versionError)
	}

	m, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		fmt.Printf("%T %#v\n", secret.Data["data"], secret.Data["data"])
		return nil, fmt.Errorf("error reading path '%s'", path)
	}

	for key, value := range m {
		type_value := reflect.TypeOf(value).Kind()

		if type_value == reflect.Slice {
			json_data, err := json.Marshal(value)
			if err != nil {
				return nil, fmt.Errorf("Error parsing field %s", key)
			}

			m[key] = string(json_data)

		}

		if type_value == reflect.Map {
			json_data, err := json.Marshal(value)
			if err != nil {
				return nil, fmt.Errorf("Error parsing field %s", key)
			}

			m[key] = string(json_data)
		}
	}

	return m, nil
}

func (c *vaultClient) QuerySecret(path string, field string, parameters ...string) (string, error) {

	var versionError string

	data, err := libraries.ParamParsingHandler(parameters)

	if _, ok := data["version"]; ok {
		versionError = " in version " + data["version"][0]
	}

	if err != nil  {
		return "", err
	}

	path, err = c.pathHandler.RenderPath(path)

	if err != nil {
		return "", err
	}

	secret, err := c.apiClient.Logical().ReadWithData(string(c.pathHandler.PathV2(path)), data)


	if secret == nil {
		return "", fmt.Errorf("secret at path '%s'%s has no field '%s'", path, versionError, field)
	}

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
			return "", fmt.Errorf("secret at path '%s'%s has no field '%s'", path, versionError, field)
		}

		typeValue := reflect.TypeOf(secretValue).Kind()

		if (typeValue == reflect.String) {
			return secretValue.(string), nil
		}
		
		json_data, err := json.Marshal(secretValue)

		if err != nil {
			return "", fmt.Errorf("Error parsing field %s", field)
		}

		return string(json_data), nil
	}

	return secretValue.(string), nil
}
