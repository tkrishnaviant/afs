package http

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/afs/url"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestManager_Run(t *testing.T) {

	testPort := 8880
	baseURL := fmt.Sprintf("http://localhost:%v", testPort)
	ctx := context.Background()
	var useCases = []struct {
		description string
		URL         string
		expect      string
		putParrot   *http.Response
		hasError    bool
	}{
		{
			description: "asset download",
			URL:         url.Join(baseURL, "/foo/bar.txt"),
			expect:      "test is test",

			putParrot: &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("test is test")),
			},
		},
		{
			description: "not found error download",
			URL:         url.Join(baseURL, "/foo/error.txt"),
			hasError:    true,
		},
	}

	parrots := map[string]*http.Response{}
	for _, useCase := range useCases {
		addURLParrots(http.MethodGet, useCase.URL, useCase.putParrot, parrots)
	}
	go startServer(testPort, parrotHandler(parrots))

	for _, useCase := range useCases {
		manager := newManager()
		response := &http.Response{}
		reader, err := manager.DownloadWithURL(ctx, useCase.URL, response)
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		data, err := ioutil.ReadAll(reader)
		_ = reader.Close()
		assert.EqualValues(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, useCase.expect, string(data), useCase.description)

	}

}
