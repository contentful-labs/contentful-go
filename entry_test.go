package contentful

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntriesService_List(t *testing.T) {
	setup()
	defer teardown()
	assert.NotNil(t, c)

	collection, err := c.Entries.List(spaceID)
	assert.NoError(t, err)
	assert.NotNil(t, collection)
	assert.NotEmpty(t, collection.Items)
}

//goos: darwin
//goarch: amd64
//pkg: github.com/contentful-labs/contentful-go
//BenchmarkEntriesService_List-12    	    4153	    287229 ns/op
func BenchmarkEntriesService_List(b *testing.B) {
	setup()
	defer teardown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Entries.List(spaceID)
	}
}

func TestEntriesService_Get(t *testing.T) {
	setup()
	defer teardown()

	assert.NotNil(t, c)
	entry, err := c.Entries.Get(spaceID, entryID)
	assert.NoError(t, err)
	assert.NotNil(t, entry.Sys)
	assert.Equal(t, entryID, entry.Sys.ID)
}

//goos: darwin
//goarch: amd64
//pkg: github.com/contentful-labs/contentful-go
//BenchmarkEntriesService_Get-12    	   10035	    115806 ns/op
func BenchmarkEntriesService_Get(b *testing.B) {
	setup()
	defer teardown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Entries.Get(spaceID, entryID)
	}
}

func TestEntriesService_Delete(t *testing.T) {
	setup()
	defer teardown()
	assert.NotNil(t, c)

	err := c.Entries.Delete(spaceID, entryID)
	assert.NoError(t, err)
}

//goos: darwin
//goarch: amd64
//pkg: github.com/contentful-labs/contentful-go
//BenchmarkEntriesService_Delete-12    	    4521	    226502 ns/op
func BenchmarkEntriesService_Delete(b *testing.B) {
	setup()
	defer teardown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Entries.Delete(spaceID, entryID)
	}
}

func TestEntriesService_Create(t *testing.T) {
	setup()
	defer teardown()
	assert.NotNil(t, c)

	err := c.Entries.Create(spaceID, &Entry{
		Fields: map[string]interface{}{
			"name": "Go Cat",
		},
	})
	assert.NoError(t, err)
}

//goos: darwin
//goarch: amd64
//pkg: github.com/contentful-labs/contentful-go
//BenchmarkEntriesService_Create-12    	    4257	    395626 ns/op
func BenchmarkEntriesService_Create(b *testing.B) {
	setup()
	defer teardown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Entries.Create(spaceID, &Entry{
			Fields: map[string]interface{}{
				"name": "Go Cat",
			},
		})
	}
}

func TestEntriesService_Publish(t *testing.T) {
	setup()
	defer teardown()
	assert.NotNil(t, c)

	happyCatEntry, err := c.Entries.Get(spaceID, entryID)
	assert.NoError(t, err)
	assert.NotNil(t, happyCatEntry)

	err = c.Entries.Publish(spaceID, &happyCatEntry)
	assert.NoError(t, err)
}

//goos: darwin
//goarch: amd64
//pkg: github.com/contentful-labs/contentful-go
//BenchmarkEntriesService_Publish-12    	    4759	    233542 ns/op
func BenchmarkEntriesService_Publish(b *testing.B) {
	setup()
	defer teardown()
	happyCatEntry, _ := c.Entries.Get(spaceID, entryID)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Entries.Publish(spaceID, &happyCatEntry)
	}
}

func TestEntriesService_UnPublish(t *testing.T) {
	setup()
	defer teardown()
	assert.NotNil(t, c)

	happyCatEntry, err := c.Entries.Get(spaceID, entryID)
	assert.NoError(t, err)
	assert.NotNil(t, happyCatEntry)

	err = c.Entries.UnPublish(spaceID, &happyCatEntry)
	assert.NoError(t, err)
}

//goos: darwin
//goarch: amd64
//pkg: github.com/contentful-labs/contentful-go
//BenchmarkEntriesService_UnPublish-12    	    4731	    234145 ns/op
func BenchmarkEntriesService_UnPublish(b *testing.B) {
	setup()
	defer teardown()
	happyCatEntry, _ := c.Entries.Get(spaceID, entryID)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.Entries.UnPublish(spaceID, &happyCatEntry)
	}
}