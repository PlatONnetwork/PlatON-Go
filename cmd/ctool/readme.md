## Compile：

Run in this directory： go build ctool.go generate or update ctool.exe file.

## Command:

##### 1.Deploy contract：
```
./ctool deploy -abi "abi json file path(must) " -code "wasm file path (must)" -config "config path(optional)"
```
##### 2.Contract call
```
./ctool invoke -addr "contract address(must) " --func "functon name and param : eg transfer("a",b,c) (must) " --abi "abi file path (must) " -config "config path(optional)"
```
##### 3.Send transaction
```
./ctool sendTransaction -from "msg sender(must)" -to "msg acceptor(must)" -value "transfer value(must)" -config "config path (optional)"
```
##### 4.Send raw transaction
```
./ctool sendRawTransaction -pk "private key file" -from "msg sender(must)" -to "msg acceptor(must)" -value "transfer value(must)" -config "config path (optional)"
```
##### 5.Query transactionReceipt
```
./ctool getTxReceipt -hash "txhash (must)" -config "config path (optional)"
```
##### 6.Prepare transaction stability test account
```
./ctool prepare -pkfile "account private key file path (must)" -size "the number of accounts " value "transfer value" -config "config path (optional)"
```

eg: ./ctool.exe pre -size 10 -pkfile "./test/privateKeys.txt" -value 0xDE0B6B3A7640000

##### 7.Make Stability test
```
./ctool stab -pkfile "account private key file path (must)" -times "send transaction times " -config "config path (optional)"
```

eg:  ./ctool.exe stab -pkfile "./test/privateKeys.txt" -times 10000

note: If the command exits normally,the next time you can continue to run with the generated accounts and the command exits abnormally, you need to re-use the pre command to generate the test accounts.

##### Config Description： The config parameter is not passed in the command, and the `config.json` file in the current directory is read by default.

The config.json file is as follows：

```
{
  "url":"http://192.168.9.73:6789",
  "gas": "0x76c0",
  "gasPrice": "0x9184e72a000",
  "from":"0xfb8c2fa47e84fbde43c97a0859557a36a5fb285b"
}
```


