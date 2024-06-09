package indexator_test

import (
	"testing"
	"twtt/internal/service/indexator"
	testhelpers "twtt/internal/test_helpers"
)

func TestService_FetchBlocks(t *testing.T) {
	container := testhelpers.GetClean(t)
	balancerEL := testhelpers.InitELBalancer(t, container.Logger)
	srv := indexator.NewService(container.Ctx, container.Logger, balancerEL, nil)
	srv.FetchBlocks()
}
