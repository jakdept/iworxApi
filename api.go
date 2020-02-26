package iworx

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/kolo/xmlrpc"
)

type NodeWorxAPI struct {
	defaultReqParams NodeWorxReqParams
	client           *xmlrpc.Client
	auth             map[string]string `xmlrpc:"apikey"`
}

type NodeWorxReqParams struct {
	Auth       map[string]string `xmlrpc:"apikey"`
	Controller string            `xmlrpc:"ctrl_name"`
	Action     string            `xmlrpc:"action"`
	Input      interface{}       `xmlrpc:"input"`
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

func (a *NodeWorxAPI) NodeWorxVersion() (string, error) {
	output := struct {
		Status      int `xmlrpc:"status"`
		RespPayload struct {
			Version string `xmlrpc:"version"`
		} `xmlrpc:"payload"`
	}{}

	err := a.Call("/nodeworx/overview", "listVersion", map[string]string{}, &output)
	return output.RespPayload.Version, err
}

// Account represents a SiteWorx account
type Account struct {
	Username     string
	Domain       string
	ContactEmail string

	Reseller      string
	HomePartition string
	Shell         string
	Package       string
	Theme         string

	BackupsEnabled        bool
	Suspended             bool
	Locked                bool
	OutgoingMailSuspended bool
	OutgoingMailHold      bool

	MailboxFormat            int8
	MaxDeferPrecent          string
	MinDeferBeforeProtection string
	MaxEmailPerHour          string

	MainIPv4 net.IP
	MainIPv6 net.IP

	EmailQuotaLimit string
	MaxAddons       string
	MaxFtp          string
	MaxMailingLists string
	MaxParked       string
	MaxPop          string
	MaxDatabases    string
	MaxSubdomains   string

	DiskLimit  int
	DiskUsed   int
	InodeLimit int
	InodeUsed  int
}

// func (a *NodeWorxAPI) ListAccounts() ([]string, error) {
// 	url := a.GenerateURL("listaccts")
// 	params := url.Query()
// 	params.Add("want", "user")
// 	url.RawQuery = params.Encode()

// 	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
// 	if err != nil {
// 		return []string{}, err
// 	}
// 	response, err := a.client.Do(request)
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	var outputData struct {
// 		Data struct {
// 			AccountList []struct {
// 				Username string `json:"user"`
// 			} `json:"acct"`
// 		} `json:"data"`
// 	}

// 	err = json.NewDecoder(response.Body).Decode(&outputData)
// 	if err != nil {
// 		return []string{}, err
// 	}

// 	userlist := []string{}
// 	for _, account := range outputData.Data.AccountList {
// 		userlist = append(userlist, account.Username)
// 	}
// 	return userlist, nil
// }
