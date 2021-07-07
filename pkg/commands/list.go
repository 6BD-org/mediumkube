package commands

import (
	"fmt"
	"mediumkube/pkg/configurations"
	"mediumkube/pkg/daemon/mesh"
	"mediumkube/pkg/models"
	"mediumkube/pkg/utils"
	"os"

	"github.com/olekukonko/tablewriter"
)

type ListHandler struct {
}

func disp(resp []models.Domain) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name", "IP", "STATE", "REASON",
	})
	for _, d := range resp {
		table.Append([]string{
			d.Name, d.IP, d.Status, d.Reason,
		})
	}
	table.Render()
}

func (handler ListHandler) Handle(args []string) {
	config := configurations.Config()
	domains, err := mesh.ListDomains(config)
	utils.CheckErr(err)
	disp(domains)
}
func (handler ListHandler) Help() {
	fmt.Println("list")
}
func (handler ListHandler) Desc() string {
	return "List nodes"
}

func init() {
	name := "list"
	CMD[name] = ListHandler{}
}
