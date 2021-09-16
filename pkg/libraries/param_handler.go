package libraries

import (
	"strings"
	"fmt"
)


func ParamParsingHandler(parameters []string) (map[string][]string, error) {

	var data map[string][]string = make(map[string][]string)

	if len(parameters) > 0 {
		paramatersSplited := strings.Split(parameters[0], string(','))

		for _, arguments := range paramatersSplited {
			arguments_data := strings.Split(arguments, string(':'))
			if len(arguments_data) == 1 || len(arguments_data) > 2 {
				return data, fmt.Errorf("Error parsing parameter '%s'", arguments_data[0])
			}

			data[arguments_data[0]] = []string{arguments_data[1]}
		}
	}

	return data, nil
}



