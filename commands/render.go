package commands

import (
	"flag"
	"io/ioutil"
	"log"
	"mediumkube/common"
	"mediumkube/utils"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

func parseConfig(configPath string) common.TemplateConfig {
	data := utils.ReadStr(configPath)

	var config common.TemplateConfig = common.TemplateConfig{}
	err := yaml.Unmarshal([]byte(data), &config)
	utils.CheckErr(err)

	config.PubKey = utils.ReadStr(config.PubKeyDir)
	config.PrivKey = utils.ReadStr(config.PrivKeyDir)
	config.HostPubKey = utils.ReadStr(config.HostPubKeyDir)

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
func Render(configPath string, templatePath string, outPath string) {
	config := parseConfig(configPath)
	templateStr := readTemplate(templatePath)

	tmpl, err := template.New("cloudInit").Funcs(sprig.TxtFuncMap()).Parse(templateStr)
	utils.CheckErr(err)

	var out *os.File
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

	configPath := handler.flagSet.String("config", "./config.yaml", "Path to config file")
	templatePath := handler.flagSet.String("template", "./cloud-init.yaml.tmpl", "Path to cloud init yaml template")
	outPath := handler.flagSet.String("out", "./cloud-init.yaml", "Path of output yaml")

	handler.flagSet.Parse(args)

	if Help(handler, args) {
		return
	}

	Render(*configPath, *templatePath, *outPath)

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
