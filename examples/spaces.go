package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/BurntSushi/toml"
	contentful "github.com/tolgaakyuz/contentful-go"
)

type Config struct {
	CDAToken string `toml:"cda"`
	CPAToken string `toml:"cpa"`
	CMAToken string `toml:"cma"`
	SpaceID  string `toml:"spaceId"`
	BaseURL  string `toml:"baseURL"`
}

var cma *contentful.Contentful
var space *contentful.Space
var config Config

func main() {
	usr, _ := user.Current()
	configFile := usr.HomeDir + "/.config/contentful.toml"
	var err error

	if _, err = toml.DecodeFile(configFile, &config); err != nil {
		fmt.Println(err)
		return
	}

	cma = contentful.NewCMA(config.CMAToken)

	getSpaces()
	// deleteAllSpaces()
}

func getSpaces() []*contentful.Space {
	collection, err := cma.Spaces.List().Next()
	if err != nil {
		log.Fatal(err)
	}

	spaces := collection.ToSpace()
	for _, space := range spaces {
		fmt.Println(space.Sys.ID, space.Name)
	}

	return spaces
}

func getSpace() *contentful.Space {
	spaceID := getSpaces()[0].Sys.ID

	space, err := cma.Spaces.Get(spaceID)
	if err != nil {
		log.Fatal(err)
	}

	return space
}

func editSpace() {
	space := getSpace()
	space.Name = "modified"
	err := cma.Spaces.Upsert(space)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteSpace() {
	space := getSpace()
	err := cma.Spaces.Delete(space)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteAllSpaces() {
	collection, err := cma.Spaces.List().Next()
	if err != nil {
		log.Fatal(err)
	}

	for _, space := range collection.ToSpace() {
		err := cma.Spaces.Delete(space)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("all spaces are deleted")
}
