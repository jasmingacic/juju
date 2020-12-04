// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charms_test

import (
	"github.com/golang/mock/gomock"
	charm "github.com/juju/charm/v8"
	"github.com/juju/charm/v8/resource"
	"github.com/juju/systems"
	"github.com/juju/systems/channel"
	"github.com/juju/version"
	gc "gopkg.in/check.v1"

	basemocks "github.com/juju/juju/api/base/mocks"
	apicommoncharms "github.com/juju/juju/api/common/charms"
	"github.com/juju/juju/apiserver/params"
	coretesting "github.com/juju/juju/testing"
)

type charmsMockSuite struct {
	coretesting.BaseSuite
	charmsCommonClient *apicommoncharms.CharmsClient
}

var _ = gc.Suite(&charmsMockSuite{})

func (s *charmsMockSuite) TestCharmInfo(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockFacadeCaller := basemocks.NewMockFacadeCaller(ctrl)

	url := "local:quantal/dummy-1"
	args := params.CharmURL{URL: url}
	info := new(params.Charm)

	params := params.Charm{
		Revision: 1,
		URL:      url,
		Config: map[string]params.CharmOption{
			"config": {
				Type:        "type",
				Description: "config-type option",
			},
		},
		LXDProfile: &params.CharmLXDProfile{
			Description: "LXDProfile",
			Devices: map[string]map[string]string{
				"tun": {
					"path": "/dev/net/tun",
					"type": "unix-char",
				},
			},
		},
		Meta: &params.CharmMeta{
			Name:           "dummy",
			Description:    "cockroachdb",
			MinJujuVersion: "2.9.0",
			Resources: map[string]params.CharmResourceMeta{
				"cockroachdb-image": {
					Type:        "oci-image",
					Description: "OCI image used for cockroachdb",
				},
			},
			Systems: []params.CharmSystem{
				{
					OS:      "ubuntu",
					Channel: "20.04/stable",
				},
			},
			Containers: map[string]params.CharmContainer{
				"cockroachdb": {
					Systems: []params.CharmSystem{
						{
							Resource: "cockroachdb-image",
						},
					},
					Mounts: []params.CharmMount{
						{
							Storage:  "database",
							Location: "/cockroach/cockroach-data",
						},
					},
				},
			},
			Platforms:     []string{"kubernetes"},
			Architectures: []string{"amd64"},
			Storage: map[string]params.CharmStorage{
				"database": {
					Type: "filesystem",
				},
			},
		},
	}

	mockFacadeCaller.EXPECT().FacadeCall("CharmInfo", args, info).SetArg(2, params).Return(nil)

	client := apicommoncharms.NewCharmsClient(mockFacadeCaller)
	got, err := client.CharmInfo(url)
	c.Assert(err, gc.IsNil)

	want := &apicommoncharms.CharmInfo{
		Revision: 1,
		URL:      url,
		Config: &charm.Config{
			Options: map[string]charm.Option{
				"config": {
					Type:        "type",
					Description: "config-type option",
				},
			},
		},
		LXDProfile: &charm.LXDProfile{
			Description: "LXDProfile",
			Config:      map[string]string{},
			Devices: map[string]map[string]string{
				"tun": {
					"path": "/dev/net/tun",
					"type": "unix-char",
				},
			},
		},
		Meta: &charm.Meta{
			Name:           "dummy",
			Description:    "cockroachdb",
			MinJujuVersion: version.MustParse("2.9.0"),
			Resources: map[string]resource.Meta{
				"cockroachdb-image": {
					Type:        resource.TypeContainerImage,
					Description: "OCI image used for cockroachdb",
				},
			},
			Systems: []systems.System{
				{
					OS: "ubuntu",
					Channel: channel.Channel{
						Name:  "20.04/stable",
						Risk:  "stable",
						Track: "20.04",
					},
				},
			},
			Containers: map[string]charm.Container{
				"cockroachdb": {
					Systems: []systems.System{
						{
							Resource: "cockroachdb-image",
						},
					},
					Mounts: []charm.Mount{
						{
							Storage:  "database",
							Location: "/cockroach/cockroach-data",
						},
					},
				},
			},
			Platforms:     []charm.Platform{"kubernetes"},
			Architectures: []charm.Architecture{"amd64"},
			Storage: map[string]charm.Storage{
				"database": {
					Type: "filesystem",
				},
			},
		},
	}
	c.Assert(got, gc.DeepEquals, want)
}
