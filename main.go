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
		VaultEndpoint  string `flag:"vault,v" env:"VAULT_ADDR" default:"https://127.0.0.1:8200" description:"Vault API endpoint. Also configurable via VAULT_ADDR."`
		VaultToken string `flag:"vault-token,f" env:"VAULT_TOKEN" description:"The vault token. Also configurable via VAULT_TOKEN."`
		TemplateFile   string `flag:"template,t" env:"TEMPLATE_FILE" description:"The template file to render. Also configurable via TEMPLATE_FILE."`
		OutputFile     string `flag:"output,o" env:"OUTPUT_FILE" description:"The output file. Also configurable via OUTPUT_FILE."`
		Environment    string `flag:"env,e" env:"ENVIRONMENT" description:"The output file. Also configurable via ENVIRONMENT."`
		Username    string `flag:"username,u" env:"USERNAME" description:"The output file. Also configurable via USERNAME."`
		Password    string `flag:"password,p" env:"PASSWORD" description:"The output file. Also configurable via PASSWORD."`
		UserPassPath    string `flag:"userpass-path,P" env:"USERPASS_PATH" description:"The output file. Also configurable via USERPASS_PATH."`

	}{}
)

func usage(msg string) {
	println(msg)
	rconfig.Usage()
	os.Exit(1)
}

func config() {
	err := rconfig.Parse(&cfg)

	if err != nil {
		log.Fatalf("Error while parsing the command line arguments: %s", err)
	}

	if cfg.VaultToken == "" {
		usage("No vault token file given")
	}

	if cfg.TemplateFile == "" {
		usage("No template file given")
	}

	if cfg.OutputFile == "" {
		usage("No output file given")
	}

	if cfg.Environment == "" {
		usage("No Environment was set")
	}


}

func main() {
	config()

	vaultToken := cfg.VaultToken

	renderer, err := template.NewVaultTemplateRenderer(string(vaultToken), cfg.VaultEndpoint, cfg.Environment)

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
