package contentful

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleAssets_Get() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"
	assetID := "asset-id"

	asset, err := cma.Assets.Get(spaceID, assetID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(asset.Sys.ID)
}

func ExampleAssets_List() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"

	collection, err := cma.Assets.List(spaceID).Next()
	if err != nil {
		log.Fatal(err)
	}

	assets := collection.ToAsset()
	for _, asset := range assets {
		fmt.Println(asset.Sys.ID)
	}
}

func ExampleAssets_Upsert_create() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"

	asset := &Asset{
		Fields: &FileFields{
			Title: "asset title",
			File: &File{
				Name:        "file name",
				ContentType: "image/jpg",
				UploadURL:   "http://desireeacres.com/gopher.jpg",
			},
		},
	}

	err := cma.Assets.Upsert(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleAssets_Upsert_update() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"
	assetID := "asset-id"

	asset, err := cma.Assets.Get(spaceID, assetID)
	if err != nil {
		log.Fatal(err)
	}

	asset.Fields.File.Name = "modified"
	err = cma.Assets.Upsert(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleAssets_Delete() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"
	assetID := "asset-id"

	asset, err := cma.Assets.Get(spaceID, assetID)
	if err != nil {
		log.Fatal(err)
	}

	err = cma.Assets.Delete(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleAssets_Delete_all() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"

	collection, err := cma.Assets.List(spaceID).Next()
	if err != nil {
		log.Fatal(err)
	}

	for _, asset := range collection.ToAsset() {
		err := cma.Assets.Delete(spaceID, asset)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleAssets_Process() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"

	asset := &Asset{
		Fields: &FileFields{
			Title: "asset title",
			File: &File{
				Name:        "file name",
				ContentType: "image/jpg",
				UploadURL:   "http://desireeacres.com/gopher.jpg",
			},
		},
	}

	err := cma.Assets.Upsert(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	err = cma.Assets.Process(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleAssets_Publish() {
	cma := NewCMA("cma-token")
	spaceID := "space-id"

	asset := &Asset{
		Fields: &FileFields{
			Title: "asset title",
			File: &File{
				Name:        "file name",
				ContentType: "image/jpg",
				UploadURL:   "http://desireeacres.com/gopher.jpg",
			},
		},
	}

	err := cma.Assets.Upsert(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}

	err = cma.Assets.Publish(spaceID, asset)
	if err != nil {
		log.Fatal(err)
	}
}

func TestAssets_MarshalJSON(t *testing.T) {}

func TestAssets_UnmarshalJSON(t *testing.T) {}

func TestAssets_List(t *testing.T) {
	var err error
	assert := assert.New(t)
	spaceID := "cfexampleapi"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(r.Method, "GET")
		assert.Equal(r.URL.Path, fmt.Sprintf("/spaces/%s/assets", spaceID))

		checkHeaders(r, assert)

		t := template.New("assets-template")
		t, _ = t.Parse(readTestData("assets.json"))

		w.WriteHeader(200)
		t.Execute(w, struct {
			SpaceID string
		}{
			SpaceID: spaceID,
		})
	})

	// test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// cma client
	cma = NewCMA(CMAToken)
	cma.BaseURL = server.URL

	collection, err := cma.Assets.List(spaceID).Next()
	assert.Nil(err)

	assets := collection.ToAsset()
	fmt.Println(assets)
	assert.Equal(3, len(assets))
	assert.Equal("1x0xpXu4pSGS4OukSyWGUK", assets[0].Sys.ID)
	assert.Equal("happycat", assets[1].Sys.ID)
}
