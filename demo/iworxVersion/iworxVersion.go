package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kingpin"
	iworx "github.com/jakdept/iworxApi"
)

var keyfile *string = kingpin.Flag("keyfile", "location to ssh key").Default("/root/.ssh/id_rsa").String()
var username *string = kingpin.Flag("username", "remote ssh user").Default("root").String()
var host *string = kingpin.Flag("host", "remote ssh host").Default("localhost").String()
var port *int = kingpin.Flag("port", "remote ssh port").Default("22").Int()

func main() {
	_ = kingpin.Parse()

	api, err := iworx.NewNodeWorxAPI(*host)
	if err != nil {
		log.Fatalf("problem creating API: %s\n", err)
	}

	sshConfig, err := iworx.InsecureSSHKeyfileConfig(*username, *keyfile)
	if err != nil {
		log.Fatalf("problem with ssh client config: %s\n", err)
	}

	err = api.SSHSessionAuthenticate(*host, *port, sshConfig)
	if err != nil {
		log.Fatalf("problem authenticating: %s\n", err)
	}

	version, err := api.NodeWorxVersion()
	if err != nil {
		log.Fatalf("problem getting iworx version: %v", err)
	}
	fmt.Printf("nodeworx is version %s\n", version)

}
