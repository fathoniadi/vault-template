package main

import (
	"github.com/Luzifer/rconfig"
	"github.com/fathoniadi/vault-template/pkg/template"
	"io/ioutil"
	"log"
	"os"
)

var (
	cfg = struct {
		VaultHost  string `flag:"host,h" env:"VAULT_HOST" default:"https://127.0.0.1:8200" description:"Vault API endpoint. Also configurable via VAULT_HOST."`
		VaultToken string `flag:"token,k" env:"VAULT_TOKEN" description:"File containt vault token. Also configurable via VAULT_TOKEN."`
		TemplateFile   string `flag:"template,t" env:"TEMPLATE_FILE" description:"The template file to render. Also configurable via TEMPLATE_FILE."`
		OutputFile     string `flag:"output,o" env:"OUTPUT_FILE" description:"The output file. Also configurable via OUTPUT_FILE."`
		DynamicPathVariable    string `flag:"path-params,p" env:"DYNAMICPATHVARIABLE" description:"Dynamic variable path templating. Also configurable via DYNAMICPATHVARIABLE. Ex. \"project=blog,environment=development\""`
		Username    string `flag:"username,U" env:"USERNAME" description:"Username to login. Also configurable via USERNAME."`
		Password    string `flag:"password,W" env:"PASSWORD" description:"Password to login. Also configurable via PASSWORD."`
		UserPassPath    string `flag:"userpass-path,P" env:"USERPASS_PATH" default:"userpass" description:"Path user was registered. Also configurable via USERPASS_PATH."`

	}{}
)

func usage(msg string) {
	println(msg)
	rconfig.Usage()
	os.Exit(1)
}

func config() (map[string]string) {
	var credentials map[string]string = make(map[string]string)

	var useUserPass  bool = true
	var useToken bool = true

	err := rconfig.Parse(&cfg)
	
	if err != nil {
		log.Fatalf("Error while parsing the command line arguments: %s", err)
	}

	if cfg.VaultToken == "" {
		useToken = false
	}

	if cfg.Username == "" || cfg.Password == "" {
		useUserPass  = false
	}

	if !useUserPass  && !useToken {
		usage("No Auth method declared")
	}

	if useUserPass {
		credentials["username"] = cfg.Username
		credentials["password"] = cfg.Password
		credentials["userpass_path"] = cfg.UserPassPath
		credentials["auth_method"] = "userpass"

	}

	if useToken {
		credentials["token"] = cfg.VaultToken
		credentials["auth_method"] = "token"
	}

	if cfg.TemplateFile == "" {
		usage("No template file given")
	}

	if cfg.OutputFile == "" {
		usage("No output file given")
	}

	return credentials
}

func main() {

	credentials := config()

	renderer, err := template.NewVaultTemplateRenderer(credentials, cfg.VaultHost, cfg.DynamicPathVariable)

	if err != nil {
		log.Fatalf("Unable to create renderer: %s", err)
	}

	templateContent, err := ioutil.ReadFile(cfg.TemplateFile)

	if err != nil {
		log.Fatalf("Unable to read template file: %s", err)
	}

	renderedContent, err := renderer.RenderTemplate(string(templateContent))

	if err != nil {
		log.Fatalf("Unable to render template: %s", err)
	}

	outputFile, err := os.Create(cfg.OutputFile)

	if err != nil {
		log.Fatalf("Unable to write output file: %s", err)
	}

	defer func() {
		err := outputFile.Close()
		if err != nil {
			log.Fatalf("Error while closing the output file: %s", err)
		}
	}()

	_, err = outputFile.Write([]byte(renderedContent))

	if err != nil {
		log.Fatalf("Error while writing the output file: %s", err)
	}
}
