// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package azure_test

import (
	"net/http"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloud"
	environscloudspec "github.com/juju/juju/environs/cloudspec"
	"github.com/juju/juju/provider/azure"
	"github.com/juju/juju/provider/azure/internal/azuretesting"
)

type AuthSuite struct {
	testing.IsolationSuite
	requests []*http.Request
}

var _ = gc.Suite(&AuthSuite{})

func (s *AuthSuite) TestAuthTokenServicePrincipalSecret(c *gc.C) {
	spec := environscloudspec.CloudSpec{
		Type:             "azure",
		Name:             "azure",
		Region:           "westus",
		Endpoint:         "https://api.azurestack.local",
		IdentityEndpoint: "https://graph.azurestack.local",
		StorageEndpoint:  "https://storage.azurestack.local",
		Credential:       fakeServicePrincipalCredential(),
	}
	senders := azuretesting.Senders{
		discoverAuthSender(),
	}
	token, tenantID, err := azure.AuthToken(spec, &senders, "https://resource")
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(token, gc.NotNil)
	c.Assert(tenantID, gc.Equals, fakeTenantId)
}

func (s *AuthSuite) TestAuthTokenInteractive(c *gc.C) {
	spec := environscloudspec.CloudSpec{
		Type:             "azure",
		Name:             "azure",
		Region:           "westus",
		Endpoint:         "https://api.azurestack.local",
		IdentityEndpoint: "https://graph.azurestack.local",
		StorageEndpoint:  "https://storage.azurestack.local",
		Credential:       fakeInteractiveCredential(),
	}
	senders := azuretesting.Senders{}
	_, _, err := azure.AuthToken(spec, &senders, "")
	c.Assert(err, gc.ErrorMatches, `auth-type "interactive" not supported`)
}

func fakeInteractiveCredential() *cloud.Credential {
	cred := cloud.NewCredential("interactive", map[string]string{
		"subscription-id": fakeSubscriptionId,
	})
	return &cred
}
