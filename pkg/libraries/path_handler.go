package libraries

import (
	"strings"
	"text/template"
	"bytes"
	"fmt"
)

type PathHandler interface {
	getDynamicPathVariable() (string)
	PathV2(path string) (string)
	RenderPath(path string) (string, error)
	GetDynamicParamParsed() (map[string]interface{}, error)
	PathParamsParsing(parameters []string) (map[string][]string, error)
}

type pathHandler struct {
	dynamiCPathVariable string
}

func NewPathHandler(dynamiCPathVariable string) (PathHandler) {

	pathHandler := &pathHandler{
		dynamiCPathVariable: dynamiCPathVariable,
	}

	return pathHandler
}

func (p *pathHandler) getDynamicPathVariable() (string) {
	return p.dynamiCPathVariable
}

func (p *pathHandler) GetDynamicParamParsed() (map[string]interface{}, error){
	var data map[string]interface{} = make(map[string]interface{})

	if (p.getDynamicPathVariable() == ""){
		return data, nil
	}

	params := strings.Split(p.getDynamicPathVariable(), string(','))

	for _, valueParams := range params {


		param := strings.Split(valueParams, string('='))

		if len(param) <= 1 {
			return nil, fmt.Errorf("Invalid dynamic path variable templating declaration %s", p.getDynamicPathVariable() )
		}

		data[param[0]] = param[1]
	}

	return data, nil

}


func (p *pathHandler) PathParamsParsing(parameters []string) (map[string][]string, error) {

	var data map[string][]string = make(map[string][]string)

	if len(parameters) > 0 {
		params := strings.Split(parameters[0], string(','))

		for _, valueParams := range params {

			param := strings.Split(valueParams, string('='))
			if len(param) <= 1 {
				return nil, fmt.Errorf("Invalid path parameter declaration %s", p.getDynamicPathVariable() )
			}
	
			data[param[0]] = []string{param[1]}
		}
	}

	return data, nil

}

func (p *pathHandler) RenderPath(path string) (string, error) {
	var renderedPath bytes.Buffer

	params, err := p.GetDynamicParamParsed()

	if err != nil {
		return "", err
	}

	if len(params) == 0 {
		return path, nil
	}

	templatePath, err := template.New("path").Parse(path)

	if err != nil {
		return "", fmt.Errorf("Cannot parse %s", path)
	} 

	err = templatePath.Execute(&renderedPath, params)

	if err != nil {
		return "", fmt.Errorf("Cannot parse %s", err.Error())
	}

	return renderedPath.String(), nil
}

func (p *pathHandler) PathV2(path string) (string) {
	pathSplited := strings.Split(path, string('/'))
	pathSplited[0] = pathSplited[0] + "/data"
	return strings.Join(pathSplited, string('/'))
}