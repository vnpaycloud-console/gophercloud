package testing

import "github.com/vnpaycloud-console/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"

const ListBGPPeersResult = `
{
  "bgp_peers": [
    {
      "auth_type": "none",
      "remote_as": 4321,
      "name": "testing-peer-1",
      "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "peer_ip": "1.2.3.4",
      "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "id": "afacc0e8-6b66-44e4-be53-a1ef16033ceb"
    },
    {
      "auth_type": "none",
      "remote_as": 4321,
      "name": "testing-peer-2",
      "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "peer_ip": "5.6.7.8",
      "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
      "id": "acd7c4a1-e243-4fe5-80f9-eba8f143ac1d"
    }
  ]
}
`

var BGPPeer1 = peers.BGPPeer{
	ID:        "afacc0e8-6b66-44e4-be53-a1ef16033ceb",
	AuthType:  "none",
	Name:      "testing-peer-1",
	TenantID:  "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	PeerIP:    "1.2.3.4",
	ProjectID: "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	RemoteAS:  4321,
}

var BGPPeer2 = peers.BGPPeer{
	AuthType:  "none",
	ID:        "acd7c4a1-e243-4fe5-80f9-eba8f143ac1d",
	Name:      "testing-peer-2",
	TenantID:  "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	PeerIP:    "5.6.7.8",
	ProjectID: "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
	RemoteAS:  4321,
}

const GetBGPPeerResult = `
{
  "bgp_peer": {
    "auth_type": "none",
    "remote_as": 4321,
    "name": "testing-peer-1",
    "tenant_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "peer_ip": "1.2.3.4",
    "project_id": "7fa3f96b-17ee-4d1b-8fbf-fe889bb1f1d0",
    "id": "afacc0e8-6b66-44e4-be53-a1ef16033ceb"
  }
}
`

const CreateRequest = `
{
  "bgp_peer": {
    "auth_type": "md5",
    "name": "gophercloud-testing-bgp-peer",
    "password": "notSoStrong",
    "peer_ip": "192.168.0.1",
    "remote_as": 20000
  }
}
`

const CreateResponse = `
{
  "bgp_peer": {
    "auth_type": "md5",
    "project_id": "52a9d4ff-81b6-4b16-a7fa-5325d3bc1c5d",
    "remote_as": 20000,
    "name": "gophercloud-testing-bgp-peer",
    "tenant_id": "52a9d4ff-81b6-4b16-a7fa-5325d3bc1c5d",
    "peer_ip": "192.168.0.1",
    "id": "b7ad63ea-b803-496a-ad59-f9ef513a5cb9"
  }
}
`

const UpdateBGPPeerRequest = `
{
  "bgp_peer": {
    "name": "test-rename-bgp-peer",
    "password": "superStrong"
  }
}
`

const UpdateBGPPeerResponse = `
{
  "bgp_peer": {
    "auth_type": "md5",
    "remote_as": 20000,
    "name": "test-rename-bgp-peer",
    "tenant_id": "52a9d4ff-81b6-4b16-a7fa-5325d3bc1c5d",
    "peer_ip": "192.168.0.1",
    "project_id": "52a9d4ff-81b6-4b16-a7fa-5325d3bc1c5d",
    "id": "b7ad63ea-b803-496a-ad59-f9ef513a5cb9"
  }
}
`
