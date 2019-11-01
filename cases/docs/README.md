# PlatON-Tests

### 测试执行

1.依赖安装 

```js
pip install -r requirements.txt
```

2.执行测试

暂时沿用pytest命令执行，新增参数参考[conftest.py](../conftest.py)。

example: 

```js
py.test 'tests/chain/test_chain_deploy.py' --nodeFile "deploy/node/test_chaininfo.yml" --accountFile "deploy/accounts.yml"
```

更多请查看[readme](../README.md)

备注：
项目管理，执行测试命令存在较多组合方式，需要输入大量的命令，考虑引入Makefile。

### 项目结构说明

1.根目录：

+ pytest.ini pytest测试配置文件，具体配置项请参考官方文件[pytest.ini](https://docs.pytest.org/en/latest/reference.html#configuration-options)。

+ conftest.py-pytest插件，主要包括命令行参数的引入，测试环境生成，测试失败用例日志下载等功能。

2.common：

+ 项目公共方法，包括日志，配置文件读取，连接linux服务器，建立web3连接，二进制包下载等。

3.conf：

+ 部署环境时本地文件路径配置，详情请查看文件备注

4.deploy：

+ 部署文件存放路径，配合conf使用

5.environment：

+ 部署环境脚本，由TestConfig, TestEnvironment, Account, Server, Node 五个类组成。
+ TestConfig 部署环境配置，测试环境全局配置。
+ Account 测试账号管理类，用于测试过程中，账号生成，转账等。
+ Server 服务器管理类，用于部署过程中，管理服务器依赖，压缩文件管理。
+ TestEnvironment 测试环境类，用于管理测试环境，启动节点，关闭节点，重启节点，获取节点，获取账号等。
+ Node 节点类，用于管理节点，具备节点基本信息，节点连接，节点运行状态，节点启停等。

6.tests：

+ 用例存放目录

7.tests/lib:

+ utils：用例公共方法
+ 配置：StakingConfig，质押基本配置；PipConfig，pip基础配置；DefaultEconomicConfig经济模型基本配置
+ genesi：Genesis genesis.json文件可转化对象
+ 其余类均是对sdk的二次封装


### 用例示例

[质押示例](../tests/example/test_staking.py)

[定制节点部署](../tests/example/test_customize_deploy.py)

### 特性说明

+ 抛弃旧框架中的部署形式，引入测试环境对象和节点对象，用于管理测试环境和节点，在节点中包含了大量节点属性，免去繁多的节点信息获取调用

+ w3连接建立引入超时机制，避免节点刚启动无法连接导致报错的问题

+ 去掉setup等，引入fixture，设计用例时，公共步骤需要设计为步骤方法，去掉重复步骤代码

+ 部署优化，牺牲环境准备时间，二次部署速度有极大提升