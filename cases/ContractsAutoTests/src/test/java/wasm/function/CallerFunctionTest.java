package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.CallerFunction;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;


/**
 *
 * @title 验证函数platon_caller
 * @description:
 * @author: liweic
 * @create: 2020/02/11
 */
public class CallerFunctionTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.CallerFunctionTest验证链上函数platon_caller",sourcePrefix = "wasm")
    public void callerFunction() {

        try {
            prepare();
            CallerFunction caller = CallerFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = caller.getContractAddress();
            String transactionHash = caller.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CallerFunctionTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("CallerFunctionTest deploy gasUsed:" + caller.getTransactionReceipt().get().getGasUsed());

            String callerAddr = caller.get_platon_caller().send();
            collector.logStepPass("getPlatONCaller函数返回值:" + callerAddr);
            collector.assertEqual(callerAddr, walletAddress.toLowerCase());


        } catch (Exception e) {
            collector.logStepFail("CallerFunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}


