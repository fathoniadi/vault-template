package api

import (
	"context"
	"encoding/json"
  "bytes"
  "strings"
	"fmt"
	"github.com/fathoniadi/vault-template/pkg/libraries"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
	"github.com/hashicorp/vault/api/auth/userpass"
)

func JSONMarshal(t interface{}) ([]byte, error) {
    buffer := &bytes.Buffer{}
    encoder := json.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    err := encoder.Encode(t)
    return buffer.Bytes(), err
}

type VaultClient interface {
	QuerySecret(path string, field string, parameters ...string) (interface{}, error)
	QuerySecretMap(path string, parameters ...string) (map[string]interface{}, error)
}

type vaultClient struct {
	apiClient   *api.Client
	pathHandler libraries.PathHandler
}

func LoginWithToken(apiClient *api.Client, token string) {
	apiClient.SetToken(strings.TrimSpace(token))
}

func LoginWithUserPass(apiClient *api.Client, credentials map[string]string) (*api.Secret, error) {
	userPassword := &userpass.Password{
		FromString: credentials["password"],
	}

	loginOption := userpass.WithMountPath(credentials["userpass_path"])

	userPassAuth, err := userpass.NewUserpassAuth(
		credentials["username"],
		userPassword,
		loginOption,
	)

	if err != nil {
		return nil, err
	}

	authInfo, err := apiClient.Auth().Login(context.TODO(), userPassAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return authInfo, nil
}

func LoginWithApprole(apiClient *api.Client, credentials map[string]string) (*api.Secret, error) {

	appRoleSecretID := &approle.SecretID{
		FromString: credentials["approlesecretid"],
	}

	appRoleAuth, err := approle.NewAppRoleAuth(
		credentials["approleid"],
		appRoleSecretID,
	)

	if err != nil {
		return nil, err
	}

	authInfo, err := apiClient.Auth().Login(context.TODO(), appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login to AppRole auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	return authInfo, nil
}

func NewVaultClient(vaultHost string, credentials map[string]string, dynamicPathVariable string) (VaultClient, error) {
	apiClient, err := api.NewClient(&api.Config{
		Address: vaultHost,
	})

	if err != nil {
		return nil, err
	}

	if credentials["auth_method"] == "token" {
		LoginWithToken(apiClient, credentials["token"])
	} else if credentials["auth_method"] == "userpass" {
		_, err := LoginWithUserPass(apiClient, credentials)
		if err != nil {
			if(strings.Contains(fmt.Sprint(err), "timeout")){
				fmt.Println("timeout")
			}
			return nil, fmt.Errorf("Invalid user or password: %s", err)
		}
	} else {
		_, err := LoginWithApprole(apiClient, credentials)
		if err != nil {
			return nil, fmt.Errorf("Invalid AppRole ID or secret: %s", err)
		}
	}

	pathHandler := libraries.NewPathHandler(dynamicPathVariable)

	vaultClient := &vaultClient{
		apiClient:   apiClient,
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

	if err != nil {
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
		jsonData, err := JSONMarshal(value)

		if err != nil {
			return nil, fmt.Errorf("Error parsing field %s", key)
		}

		m[key] = string(jsonData)
	}

	return m, nil
}

func (c *vaultClient) QuerySecret(path string, field string, parameters ...string) (interface{}, error) {
	var versionError string

	data, err := c.pathHandler.PathParamsParsing(parameters)

	if _, ok := data["version"]; ok {
		versionError = " in version " + data["version"][0]
	}
	if err != nil {
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
		return nil, fmt.Errorf("secret at path '%s'%s has no field '%s'", path, versionError, field)
	}

	secretValue, ok := secret.Data[field]
	if !ok {
		m, ok := secret.Data["data"].(map[string]interface{})

		if !ok {
			fmt.Printf("%T %#v\n", secret.Data["data"], secret.Data["data"])
			return nil, fmt.Errorf("error reading path '%s'", path)
		}

		secretValue, ok := m[field]

		if !ok {
			return nil, fmt.Errorf("secret at path '%s'%s has no field '%s'", path, versionError, field)
		}

		jsonData, err := JSONMarshal(secretValue)


		if err != nil {
			return nil, fmt.Errorf("Error parsing field %s", field)
		}
		return string(jsonData), nil
	}

	jsonData, err := JSONMarshal(secretValue)

	if err != nil {
		return nil, fmt.Errorf("Error parsing field %s", field)
	}
	return string(jsonData), nil
}
