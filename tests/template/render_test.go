package template

import (
	"mediumkube/pkg/commands"
	"mediumkube/pkg/utils"
	"os"
	"testing"
)

func TestRender(t *testing.T) {
	os.Remove("./test.yaml")
	commands.Render("test.yaml.tmpl", "test.yaml")
	if utils.ReadStr("./test.gold.yaml") != utils.ReadStr("./test.yaml") {
		t.Fail()
	}
	os.Remove("./test.yaml")
}
