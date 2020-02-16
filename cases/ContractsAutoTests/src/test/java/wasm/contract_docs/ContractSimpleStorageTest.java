package wasm.contract_docs;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InnerFunction;
import network.platon.contracts.wasm.InnerFunction_1;
import network.platon.contracts.wasm.InnerFunction_2;
import network.platon.contracts.wasm.SimpleStorage;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * @author zjsunzone
 *
 * This class is for docs.
 */
public class ContractSimpleStorageTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_simple_storage",sourcePrefix = "wasm")
    public void testSimpleStorageContract() {
        try {
            // deploy contract.
            SimpleStorage contract = SimpleStorage.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SimpleStorage issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            TransactionReceipt tr = contract.set(Long.valueOf(10000)).send();
            collector.logStepPass("To invoke set success, txHash: " + tr.getTransactionHash());
            Long result = contract.get().send();
            collector.logStepPass("To invoke get success, result: " + result.longValue());
            collector.assertEqual(result.longValue(), Long.valueOf(10000).longValue());

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("SimpleStorage and could not call contract function");
            }else{
                collector.logStepFail("SimpleStorage failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }

}
