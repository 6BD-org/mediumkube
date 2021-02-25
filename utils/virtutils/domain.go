package virtutils

import (
	"bytes"
	"fmt"
	"mediumkube/common"
	"text/template"

	"github.com/Masterminds/sprig"
)

const domainTmplKvm = `
<domain type='{{.Type}}'>
	<name>{{.Name}}</name>
	<title>{{.Name}} - Generated by MediumKube</title>
	<os>
    	<type arch='x86_64'>hvm</type>
    	<boot dev='hd'/>
  	</os>

	<!--Resource Allocation-->
	<vcpu placement='static'>{{.CPU}}</vcpu>
	<memory unit='{{.MemoryUnit}}'>{{.MemoryMagnitude}}</memory>

	<features>
    	<acpi/>
    	<apic/>
  	</features>


	<on_poweroff>destroy</on_poweroff>
	<on_reboot>restart</on_reboot>
	<on_crash>destroy</on_crash>
  
	<devices>
		<!-- System disk -->
		<disk type='file' device='disk'> 
			<driver name='qemu' type='qcow2'/>
			<source file='{{.OSImage}}' />
			<target dev='hda' bus='ide' />
			<address type='drive' controller='0' bus='0' target='0' unit='0'/>
		</disk>

		<!-- Cloud init disk -->
		<disk type='file' device='cdrom'>
			<driver name='qemu' type='raw'/>
			<source file='{{.CloudInitImage}}' />
			<target dev='hdb' bus='ide'/>
			<address type='drive' controller='0' bus='0' target='0' unit='1'/>
		</disk>

		<controller type='usb' index='0' model='ich9-ehci1'>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x7'/>
		</controller>
		<controller type='usb' index='0' model='ich9-uhci1'>
			<master startport='0'/>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0' multifunction='on'/>
		</controller>
		<controller type='usb' index='0' model='ich9-uhci2'>
			<master startport='2'/>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x1'/>
		</controller>
		<controller type='usb' index='0' model='ich9-uhci3'>
			<master startport='4'/>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x2'/>
		</controller>
		<controller type='pci' index='0' model='pci-root'/>
		<controller type='ide' index='0'>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x01' function='0x1'/>
		</controller>
		<controller type='virtio-serial' index='0'>
			<address type='pci' domain='0x0000' bus='0x00' slot='0x05' function='0x0'/>
		</controller>
  

		<interface type='bridge'>
			<source bridge='{{.Bridge}}'/>
		</interface>
		<serial type='pty'>
			<target type='isa-serial' port='0'>
			<model name='isa-serial'/>
			</target>
		</serial>
		<console type='pty'>
			<target type='serial' port='0'/>
		</console>
	
	</devices>

</domain>
`

// GetDeploymentConfig Render a domain deployment config
func GetDeploymentConfig(param common.DomainCreationParam) (string, error) {
	var tmplStr string
	if param.Type == "kvm" {
		tmplStr = domainTmplKvm
	} else {
		return "", fmt.Errorf("VM type not supported")
	}

	tmpl, err := template.New("domain").Funcs(sprig.TxtFuncMap()).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("Unable to generate template")
	}
	buffer := bytes.NewBuffer(make([]byte, 0))
	err = tmpl.Execute(buffer, param)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
