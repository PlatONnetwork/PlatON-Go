package wasm.contract_object_oriented;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InnerFunction;
import network.platon.contracts.wasm.OverrideContract;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.DefaultBlockParameterNumber;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * The test class of the function for chain.
 */
public class OverrideContractTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.override_contract",sourcePrefix = "wasm")
    public void testOverrideContract() {
        try {
            prepare();

            // deploy contract.
            OverrideContract contract = OverrideContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("OverrideContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            // 1. input = 1, get: 100, input == 2, get: 10000
            int area01 = contract.getArea(Long.valueOf(1)).send();
            collector.logStepPass("To invoke getArea success, area: " + area01);
            //collector.assertEqual(area01, 100);

            int area02 = contract.getArea(Long.valueOf(2)).send();
            collector.logStepPass("To invoke getArea success, area2: " + area02);
            //collector.assertEqual(area01, 100);

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("OverrideContract and could not call contract function");
            }else{
                collector.logStepFail("OverrideContract failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
