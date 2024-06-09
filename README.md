## Overview
```shell
make lint && make test
```

### prepare to run
```shell
cp configs/sample.app_conf.yml configs/app_conf.yml
```
populate `configs/app_conf.yml` with your eth rpc urls. you can us same urls as fallback for simplification.


usage:
```shell
curl http://127.0.0.1:8000/get_current_block
```

```shell
curl http://127.0.0.1:8000/subscribe?address=0xdAC17F958D2ee523a2206206994597C13D831ec7
```

```shell
curl http://127.0.0.1:8000/get_transactions?address=0xdAC17F958D2ee523a2206206994597C13D831ec7
```

### implementation
* see `internal/service/indexator`
  * `background_block_sync.go` - sync blocks in background
  * `background_transaction_sync.go` - parse transactions in background for subscribed addresses
* see `internal/service/eth_nodes` for rpc calls for eth nodes
  * there are libs only for quick decode strings `0x123` into big.Int. all request without any library
  * `internal/service/eth_nodes/el.go` - methods for fetching blocks in `*ethclient.Client` style
    * `internal/service/eth_nodes/el_test.go` - verify that pure rpc calls fetch correct data