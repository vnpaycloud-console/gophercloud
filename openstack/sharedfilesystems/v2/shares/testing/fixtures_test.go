package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	fake "github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

const (
	shareEndpoint = "/shares"
	shareID       = "011d21e2-fbc3-4e4a-9993-9ea223f73264"
)

var createRequest = `{
		"share": {
			"name": "my_test_share",
			"size": 1,
			"share_proto": "NFS",
                        "scheduler_hints": {
                            "same_host": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
                            "different_host": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef"
                        }
		}
	}`

var createResponse = `{
		"share": {
			"name": "my_test_share",
			"share_proto": "NFS",
			"size": 1,
			"status": null,
			"share_server_id": null,
			"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
			"share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
			"share_type_name": "default",
			"availability_zone": null,
			"created_at": "2015-09-18T10:25:24.533287",
			"export_location": null,
			"links": [
				{
					"href": "http://172.18.198.54:8786/v1/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
					"rel": "self"
				},
				{
					"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
					"rel": "bookmark"
				}
			],
			"share_network_id": null,
			"export_locations": [],
			"host": null,
			"access_rules_status": "active",
			"has_replicas": false,
			"replication_type": null,
			"task_state": null,
			"snapshot_support": true,
			"create_share_from_snapshot_support": true,
			"consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
			"source_cgsnapshot_member_id": null,
			"volume_type": "default",
			"snapshot_id": null,
			"is_public": true,
			"metadata": {
				"project": "my_app",
				"aim": "doc",
                                "__affinity_same_host": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
                                "__affinity_different_host": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef"
			},
			"id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
			"description": "My custom share London"
		}
	}`

// MockCreateResponse creates a mock response
func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, createRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, createResponse)
	})
}

// MockDeleteResponse creates a mock delete response
func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

var updateRequest = `{
		"share": {
			"display_name": "my_new_test_share",
			"display_description": "",
			"is_public": false
		}
	}`

var updateResponse = `
{
	"share": {
		"links": [
			{
				"href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
				"rel": "self"
			},
			{
				"href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
				"rel": "bookmark"
			}
		],
		"availability_zone": "nova",
		"share_network_id": "713df749-aac0-4a54-af52-10f6c991e80c",
		"export_locations": [],
		"share_server_id": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
		"share_group_id": null,
		"snapshot_id": null,
		"id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
		"size": 1,
		"share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
		"share_type_name": "default",
		"export_location": null,
		"project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
		"metadata": {
			"project": "my_app",
			"aim": "doc"
		},
		"status": "error",
		"description": "",
		"host": "manila2@generic1#GENERIC1",
		"task_state": null,
		"is_public": false,
		"snapshot_support": true,
		"create_share_from_snapshot_support": true,
		"name": "my_new_test_share",
		"created_at": "2015-09-18T10:25:24.000000",
		"share_proto": "NFS",
		"volume_type": "default"
	}
}
`

func MockUpdateResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, updateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, updateResponse)
	})
}

var getResponse = `{
    "share": {
        "links": [
            {
                "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
                "rel": "self"
            },
            {
                "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
                "rel": "bookmark"
            }
        ],
        "availability_zone": "nova",
        "share_network_id": "713df749-aac0-4a54-af52-10f6c991e80c",
        "share_server_id": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
        "snapshot_id": null,
        "id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
        "size": 1,
        "share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
        "share_type_name": "default",
        "consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
        "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
        "metadata": {
            "project": "my_app",
            "aim": "doc"
        },
        "status": "available",
        "description": "My custom share London",
        "host": "manila2@generic1#GENERIC1",
        "has_replicas": false,
        "replication_type": null,
        "task_state": null,
        "is_public": true,
        "snapshot_support": true,
        "create_share_from_snapshot_support": true,
        "name": "my_test_share",
        "created_at": "2015-09-18T10:25:24.000000",
        "share_proto": "NFS",
        "volume_type": "default",
        "source_cgsnapshot_member_id": null
    }
}`

// MockGetResponse creates a mock get response
func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, getResponse)
	})
}

var listDetailResponse = `{
		"shares": [
			{
		        "links": [
		            {
		                "href": "http://172.18.198.54:8786/v2/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
		                "rel": "self"
		            },
		            {
		                "href": "http://172.18.198.54:8786/16e1ab15c35a457e9c2b2aa189f544e1/shares/011d21e2-fbc3-4e4a-9993-9ea223f73264",
		                "rel": "bookmark"
		            }
		        ],
		        "availability_zone": "nova",
		        "share_network_id": "713df749-aac0-4a54-af52-10f6c991e80c",
		        "share_server_id": "e268f4aa-d571-43dd-9ab3-f49ad06ffaef",
		        "snapshot_id": null,
		        "id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
		        "size": 1,
		        "share_type": "25747776-08e5-494f-ab40-a64b9d20d8f7",
		        "share_type_name": "default",
		        "consistency_group_id": "9397c191-8427-4661-a2e8-b23820dc01d4",
		        "project_id": "16e1ab15c35a457e9c2b2aa189f544e1",
		        "metadata": {
		            "project": "my_app",
		            "aim": "doc"
		        },
		        "status": "available",
		        "description": "My custom share London",
		        "host": "manila2@generic1#GENERIC1",
		        "has_replicas": false,
		        "replication_type": null,
		        "task_state": null,
		        "is_public": true,
		        "snapshot_support": true,
		        "create_share_from_snapshot_support": true,
		        "name": "my_test_share",
		        "created_at": "2015-09-18T10:25:24.000000",
		        "share_proto": "NFS",
		        "volume_type": "default",
		        "source_cgsnapshot_member_id": null
		    }
		]
	}`

var listDetailEmptyResponse = `{"shares": []}`

// MockListDetailResponse creates a mock detailed-list response
func MockListDetailResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("offset")

		switch marker {
		case "":
			fmt.Fprint(w, listDetailResponse)
		default:
			fmt.Fprint(w, listDetailEmptyResponse)
		}
	})
}

var listExportLocationsResponse = `{
    "export_locations": [
        {
		"path": "127.0.0.1:/var/lib/manila/mnt/share-9a922036-ad26-4d27-b955-7a1e285fa74d",
        	"share_instance_id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
		"is_admin_only": false,
        	"id": "80ed63fc-83bc-4afc-b881-da4a345ac83d",
		"preferred": false
	}
    ]
}`

// MockListExportLocationsResponse creates a mock get export locations response
func MockListExportLocationsResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/export_locations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, listExportLocationsResponse)
	})
}

var getExportLocationResponse = `{
    "export_location": {
	"path": "127.0.0.1:/var/lib/manila/mnt/share-9a922036-ad26-4d27-b955-7a1e285fa74d",
	"share_instance_id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
	"is_admin_only": false,
	"id": "80ed63fc-83bc-4afc-b881-da4a345ac83d",
	"preferred": false
    }
}`

// MockGetExportLocationResponse creates a mock get export location response
func MockGetExportLocationResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/export_locations/80ed63fc-83bc-4afc-b881-da4a345ac83d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, getExportLocationResponse)
	})
}

var grantAccessRequest = `{
		"allow_access": {
			"access_type": "ip",
			"access_to": "0.0.0.0/0",
			"access_level": "rw"
		}
	}`

var grantAccessResponse = `{
    "access": {
	"share_id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
	"access_type": "ip",
	"access_to": "0.0.0.0/0",
	"access_key": "",
	"access_level": "rw",
	"state": "new",
	"id": "a2f226a5-cee8-430b-8a03-78a59bd84ee8"
    }
}`

// MockGrantAccessResponse creates a mock grant access response
func MockGrantAccessResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, grantAccessRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, grantAccessResponse)
	})
}

var revokeAccessRequest = `{
	"deny_access": {
		"access_id": "a2f226a5-cee8-430b-8a03-78a59bd84ee8"
	}
}`

// MockRevokeAccessResponse creates a mock revoke access response
func MockRevokeAccessResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, revokeAccessRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var listAccessRightsRequest = `{
		"access_list": null
	}`

var listAccessRightsResponse = `{
		"access_list": [
			{
				"share_id": "011d21e2-fbc3-4e4a-9993-9ea223f73264",
				"access_type": "ip",
				"access_to": "0.0.0.0/0",
				"access_key": "",
				"access_level": "rw",
				"state": "new",
				"id": "a2f226a5-cee8-430b-8a03-78a59bd84ee8"
			}
		]
	}`

// MockListAccessRightsResponse creates a mock list access response
func MockListAccessRightsResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, listAccessRightsRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, listAccessRightsResponse)
	})
}

var extendRequest = `{
		"extend": {
			"new_size": 2
		}
	}`

// MockExtendResponse creates a mock extend share response
func MockExtendResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, extendRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var shrinkRequest = `{
		"shrink": {
			"new_size": 1
		}
	}`

// MockShrinkResponse creates a mock shrink share response
func MockShrinkResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, shrinkRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var getMetadataResponse = `{
		"metadata": {
			"foo": "bar"
		}
	}`

// MockGetMetadataResponse creates a mock get metadata response
func MockGetMetadataResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/metadata", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, getMetadataResponse)
	})
}

var getMetadatumResponse = `{
		"meta": {
			"foo": "bar"
		}
	}`

// MockGetMetadatumResponse creates a mock get metadatum response
func MockGetMetadatumResponse(t *testing.T, key string) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/metadata/"+key, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, getMetadatumResponse)
	})
}

var setMetadataRequest = `{
		"metadata": {
			"foo": "bar"
		}
	}`

var setMetadataResponse = `{
		"metadata": {
			"foo": "bar"
		}
	}`

// MockSetMetadataResponse creates a mock set metadata response
func MockSetMetadataResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/metadata", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, setMetadataRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, setMetadataResponse)
	})
}

var updateMetadataRequest = `{
		"metadata": {
			"foo": "bar"
		}
	}`

var updateMetadataResponse = `{
		"metadata": {
			"foo": "bar"
		}
	}`

// MockUpdateMetadataResponse creates a mock update metadata response
func MockUpdateMetadataResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/metadata", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, updateMetadataRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, updateMetadataResponse)
	})
}

// MockDeleteMetadatumResponse creates a mock unset metadata response
func MockDeleteMetadatumResponse(t *testing.T, key string) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/metadata/"+key, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusOK)
	})
}

var revertRequest = `{
		"revert": {
			"snapshot_id": "ddeac769-9742-497f-b985-5bcfa94a3fd6"
		}
	}`

// MockRevertResponse creates a mock revert share response
func MockRevertResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, revertRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var resetStatusRequest = `{
		"reset_status": {
			"status": "error"
		}
	}`

// MockResetStatusResponse creates a mock reset status share response
func MockResetStatusResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, resetStatusRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var forceDeleteRequest = `{
                "force_delete": null
        }`

// MockForceDeleteResponse creates a mock force delete share response
func MockForceDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, forceDeleteRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}

var unmanageRequest = `{
                "unmanage": null
        }`

// MockUnmanageResponse creates a mock unmanage share response
func MockUnmanageResponse(t *testing.T) {
	th.Mux.HandleFunc(shareEndpoint+"/"+shareID+"/action", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, unmanageRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
	})
}
