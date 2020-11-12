package network.platon.test.wasm.function;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.CallValueFunction;
import org.junit.Test;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;


/**
 *
 * @title 验证函数platon_call_value
 * @description:
 * @author: liweic
 * @create: 2020/02/11
 */
public class CallValueFunctionTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.CallValueFunctionTest验证链上函数platon_call_value",sourcePrefix = "wasm")
    public void CallValuefunction() {

        try {
            prepare();
            CallValueFunction callvalue = CallValueFunction.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = callvalue.getContractAddress();
            String transactionHash = callvalue.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CallerFunctionTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("CallerFunctionTest deploy gasUsed:" + callvalue.getTransactionReceipt().get().getGasUsed());

            Uint8 callval = callvalue.get_platon_call_value().send();
            collector.logStepPass("getPlatONCallValue函数返回值:" + callval.value);
            collector.assertEqual(callval.value.toString(), "0");


        } catch (Exception e) {
            collector.logStepFail("CallValueFunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}



