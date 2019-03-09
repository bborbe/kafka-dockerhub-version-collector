// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"context"

	"github.com/bborbe/kafka-dockerhub-version-collector/avro"
)

type fetcherList struct {
	fetcher []Fetcher
}

func NewFetcherList(
	fetcher []Fetcher,
) Fetcher {
	return &fetcherList{
		fetcher: fetcher,
	}
}

func (f *fetcherList) Fetch(ctx context.Context, versions chan<- avro.ApplicationVersionAvailable) error {
	for _, fetcher := range f.fetcher {
		if err := fetcher.Fetch(ctx, versions); err != nil {
			return err
		}
	}
	return nil
}
