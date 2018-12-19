## Compile：

Run in this directory： go build cli.go Generate (update) cli.exe file.

## Command:

##### 1.Deploy contract：
```
./cli -cmd deploy -abi "abi json file path(must) " -code "wasm file path (must)" -config "config path(optional)"
```
##### 2.Contract call
```
./cli -cmd invoke  -addr "contract address(must) " --func "functon name and param : eg transfer("a",b,c) (must) " --abi "abi file path (must) " -config "config path(optional)"
```
##### 3.Query transactionReceipt
```
./cli -cmd getTxReceipt -hash "txhash (must)" -config "config path (optional)"
```

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


### Notes and Exceptions：

1.Comments are not supported in the configuration file config.json. Comments may cause errors.

2.Please ensure that the node starts normally and starts mining, otherwise the transaction cannot be packaged.

3.The ip and port in the configuration file are guaranteed to be correct, otherwise an exception will be thrown.
```
panic: runtime error: invalid memory address or nil pointer dereference
```

4.Deploy contract: must specify the full path of the contract abi and wasm files.

5.Contract call：
```
  the contract address is not exist ...
```
  The contract was not successfully deploy, the code of the contract could not be obtained, and the contract was successfully issued.
  The contract address parameter in the command was correct.






