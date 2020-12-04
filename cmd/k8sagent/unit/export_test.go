// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// +build !windows

package unit

import (
	"github.com/juju/cmd"
	"github.com/juju/names/v4"
	"github.com/juju/utils/v2/voyeur"

	"github.com/juju/juju/agent"
	"github.com/juju/juju/cmd/jujud/agent/agentconf"
	"github.com/juju/juju/cmd/k8sagent/utils"
	"github.com/juju/juju/worker/logsender"
)

type (
	ManifoldsConfig = manifoldsConfig
	K8sUnitAgent    = k8sUnitAgent
)

type K8sUnitAgentTest interface {
	cmd.Command
	DataDir() string
	SetAgentConf(cfg agentconf.AgentConf)
	ChangeConfig(change agent.ConfigMutator) error
	CurrentConfig() agent.Config
	Tag() names.UnitTag
	CharmModifiedVersion() int
}

func NewForTest(
	ctx *cmd.Context,
	bufferedLogger *logsender.BufferedLogWriter,
	configChangedVal *voyeur.Value,
	fileReaderWriter utils.FileReaderWriter,
	environment utils.Environment,
) K8sUnitAgentTest {
	return &k8sUnitAgent{
		ctx:              ctx,
		AgentConf:        agentconf.NewAgentConf(""),
		bufferedLogger:   bufferedLogger,
		configChangedVal: configChangedVal,
		fileReaderWriter: fileReaderWriter,
		environment:      environment,
	}
}

func (c *k8sUnitAgent) SetAgentConf(cfg agentconf.AgentConf) {
	c.AgentConf = cfg
}

func (c *k8sUnitAgent) DataDir() string {
	return c.AgentConf.DataDir()
}
