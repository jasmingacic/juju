// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package series

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

const distroInfoContents = `version,codename,series,created,release,eol,eol-server
10.04,Firefox,firefox,2009-10-13,2010-04-26,2016-04-26
12.04 LTS,Precise Pangolin,precise,2011-10-13,2012-04-26,2017-04-26
99.04,Focal,focal,2020-04-25,2020-10-17,2365-07-17
`

type SupportedSeriesSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&SupportedSeriesSuite{})

func (s *SupportedSeriesSuite) TestSeriesForTypes(c *gc.C) {
	tmpFile, close := makeTempFile(c, distroInfoContents)
	defer close()

	now := time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)

	info, err := seriesForTypes(tmpFile.Name(), now, "", "")
	c.Assert(err, jc.ErrorIsNil)

	ctrlSeries := info.ControllerSeries()
	c.Assert(ctrlSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty"})

	wrkSeries := info.WorkloadSeries()
	c.Assert(wrkSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty", "centos7", "centos8", "genericlinux", "kubernetes", "opensuseleap", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81"})
}

func (s *SupportedSeriesSuite) TestSeriesForTypesUsingImageStream(c *gc.C) {
	tmpFile, close := makeTempFile(c, distroInfoContents)
	defer close()

	now := time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)

	info, err := seriesForTypes(tmpFile.Name(), now, "focal", "daily")
	c.Assert(err, jc.ErrorIsNil)

	ctrlSeries := info.ControllerSeries()
	c.Assert(ctrlSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty"})

	wrkSeries := info.WorkloadSeries()
	c.Assert(wrkSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty", "centos7", "centos8", "genericlinux", "kubernetes", "opensuseleap", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81"})
}

func (s *SupportedSeriesSuite) TestSeriesForTypesUsingInvalidImageStream(c *gc.C) {
	tmpFile, close := makeTempFile(c, distroInfoContents)
	defer close()

	now := time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)

	info, err := seriesForTypes(tmpFile.Name(), now, "focal", "turtle")
	c.Assert(err, jc.ErrorIsNil)

	ctrlSeries := info.ControllerSeries()
	c.Assert(ctrlSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty"})

	wrkSeries := info.WorkloadSeries()
	c.Assert(wrkSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty", "centos7", "centos8", "genericlinux", "kubernetes", "opensuseleap", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81"})
}

func (s *SupportedSeriesSuite) TestSeriesForTypesUsingInvalidSeries(c *gc.C) {
	tmpFile, close := makeTempFile(c, distroInfoContents)
	defer close()

	now := time.Date(2020, 3, 16, 0, 0, 0, 0, time.UTC)

	info, err := seriesForTypes(tmpFile.Name(), now, "firewolf", "daily")
	c.Assert(err, jc.ErrorIsNil)

	ctrlSeries := info.ControllerSeries()
	c.Assert(ctrlSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty"})

	wrkSeries := info.WorkloadSeries()
	c.Assert(wrkSeries, jc.DeepEquals, []string{"groovy", "focal", "bionic", "xenial", "trusty", "centos7", "centos8", "genericlinux", "kubernetes", "opensuseleap", "win10", "win2008r2", "win2012", "win2012hv", "win2012hvr2", "win2012r2", "win2016", "win2016hv", "win2016nano", "win2019", "win7", "win8", "win81"})
}

func makeTempFile(c *gc.C, content string) (*os.File, func()) {
	tmpfile, err := ioutil.TempFile("", "distroinfo")
	if err != nil {
		c.Assert(err, jc.ErrorIsNil)
	}

	_, err = tmpfile.Write([]byte(content))
	c.Assert(err, jc.ErrorIsNil)

	// Reset the file for reading.
	_, err = tmpfile.Seek(0, 0)
	c.Assert(err, jc.ErrorIsNil)

	return tmpfile, func() {
		err := os.Remove(tmpfile.Name())
		c.Assert(err, jc.ErrorIsNil)
	}
}
