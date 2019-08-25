# Try it out

```sh
nscli rest-server --chain-id my-chain --trust-node
```

```sh

nscli tx send $(nscli keys show jack --address) $(nscli keys show alice --address) 12nametoken

nscli query tx 2550AE955AC752E9B209AFED0D733597AEB96E77F88BC22199D3C809973A8975

nscli query account $(nscli keys show jack --address)
nscli query account $(nscli keys show alice --address)
```
