package contentful

import (
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryInclude(t *testing.T) {
	q := NewQuery().Include(5)
	expected := url.Values{}
	expected.Set("include", "5")
	assert.Equal(t, expected.Encode(), q.String())

	assert.Panics(t, func() {
		q := NewQuery().Include(11)
		q.String()
	}, "out of range `include` should panic")
}

func TestQueryContentType(t *testing.T) {
	q := NewQuery().ContentType("content_type")
	expected := url.Values{}
	expected.Set("content_type", "content_type")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQuerySelect(t *testing.T) {
	q := NewQuery().
		ContentType("ct").
		Select([]string{"field1", "field2"})

	expected := url.Values{}
	expected.Set("content_type", "ct")
	expected.Set("select", "field1,field2")
	assert.Equal(t, expected.Encode(), q.String())

	assert.Panics(t, func() {
		q := NewQuery().Select([]string{"field1", "field2"})
		expected := url.Values{}
		expected.Set("select", "field1,field2")
		assert.Equal(t, expected.Encode(), q.String())
	}, "select needs content_type")

	assert.Panics(t, func() {
		fields := []string{}
		for i := 0; i < 110; i++ {
			fields = append(fields, "field"+strconv.Itoa(i))
		}

		q := NewQuery().Select(fields)
		q.String()
	}, "select accepts 100 fields max")

	assert.Panics(t, func() {
		q := NewQuery().Select([]string{"field1", "field2.d1", "field3.d2.d3"})
		q.String()
	}, "select accepts depths 3 max")
}

func TestQueryEqual(t *testing.T) {
	q := NewQuery().Equal("field1", 10)
	expected := url.Values{}
	expected.Set("field1", "10")
	assert.Equal(t, expected.Encode(), q.String())

	q = q.Equal("field1", "11")
	expected.Set("field1", "11")
	assert.Equal(t, expected.Encode(), q.String())

	q = q.Equal("field1", time.Now())
	expected.Del("field1")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryNotEqual(t *testing.T) {
	q := NewQuery().NotEqual("field1", 10)
	expected := url.Values{}

	expected.Set("field1[ne]", "10")
	assert.Equal(t, expected.Encode(), q.String())

	q = q.NotEqual("field1", "11")
	expected.Set("field1[ne]", "11")
	assert.Equal(t, expected.Encode(), q.String())

	q = q.NotEqual("field1", time.Now())
	expected.Del("field1[ne]")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryAll(t *testing.T) {
	q := NewQuery().All("field1", []string{"10", "test"})
	expected := url.Values{}
	expected.Set("field1[all]", "10,test")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryIn(t *testing.T) {
	q := NewQuery().In("sys.id", []string{"test", "test2"})
	expected := url.Values{}
	expected.Set("sys.id[in]", "test,test2")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryNotIn(t *testing.T) {
	q := NewQuery().NotIn("sys.id", []string{"test3"})
	expected := url.Values{}
	expected.Set("sys.id[nin]", "test3")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryExists(t *testing.T) {
	q := NewQuery().Exists("sys.id")
	expected := url.Values{}
	expected.Set("sys.id[exists]", "true")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryNotExists(t *testing.T) {
	q := NewQuery().NotExists("sys.id")
	expected := url.Values{}
	expected.Set("sys.id[exists]", "false")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryLessThan(t *testing.T) {
	q := NewQuery().LessThan("fields.date", 10)
	expected := url.Values{}
	expected.Set("fields.date[lt]", "10")
	assert.Equal(t, expected.Encode(), q.String())

	now := time.Now()
	q = NewQuery().LessThan("fields.date", now)
	expected = url.Values{}
	expected.Set("fields.date[lt]", now.Format("2006-01-02 15:04:05"))
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryLessThanOrEqual(t *testing.T) {
	q := NewQuery().LessThanOrEqual("fields.date", 10)
	expected := url.Values{}
	expected.Set("fields.date[lte]", "10")
	assert.Equal(t, expected.Encode(), q.String())

	now := time.Now()
	q = NewQuery().LessThanOrEqual("fields.date", now)
	expected = url.Values{}
	expected.Set("fields.date[lte]", now.Format("2006-01-02 15:04:05"))
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryGreaterThan(t *testing.T) {
	q := NewQuery().GreaterThan("fields.date", 10)
	expected := url.Values{}
	expected.Set("fields.date[gt]", "10")
	assert.Equal(t, expected.Encode(), q.String())

	now := time.Now()
	q = NewQuery().GreaterThan("fields.date", now)
	expected = url.Values{}
	expected.Set("fields.date[gt]", now.Format("2006-01-02 15:04:05"))
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryGreaterThanOrEqual(t *testing.T) {
	q := NewQuery().GreaterThanOrEqual("fields.date", 10)
	expected := url.Values{}
	expected.Set("fields.date[gte]", "10")
	assert.Equal(t, expected.Encode(), q.String())

	now := time.Now()
	q = NewQuery().GreaterThanOrEqual("fields.date", now)
	expected = url.Values{}
	expected.Set("fields.date[gte]", now.Format("2006-01-02 15:04:05"))
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryQuery(t *testing.T) {
	q := NewQuery().Query("query_str")
	expected := url.Values{}
	expected.Set("query", "query_str")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryMatch(t *testing.T) {
	q := NewQuery().Match("field1", "match_query")
	expected := url.Values{}
	expected.Set("field1[match]", "match_query")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryNear(t *testing.T) {
	q := NewQuery().Near("field1", 38, -120)
	expected := url.Values{}
	expected.Set("field1[near]", "38,-120")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryWithin(t *testing.T) {
	q := NewQuery().Within("field1", 38, -120, 10, 120)
	expected := url.Values{}
	expected.Set("field1[within]", "38,-120,10,120")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryWithinRadius(t *testing.T) {
	q := NewQuery().WithinRadius("field1", 38, -120, 22)
	expected := url.Values{}
	expected.Set("field1[within]", "38,-120,22")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryOrder(t *testing.T) {
	q := NewQuery().ContentType("ct").Order("field1", false)
	expected := url.Values{}
	expected.Set("content_type", "ct")
	expected.Set("order", "field1")
	assert.Equal(t, expected.Encode(), q.String())

	q = NewQuery().ContentType("ct").Order("field1", true)
	expected = url.Values{}
	expected.Set("content_type", "ct")
	expected.Set("order", "-field1")
	assert.Equal(t, expected.Encode(), q.String())

	q = NewQuery().
		ContentType("ct").
		Order("field1", true).
		Order("field2", false).
		Order("field3", false)

	expected = url.Values{}
	expected.Set("content_type", "ct")
	expected.Set("order", "-field1,field2,field3")
	assert.Equal(t, expected.Encode(), q.String())

	// assert.Panics(t, func() {
	// q := NewQuery().Order("field1", false)
	// q.String()
	// }, "out of range limit should panic")
}

func TestQueryLimit(t *testing.T) {
	q := NewQuery().Limit(10)
	expected := url.Values{}
	expected.Set("limit", "10")
	assert.Equal(t, expected.Encode(), q.String())

	assert.Panics(t, func() {
		q := NewQuery().Limit(3000)
		q.String()
	}, "out of range limit should panic")
}

func TestQuerySkip(t *testing.T) {
	q := NewQuery().Skip(10)
	expected := url.Values{}
	expected.Set("skip", "10")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryMimeType(t *testing.T) {
	q := NewQuery().MimeType("image")
	expected := url.Values{}
	expected.Set("mimetype_group", "image")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryLinksToEntry(t *testing.T) {
	q := NewQuery().LinksToEntry("21tWpTXe2XBHtS1Ytg7n5C")
	expected := url.Values{}
	expected.Set("links_to_entry", "21tWpTXe2XBHtS1Ytg7n5C")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQueryLinksToAsset(t *testing.T) {
	q := NewQuery().LinksToAsset("21tWpTXe2XBHtS1Ytg7n5C")
	expected := url.Values{}
	expected.Set("links_to_asset", "21tWpTXe2XBHtS1Ytg7n5C")
	assert.Equal(t, expected.Encode(), q.String())
}

func TestQuery(t *testing.T) {
	q := NewQuery().
		Equal("cat.name", "catname").
		NotEqual("cat.name", "dogname").
		In("sys.id", []string{"test", "test2"}).
		NotIn("sys.id", []string{"test3"}).
		LessThan("fields.cat", 4)

	expected := url.Values{}
	expected.Set("cat.name", "catname")
	expected.Set("cat.name[ne]", "dogname")
	expected.Set("sys.id[in]", "test,test2")
	expected.Set("sys.id[nin]", "test3")
	expected.Set("fields.cat[lt]", "4")

	assert.Equal(t, expected.Encode(), q.String())
}
