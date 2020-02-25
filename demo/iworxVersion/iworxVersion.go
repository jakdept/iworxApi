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

	log.Println("authenticating")
	err = api.AuthViaInsecureSSHKeyfile(*host, *username, *keyfile, *port)
	if err != nil {
		log.Fatalf("problem authenticating: %s\n", err)
	}

	version, err := api.NodeWorxVersion()
	if err != nil {
		log.Fatalf("problem getting iworx version: %v", err)
	}
	fmt.Printf("nodeworx is version %s\n", version)

	log.Println("i'm here")

	// accounts, err := api.ListAccounts()
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// fmt.Println("All accounts on server:")
	// for _, each := range accounts {
	// 	fmt.Println(each)
	// }
}
