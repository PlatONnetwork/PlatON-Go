package wasm.function;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Sha3Function;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;


/**
 *
 * @title 验证函数platon_sha3
 * @description:
 * @author: liweic
 * @create: 2020/02/11
 */
public class Sha3FunctionTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.SpecialFunctionsA验证链上函数platon_sha3",sourcePrefix = "wasm")
    public void Sha3function() {

        try {
            prepare();
            Sha3Function shafunction = Sha3Function.deploy(web3j, transactionManager, provider).send();
            String contractAddress = shafunction.getContractAddress();
            String transactionHash = shafunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CallerFunctionTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            Integer sha3 = shafunction.Sha3Result().send();
            collector.logStepPass("sha3函数返回值:" + sha3);
            collector.assertEqual(sha3, new Integer("114259850"));


        } catch (Exception e) {
            collector.logStepFail("Sha3FunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}



