package models

type Domain struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	IP       string `json:"ip"`
	Reason   string `json:"reason"`
	Location string `json:"location"`
}
