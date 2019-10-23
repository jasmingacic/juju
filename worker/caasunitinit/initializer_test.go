// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasunitinit_test

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/golang/mock/gomock"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/names.v3"

	"github.com/juju/juju/caas"
	"github.com/juju/juju/caas/kubernetes/provider/exec"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/juju/worker/caasoperator"
	"github.com/juju/juju/worker/caasoperator/mocks"
	"github.com/juju/juju/worker/caasunitinit"
)

type UnitInitializerSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&UnitInitializerSuite{})

func (s *UnitInitializerSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
}

func (s *UnitInitializerSuite) TestInitialize(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockExecClient := mocks.NewMockExecutor(ctrl)

	params := caasunitinit.InitializeUnitParams{
		UnitTag: names.NewUnitTag("gitlab/0"),
		Logger:  loggo.GetLogger("test"),
		Paths: caasoperator.Paths{
			State: caasoperator.StatePaths{
				CharmDir: "dir/charm",
			},
		},
		ExecClient: mockExecClient,
		OperatorInfo: caas.OperatorInfo{
			CACert: "ca-cert",
		},
		UnitProviderIDFunc: func(unit names.UnitTag) (string, error) {
			return "gitlab-ffff", nil
		},
		TempDir: func(dir string, prefix string) (string, error) {
			return filepath.Join(dir, prefix+"-random"), nil
		},
		WriteFile: func(path string, data []byte, perm os.FileMode) error {
			return nil
		},
	}

	gomock.InOrder(
		mockExecClient.EXPECT().Exec(exec.ExecParams{
			Commands:      []string{"mkdir", "-p", filepath.Join(os.TempDir(), "unit-gitlab-0-random")},
			PodName:       "gitlab-ffff",
			ContainerName: "juju-pod-init",
			Stdout:        &bytes.Buffer{},
			Stderr:        &bytes.Buffer{},
		}, gomock.Any()).Return(nil),
		mockExecClient.EXPECT().Copy(exec.CopyParam{
			Src: exec.FileResource{
				Path: "dir/charm",
			},
			Dest: exec.FileResource{
				Path:          filepath.Join(os.TempDir(), "unit-gitlab-0-random"),
				PodName:       "gitlab-ffff",
				ContainerName: "juju-pod-init",
			},
		}, gomock.Any()).Return(nil),
		mockExecClient.EXPECT().Copy(exec.CopyParam{
			Src: exec.FileResource{
				Path: filepath.Join(os.TempDir(), "unit-gitlab-0-random/ca.crt"),
			},
			Dest: exec.FileResource{
				Path:          filepath.Join(os.TempDir(), "unit-gitlab-0-random/ca.crt"),
				PodName:       "gitlab-ffff",
				ContainerName: "juju-pod-init",
			},
		}, gomock.Any()).Return(nil),
		mockExecClient.EXPECT().Copy(exec.CopyParam{
			Src: exec.FileResource{
				Path: "/var/lib/juju/agents/unit-gitlab-0/operator-client-cache.yaml",
			},
			Dest: exec.FileResource{
				Path:          filepath.Join(os.TempDir(), "unit-gitlab-0-random/operator-client-cache.yaml"),
				PodName:       "gitlab-ffff",
				ContainerName: "juju-pod-init",
			},
		}, gomock.Any()).Return(nil),
		mockExecClient.EXPECT().Exec(exec.ExecParams{
			Commands: []string{"/var/lib/juju/tools/jujud", "caas-unit-init",
				"--send", "--unit", "unit-gitlab-0",
				"--charm-dir",
				filepath.Join(os.TempDir(), "unit-gitlab-0-random/charm"),
				"--operator-file",
				filepath.Join(os.TempDir(), "unit-gitlab-0-random/operator-client-cache.yaml"),
				"--operator-ca-cert-file",
				filepath.Join(os.TempDir(), "unit-gitlab-0-random/ca.crt"),
			},
			WorkingDir:    "/var/lib/juju",
			PodName:       "gitlab-ffff",
			ContainerName: "juju-pod-init",
			Stdout:        &bytes.Buffer{},
			Stderr:        &bytes.Buffer{},
		}, gomock.Any()).Return(nil),
	)

	cancel := make(chan struct{})
	err := caasunitinit.InitializeUnit(params, cancel)
	c.Assert(err, jc.ErrorIsNil)
}

func (s *UnitInitializerSuite) TestInitializeUnitMissing(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockExecClient := mocks.NewMockExecutor(ctrl)

	params := caasunitinit.InitializeUnitParams{
		UnitTag: names.NewUnitTag("gitlab/0"),
		Logger:  loggo.GetLogger("test"),
		Paths: caasoperator.Paths{
			State: caasoperator.StatePaths{
				CharmDir: "dir/charm",
			},
		},
		ExecClient: mockExecClient,
		OperatorInfo: caas.OperatorInfo{
			CACert: "ca-cert",
		},
		UnitProviderIDFunc: func(unit names.UnitTag) (string, error) {
			return "", errors.NotFoundf("unit")
		},
		TempDir: func(dir string, prefix string) (string, error) {
			return filepath.Join(dir, prefix+"-random"), nil
		},
		WriteFile: func(path string, data []byte, perm os.FileMode) error {
			return nil
		},
	}

	gomock.InOrder()

	cancel := make(chan struct{})
	err := caasunitinit.InitializeUnit(params, cancel)
	c.Assert(err, gc.ErrorMatches, "unit not found")
}

func (s *UnitInitializerSuite) TestInitializeContainerMissing(c *gc.C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	mockExecClient := mocks.NewMockExecutor(ctrl)

	params := caasunitinit.InitializeUnitParams{
		UnitTag: names.NewUnitTag("gitlab/0"),
		Logger:  loggo.GetLogger("test"),
		Paths: caasoperator.Paths{
			State: caasoperator.StatePaths{
				CharmDir: "dir/charm",
			},
		},
		ExecClient: mockExecClient,
		OperatorInfo: caas.OperatorInfo{
			CACert: "ca-cert",
		},
		UnitProviderIDFunc: func(unit names.UnitTag) (string, error) {
			return "gitlab-ffff", nil
		},
		TempDir: func(dir string, prefix string) (string, error) {
			return filepath.Join(dir, prefix+"-random"), nil
		},
		WriteFile: func(path string, data []byte, perm os.FileMode) error {
			return nil
		},
	}

	gomock.InOrder(
		mockExecClient.EXPECT().Exec(exec.ExecParams{
			Commands:      []string{"mkdir", "-p", filepath.Join(os.TempDir(), "unit-gitlab-0-random")},
			PodName:       "gitlab-ffff",
			ContainerName: "juju-pod-init",
			Stdout:        &bytes.Buffer{},
			Stderr:        &bytes.Buffer{},
		}, gomock.Any()).Return(nil),
		mockExecClient.EXPECT().Copy(exec.CopyParam{
			Src: exec.FileResource{
				Path: "dir/charm",
			},
			Dest: exec.FileResource{
				Path:          filepath.Join(os.TempDir(), "unit-gitlab-0-random"),
				PodName:       "gitlab-ffff",
				ContainerName: "juju-pod-init",
			},
		}, gomock.Any()).Return(errors.NotFoundf("container")),
	)

	cancel := make(chan struct{})
	err := caasunitinit.InitializeUnit(params, cancel)
	c.Assert(err, gc.ErrorMatches, "container not found")
}
