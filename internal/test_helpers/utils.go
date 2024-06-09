package testhelpers

import (
	"os"
	"strings"
	"testing"
	"time"
	"twtt/internal/logger"
	ethnodes "twtt/internal/service/eth_nodes"
	"twtt/internal/utils"

	"golang.org/x/exp/rand"
)

var MainetRPC *utils.RoundRobinBalancer[string]

func InitELBalancer(t *testing.T, log logger.AppLogger) *ethnodes.ELBalancer {
	InitMainetRPC(t)
	return ethnodes.NewELBalancer(log, MainetRPC.Values(), MainetRPC.Values(), 3, nil)
}

//nolint:gocritic
func InitMainetRPC(t *testing.T) {
	rawRPC := os.Getenv("RPC_URLS")
	if rawRPC == "" {
		t.Fatal("RPC_URLS env variable is not set")
	}
	rawRPC = "https://rpc.ankr.com/eth/660eb7902dae14909a7f4def2456a39b636b73165f79870185036990e866da18"
	rpcList := strings.Split(rawRPC, ",")
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano()))) // nolint:gosec
	r.Shuffle(len(rpcList), func(i, j int) { rpcList[i], rpcList[j] = rpcList[j], rpcList[i] })
	for i := range rpcList {
		rpcList[i] = strings.TrimSpace(rpcList[i])
	}
	MainetRPC = utils.NewRoundRobinBalancer(rpcList)
}
