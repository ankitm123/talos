// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package aws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

// MetadataConfig represents a metadata AWS instance.
type MetadataConfig struct {
	Hostname     string `json:"hostname,omitempty"`
	InstanceID   string `json:"instance-id,omitempty"`
	InstanceType string `json:"instance-type,omitempty"`
	PublicIPv4   string `json:"public-ipv4,omitempty"`
	PublicIPv6   string `json:"ipv6,omitempty"`
	Region       string `json:"region,omitempty"`
	Zone         string `json:"zone,omitempty"`
}

//nolint:gocyclo
func (a *AWS) getMetadata(ctx context.Context) (*MetadataConfig, error) {
	getMetadataKey := func(key string) (v string, err error) {
		v, err = a.metadataClient.GetMetadataWithContext(ctx, key)
		if err != nil {
			if awsErr, ok := err.(awserr.RequestFailure); ok {
				if awsErr.StatusCode() == http.StatusNotFound {
					return "", nil
				}
			}

			return "", fmt.Errorf("failed to fetch %q from IMDS: %w", key, err)
		}

		return v, nil
	}

	var (
		metadata MetadataConfig
		err      error
	)

	if metadata.Hostname, err = getMetadataKey("hostname"); err != nil {
		return nil, err
	}

	if metadata.InstanceType, err = getMetadataKey("instance-type"); err != nil {
		return nil, err
	}

	if metadata.InstanceID, err = getMetadataKey("instance-id"); err != nil {
		return nil, err
	}

	if metadata.PublicIPv4, err = getMetadataKey("public-ipv4"); err != nil {
		return nil, err
	}

	if metadata.PublicIPv6, err = getMetadataKey("ipv6"); err != nil {
		return nil, err
	}

	if metadata.Region, err = getMetadataKey("placement/region"); err != nil {
		return nil, err
	}

	if metadata.Zone, err = getMetadataKey("placement/availability-zone"); err != nil {
		return nil, err
	}

	return &metadata, nil
}
