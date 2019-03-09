// Copyright (c) 2018 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bborbe/kafka-dockerhub-version-collector/avro"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

//go:generate counterfeiter -o ../mocks/fetcher.go --fake-name Fetcher . Fetcher
type Fetcher interface {
	Fetch(ctx context.Context, versions chan<- avro.ApplicationVersionAvailable) error
}

func NewFetcher(
	httpClient *http.Client,
	url string,
	repository string,
) Fetcher {
	return &fetcher{
		httpClient: httpClient,
		url:        url,
		repository: repository,
	}
}

type fetcher struct {
	httpClient *http.Client
	url        string
	repository string
}

func (f *fetcher) Fetch(ctx context.Context, versions chan<- avro.ApplicationVersionAvailable) error {
	url := fmt.Sprintf("%s/v2/repositories/%s/tags/", f.url, f.repository)
	for {
		var result struct {
			Next    string `json:"next"`
			Results []struct {
				Name string `json:"name"`
			} `json:"results"`
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return errors.Wrap(err, "build request failed")
		}
		glog.V(1).Infof("%s %s", req.Method, req.URL.String())
		resp, err := f.httpClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "request failed")
		}
		if resp.StatusCode/100 != 2 {
			return errors.New("request status code != 2xx")
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return errors.Wrap(err, "decode json failed")
		}
		resp.Body.Close()
		for _, tag := range result.Results {
			version := avro.ApplicationVersionAvailable{
				App:     "docker.io/" + f.repository,
				Version: tag.Name,
			}
			select {
			case <-ctx.Done():
				glog.Infof("context done => return")
				return nil
			case versions <- version:
				glog.V(2).Infof("send %s:%s", version.App, version.Version)
			}
		}
		if result.Next == "" {
			return nil
		}
		url = result.Next
	}
}
