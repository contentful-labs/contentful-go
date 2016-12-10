package contentful

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpace(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space := c.NewSpace()
	space.Sys.ID = "newspace"
	space.Name = "new space"
	req, err := space.getSaveReq()
	assert.Nil(err)
	assert.Equal("POST", req.Method)
	assert.Equal("/spaces", req.URL.Path)
	assert.Equal("1", req.Header.Get("X-Contentful-Version"))
	space.Save()

	nSpace, _ := c.GetSpace(space.ID())
	assert.Equal(nSpace.Sys.ID, space.Sys.ID)
	assert.Equal(nSpace.Name, space.Name)
}

func TestUpdateSpace(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, _ := c.GetSpace(spaceID)
	space.Name = "changed-space-name"
	req, err := space.getSaveReq()
	assert.Nil(err)
	assert.Equal("PUT", req.Method)
	assert.Equal("/spaces/"+spaceID, req.URL.Path)
	version := strconv.Itoa(space.Sys.Version)
	assert.Equal(version, req.Header.Get("X-Contentful-Version"))
	space.Save()
}

func TestGetLocales(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, _ := c.GetSpace(spaceID)
	col := space.GetLocales()
	assert.IsType(Collection{}, *col)
	col, err := col.Next()
	assert.Nil(err)
	locales := col.ToLocale()
	assert.Equal(1, len(locales))
	assert.Equal(space, locales[0].s)
	assert.Equal(c, locales[0].c)
	assert.Equal("English (United States)", locales[0].Name)
	assert.Equal("en-US", locales[0].Code)
	assert.Equal("en-US", locales[0].FallbackCode)
	assert.Equal(true, locales[0].CDA)
	assert.Equal(true, locales[0].CMA)
	assert.Equal(false, locales[0].Default)
	assert.Equal(false, locales[0].Optional)
	assert.Equal("34N35DoyUQAtaKwWTgZs34", locales[0].Sys.ID)
	assert.Equal(0, locales[0].Sys.Version)
}

func TestGetAssets(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, _ := c.GetSpace(spaceID)
	col := space.GetAssets()
	assert.IsType(Collection{}, *col)
	col, err := col.Next()
	assert.Nil(err)
	assets := col.ToAsset()
	assert.Equal(5, len(assets))
	for _, asset := range assets {
		assert.Equal("Asset", asset.Sys.Type)
	}
	assert.Equal("Doge", assets[0].Fields.Title)
	assert.Equal("//images.contentful.com/cfexampleapi/1x0xpXu4pSGS4OukSyWGUK/cc1239c6385428ef26f4180190532818/doge.jpg", assets[0].Fields.File.URL)
	assert.Equal(522943, assets[0].Fields.File.Detail.Size)
	assert.Equal(5800, assets[0].Fields.File.Detail.Image.Width)
	assert.Equal(4350, assets[0].Fields.File.Detail.Image.Height)

	assert.Equal("hehehe", assets[1].Fields.Title)
	assert.Equal("//images.flinkly.com/222ru4k10hm8/3HNzx9gvJScKku4UmcekYw/997663077456077dde5b5be9bd3c1386/d3b8dad44e5066cfb805e2357469ee64.png", assets[1].Fields.File.URL)
	assert.Equal(6198, assets[1].Fields.File.Detail.Size)
	assert.Equal(206, assets[1].Fields.File.Detail.Image.Width)
	assert.Equal(79, assets[1].Fields.File.Detail.Image.Height)
}

func TestGetAsset(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, _ := c.GetSpace(spaceID)
	asset, err := space.GetAsset("1x0xpXu4pSGS4OukSyWGUK")
	assert.Nil(err)
	assert.IsType(Asset{}, *asset)
	assert.Equal("Asset", asset.Sys.Type)
	assert.Equal("Doge", asset.Fields.Title)
	assert.Equal("//images.contentful.com/cfexampleapi/1x0xpXu4pSGS4OukSyWGUK/cc1239c6385428ef26f4180190532818/doge.jpg", asset.Fields.File.URL)
	assert.Equal(522943, asset.Fields.File.Detail.Size)
	assert.Equal(5800, asset.Fields.File.Detail.Image.Width)
	assert.Equal(4350, asset.Fields.File.Detail.Image.Height)
}

func TestGetAssetWithLocale(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)
	space, _ := c.GetSpace(spaceID)
	asset, err := space.GetAssetWithLocale("3HNzx9gvJScKku4UmcekYw", "de")
	assert.Nil(err)
	assert.IsType(Asset{}, *asset)
	assert.Equal("Asset", asset.Sys.Type)
	assert.Equal("hehehe-de", asset.Fields.Title)
	assert.Equal("//images.flinkly.com/222ru4k10hm8/3HNzx9gvJScKku4UmcekYw/997663077456077dde5b5be9bd3c1386/d3b8dad44e5066cfb805e2357469ee64.png-de", asset.Fields.File.URL)
	assert.Equal(6198, asset.Fields.File.Detail.Size)
	assert.Equal(206, asset.Fields.File.Detail.Image.Width)
	assert.Equal(79, asset.Fields.File.Detail.Image.Height)
}

func TestNewLocale(t *testing.T) {
	setup()
	defer teardown()

	assert := assert.New(t)

	space, _ := c.GetSpace(spaceID)
	locale := space.NewLocale()
	assert.IsType(Locale{}, *locale)
	assert.Equal("Locale", locale.Sys.Type)
	assert.Equal(0, locale.Sys.Version)
	assert.Equal(space, locale.s)
	assert.Equal(c, locale.c)
}
