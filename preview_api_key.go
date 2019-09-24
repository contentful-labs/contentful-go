package contentful

import (
	"fmt"
)

type PreviewAPIKeyService service

func (service *PreviewAPIKeyService) Get(spaceID string, apiKey *APIKey) (*APIKey, error) {
	previewKeyId := apiKey.PreviewAPIKey.Sys.ID
	path := fmt.Sprintf("/spaces/%s/preview_api_keys/%s", spaceID, previewKeyId)
	req, err := service.c.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	var previewKey APIKey
	if err := service.c.do(req, &previewKey); err != nil {
		return nil, err
	}

	return &previewKey, nil
}
