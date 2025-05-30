package testing

import (
	"net/http"
	"testing"

	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	fake "github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

// HandleGetAccountSuccessfully creates an HTTP handler at `/` on the test handler mux that
// responds with a `Get` response.
func HandleGetAccountSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("X-Account-Container-Count", "2")
		w.Header().Set("X-Account-Object-Count", "5")
		w.Header().Set("X-Account-Meta-Quota-Bytes", "42")
		w.Header().Set("X-Account-Bytes-Used", "14")
		w.Header().Set("X-Account-Meta-Subject", "books")
		w.Header().Set("Date", "Fri, 17 Jan 2014 16:09:56 UTC")
		w.Header().Set("X-Account-Meta-Temp-URL-Key", "testsecret")

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleGetAccountNoQuotaSuccessfully creates an HTTP handler at `/` on the
// test handler mux that responds with a `Get` response.
func HandleGetAccountNoQuotaSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "HEAD")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("X-Account-Container-Count", "2")
		w.Header().Set("X-Account-Object-Count", "5")
		w.Header().Set("X-Account-Bytes-Used", "14")
		w.Header().Set("X-Account-Meta-Subject", "books")
		w.Header().Set("Date", "Fri, 17 Jan 2014 16:09:56 UTC")

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateAccountSuccessfully creates an HTTP handler at `/` on the test handler mux that
// responds with a `Update` response.
func HandleUpdateAccountSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "X-Account-Meta-Gophercloud-Test", "accounts")
		th.TestHeader(t, r, "X-Remove-Account-Meta-Gophercloud-Test-Remove", "remove")
		th.TestHeader(t, r, "Content-Type", "")
		th.TestHeader(t, r, "X-Detect-Content-Type", "false")
		th.TestHeaderUnset(t, r, "X-Account-Meta-Temp-URL-Key")

		w.Header().Set("Date", "Fri, 17 Jan 2014 16:09:56 UTC")
		w.WriteHeader(http.StatusNoContent)
	})
}
