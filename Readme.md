# `vault-template`

Render templated config files with secrets from [HashiCorp Vault](https://www.vaultproject.io/). Inspired by [vaultenv](https://github.com/channable/vaultenv).

* Define a template for your config file which contains secrets at development time.
* Use `vault-template` to render your config file template by fetching secrets from Vault at runtime.

## Usage

```text
Usage of ./vault-template:
  -o, --output string             The output file.
                                  Also configurable via OUTPUT_FILE.
  -t, --template string           The template file to render.
                                  Also configurable via TEMPLATE_FILE.
  -v, --vault string              Vault API endpoint.
                                  Also configurable via VAULT_ADDR.
                                  (default "http://127.0.0.1:8200")
  -f, --vault-token string        The vault token.
                                  Also configurable via VAULT_TOKEN_FILE.
  
  -u, --username string           Username to login
                                  Also configurable via USERNAME

  -e, --environment string        Path environment templating
                                  Also configurable via ENVIRONMENT

  -p, --password string           Password to login
                                  Also configurable via PASSWORD

  -P --userpass-path string       Path user was registered. 
                                  Also configurable via USERPASS_PATH.
                                  (default "userpass")
```

A [docker image is availabe on Dockerhub.](https://hub.docker.com/r/rplan/vault-template)

## Template

First of all, suppose that the secret was created with `vault write secret/mySecret name=john password=secret`.

The templates will be rendered using the [Go template](https://golang.org/pkg/text/template/) mechanism.

Currently vault-template can render two functions:
- `vault`
- `vaultMap`

The `vault` function takes two string parameters which specify the path to the secret and the field inside to return.

```gotemplate
mySecretName = {{ vault "secret/mySecret" "name" }}
mySecretPassword = {{ vault "secret/mySecret" "password" }}
```

```text
mySecretName = john
mySecretPassword = secret
```

with specific version of secret:

```gotemplate
mySecretName = {{ vault "secret/mySecret" "name" "version:1" }}
mySecretPassword = {{ vault "secret/mySecret" "password" }}
```

```text
mySecretName = johni
mySecretPassword = secret
```


The `vaultMap` function takes one string parameter which specify the path to the secret to return.

```gotemplate
{{ range $name, $secret := vaultMap "secret/mySecret"}}
{{ $name }}: {{ $secret }}
{{- end }}
```

```text
name: john
password: secret
```

with specific version:
```gotemplate
{{ range $name, $secret := vaultMap "secret/mySecret" "version:1"}}
{{ $name }}: {{ $secret }}
{{- end }}
```

```text
name: johni
password: secret
```

More real example:

```gotemplate
---

{{ range $name, $secret := vaultMap "thoni-website/{{ .Environment }}/database" "version:1"}}
{{ $name }}: {{ $secret }}
{{- end }}

```

And command that use this template in kubernetes:
```
vault-template -o values.yaml -t values.tmpl -v "http://vault.default.svc.cluster.local:8200" -f "$(token)" -e development
```
