package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.IntegerDataTypeContract;
import network.platon.contracts.wasm.ReferenceDataTypeContract;
import org.junit.Test;
import org.web3j.abi.datatypes.Address;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class is used to test data type of reference.
 */
public class ReferenceDataTypeContractTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.reference_data_type",sourcePrefix = "wasm")
    public void testIntegerTypeContract() {

        try {
            prepare();

            // deploy contract.
            ReferenceDataTypeContract contract = ReferenceDataTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ReferenceDataTypeContract deploy successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);



        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ReferenceDataTypeContract and could not call contract function");
            }else{
                collector.logStepFail("ReferenceDataTypeContract failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
