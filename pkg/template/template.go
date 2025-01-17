package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"github.com/fathoniadi/vault-template/pkg/api"
	"os"
	"strings"
	"text/template"
)

type VaultTemplateRenderer struct {
	vaultClient api.VaultClient
	
}

func NewVaultTemplateRenderer(credentials map[string]string, vaultEndpoint string, dynamicPathVariable string) (*VaultTemplateRenderer, error) {

	vaultClient, err := api.NewVaultClient(vaultEndpoint, credentials, dynamicPathVariable)

	if err != nil {
		return nil, err
	}

	return &VaultTemplateRenderer{
		vaultClient: vaultClient,
	}, nil
}

func (v *VaultTemplateRenderer) RenderTemplate(templateContent string) (string, error) {
	funcMap := template.FuncMap{
		"vault":    v.vaultClient.QuerySecret,
		"vaultMap": v.vaultClient.QuerySecretMap,
	}

	tmpl, err := template.
		New("template").
		Funcs(sprig.TxtFuncMap()).
		Funcs(funcMap).
		Parse(templateContent)

	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer

	envMap := envToMap()
	if err := tmpl.Execute(&outputBuffer, envMap); err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}

func envToMap() map[string]string {
	envMap := map[string]string{}

	for _, v := range os.Environ() {
		splitV := strings.Split(v, "=")
		envMap[splitV[0]] = splitV[1]
	}

	return envMap
}
