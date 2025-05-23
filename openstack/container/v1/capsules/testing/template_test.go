package testing

import (
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/container/v1/capsules"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestTemplateParsing(t *testing.T) {
	templateJSON := new(capsules.Template)
	templateJSON.Bin = []byte(ValidJSONTemplate)
	err := templateJSON.Parse()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ValidJSONTemplateParsed, templateJSON.Parsed)

	templateYAML := new(capsules.Template)
	templateYAML.Bin = []byte(ValidYAMLTemplate)
	err = templateYAML.Parse()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ValidYAMLTemplateParsed, templateYAML.Parsed)

	templateInvalid := new(capsules.Template)
	templateInvalid.Bin = []byte("Keep Austin Weird")
	err = templateInvalid.Parse()
	if err == nil {
		t.Error("Template parsing did not catch invalid template")
	}
}
