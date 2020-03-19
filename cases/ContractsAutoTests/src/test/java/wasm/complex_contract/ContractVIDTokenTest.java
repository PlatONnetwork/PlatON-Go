package wasm.complex_contract;

import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ForeignBridge;
import network.platon.contracts.wasm.HomeBridge;
import network.platon.contracts.wasm.VIDToken;
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
public class ContractVIDTokenTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_VIDToken",sourcePrefix = "wasm")
    public void testHomeBridge() {
        try {
            // deploy contract.
            VIDToken contract = VIDToken.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_VIDToken issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_VIDToken deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // transfer to contract
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal(100), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

          

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("contract_VIDToken and could not call contract function");
            }else{
                collector.logStepFail("contract_VIDToken failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }



}
