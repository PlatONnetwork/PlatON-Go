## 编译：
在本目录下运行： go build cli.go 生成（更新）cli.exe文件

## 命令:
##### 1.发布合约：
```
./cli -cmd deploy -abi "abi json file path(must) " -code "wasm file path (must)" -config "config path(optional)"
```
##### 2.合约调用
```
./cli -cmd invoke  -addr "contract address(must) " --func "functon name and param : eg transfer("a",b,c) (must) " --abi "abi file path (must) " -config "config path(optional)"
```
##### 3.查询交易receipt
```
./cli -cmd getTxReceipt -hash "txhash (must)" -config "config path (optional)"
```

##### config说明： 命令中不传递config参数，则默认读取本当前目录下的config.json文件

config.json文件如下：

```
{
  "url":"http://192.168.9.73:6789",
  "gas": "0x76c0",
  "gasPrice": "0x9184e72a000",
  "from":"0xfb8c2fa47e84fbde43c97a0859557a36a5fb285b"
}
```
