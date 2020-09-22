package network.platon.test.evm.versioncompatible.v0_4_25;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.DeprecatedFunctions;
import org.junit.Test;
import org.web3j.tuples.generated.Tuple2;
import org.web3j.tx.exceptions.ContractCallException;

/**
 * @title 0.5.0版本弃用但0.4.25版本仍生效函数测试
 * 1.callcode()（0.5.0版本已弃用，使用delegatecall()函数代替） 验证
 * 2.suicide() （0.5.0版本已弃用，使用selfdestruct()函数替用）验证
 * 3.sha3() （0.5.0版本已弃用，使用keccak256()函数代替）验证
 * 4.throw （0.5.0版本已弃用，使用异常函数验证）验证
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class DeprecatedFunctionsTest extends ContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testFunctionCheck",
            author = "albedo", showName = "network.platon.test.evm.DeprecatedFunctionsTest-0.4.25版本仍生效函数", sourcePrefix = "evm")
    public void testFunctionCheck() {
        try {
            prepare();
            DeprecatedFunctions deprecatedFunctions = DeprecatedFunctions.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = deprecatedFunctions.getContractAddress();
            String transactionHash = deprecatedFunctions.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("DeprecatedFunctions issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + deprecatedFunctions.getTransactionReceipt().get().getGasUsed());
            Tuple2<Boolean, byte[]> result = deprecatedFunctions.functionCheck().send();
            collector.assertEqual(result.getValue1(), Boolean.TRUE, "checkout deprecated function cast result");

        } catch (Exception e) {
            collector.logStepFail("DeprecatedFunctionsTest testFunctionCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testThrowCheck",
            author = "albedo", showName = "DeprecatedFunctionsTest-验证throw关键字", sourcePrefix = "evm")
    public void testThrowCheck() {
        try {
            prepare();
            DeprecatedFunctions deprecatedFunctions = DeprecatedFunctions.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = deprecatedFunctions.getContractAddress();
            String transactionHash = deprecatedFunctions.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("DeprecatedFunctions issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + deprecatedFunctions.getTransactionReceipt().get().getGasUsed());
            try {
                deprecatedFunctions.throwCheck(false).send();
            } catch (ContractCallException e) {
                collector.assertEqual(e.getMessage(), "Empty value (0x) returned from contract", "checkout throw result");
            }
            Boolean result = deprecatedFunctions.throwCheck(true).send();
            collector.assertEqual(result, Boolean.TRUE, "checkout deprecated function result");
        } catch (Exception e) {
            collector.logStepFail("DeprecatedFunctionsTest testThrowCheck failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "testSuicide",
            author = "albedo", showName = "DeprecatedFunctionsTest-验证suicide函数", sourcePrefix = "evm")
    public void testSuicide() {
        try {
            prepare();
            DeprecatedFunctions deprecatedFunctions = DeprecatedFunctions.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = deprecatedFunctions.getContractAddress();
            String transactionHash = deprecatedFunctions.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("DeprecatedFunctions issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + deprecatedFunctions.getTransactionReceipt().get().getGasUsed());
            deprecatedFunctions.kill().send();
            try {
                deprecatedFunctions.functionCheck().send();
            } catch (Exception e) {
                collector.logStepPass("checkout suicide function success ");
            }
        } catch (Exception e) {
            collector.logStepFail("DeprecatedFunctionsTest testSuicide failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
