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

func LoginWithToken(apiClient *api.Client, token string) {
	apiClient.SetToken(strings.TrimSpace(token))
} 


func LoginWithUserPass(apiClient *api.Client, credentials map[string]string) (string, error) {
	options := map[string]interface{}{
		"password": credentials["password"],
	}

	path := fmt.Sprintf("auth/%s/login/%s", credentials["userpass_path"] , credentials["username"])

	secret, err := apiClient.Logical().Write(path, options)

	if err != nil {
		return "", err
	}

	token := secret.Auth.ClientToken

	return token, nil
}


func NewVaultClient(vaultHost string, credentials map[string]string, PathParams string) (VaultClient, error) {
	apiClient, err := api.NewClient(&api.Config{
		Address: vaultHost,
	})

	if err != nil {
		return nil, err
	}

	if(credentials["auth_method"] == "token"){
		LoginWithToken(apiClient, credentials["token"])
	} else {
		token, err := LoginWithUserPass(apiClient, credentials)
		LoginWithToken(apiClient, token)
		if err != nil {
			return nil, fmt.Errorf("Invalid user or password")
		}
	}

	pathHandler := libraries.NewPathHandler(PathParams)
	
	vaultClient := &vaultClient{
		apiClient: apiClient,
		pathHandler: pathHandler,
	}
	
	return vaultClient, nil
}


func (c *vaultClient) QuerySecretMap(path string, parameters ...string) (map[string]interface{}, error) {
	var versionError string
	
	data, err := c.pathHandler.PathParamsParsing(parameters)
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
	
		if type_value == reflect.Slice || type_value == reflect.Map {
	
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

	data, err := c.pathHandler.PathParamsParsing(parameters)

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
