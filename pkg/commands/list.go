package commands

import (
	"fmt"
	"mediumkube/pkg/models"
	"mediumkube/pkg/services"
	"mediumkube/pkg/utils"
	"os"

	"github.com/olekukonko/tablewriter"
)

type ListHandler struct {
}

func disp(resp []models.Domain) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Name", "IP", "STATE", "REASON", "Location",
	})
	for _, d := range resp {
		table.Append([]string{
			d.Name, d.IP, d.Status, d.Reason, d.Location,
		})
	}
	table.Render()
}

func (handler ListHandler) Handle(args []string) {
	domains, err := services.GetMeshService().ListDomains()
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
