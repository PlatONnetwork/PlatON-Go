package wasm.contract_object_oriented;

import com.platon.rlp.datatypes.Uint32;
import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.OverrideContract;
import org.junit.Test;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * @author zjsunzone
 *
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
            OverrideContract contract = OverrideContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("OverrideContract issued successfully.contractAddress:"
                    + contractAddress + ", hash:" + transactionHash
                    + " gasUsed:" + contract.getTransactionReceipt().get().getGasUsed().toString());

            // 1. input = 1, get: 100, input == 2, get: 10000
            Uint32 area01 = contract.getArea(Uint64.of(1)).send();
            collector.logStepPass("To invoke getArea success, area: " + area01.getValue().toString());
            collector.assertEqual(area01.getValue(), BigInteger.valueOf(100));

            Uint32 area02 = contract.getArea(Uint64.of(2)).send();
            collector.logStepPass("To invoke getArea success, area2: " + area02.getValue().toString());
            collector.assertEqual(area02.getValue(), BigInteger.valueOf(10000));

        } catch (Exception e) {
            collector.logStepFail("OverrideContract failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
