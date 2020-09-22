package network.platon.test.evm.exceptionhandle;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.RequireHandle;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.exceptions.TransactionException;

import java.math.BigInteger;

/**
 * @title require(bool condition)函数测试
 * 如果条件不满足则撤销状态更改，用于检查由输入或者外部组件引起的错误
 *（1）消息调用函数过程中没有正确结束异常验证
 *（2）new创建合约没有正确返回产生异常验证
 *（3）调用外部函数，被调用的对象不包含代码异常验证
 *（4）合约没有payable修饰符的public的函数在接收主币时（包括构造函数，和回退函数）异常验证
 *（5）合约通过一个public的getter函数（public getter funciton）接收主币异常验证
 *（6）.transfer()函数执行失败异常验证
 * @description:
 * @author: albedo
 * @create: 2019/12/31
 */
public class RequireHandleTest extends ContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "functionCallException",
            author = "albedo", showName = "network.platon.test.evm.exceptionhandle.RequireHandle-消息调用函数过程中没有正确结束异常", sourcePrefix = "evm")
    public void testFunctionCallException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            try {
                handle.functionCallException(new BigInteger("1000")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail(" RequireHandleTest testFunctionCallException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "newContractException",
            author = "albedo", showName = "exceptionhandle.RequireHandle-new创建合约没有正确返回", sourcePrefix = "evm")
    public void testNewContractException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + handle.getTransactionReceipt().get().getGasUsed());
            try {
                handle.newContractException().send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireHandleTest testNewContractException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "outFunctionCallException",
            author = "albedo", showName = "exceptionhandle.RequireHandle-调用外部函数，被调用的对象不包含代码", sourcePrefix = "evm")
    public void testOutFunctionCallException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            try {
                handle.outFunctionCallException(new BigInteger("1000")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireHandleTest testOutFunctionCallException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "nonPayableReceiveEthException",
            author = "albedo", showName = "exceptionhandle.RequireHandle-合约在没有payable修饰符的public的函数中接收主币", sourcePrefix = "evm")
    public void testNonPayableReceiveEthException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            try {
                handle.nonPayableReceiveEthException(new BigInteger("1000")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireHandleTest testNonPayableReceiveEthException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "publicGetterReceiveEthException",
            author = "albedo", showName = "exceptionhandle.RequireHandle-合约通过一个public的getter函数接收主币", sourcePrefix = "evm")
    public void testPublicGetterReceiveEthException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            try {
                handle.publicGetterReceiveEthException(new BigInteger("1000")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireHandleTest testPublicGetterReceiveEthException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "transferCallException",
            author = "albedo", showName = "exceptionhandle.RequireHandle-transfer()函数执行失败", sourcePrefix = "evm")
    public void testTransferCallException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            try {
                handle.transferCallException(new BigInteger("100000000000000000000000000000000000000000"),new BigInteger("100000000000")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout assert throw exception:" + e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireHandleTest testTransferCallException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "paramException",
            author = "albedo", showName = "exceptionhandle.RequireHandle-检测异常", sourcePrefix = "evm")
    public void testParamException() {
        try {
            prepare();
            RequireHandle handle = RequireHandle.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = handle.getContractAddress();
            String transactionHash = handle.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("RequireHandle issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            TransactionReceipt receipt = handle.paramException(new BigInteger("5")).send();
            collector.logStepPass("checkout require normal,transactionHash="+receipt.getTransactionHash());
            try {
                handle.paramException(new BigInteger("11")).send();
            } catch (TransactionException e) {
                collector.logStepPass("checkout require throw exception:"+e.getMessage());
            }
        } catch (Exception e) {
            collector.logStepFail("RequireHandleTest testParamException failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
