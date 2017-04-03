package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/BurntSushi/toml"
	contentful "github.com/contentful-labs/contentful-go"
)

type Config struct {
	CDAToken string `toml:"cda"`
	CPAToken string `toml:"cpa"`
	CMAToken string `toml:"cma"`
	SpaceID  string `toml:"spaceId"`
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

	getContentTypes()
	getContentType()
	createContentType()
	activateContentType()
	deactivateContentType()
	deleteContentType()
	deleteAllDraftContentTypes()
}

func getContentTypes() []*contentful.ContentType {
	collection, err := cma.ContentTypes.List(config.SpaceID).Next()
	if err != nil {
		log.Fatal(err)
	}

	contentTypes := collection.ToContentType()
	for _, contentType := range contentTypes {
		fmt.Println(contentType.Sys.ID, contentType.Sys.PublishedAt)
	}

	return contentTypes
}

func getContentType() *contentful.ContentType {
	contentType, err := cma.ContentTypes.Get(config.SpaceID, "contentTypeTest1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(contentType.Sys.ID, contentType.Name)

	return contentType
}

func createContentType() (*contentful.ContentType, error) {
	contentType := &contentful.ContentType{
		Name:         "test content type",
		DisplayField: "field1_id",
		Description:  "content type description",
		Fields: []*contentful.Field{
			&contentful.Field{
				ID:       "field1_id",
				Name:     "field1",
				Type:     "Symbol",
				Required: false,
				Disabled: false,
			},
			&contentful.Field{
				ID:       "field2_id",
				Name:     "field2",
				Type:     "Symbol",
				Required: false,
				Disabled: true,
			},
		},
	}

	err := cma.ContentTypes.Upsert(config.SpaceID, contentType)
	if err != nil {
		return nil, err
	}

	fmt.Println(contentType.Sys.ID, contentType.Name)

	return contentType, nil
}

func activateContentType() *contentful.ContentType {
	contentType, _ := createContentType()

	err := cma.ContentTypes.Activate(config.SpaceID, contentType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is activated :", contentType.Name)

	return contentType
}

func deactivateContentType() *contentful.ContentType {
	contentType := activateContentType()

	err := cma.ContentTypes.Deactivate(config.SpaceID, contentType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is deactivated :", contentType.Name)

	return contentType
}

func deleteContentType() {
	contentType := deactivateContentType()

	err := cma.ContentTypes.Delete(config.SpaceID, contentType)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is deleted")
}

func deleteAllDraftContentTypes() {
	contentTypes := getContentTypes()

	for _, contentType := range contentTypes {
		if contentType.Sys.PublishedAt == "" {
			err := cma.ContentTypes.Delete(config.SpaceID, contentType)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("all draft content types are deleted")
}
