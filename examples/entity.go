package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/BurntSushi/toml"
	contentful "github.com/tolgaakyuz/contentful.go"
)

var cma *contentful.Contentful
var space *contentful.Space

func main() {
	usr, _ := user.Current()
	configFile := usr.HomeDir + "/.config/contentful.toml"
	var err error

	config := struct {
		CDAToken string `toml:"cda"`
		CPAToken string `toml:"cpa"`
		CMAToken string `toml:"cma"`
		SpaceID  string `toml:"spaceId"`
	}{}

	if _, err = toml.DecodeFile(configFile, &config); err != nil {
		fmt.Println(err)
		return
	}

	cma = contentful.NewCMA(config.CMAToken)
	space, err = cma.GetSpace(config.SpaceID)
	if err != nil {
		log.Fatal(err)
	}

	getEntries()
}

func getEntries() []*contentful.Entry {}
