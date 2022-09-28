# `vault-template`

Render templated config files with secrets from [HashiCorp Vault](https://www.vaultproject.io/). Inspired by [vaultenv](https://github.com/channable/vaultenv).

* Define a template for your config file which contains secrets at development time.
* Use `vault-template` to render your config file template by fetching secrets from Vault at runtime.

## Usage

```text
Usage of vault-template:
  -i, --approleid string       AppRole ID. Also configurable via AppRoleID.
  -s, --approlesecret string   AppRole Secret ID. Also via AppRoleSecret.
  -h, --host string            Vault API endpoint. Also configurable via VAULT_HOST. (default "https://127.0.0.1:8200")
  -o, --output string          The output file. Also configurable via OUTPUT_FILE.
  -W, --password string        Password to login. Also configurable via PASSWORD.
  -p, --path-params string     Dynamic variable path templating. Also configurable via DYNAMICPATHVARIABLE. Ex. "project=blog,environment=development"
  -t, --template string        The template file to render. Also configurable via TEMPLATE_FILE.
  -k, --token string           File containt vault token. Also configurable via VAULT_TOKEN.
  -U, --username string        Username to login. Also configurable via USERNAME.
  -P, --userpass-path string   Path user was registered. Also configurable via USERPASS_PATH. (default "userpass")
```

[docker image is availabe on Dockerhub.](https://hub.docker.com/r/fathoniadi/vault-template)

## Auth

The `vault-template` support two mechanism vault's authentication, using token or username and password.

Example using token:

```
vault-template -t ./env -o ./env.ready -h "http://localhost:8200" -k "$(cat ./token)"
```

Example using username and password

```
vault-template -t ./env -o ./env.ready -h "http://localhost:8200" -U fathoniadi -W "$(cat ./password)"
```

Example using AppRole
```
./vault-template -t ./env -o ./env.ready --approleid 00108aa6-1234-1234-1234-66efefd00000 --approlesecret 0056cccc-b017-1234-bbbb-be2febe00000 
```


## Template

First of all, suppose that the secret was created with `vault write secret/mySecret name=john password=secret`.

The templates will be rendered using the [Go template](https://golang.org/pkg/text/template/) mechanism.

Currently vault-template can render two functions:
- `vault`
- `vaultMap`

The `vault` function takes three string parameters which specify the path to the secret, the field inside to return and the version of secret.

```gotemplate
mySecretName = {{ vault "secret/mySecret" "name" }}
mySecretPassword = {{ vault "secret/mySecret" "password" }}
```

```text
mySecretName = john
mySecretPassword = secret
```


Note:
If you don't specify the version of secret, the `vault` function will return value from the latest version of secret


with specific version of secret:

```gotemplate
mySecretName = {{ vault "secret/mySecret" "name" "version=1" }}
mySecretPassword = {{ vault "secret/mySecret" "password" }}
```

```text
mySecretName = johni
mySecretPassword = secret
```


The `vaultMap` function takes two string parameter which specify the path to the secret to return and the version of secret.



```gotemplate
{{ range $name, $secret := vaultMap "secret/mySecret"}}
{{ $name }}: {{ $secret }}
{{- end }}
```

```text
name: john
password: secret
```

Note:
If you don't specify the version of secret, the `vaultMap` function will return value from the latest version of secret

with specific version:
```gotemplate
{{ range $name, $secret := vaultMap "secret/mySecret" "version=1"}}
{{ $name }}: {{ $secret }}
{{- end }}
```

```text
name: johni
password: secret
```

### Path Paramaterized 

```gotemplate
---

{{ range $name, $secret := vaultMap "{{ .project }}/{{ .environment }}/database" "version=1"}}
{{ $name }}: {{ $secret }}
{{- end }}

```

And command that use this feature:

```
vault-template -o values.yaml -t values.tmpl -v "http://localhost:8200" -k "$(token)" -p "project=devops,environment=development"
```


Forked from `vault-template` by [Actano GmbH](https://github.com/actano) and [Minh-Danh](https://github.com/minhdanh). Thanks for your great work.
