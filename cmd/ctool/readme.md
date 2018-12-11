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


### 注意事项和异常说明：
1.配置文件config.json中不支持注释，有注释可能会引起错误

2.请确保节点正常启动,并开启挖矿,否则交易无法打包上链

3.配置文件中ip和端口确保正确，否则会抛出一下异常
```
panic: runtime error: invalid memory address or nil pointer dereference
```

4.发布合约：必须指定合约abi和wasm文件的全路径

5.合约调用：
```
  the contract address is not exist ...
```
  合约没有发布成功，获取不到合约的code，检查合约是否发布成功，命令中合约地址参数是否正确






