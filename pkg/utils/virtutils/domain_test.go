package virtutils

import (
	"log"
	"mediumkube/pkg/common"
	"testing"
)

func TestRenderDomain(t *testing.T) {
	param := common.NewDomainCreationParam(
		"test-domain",
		"2", "2G", "./a.img", "./c.img", "br0",
	)

	domainXML, err := GetDeploymentConfig(param)
	if err != nil {
		t.Fail()
	}

	expectation := `
<domain type='kvm'>
	<name>test-domain</name>
	<title>test-domain - Generated by MediumKube</title>
	
	<!--Resource Allocation-->
	<vcpu>2</vcpu>
	<maxMemory unit='G'>2</maxMemory>
	<memory unit='G'>2</memory>
	<currentMemory unit='G'>1</currentMemory>

	<devices>
		<!-- System disk -->
		<disk type='file' device='disk'> 
			<source file='./a.img' />
		</disk>

		<!-- Cloud init disk -->
		<disk type='file' device='cdrom'>
			<source file='./c.img' />
		</disk>

		<interface type='bridge'>
			<source bridge='br0'>
		</interface>
	</devices>

</domain>
`

	log.Println(domainXML)
	log.Println(expectation)
}
