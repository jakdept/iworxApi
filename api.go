package iworx

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/kolo/xmlrpc"
)

type NodeWorxAPI struct {
	client *xmlrpc.Client
	auth   map[string]string `xmlrpc:"apikey"`
}

// auth object may be:
//     map[string]string{"sessionid": "sessid"}
//     map[string]string{"apikey":"key"}
//     map[string]string{"email":"username@domain.com", "password":"hunter2"}
// For SiteWorx for all three options add another "domain" key

const NodeWorxAPIRoute = "iworx.route"

func NewNodeWorxAPI(hostname string) (*NodeWorxAPI, error) {
	client, err := xmlrpc.NewClient(fmt.Sprintf("https://%s:2443/xmlrpc", hostname),
		&http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, //nolint
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return &NodeWorxAPI{client: client}, nil
}

func (a *NodeWorxAPI) Call(
	controller string,
	action string,
	input interface{},
	output interface{},
) error {

	if len(a.auth) == 0 {
		return errors.New("API not authenticated")
	}

	err := a.client.Call(NodeWorxAPIRoute, []interface{}{
		a.auth, controller, action, input,
	}, output)
	return err
}
