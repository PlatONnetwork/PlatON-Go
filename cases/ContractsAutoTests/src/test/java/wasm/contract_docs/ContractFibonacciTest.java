package wasm.contract_docs;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Fibonacci;
import network.platon.contracts.wasm.InnerFunction;
import network.platon.contracts.wasm.InnerFunction_1;
import network.platon.contracts.wasm.InnerFunction_2;
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
 * This class exists for docs.
 */
public class ContractFibonacciTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_fibonacci",sourcePrefix = "wasm")
    public void testFibonacciContract() {

        try {
            // deploy contract.
            Fibonacci contract = Fibonacci.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("Fibonacci deploy successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            Long result = contract.fibonacci(Long.valueOf(10)).send();
            collector.logStepPass("To invoke fibonacci success, result: " + result);



        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("Fibonacci and could not call contract function");
            }else{
                collector.logStepFail("Fibonacci failure,exception msg:" , e.getMessage());
            }
            e.printStackTrace();
        }
    }


}
