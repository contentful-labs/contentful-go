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

	space, err = cma.GetSpace(config.SpaceID)
	if err != nil {
		log.Fatal(err)
	}

	getContentTypes()
	getContentType()
	createContentType()
	activateContentType()
	deactivateContentType()
	deleteContentType()
	deleteAllDraftContentTypes()
}

func getContentTypes() []*contentful.ContentType {
	collection, err := space.GetContentTypes().Next()
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
	contentType, err := space.GetContentType("contentTypeTest1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(contentType.Sys.ID, contentType.Name)

	return contentType
}

func createContentType() (*contentful.ContentType, error) {
	contentType := space.NewContentType()
	contentType.Name = "test content type"
	contentType.DisplayField = "field1_id"
	contentType.Description = "content type description"
	contentType.Fields = []*contentful.Field{
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
	}

	err := contentType.Save()
	if err != nil {
		return nil, err
	}

	fmt.Println(contentType.Sys.ID, contentType.Name)

	return contentType, nil
}

func activateContentType() error {
	contentType, _ := createContentType()

	err := contentType.Activate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is activated :", contentType.Name)

	return nil
}

func deactivateContentType() error {
	contentType, _ := createContentType()

	err := contentType.Activate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is activated :", contentType.Name)

	err = contentType.Deactivate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is deactivated :", contentType.Name)

	return nil
}

func deleteContentType() {
	var err error
	contentType, _ := createContentType()

	err = contentType.Activate()
	if err != nil {
		log.Fatal(err)
	}

	err = contentType.Delete()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("content type is deleted")
}

func deleteAllDraftContentTypes() {
	collection, err := space.GetContentTypes().Next()
	if err != nil {
		log.Fatal(err)
	}

	contentTypes := collection.ToContentType()

	for _, contentType := range contentTypes {
		if contentType.Sys.PublishedAt == "" {
			err := contentType.Delete()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fmt.Println("all draft content types are deleted")
}
