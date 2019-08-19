package http

import (
	"context"
	"fmt"
	"github.com/viant/afs/storage"
	"io"
	"net/http"
)

//Download downloads asset for supplied object
func (s *manager) Download(ctx context.Context, object storage.Object, options ...storage.Option) (io.ReadCloser, error) {
	return s.DownloadWithURL(ctx, object.URL(), options...)
}

//Download downloads asset for supplied object
func (s *manager) DownloadWithURL(ctx context.Context, URL string, options ...storage.Option) (io.ReadCloser, error) {
	request, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := s.run(ctx, URL, request, options...)
	if err != nil {
		return nil, err
	}
	if IsStatusOK(response) {
		return response.Body, nil
	}
	return nil, fmt.Errorf("invalid status code: %v", response.StatusCode)
}
