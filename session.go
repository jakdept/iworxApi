package iworx

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/tiaguinho/gosoap"
	"golang.org/x/crypto/ssh"
)

type NodeWorxAPI struct {
	hostname string
	session  string
	client   *gosoap.Client
	setup    bool
}

// TODO figure out how to standardize ssh.Session and exec.Command

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
	a.session = strings.TrimSpace(string(output))

	a.NodeWorxSessionAuthenticate(a.session)
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

func NewNodeWorxAPI(hostname string) (*NodeWorxAPI, error) {
	newClient, err := gosoap.SoapClient("https://" + hostname + ":2443/soap?wsdl")
	if err != nil {
		return nil, nil
	}
	return &NodeWorxAPI{
		hostname: hostname,
		client:   newClient,
	}, nil
}

func (a *NodeWorxAPI) NodeWorxSessionAuthenticate(session string) error {
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

	params := gosoap.Params{
		"apikey": gosoap.Params{
			"sessionid": session,
		},
		"ctrl_name": "Index",
		"action":    "ssoCommit",
		"input": map[string]string{
			"sid": session,
		},
	}

	resp, err := a.client.Call("Index", params)
	if err != nil {
		return err
	}
	spew.Dump(resp.Header)
	spew.Dump(resp.Body)
	spew.Dump(resp.Payload)

	return nil
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
