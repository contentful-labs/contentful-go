package main

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/BurntSushi/toml"
	contentful "github.com/tolgaakyuz/contentful.go"
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

	getAssets()
	getAsset()
	createAsset()
	createAndEditAsset()
	editExistingAsset()
	deleteExistingAsset()
	createAndProcessAsset()
	publishAsset()
}

func getAssets() []*contentful.Asset {
	collection, err := space.GetAssets().Next()
	if err != nil {
		log.Fatal(err)
	}

	assets := collection.ToAsset()

	/* for _, asset := range assets {
		fmt.Println(asset.ID(), asset.Fields.Title)
	} */

	return assets
}

func getAsset() *contentful.Asset {
	assetID := getAssets()[0].ID()

	asset, err := space.GetAsset(assetID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(asset.Sys.Type)

	return asset
}

func createAsset() *contentful.Asset {
	asset := space.NewAsset()
	asset.Fields.Title = "asset title"
	asset.Fields.File = &contentful.File{
		Name:        "file name",
		ContentType: "image/jpg",
		UploadURL:   "http://desireeacres.com/gopher.jpg",
	}

	err := asset.Save()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset is created with id: " + asset.Sys.ID)

	return asset
}

func createAndEditAsset() {
	asset := space.NewAsset()
	asset.Fields.Title = "asset title"
	asset.Fields.File = &contentful.File{
		Name:        "file name",
		ContentType: "image/jpg",
		UploadURL:   "http://desireeacres.com/gopher.jpg",
	}

	err := asset.Save()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset is created with id: " + asset.Sys.ID)
	fmt.Println(asset.Fields.Title)

	asset.Fields.Title = "updated title for asset"
	err = asset.Save()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset updated with id: " + asset.Sys.ID)
	fmt.Println(asset.Fields.Title)
}

func editExistingAsset() {
	asset := getAssets()[0]
	asset.Fields.Title = "existing asset is updated"

	fmt.Println("updating asset with id: " + asset.ID())

	err := asset.Save()
	if err != nil {
		log.Fatal(err)
	}

	updatedAsset, err := space.GetAsset(asset.ID())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new asset with id: " + updatedAsset.ID())
	fmt.Println(updatedAsset.Fields.Title)
}

func deleteExistingAsset() {
	asset := getAssets()[0]
	fmt.Println("deleting asset with id: " + asset.ID())

	err := asset.Delete()
	if err != nil {
		log.Fatal(err)
	}
}

func createAndProcessAsset() *contentful.Asset {
	asset := space.NewAsset()
	asset.Fields.Title = "asset title"
	asset.Fields.File = &contentful.File{
		Name:        "file name",
		ContentType: "image/jpg",
		UploadURL:   "http://desireeacres.com/gopher.jpg",
	}

	err := asset.Save()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset is created with id: " + asset.Sys.ID)
	fmt.Println(asset.Fields.Title)

	err = asset.Process()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("processing..." + asset.Sys.ID)

	return asset
}

func publishAsset() {
	asset := createAndProcessAsset()
	time.Sleep(3 * time.Second)

	err := asset.Publish()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("published with id : " + asset.ID())
	fmt.Println("publish date: " + asset.Sys.PublishedAt)
}
