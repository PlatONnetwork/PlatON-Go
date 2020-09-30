# 使用手冊

​		合约自动化测试代码是一个基于`client.java.sdk`并通过`maven`来管理的工程。因此使用时只需要配置`maven`，并通过`maven client compile` 就可以生成一个`java project`并导入`IDE`，然后就可以开始愉快的编写合约测试

## 1 合约编写

### 1.1 合约源码存放路径

`src/test/resources/contracts`，目前还未集成编译器，所以暂时无法把合约源码编译成二进制和`abi`文件，计划12.20前支持编译。

### 1.2 合约二进制和`abi`文件存放路径

`src/test/resources/contracts/build`

### 1.3 包装类生成

* EVM:将二进制和`abi`文件放在`src/test/resources/contracts/evm/build`路径下，再执行`src/test/java/evm/beforetest`路径下的`GeneratorPreTest`类的`junit`方法（右键直接使用`junit`插件即可执行）
* WASM:将二进制和`abi`文件放在`src/test/resources/contracts/wasm/build`路径下，再执行`src/test/java/wasm/beforetest`路径下的`WASMGeneratorPreTest`类的`junit`方法（右键直接使用`junit`插件即可执行）

## 2 合约测试脚本编写

### 2.1 示例脚本

```java
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.junit.rules.AssertCollector;
import network.platon.autotest.junit.rules.DriverService;
import network.platon.contracts.evm.HumanStandardToken;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;
import java.math.BigInteger;


/**
 * @title 代币转移
 * @description:
 * @author: qcxiao
 * @create: 2019/12/16 13:39
 **/
public class TokenTransferTest {
    @Rule
    public AssertCollector collector = new AssertCollector();

    @Rule
    public DriverService driverService = new DriverService();
    // 底层链ID
    private long chainId;
    // 每次转移的代币数量
    private String transferAmount;
    // 发行代币的总额
    private String ownerAmount;
    // 发行代币的地址
    private final static String transferFrom = "0x03f0e0a226f081a5daecfda222cafc959ed7b800";
    // 接收代币的地址
    private final static String transferTo = "0x8d2b8b62d2ff5e7d17f91cf821cafee8e1fe4584";
    // 代币名称
    private String tokenName;

    @Before
    public void before() {
        chainId = Integer.valueOf(driverService.param.get("chainId"));
        ownerAmount = driverService.param.get("ownerAmount");
        transferAmount = driverService.param.get("transferAmount");
        tokenName = driverService.param.get("tokenName");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "complexcontracts.TokenTransferTest-代币转移")
    public void testTransfer() {
        Web3j web3j = null;
        Credentials credentials = null;
        try {
            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
            credentials = Credentials.create(driverService.param.get("privateKey"));
            collector.logStepPass("currentBlockNumber:" + web3j.platonBlockNumber().send().getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("The node is unable to connect", e.toString());
            e.printStackTrace();
        }

        ContractGasProvider provider = new ContractGasProvider(new BigInteger("50000000000"), new BigInteger("3000000"));
        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);

        try {
            HumanStandardToken token = HumanStandardToken.deploy(web3j, transactionManager, provider,
                    new BigInteger(ownerAmount), tokenName, BigInteger.valueOf(18), "USDT").send();
            String contractAddress = token.getContractAddress();
            TransactionReceipt tx = token.getTransactionReceipt().get();
            collector.logStepPass("Token issued successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash() +
                    ", tokenName:" + token.name().send() + ", symbol:" + token.symbol().send());
            collector.assertEqual(tokenName, token.name().send(), "checkout tokenName");
            collector.logStepPass("5次循环调用...");
            for (int i = 1; i < 6; i++) {
                TransactionReceipt transactionReceipt = HumanStandardToken.load(contractAddress, web3j, transactionManager, provider, chainId)
                        .transfer(transferTo, new BigInteger(transferAmount)).send();
                BigInteger toBalance = token.balanceOf(transferTo).send();
                BigInteger fromBalance = token.balanceOf(transferFrom).send();
                collector.logStepPass("Token transfer successful.transactionHash:" + transactionReceipt.getTransactionHash());
                collector.logStepPass("currentTransferTime:" + i + ", currentBlockNumber:" + transactionReceipt.getBlockNumber());
                collector.logStepPass("transferToBalance:" + toBalance + ", transferFromBalance:" + fromBalance);
                // 累计转移的数量
                BigInteger amount = new BigInteger(transferAmount).multiply(BigInteger.valueOf(i));
                // 判断代币接收地址的余额是否正确
                collector.assertEqual(amount, toBalance, "checkout every time transferTo balance.");
                // 判断代币转出地址余额是否正确
                collector.assertEqual((new BigInteger(ownerAmount)).subtract(amount), fromBalance, "checkout every time transferFrom balance.");
            }
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
```

### 2.2 类名称说明

类名称即测试套件名称，必须以`Test`单词结尾

### 2.3 数据驱动

#### 2.3.1 类成员变量声明

```java
@Rule
public AssertCollector collector = new AssertCollector();
@Rule
public DriverService driverService = new DriverService();
```

#### 2.3.2 指定对应的测试数据

```java
@DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "complexcontracts.TokenTransferTest-代币转移")
```

`type = DataSourceType.EXCEL`：表示数据通过`Excel`表管理

`file = "test.xls"`：表示文件名称

`sheetName = "Sheet1"`：指定`EXCEL`表的`Sheet`

`author = "qcxiao"`：表示作者

`showName = "complexcontracts.TokenTransferTest-代币转移"`：日志中的名称显示

#### 2.3.3 测试数据文件

文件路径：`src/test/resources/` + 对应测试类(测试脚本)的名称，如`src/test/resources/complexcontracts.TokenTransferTest`

文件名称：`src/test/resources/complexcontracts.TokenTransferTest`下的文件名称与`SheetName`需要和2.3.2中的文件名称一致

示例：

| `caseName`      | `caseDescription` | `caseRun` | `casePriority` | `ownerAmount`                | `transferAmount` | `tokenName` |
| --------------- | ----------------- | --------- | -------------- | ---------------------------- | ---------------- | ----------- |
| `ERC20代币转移` | 普通转账          | Y         | `P1`           | 1000000000000000000000000000 | 1000             | `qcxiao`    |

不同的测试点，数据可以有多行，执行的时候框架会自动执行到每行数据。

#### 2.3.4 测试数据读取

`driverService.param.get("ownerAmount")`

### 2.4 日志

#### 2.4.1 成功日志

```java
collector.logStepPass("currentBlockNumber:" + web3j.platonBlockNumber().send().getBlockNumber());
```

#### 2.4.2 失败日志

```java
collector.logStepFail("The node is unable to connect", e.toString());
```

### 2.5 断言

```java
// 判断代币转出地址余额是否正确
collector.assertEqual((new BigInteger(ownerAmount)).subtract(amount), fromBalance, "checkout every time transferFrom balance.");
```

## 3 执行测试

### 3.1 `Junit`插件方式

测试类中使用右键`junit`方式直接测试

### 3.2 `maven`方式

```java
mvn clean test
```

## 4 报告生成

#### 4.1 报告路径配置

`src/test/resources/test.properties`文件中的`logDir = C\:\\autotest_log\\`字段用以配置

#### 4.2 报告示例

![image-20191217094605262](https://github.com/qcblockchain/PlatON-Go/blob/feature/wasm/cases/ContractsAutoTests/src/main/resources/templates/images/image-20191217094605262.png)

#### 4.3 编码说明

使用`mvn clean test`测试时，因为`windows`默认采用`GBK`编码，所以`src/main/resources/templates`里面的文件编码需要调整成`GBK`，否则报告会乱码



## 5 补充说明

目前测试用例版本将包含EVM和WASM两个版本，脚本数和用例数已经变得庞大，为了规范两种类别的自动化测试用例，接下来将区分测试代码和测试数据：

### 5.1 测试方法改造

测试方法的注解需要带`sourcePrefix = "evm"`或者`sourcePrefix = "wasm"`

完整信息：`@DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",        author = "qcxiao", showName = "complexcontracts.TokenTransferTest-代币转移", sourcePrefix = "evm")`

### 5.2 测试数据转移

测试数据均需要放到`evm`或者`wasm`文件夹下面

### 5.2 测试优先级

因为当前EVM测试需要半小时，如何全部执行时只执行WASM的用例，需要在Excel表中的casePriority字段标识成WASM，执行测试时通过命令：`mvn test -DcasePriority=WASM`，此时将只执行WASM的用例

### 5.3 测试代码结构图

![结构](https://github.com/qcblockchain/PlatON-Go/blob/feature/wasm/cases/ContractsAutoTests/src/main/resources/templates/images/%E7%BB%93%E6%9E%84.png)
