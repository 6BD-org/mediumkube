package handlers

import (
	"fmt"
	"mediumkube/pkg/models"
	"mediumkube/pkg/services"
	"mediumkube/pkg/utils"
	"os"

	"github.com/olekukonko/tablewriter"
)

type ListNodeHandler struct {
}

func (h ListNodeHandler) Desc() string {
	return "List nodes in cluster"
}

func (h ListNodeHandler) Help() {
	fmt.Println(h.Desc())
}

func dispLeases(peers []models.PeerLease) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Cidr"})
	for _, peer := range peers {
		table.Append([]string{peer.Id, peer.Cidr})
	}
	table.Render()
}

func (h ListNodeHandler) Handle(args []string) {
	leases, err := services.GetMeshService().ListLeases()
	utils.CheckErr(err)
	dispLeases(leases)

}

func init() {
	name := "leases"
	CMD[name] = ListNodeHandler{}
}
