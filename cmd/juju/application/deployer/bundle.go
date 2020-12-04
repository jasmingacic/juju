// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package deployer

import (
	"fmt"
	"strings"

	"github.com/juju/charm/v8"
	"github.com/juju/cmd"
	"github.com/juju/errors"

	"github.com/juju/juju/api/application"
	commoncharm "github.com/juju/juju/api/common/charm"
	"github.com/juju/juju/cmd/juju/application/bundle"
	"github.com/juju/juju/cmd/juju/application/store"
	"github.com/juju/juju/cmd/juju/application/utils"
	"github.com/juju/juju/core/constraints"
	"github.com/juju/juju/core/devices"
	"github.com/juju/juju/resource/resourceadapters"
	"github.com/juju/juju/storage"
)

type deployBundle struct {
	model ModelCommand
	steps []DeployStep

	dryRun bool
	force  bool
	trust  bool

	bundleDataSource  charm.BundleDataSource
	bundleDir         string
	bundleURL         *charm.URL
	bundleOverlayFile []string
	origin            commoncharm.Origin

	resolver             Resolver
	authorizer           store.MacaroonGetter
	newConsumeDetailsAPI func(url *charm.OfferURL) (ConsumeDetails, error)
	deployResources      resourceadapters.DeployResourcesFunc

	useExistingMachines bool
	bundleMachines      map[string]string
	bundleStorage       map[string]map[string]storage.Constraints
	bundleDevices       map[string]map[string]devices.Constraints

	targetModelName string
	targetModelUUID string
	controllerName  string
	accountUser     string
}

// deploy is the business logic of deploying a bundle after
// it's been prepared.
func (d *deployBundle) deploy(
	ctx *cmd.Context,
	deployAPI DeployerAPI,
	resolver Resolver,
	macaroonGetter store.MacaroonGetter,
) (rErr error) {
	d.resolver = resolver
	d.authorizer = macaroonGetter
	bakeryClient, err := d.model.BakeryClient()
	if err != nil {
		return errors.Trace(err)
	}

	var ok bool
	if d.targetModelUUID, ok = deployAPI.ModelUUID(); !ok {
		return errors.New("API connection is controller-only (should never happen)")
	}

	if d.targetModelName, _, err = d.model.ModelDetails(); err != nil {
		return errors.Annotatef(err, "could not retrieve model name")
	}

	controllerName, err := d.model.ControllerName()
	if err != nil {
		return errors.Trace(err)
	}
	d.controllerName = controllerName
	accountDetails, err := d.model.CurrentAccountDetails()
	if err != nil {
		return errors.Trace(err)
	}
	d.accountUser = accountDetails.User

	// Compose bundle to be deployed and check its validity before running
	// any pre/post checks.
	var bundleData *charm.BundleData
	if bundleData, err = bundle.ComposeAndVerifyBundle(d.bundleDataSource, d.bundleOverlayFile); err != nil {
		return errors.Annotatef(err, "cannot deploy bundle")
	}
	d.bundleDir = d.bundleDataSource.BasePath()

	// Short-circuit trust checks if the operator specifies '--force'
	if !d.trust {
		if tl := appsRequiringTrust(bundleData.Applications); len(tl) != 0 && !d.force {
			return errors.Errorf(`Bundle cannot be deployed without trusting applications with your cloud credentials.
Please repeat the deploy command with the --trust argument if you consent to trust the following application(s):
  - %s`, strings.Join(tl, "\n  - "))
		}
	}

	for app, applicationSpec := range bundleData.Applications {
		if applicationSpec.Plan != "" {
			for _, step := range d.steps {
				s := step

				charmURL, err := resolveCharmURL(applicationSpec.Charm)
				if err != nil {
					return errors.Trace(err)
				}
				cons, err := constraints.Parse(applicationSpec.Constraints)
				if err != nil {
					return errors.Trace(err)
				}

				platform, err := utils.DeducePlatform(cons, applicationSpec.Series)
				if err != nil {
					return errors.Trace(err)
				}

				origin, err := utils.DeduceOrigin(charmURL, d.origin.CoreChannel(), platform)
				if err != nil {
					return errors.Trace(err)
				}

				deployInfo := DeploymentInfo{
					CharmID: application.CharmID{
						URL:    charmURL,
						Origin: origin,
					},
					ApplicationName: app,
					ApplicationPlan: applicationSpec.Plan,
					ModelUUID:       d.targetModelUUID,
					Force:           d.force,
				}

				err = s.RunPre(deployAPI, bakeryClient, ctx, deployInfo)
				if err != nil {
					return errors.Trace(err)
				}

				defer func() {
					err = errors.Trace(s.RunPost(deployAPI, bakeryClient, ctx, deployInfo, rErr))
					if err != nil {
						rErr = err
					}
				}()
			}
		}
	}
	spec := d.makeBundleDeploySpec(ctx, deployAPI)

	// TODO(ericsnow) Do something with the CS macaroons that were returned?
	// Deploying bundles does not allow the use force, it's expected that the
	// bundle is correct and therefore the charms are also.
	if _, err := bundleDeploy(bundleData, spec); err != nil {
		return errors.Annotate(err, "cannot deploy bundle")
	}
	return nil
}

func (d *deployBundle) makeBundleDeploySpec(ctx *cmd.Context, apiRoot DeployerAPI) bundleDeploySpec {
	// set the consumer details API factory method on the spec, so it makes it
	// possible to communicate with other controllers, that are found within
	// the local cache.
	// If no controller is found within the local cache, an error will be raised
	// which should ask the user to login.
	getConsumeDetails := func(url *charm.OfferURL) (ConsumeDetails, error) {
		// Ensure that we have a url source when querying the controller.
		if url.Source == "" {
			url.Source = d.controllerName
		}
		return d.newConsumeDetailsAPI(url)
	}

	return bundleDeploySpec{
		ctx:                  ctx,
		filesystem:           d.model.Filesystem(),
		dryRun:               d.dryRun,
		force:                d.force,
		trust:                d.trust,
		bundleDataSource:     d.bundleDataSource,
		bundleDir:            d.bundleDir,
		bundleURL:            d.bundleURL,
		bundleOverlayFile:    d.bundleOverlayFile,
		origin:               d.origin,
		deployAPI:            apiRoot,
		bundleResolver:       d.resolver,
		authorizer:           d.authorizer,
		getConsumeDetailsAPI: getConsumeDetails,
		deployResources:      d.deployResources,
		useExistingMachines:  d.useExistingMachines,
		bundleMachines:       d.bundleMachines,
		bundleStorage:        d.bundleStorage,
		bundleDevices:        d.bundleDevices,
		targetModelName:      d.targetModelName,
		targetModelUUID:      d.targetModelUUID,
		controllerName:       d.controllerName,
		accountUser:          d.accountUser,
	}
}

type localBundle struct {
	deployBundle
}

// String returns a string description of the deployer.
func (d *localBundle) String() string {
	str := fmt.Sprintf("deploy local bundle from: %s", d.bundleDir)
	if isEmptyOrigin(d.origin, commoncharm.OriginLocal) {
		return str
	}
	var channel string
	if ch := d.origin.CoreChannel().String(); ch != "" {
		channel = fmt.Sprintf(" from channel %s", ch)
	}
	return fmt.Sprintf("%s%s", str, channel)
}

// PrepareAndDeploy deploys a local bundle, no further preparation is needed.
func (d *localBundle) PrepareAndDeploy(ctx *cmd.Context, deployAPI DeployerAPI, resolver Resolver, macaroonGetter store.MacaroonGetter) error {
	return d.deploy(ctx, deployAPI, resolver, macaroonGetter)
}

type charmstoreBundle struct {
	deployBundle
}

// String returns a string description of the deployer.
func (d *charmstoreBundle) String() string {
	str := fmt.Sprintf("deploy charm store bundle: %s", d.bundleURL.String())
	if isEmptyOrigin(d.origin, commoncharm.OriginCharmStore) {
		return str
	}
	var channel string
	if ch := d.origin.CoreChannel().String(); ch != "" {
		channel = fmt.Sprintf(" from channel %s", ch)
	}
	return fmt.Sprintf("%s%s", str, channel)
}

// PrepareAndDeploy deploys a local bundle, no further preparation is needed.
func (d *charmstoreBundle) PrepareAndDeploy(ctx *cmd.Context, deployAPI DeployerAPI, resolver Resolver, macaroonGetter store.MacaroonGetter) error {
	ctx.Infof("Located bundle %q", d.bundleURL)
	return d.deploy(ctx, deployAPI, resolver, macaroonGetter)
}
