package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/image/v2/imageimport"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	fakeclient "github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/info/import", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ImportGetResult)
	})

	validImportMethods := []string{
		string(imageimport.GlanceDirectMethod),
		string(imageimport.WebDownloadMethod),
	}

	s, err := imageimport.Get(context.TODO(), fakeclient.ServiceClient()).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.ImportMethods.Description, "Import methods available.")
	th.AssertEquals(t, s.ImportMethods.Type, "array")
	th.AssertDeepEquals(t, s.ImportMethods.Value, validImportMethods)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/images/da3b75d9-3f4a-40e7-8a2c-bfab23927dea/import", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fakeclient.TokenID)
		th.TestJSONRequest(t, r, ImportCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{}`)
	})

	opts := imageimport.CreateOpts{
		Name: imageimport.WebDownloadMethod,
		URI:  "http://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img",
	}
	err := imageimport.Create(context.TODO(), fakeclient.ServiceClient(), "da3b75d9-3f4a-40e7-8a2c-bfab23927dea", opts).ExtractErr()
	th.AssertNoErr(t, err)
}
