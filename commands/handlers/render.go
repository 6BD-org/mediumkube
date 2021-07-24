package handlers

import (
	"flag"
	"io/ioutil"
	"log"
	"mediumkube/pkg/common"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

func parseConfig(configPath string) common.OverallConfig {
	data := utils.ReadStr(configPath)

	var config common.OverallConfig = common.OverallConfig{}
	err := yaml.Unmarshal([]byte(data), &config)
	utils.CheckErr(err)

	return config
}

func readTemplate(templatePath string) string {

	data, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Panic("Error reading file")
	}

	return string(data)
}

// Render render the template using config object
func Render(templatePath string, outPath string) {
	config := configurations.Config()
	templateStr := readTemplate(templatePath)

	tmpl, err := template.New("cloudInit").Funcs(sprig.TxtFuncMap()).Parse(templateStr)
	utils.CheckErr(err)

	var out *os.File
	os.Remove(outPath)
	out, err = os.Create(outPath)
	utils.CheckErr(err)
	err = tmpl.Execute(out, config)
	utils.CheckErr(err)
}

// RenderHandler This handles render command
type RenderHandler struct {
	flagSet *flag.FlagSet
}

func (handler RenderHandler) Handle(args []string) {

	templatePath := handler.flagSet.String("template", "./cloud-init.yaml.tmpl", "Path to cloud init yaml template")
	outPath := handler.flagSet.String("out", "./cloud-init.yaml", "Path of output yaml")

	handler.flagSet.Parse(args[1:])

	if Help(handler, args) {
		return
	}

	Render(*templatePath, *outPath)

}

func (handler RenderHandler) Desc() string {
	return "render cloud init template"
}

func (handler RenderHandler) Help() {
	handler.flagSet.Usage()
}

func init() {
	var name = "render"
	// Register to root
	CMD[name] = RenderHandler{
		flagSet: flag.NewFlagSet(name, flag.ExitOnError),
	}
}
