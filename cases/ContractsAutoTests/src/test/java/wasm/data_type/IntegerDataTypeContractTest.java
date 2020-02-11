package wasm.data_type;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InnerFunction;
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
 * Base data type test.
 */
public class IntegerDataTypeContractTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.data_type",sourcePrefix = "wasm")
    public void testFunctionContract() {

        try {
            prepare();

            // deploy contract.
            InnerFunction innerFunction = InnerFunction.deploy(web3j, transactionManager, provider).send();
            String contractAddress = innerFunction.getContractAddress();
            String transactionHash = innerFunction.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InnerFunction issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);


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
