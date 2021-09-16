package libraries

import (
	"strings"
	"text/template"
	"bytes"
	"fmt"
)

type PathHandler interface {
	GetEnv() (string)
	PathV2(path string) (string)
	RenderPath(path string) (string, error)
}

type pathHandler struct {
	Environment string
}

type pathHandler2 struct {
	Environment string
}

func NewPathHandler(environment string) (PathHandler) {
	pathHandler := &pathHandler{
		Environment: environment,
	}

	return pathHandler
}

func (p *pathHandler) GetEnv() (string) {
	return p.Environment
}

func (p *pathHandler) RenderPath(path string) (string, error) {
	var renderedPath bytes.Buffer

	templatePath, err := template.New("path").Parse(path)

	if err != nil {
		return "", fmt.Errorf("Cannot parse %s", path)
	} 

	err = templatePath.Execute(&renderedPath, p)

	if err != nil {
		panic(err)
	}

	return renderedPath.String(), nil
}

func (p *pathHandler) PathV2(path string) (string) {
	pathSplited := strings.Split(path, string('/'))
	pathSplited[0] = pathSplited[0] + "/data"
	return strings.Join(pathSplited, string('/'))
}