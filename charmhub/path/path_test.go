// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package path

import (
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type PathSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&PathSuite{})

func (s *PathSuite) TestJoin(c *gc.C) {
	rawURL := MustParseURL(c, "http://foobar/v1/path/")

	path := MakePath(rawURL)
	appPath, err := path.Join("entity", "app")
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(appPath.String(), gc.Equals, "http://foobar/v1/path/entity/app")
}

func (s *PathSuite) TestJoinMultipleTimes(c *gc.C) {
	rawURL := MustParseURL(c, "http://foobar/v1/path/")

	path := MakePath(rawURL)
	entityPath, err := path.Join("entity")
	c.Assert(err, jc.ErrorIsNil)

	appPath, err := entityPath.Join("app")
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(appPath.String(), gc.Equals, "http://foobar/v1/path/entity/app")
}

func (s *PathSuite) TestQuery(c *gc.C) {
	rawURL := MustParseURL(c, "http://foobar/v1/path")

	path := MakePath(rawURL)

	newPath, err := path.Query("q", "foo")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(path.String(), gc.Equals, "http://foobar/v1/path")
	c.Assert(newPath.String(), gc.Equals, "http://foobar/v1/path?q=foo")
}

func (s *PathSuite) TestQueryEmptyValue(c *gc.C) {
	rawURL := MustParseURL(c, "http://foobar/v1/path")

	path := MakePath(rawURL)

	newPath, err := path.Query("q", "")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(path.String(), gc.Equals, newPath.String())
	//c.Assert(newPath.String(), gc.Equals, "http://foobar/v1/path?q=foo")
}

func (s *PathSuite) TestJoinQuery(c *gc.C) {
	rawURL := MustParseURL(c, "http://foobar/v1/path")

	path := MakePath(rawURL)
	entityPath, err := path.Join("entity")
	c.Assert(err, jc.ErrorIsNil)

	appPath, err := entityPath.Join("app")
	c.Assert(err, jc.ErrorIsNil)

	newPath, err := appPath.Query("q", "foo")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(appPath.String(), gc.Equals, "http://foobar/v1/path/entity/app")
	c.Assert(newPath.String(), gc.Equals, "http://foobar/v1/path/entity/app?q=foo")
}

func (s *PathSuite) TestMultipleQueries(c *gc.C) {
	rawURL := MustParseURL(c, "http://foobar/v1/path")

	path := MakePath(rawURL)

	newPath, err := path.Query("q", "foo1")
	c.Assert(err, jc.ErrorIsNil)

	newPath, err = newPath.Query("q", "foo2")
	c.Assert(err, jc.ErrorIsNil)

	newPath, err = newPath.Query("x", "bar")
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(path.String(), gc.Equals, "http://foobar/v1/path")
	c.Assert(newPath.String(), gc.Equals, "http://foobar/v1/path?q=foo1&q=foo2&x=bar")
}
