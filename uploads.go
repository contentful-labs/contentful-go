package contentful

import (
	"encoding/json"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
	"strings"
)

//UploadService service
type UploadService service

//Upload model
type Upload struct {
	Sys *Sys `json:"sys"`
}

//UploadFrom model
type UploadFrom struct {
	Sys *Sys `json:"sys"`
}

// MarshalJSON for custom json marshaling
func (upload *Upload) MarshalJSON() ([]byte, error) {
	payload := map[string]interface{}{
		"sys": "",
	}

	payload["sys"] = upload.Sys
	return json.Marshal(payload)
}

// UnmarshalJSON for custom json unmarshaling
func (upload *Upload) UnmarshalJSON(data []byte) error {
	type Alias *Upload

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}

	upload.Sys = &Sys{}
	b, _ := json.Marshal(payload["sys"])
	if err := json.Unmarshal(b, upload.Sys); err != nil {
		return err
	}

	if err := json.Unmarshal(data, Alias(upload)); err != nil {
		return err
	}

	return nil
}

//UploadFile for upload a file in your space
func (s *UploadService) UploadFile(spaceId, fpath string) (up *Upload, err error) {

	var v Upload

	url := fmt.Sprintf("/spaces/%s/uploads", spaceId)

	//Read file
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	//Read Stat
	fstat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	reader := io.LimitReader(f, fstat.Size())
	tmpl := `{{ magenta "{prefix}"}} {{ bar . (magenta "[") "◼" (cycle . "□" ) "□" "]"}} {{speed . | magenta }} {{percent . | magenta}}`
	tmpl = strings.Replace(tmpl, "{prefix}", f.Name(), -1)
	bar := pb.ProgressBarTemplate(tmpl).Start64(fstat.Size())
	bar.Set(pb.Bytes, true)
	bar.Set(pb.SIBytesPrefix, true)
	barReader := bar.NewProxyReader(reader)

	req, err := s.c.newRequest("POST", url, nil, barReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	err = s.c.do(req, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

//RetrievingUpload returns info of item upload
func (s *UploadService) RetrievingUpload(spaceID, uploadID string) (*Upload, error) {
	var v Upload

	url := fmt.Sprintf("/spaces/%s/uploads/%s", spaceID, uploadID)

	req, err := s.c.newRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}

	err = s.c.do(req, &v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *UploadService) DeleteUpload(spaceId, uploadId string) error {

	url := fmt.Sprintf("/spaces/%s/uploads/%s", spaceId, uploadId)

	req, err := s.c.newRequest("DELETE", url, nil, nil)
	if err != nil {
		return err
	}

	err = s.c.do(req, nil)
	if err != nil {
		return err
	}

	return nil
}
