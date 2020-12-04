// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"time"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/gnuflag"
	"github.com/juju/names/v4"

	"github.com/juju/juju/apiserver/params"
	jujucmd "github.com/juju/juju/cmd"
	"github.com/juju/juju/cmd/modelcmd"
	"github.com/juju/juju/cmd/plugins/juju-wait-for/api"
	"github.com/juju/juju/cmd/plugins/juju-wait-for/query"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/core/status"
)

func newApplicationCommand() cmd.Command {
	cmd := &applicationCommand{}
	cmd.newWatchAllAPIFunc = func() (api.WatchAllAPI, error) {
		client, err := cmd.NewAPIClient()
		if err != nil {
			return nil, errors.Trace(err)
		}
		return watchAllAPIShim{
			Client: client,
		}, nil
	}
	return modelcmd.Wrap(cmd)
}

const applicationCommandDoc = `
Wait for a given application to reach a goal state.
arguments:
name
   application name identifier
options:
--query (= 'life=="alive" && status=="active"')
   query represents the goal state of a given application
`

// applicationCommand defines a command for waiting for applications.
type applicationCommand struct {
	waitForCommandBase

	name    string
	query   string
	timeout time.Duration

	found   bool
	appInfo params.ApplicationInfo
}

// Info implements Command.Info.
func (c *applicationCommand) Info() *cmd.Info {
	return jujucmd.Info(&cmd.Info{
		Name:    "application",
		Args:    "[<name>]",
		Purpose: "wait for an application to reach a goal state",
		Doc:     applicationCommandDoc,
	})
}

// SetFlags implements Command.SetFlags.
func (c *applicationCommand) SetFlags(f *gnuflag.FlagSet) {
	c.waitForCommandBase.SetFlags(f)
	f.StringVar(&c.query, "query", `life=="alive" && status=="active"`, "query the goal state")
	f.DurationVar(&c.timeout, "timeout", time.Minute*10, "how long to wait, before timing out")
}

// Init implements Command.Init.
func (c *applicationCommand) Init(args []string) (err error) {
	if len(args) == 0 {
		return errors.New("application name must be supplied when waiting for an application")
	}
	if len(args) != 1 {
		return errors.New("only one application name can be supplied as an argument to this command")
	}
	if ok := names.IsValidApplication(args[0]); !ok {
		return errors.Errorf("%q is not valid application name", args[0])
	}
	c.name = args[0]

	return nil
}

func (c *applicationCommand) Run(ctx *cmd.Context) error {
	strategy := &Strategy{
		ClientFn: c.newWatchAllAPIFunc,
		Timeout:  c.timeout,
	}
	err := strategy.Run(c.name, c.query, c.waitFor)
	return errors.Trace(err)
}

func (c *applicationCommand) waitFor(name string, deltas []params.Delta, q query.Query) (bool, error) {
	for _, delta := range deltas {
		logger.Tracef("delta %T: %v", delta.Entity, delta.Entity)

		switch entityInfo := delta.Entity.(type) {
		case *params.ApplicationInfo:
			if entityInfo.Name == name {
				if delta.Removed {
					return false, errors.Errorf("application %v removed", name)
				}

				scope := MakeApplicationScope(entityInfo)
				if done, err := runQuery(q, scope); err != nil {
					return false, errors.Trace(err)
				} else if done {
					return true, nil
				}

				c.found = entityInfo.Life != life.Dead
				c.appInfo = *entityInfo
			}
		}
	}

	if !c.found {
		logger.Infof("application %q not found, waiting...", name)
		return false, nil
	}

	currentStatus := c.appInfo.Status.Current

	units := make(map[string]*params.UnitInfo)
	for _, delta := range deltas {
		switch entityInfo := delta.Entity.(type) {
		case *params.UnitInfo:
			if delta.Removed {
				delete(units, entityInfo.Name)
			}
			if entityInfo.Application == name {
				units[entityInfo.Name] = entityInfo
			}
		}
	}

	logOutput := currentStatus.String() != "unset" && len(units) > 0

	appInfo := c.appInfo
	appInfo.Status.Current = deriveApplicationStatus(currentStatus, units)

	scope := MakeApplicationScope(&appInfo)
	if done, err := runQuery(q, scope); err != nil {
		return false, errors.Trace(err)
	} else if done {
		return true, nil
	}

	if logOutput {
		logger.Infof("application %q found with %q, waiting for goal state", name, currentStatus)
	}

	return false, nil
}

// ApplicationScope allows the query to introspect a application entity.
type ApplicationScope struct {
	ApplicationInfo *params.ApplicationInfo
}

// MakeApplicationScope creates an ApplicationScope from an ApplicationInfo
func MakeApplicationScope(info *params.ApplicationInfo) ApplicationScope {
	return ApplicationScope{
		ApplicationInfo: info,
	}
}

// GetIdents returns the identifiers with in a given scope.
func (m ApplicationScope) GetIdents() []string {
	return getIdents(m.ApplicationInfo)
}

// GetIdentValue returns the value of the identifier in a given scope.
func (m ApplicationScope) GetIdentValue(name string) (query.Box, error) {
	switch name {
	case "name":
		return query.NewString(m.ApplicationInfo.Name), nil
	case "life":
		return query.NewString(string(m.ApplicationInfo.Life)), nil
	case "exposed":
		return query.NewBool(m.ApplicationInfo.Exposed), nil
	case "charm-url":
		return query.NewString(m.ApplicationInfo.CharmURL), nil
	case "min-units":
		return query.NewInteger(int64(m.ApplicationInfo.MinUnits)), nil
	case "subordinate":
		return query.NewBool(m.ApplicationInfo.Subordinate), nil
	case "status":
		return query.NewString(string(m.ApplicationInfo.Status.Current)), nil
	case "workload-version":
		return query.NewString(m.ApplicationInfo.WorkloadVersion), nil
	}
	return nil, errors.Annotatef(query.ErrInvalidIdentifier(name), "Runtime Error: identifier %q not found on ApplicationInfo", name)
}

func deriveApplicationStatus(currentStatus status.Status, units map[string]*params.UnitInfo) status.Status {
	// If the application is unset, then derive it from the units.
	if currentStatus.String() != "unset" {
		return currentStatus
	}

	statuses := make([]status.StatusInfo, 0)
	for _, unit := range units {
		agentStatus := unit.WorkloadStatus
		statuses = append(statuses, status.StatusInfo{
			Status: agentStatus.Current,
		})
	}

	derived := status.DeriveStatus(statuses)
	return derived.Status
}
