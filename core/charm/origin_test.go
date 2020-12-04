// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm_test

import (
	"github.com/juju/testing"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/charm"
)

type platformSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&platformSuite{})

func (s platformSuite) TestParsePlatform(c *gc.C) {
	tests := []struct {
		Name        string
		Value       string
		Expected    charm.Platform
		ExpectedErr string
	}{{
		Name:        "empty",
		Value:       "",
		ExpectedErr: "platform cannot be empty",
	}, {
		Name:        "empty components",
		Value:       "//",
		ExpectedErr: `architecture in platform "//" not valid`,
	}, {
		Name:        "too many components",
		Value:       "////",
		ExpectedErr: `platform is malformed and has too many components "////"`,
	}, {
		Name:  "architecture",
		Value: "amd64",
		Expected: charm.Platform{
			Architecture: "amd64",
		},
	}, {
		Name:  "architecture and series",
		Value: "amd64/series",
		Expected: charm.Platform{
			Architecture: "amd64",
			Series:       "series",
		},
	}, {
		Name:  "architecture, os and series",
		Value: "amd64/os/series",
		Expected: charm.Platform{
			Architecture: "amd64",
			OS:           "os",
			Series:       "series",
		},
	}, {
		Name:  "architecture, unknown os and series",
		Value: "amd64/unknown/series",
		Expected: charm.Platform{
			Architecture: "amd64",
			OS:           "",
			Series:       "series",
		},
	}, {
		Name:  "architecture, unknown os and unknown series",
		Value: "amd64/unknown/unknown",
		Expected: charm.Platform{
			Architecture: "amd64",
			OS:           "",
			Series:       "",
		},
	}, {
		Name:  "architecture and unknown series",
		Value: "amd64/unknown",
		Expected: charm.Platform{
			Architecture: "amd64",
			OS:           "",
			Series:       "",
		},
	}}
	for k, test := range tests {
		c.Logf("test %q at %d", test.Name, k)
		ch, err := charm.ParsePlatformNormalize(test.Value)
		if test.ExpectedErr != "" {
			c.Assert(err, gc.ErrorMatches, test.ExpectedErr)
		} else {
			c.Assert(ch, gc.DeepEquals, test.Expected)
			c.Assert(err, gc.IsNil)
		}
	}
}

func (s platformSuite) TestString(c *gc.C) {
	tests := []struct {
		Name     string
		Value    string
		Expected string
	}{{
		Name:     "architecture",
		Value:    "amd64",
		Expected: "amd64",
	}, {
		Name:     "architecture and series",
		Value:    "amd64/series",
		Expected: "amd64/series",
	}, {
		Name:     "architecture, os and series",
		Value:    "amd64/os/series",
		Expected: "amd64/os/series",
	}}
	for k, test := range tests {
		c.Logf("test %q at %d", test.Name, k)
		platform, err := charm.ParsePlatformNormalize(test.Value)
		c.Assert(err, gc.IsNil)
		c.Assert(platform.String(), gc.DeepEquals, test.Expected)
	}
}
