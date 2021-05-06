package common

import "mediumkube/utils"

// DomainCreationParam contains all necessary information
// to create a domain
type DomainCreationParam struct {
	Name           string
	CPU            string
	Memory         string
	OSImage        string
	CloudInitImage string
	Bridge         string
	Type           string // Need to be determined at rumtime

}

// NewDomainCreationParam Nothing interesting
func NewDomainCreationParam(name string, cpu string, memory string, osImage string, cloudInitImage string, bridge string) DomainCreationParam {
	return DomainCreationParam{
		Name:           name,
		CPU:            cpu,
		Memory:         memory,
		OSImage:        osImage,
		CloudInitImage: cloudInitImage,
		Bridge:         bridge,
		Type:           "kvm",
	}
}

// MemoryUnit as string
func (param DomainCreationParam) MemoryUnit() string {
	_, unit, err := utils.GetMagnitudeAndUnitStr(param.Memory)
	utils.CheckErr(err)
	return unit
}

// MemoryMagnitude returns memory magnitude. Must be used with Unit
func (param DomainCreationParam) MemoryMagnitude() float64 {
	mag, _, err := utils.GetMagnitudeAndUnitStr(param.Memory)
	utils.CheckErr(err)
	return mag
}

// CurrentMemory returns initial memory allocated to machine when it's
// created. This value is half of machine's max memory by default
// Refer to https://libvirt.org/formatdomain.html for more details
func (param DomainCreationParam) CurrentMemory() float64 {
	return param.MemoryMagnitude() / 2
}
