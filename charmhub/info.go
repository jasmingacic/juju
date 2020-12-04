// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmhub

import (
	"context"
	"net/http"
	"strings"

	"github.com/juju/errors"
	"github.com/kr/pretty"

	"github.com/juju/juju/charmhub/path"
	"github.com/juju/juju/charmhub/transport"
)

// InfoClient defines a client for info requests.
type InfoClient struct {
	path   path.Path
	client RESTClient
	logger Logger
}

// NewInfoClient creates a InfoClient for requesting
func NewInfoClient(path path.Path, client RESTClient, logger Logger) *InfoClient {
	return &InfoClient{
		path:   path,
		client: client,
		logger: logger,
	}
}

// Info requests the information of a given charm. If that charm doesn't exist
// an error stating that fact will be returned.
func (c *InfoClient) Info(ctx context.Context, name string) (transport.InfoResponse, error) {
	c.logger.Tracef("Info(%s)", name)
	var resp transport.InfoResponse
	path, err := c.path.Join(name)
	if err != nil {
		return resp, errors.Trace(err)
	}

	path, err = path.Query("fields", defaultInfoFilter())
	if err != nil {
		return resp, errors.Trace(err)
	}

	restResp, err := c.client.Get(ctx, path, &resp)
	if err != nil {
		return resp, errors.Trace(err)
	}

	if resultErr := resp.ErrorList.Combine(); resultErr != nil {
		if restResp.StatusCode == http.StatusNotFound {
			return resp, errors.NewNotFound(resultErr, "")
		}
		return resp, resultErr
	}

	switch resp.Type {
	case "charm", "bundle":
	default:
		return resp, errors.Errorf("unexpected response type %q, expected charm or bundle", resp.Type)
	}

	c.logger.Tracef("Info() unmarshalled: %s", pretty.Sprint(resp))
	return resp, nil
}

// defaultInfoFilter returns a filter string to retrieve all data
// necessary to fill the transport.InfoResponse.  Without it, we'd
// receive the Name, ID and Type.
func defaultInfoFilter() string {
	filter := defaultResultFilter
	filter = append(filter, appendFilterList("default-release.revision", defaultDownloadFilter)...)
	filter = append(filter, appendFilterList("default-release", infoRevisionFilter)...)
	filter = append(filter, appendFilterList("default-release", defaultChannelFilter)...)
	filter = append(filter, appendFilterList("channel-map.revision", defaultDownloadFilter)...)
	filter = append(filter, appendFilterList("channel-map", infoRevisionFilter)...)
	filter = append(filter, appendFilterList("channel-map", defaultChannelFilter)...)
	filter = append(filter, appendFilterList("default-release.resources", resourceFilter)...)
	filter = append(filter, appendFilterList("channel-map.resources", resourceFilter)...)
	return strings.Join(filter, ",")
}

var infoRevisionFilter = []string{
	"revision.config-yaml",
	"revision.created-at",
	"revision.metadata-yaml",
	"revision.platforms.architecture",
	"revision.platforms.os",
	"revision.platforms.series",
	"revision.revision",
	"revision.version",
}
