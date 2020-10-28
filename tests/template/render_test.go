package template

import (
	"mediumkube/commands"
	"mediumkube/utils"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	os.Remove("./test.yaml")
	commands.Render("./test-config.yaml", "test.yaml.tmpl", "test.yaml")
	if utils.ReadStr("./test.gold.yaml") != utils.ReadStr("./test.yaml") {
		t.Fail()
	}
}
