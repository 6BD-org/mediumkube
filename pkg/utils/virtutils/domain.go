package virtutils

import (
	"bytes"
	"fmt"
	"mediumkube/pkg/common"
	"text/template"

	"github.com/Masterminds/sprig"
)

// GetDeploymentConfig Render a domain deployment config
func GetDeploymentConfig(param common.DomainCreationParam, tmplStr string) (string, error) {

	tmpl, err := template.New("domain").Funcs(sprig.TxtFuncMap()).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("Unable to generate template")
	}
	buffer := bytes.NewBuffer(make([]byte, 0))
	err = tmpl.Execute(buffer, param)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
