package wasm.function;

import com.platon.rlp.datatypes.Uint32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Sha3Function;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;


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
            author = "liweic", showName = "wasm.Sha3FunctionTest验证链上函数platon_sha3",sourcePrefix = "wasm")
    public void Sha3function() {

        try {
            prepare();
            Sha3Function shafunction = Sha3Function.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = shafunction.getContractAddress();
            String transactionHash = shafunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("CallerFunctionTest issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("CallerFunctionTest deploy gasUsed:" + shafunction.getTransactionReceipt().get().getGasUsed());

            Uint32 sha3 = shafunction.Sha3Result().send();
            collector.logStepPass("sha3函数返回值:" + sha3.value);
            collector.assertEqual(sha3.value, new BigInteger("114259850"));


        } catch (Exception e) {
            collector.logStepFail("Sha3FunctionTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}



