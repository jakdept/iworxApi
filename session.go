package iworx

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/kolo/xmlrpc"
	"golang.org/x/crypto/ssh"
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

func (a *NodeWorxAPI) NodeWorxSessionAuthenticate(session string, domain string) {
	// $key = array( 'sessionid' => '3c8ae9d982edd507428d8fdd53855a77' );
	// $input = array();
	// $params = array( 'apikey'    => $key,
	//                  'ctrl_name' => $api_controller,
	//                  'action'    => $action,
	//                  'input'     => $input );
	// // You can connect using XMLRPC, like this:
	// // NOTE: This example makes use of the Zend Framework's XMLRPC library.
	// $client = new Zend_XmlRpc_Client( 'https://license-api.interworx.com:2443/xmlrpc' );
	// $result = $client->call( 'iworx.route', $params );

	a.auth = map[string]string{
		"sessionid": session,
	}
	if domain != "" {
		a.auth["domain"] = domain
	}
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

func (a *NodeWorxAPI) AuthViaInsecureSSHKeyfile(
	hostname, username, keyFile string, port int) error {
	creds, err := SSHKeyfileInsecureRemote(username, keyFile)
	if err != nil {
		return err
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), &creds)
	if err != nil {
		return err
	}
	session, err := conn.NewSession()
	if err != nil {
		conn.Close()
		return err
	}

	// nodeworx -nu --controller Index --action getSession
	cmd := "nodeworx -nu --controller Index --action getSession"

	output, err := session.Output(cmd)
	if err != nil {
		return err
	}

	a.NodeWorxSessionAuthenticate(strings.TrimSpace(string(output)), "")
	return nil
}

func SSHKeyfileInsecureRemote(username, keyFile string) (ssh.ClientConfig, error) {
	// read the keyfile
	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return ssh.ClientConfig{}, err
	}

	return ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint
	}, nil
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
