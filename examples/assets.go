package main

import (
	"fmt"
	"log"
	"os/user"
	"time"

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
	collection, err := cma.Assets.List(config.SpaceID).Next()
	if err != nil {
		log.Fatal(err)
	}

	return collection.ToAsset()
}

func getAsset() *contentful.Asset {
	assetID := getAssets()[0].Sys.ID
	asset, err := cma.Assets.Get(config.SpaceID, assetID)
	if err != nil {
		log.Fatal(err)
	}

	return asset
}

func createAsset() *contentful.Asset {
	asset := &contentful.Asset{
		Fields: &contentful.FileFields{
			Title: "asset title",
			File: &contentful.File{
				Name:        "file name",
				ContentType: "image/jpg",
				UploadURL:   "http://desireeacres.com/gopher.jpg",
			},
		},
	}

	err := cma.Assets.Upsert(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset is created with id: " + asset.Sys.ID)

	return asset
}

func createAndEditAsset() {
	asset := &contentful.Asset{
		Fields: &contentful.FileFields{
			Title: "asset title",
			File: &contentful.File{
				Name:        "file name",
				ContentType: "image/jpg",
				UploadURL:   "http://desireeacres.com/gopher.jpg",
			},
		},
	}

	err := cma.Assets.Upsert(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset is created with id: " + asset.Sys.ID)
	fmt.Println(asset.Fields.Title)

	asset.Fields.Title = "updated title for asset"
	err = cma.Assets.Upsert(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("asset updated with id: " + asset.Sys.ID)
	fmt.Println(asset.Fields.Title)
}

func editExistingAsset() {
	asset := getAssets()[0]
	asset.Fields.Title = "existing asset is updated"

	fmt.Println("updating asset with id: " + asset.Sys.ID)

	err := cma.Assets.Upsert(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	updatedAsset, err := cma.Assets.Get(config.SpaceID, asset.Sys.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new asset with id: " + updatedAsset.Sys.ID)
	fmt.Println(updatedAsset.Fields.Title)
}

func deleteExistingAsset() {
	asset := getAssets()[0]
	fmt.Println("deleting asset with id: " + asset.Sys.ID)

	err := cma.Assets.Delete(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}
}

func createAndProcessAsset() *contentful.Asset {
	asset := createAsset()

	fmt.Println("asset is created with id: " + asset.Sys.ID)
	fmt.Println(asset.Fields.Title)

	err := cma.Assets.Process(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("processing..." + asset.Sys.ID)

	return asset
}

func publishAsset() {
	asset := createAndProcessAsset()
	time.Sleep(3 * time.Second)

	err := cma.Assets.Publish(config.SpaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("published with id : " + asset.Sys.ID)
	fmt.Println("publish date: " + asset.Sys.PublishedAt)
}
