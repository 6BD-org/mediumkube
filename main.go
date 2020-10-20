package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v3"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// TemplateConfig this is parsed from config yaml
type TemplateConfig struct {
	HTTPSProxy    string `yaml:"https-proxy,omitempty"`
	HTTPProxy     string `yaml:"http-proxy,omitempty"`
	PubKeyDir     string `yaml:"pub-key-dir"`
	PrivKeyDir    string `yaml:"priv-key-dir"`
	HostPubKeyDir string `yaml:"host-pub-key-dir"`

	PubKey     string
	PrivKey    string
	HostPubKey string
}

func readStr(path string) string {
	data, err := ioutil.ReadFile(path)
	check(err)
	return string(data)
}

func parseConfig(configPath string) TemplateConfig {
	data, err := ioutil.ReadFile(configPath)
	check(err)

	var config TemplateConfig = TemplateConfig{}
	err = yaml.Unmarshal(data, &config)
	check(err)

	config.PubKey = readStr(config.PubKeyDir)
	config.PrivKey = readStr(config.PrivKeyDir)
	config.HostPubKey = readStr(config.HostPubKeyDir)

	return config
}

func readTemplate(templatePath string) string {

	data, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Panic("Error reading file")
	}

	return string(data)
}

func testMarshal() {
	config := TemplateConfig{
		HTTPProxy:  "http",
		HTTPSProxy: "https",
	}
	data, _ := yaml.Marshal(config)
	fmt.Print(string(data))
}

func main() {

	templatePath := flag.String("template", "./cloud-init.yaml.tmpl", "Path to cloud init yaml template")
	configPath := flag.String("config", "./config.yaml", "Path to config file")
	outPath := flag.String("out", "./cloud-init.yaml", "Path of output yaml")

	config := parseConfig(*configPath)
	templateStr := readTemplate(*templatePath)

	tmpl, err := template.New("cloudInit").Funcs(sprig.TxtFuncMap()).Parse(templateStr)
	check(err)

	var out *os.File
	out, err = os.Create(*outPath)
	check(err)

	err = tmpl.Execute(out, config)
	check(err)

}

func init() {

}
