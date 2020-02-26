package iworx

import (
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

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
