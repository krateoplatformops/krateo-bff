//go:build integration
// +build integration

package resolvers_test

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo-bff/apis/core"
	"github.com/krateoplatformops/krateo-bff/internal/resolvers"
)

func TestResolveCardTemplateWithEval(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	opts := resolvers.CardTemplateGetOneOpts{
		RESTConfig: rc,
		AuthnNS:    "default",
		Username:   "demo",
	}
	nfo, err := resolvers.CardTemplateGetOne(context.TODO(),
		&core.Reference{
			Name: "card-dev", Namespace: "dev-system",
		}, opts)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(nfo)
}

func TestResolveCardTemplateWithoutEval(t *testing.T) {
	rc, err := newRestConfig()
	if err != nil {
		t.Fatal(err)
	}

	opts := resolvers.CardTemplateGetOneOpts{
		RESTConfig: rc,
		AuthnNS:    "default",
		Username:   "demo",
	}

	nfo, err := resolvers.CardTemplateGetOne(context.TODO(),
		&core.Reference{
			Name: "card-dev", Namespace: "dev-system",
		}, opts)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(nfo)
}
