package template

import (
	"fmt"
	"mediumkube/commands/handlers"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/utils"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	os.Remove("./test.yaml")
	config := configurations.LoadConfigFromFile("./test-config.yaml")
	handlers.DoRender("test.yaml.tmpl", "test.yaml", config)
	if utils.ReadStr("./test.gold.yaml") != utils.ReadStr("./test.yaml") {
		fmt.Println(utils.ReadStr("./test.gold.yaml"))
		fmt.Println(utils.ReadStr("./test.yaml"))
		t.Fail()
	}
	os.Remove("./test.yaml")
}
