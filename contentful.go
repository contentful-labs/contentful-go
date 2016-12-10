package contentful

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/moul/http2curl"
)

// Contentful model
type Contentful struct {
	client      *http.Client
	api         string
	token       string
	Debug       bool
	QueryParams map[string]string
	Headers     map[string]string
	BaseURL     string
}

// NewCMA returns a CMA client
func NewCMA(token string) *Contentful {
	return &Contentful{
		client: http.DefaultClient,
		api:    "CMA",
		token:  token,
		Debug:  false,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
			"Content-Type":  "application/vnd.contentful.management.v1+json",
		},
		BaseURL: "https://api.contentful.com",
	}
}

// NewCDA returns a CDA client
func NewCDA(token string) *Contentful {
	return &Contentful{
		client: http.DefaultClient,
		api:    "CDA",
		token:  token,
		Debug:  false,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		BaseURL: "https://cda.contentful.com",
	}
}

// NewCPA returns a CPA client
func NewCPA(token string) *Contentful {
	return &Contentful{
		client: http.DefaultClient,
		Debug:  false,
		api:    "CPA",
		token:  token,
		Headers: map[string]string{
			"Authorization": "Bearer " + token,
		},
		BaseURL: "https://preview.contentful.com",
	}
}

// SetOrganization sets the given organization id
func (c *Contentful) SetOrganization(organizationID string) *Contentful {
	c.Headers["X-Contentful-Organization"] = organizationID

	return c
}

// NewSpace creates a new space for the given token
func (c *Contentful) NewSpace() *Space {
	return &Space{
		c: c,
		Sys: &Sys{
			Type:    "Space",
			Version: 1,
		},
	}
}

// GetSpace returns a space
func (c *Contentful) GetSpace(spaceID string) (*Space, error) {
	req, err := c.newRequest("GET", "/spaces/"+spaceID, nil, nil)
	if err != nil {
		return &Space{}, err
	}

	var space Space
	if ok := c.do(req, &space); ok != nil {
		return &Space{}, ok
	}

	space.c = c
	return &space, nil
}

// GetSpaces creates a collection
func (c *Contentful) GetSpaces() *Collection {
	req, _ := c.newRequest("GET", "/spaces", nil, nil)

	col := NewCollection(&CollectionOptions{})
	col.c = c
	col.req = req

	return col
}

func (c *Contentful) newRequest(method, path string, query url.Values, body io.Reader) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, err
	}

	// set query params
	for key, value := range c.QueryParams {
		query.Set(key, value)
	}

	u.Path = path
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// set headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (c *Contentful) do(req *http.Request, v interface{}) error {
	if c.Debug == true {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Println(command)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return c.handleError(req, res)
	}

	if v != nil {
		defer res.Body.Close()
		err = json.NewDecoder(res.Body).Decode(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Contentful) handleError(req *http.Request, res *http.Response) error {
	var e ErrorResponse
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&e)
	if err != nil {
		return err
	}

	if res.StatusCode == 500 {
		fmt.Println("500 response")
	}

	if e.Sys.ID == "NotFound" {
		return NotFoundError{
			APIError{
				req: req,
				res: res,
				err: &e,
			},
		}
	}

	if e.Sys.ID == "AccessTokenInvalid" {
		return AccessTokenInvalidError{
			APIError{
				req: req,
				res: res,
				err: &e,
			},
		}
	}

	if e.Sys.ID == "VersionMismatch" || e.Sys.ID == "Conflict" {
		return VersionMismatchError{
			APIError{
				req: req,
				res: res,
				err: &e,
			},
		}
	}

	if e.Sys.ID == "ValidationFailed" {
		return ValidationFailedError{
			APIError{
				req: req,
				res: res,
				err: &e,
			},
		}
	}

	fmt.Println(e.Sys.ID)

	return e
}
