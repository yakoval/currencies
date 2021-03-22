package update

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type httpSource struct {
	uri string
}

func (es *httpSource) content() (io.Reader, error) {
	req, err := http.NewRequest("GET", es.uri, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Error getting content from URL, code " + resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(content), nil
}
