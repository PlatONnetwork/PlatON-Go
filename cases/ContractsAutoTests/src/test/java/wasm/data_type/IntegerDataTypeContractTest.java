package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InnerFunction;
import network.platon.contracts.wasm.IntegerDataTypeContract;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.DefaultBlockParameterNumber;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Bytes;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * Base data type test.
 */
public class IntegerDataTypeContractTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.integer_data_type",sourcePrefix = "wasm")
    public void testIntegerTypeContract() {

        try {
            prepare();

            // deploy contract.
            IntegerDataTypeContract contract = IntegerDataTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("IntegerDataTypeContract issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            // test: int8
            int int8 = contract.int8().send();
            collector.logStepPass("To invoke int8 success, int8: " + int8);

            // test: int32
            int int32 = contract.int32().send();
            collector.logStepPass("To invoke int32 success, int8: " + int32);

            // test: int64
            int int64 = contract.int64().send();
            collector.logStepPass("To invoke int64 success, int64: " + int8);

            // test: uint8
            byte uint8 = contract.uint8t(Byte.valueOf((byte)1)).send();
            collector.logStepPass("To invoke uint8t success, uint8: " + uint8);


        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("IntegerDataTypeContract and could not call contract function");
            }else{
                collector.logStepFail("IntegerDataTypeContract failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
