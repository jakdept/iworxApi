package iworx

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (a *NodeWorxAPI) NodeWorxSessionAuthenticate(session, domain string) {
	// set up the required object for session based authentication
	a.auth = map[string]string{
		"sessionid": session,
	}
	if domain != "" {
		a.auth["domain"] = domain
	}
}

func (a *NodeWorxAPI) APIKeyAuthenticate(key, domain string) {
	// set up the required object for apikey based authentication
	a.auth = map[string]string{
		"apikey": strings.TrimSpace(key),
	}
	if domain != "" {
		a.auth["domain"] = domain
	}
}

func (a *NodeWorxAPI) UserAuthenticate(username, password, domain string) {
	// set up the required object for user based authentication
	a.auth = map[string]string{
		"email":    username,
		"password": password,
	}
	if domain != "" {
		a.auth["domain"] = domain
	}
}

func (a *NodeWorxAPI) SSHSessionAuthenticate(
	hostname string,
	port int,
	config ssh.ClientConfig,
) error {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, port), &config)
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

func (a *NodeWorxAPI) LocalSessionAuthenticate() error {
	// nodeworx -nu --controller Index --action getSession
	cmd := exec.Command(
		"/usr/bin/nodeworx",
		"-nu",
		"--controller",
		"Index",
		"--action",
		"getSession",
	)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	a.NodeWorxSessionAuthenticate(strings.TrimSpace(string(output)), "")
	return nil
}

func InsecureSSHKeyfileConfig(username, keyFile string) (ssh.ClientConfig, error) {
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
