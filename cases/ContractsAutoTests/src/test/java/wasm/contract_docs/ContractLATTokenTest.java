package wasm.contract_docs;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.LATToken;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @author zjsunzone
 *
 * This class exists for docs.
 */
public class ContractLATTokenTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_LATToken",sourcePrefix = "wasm")
    public void testToken() {

        try {
            // deploy contract.
            // TransactionManager transactionManager, GasProvider contractGasProvider,
            // Uint64 _initialAmount, String _tokenName, Uint8 _decimalUnits, String _tokenSymbol
            LATToken contract = LATToken.deploy(web3j, transactionManager, provider,
                    Uint64.of("100000000000000"),
                    "LTT Token",
                    Uint8.of(6),
                    "LTT").send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_LATToken deploy successfully. contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_LATToken deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // the info of contract
            String name = contract.getName().send();
            collector.logStepPass("Call getName of LATToken, name: " + name);
            String symbol = contract.getSymbol().send();
            collector.logStepPass("Call getSymbol of LATToken, symbol: " + symbol);
            Uint64 totalSUpply = contract.getTotalSupply().send();
            collector.logStepPass("Call totalSUpply of LATToken, totalSUpply: " + totalSUpply.getValue());
            Uint8 decimal = contract.getDecimals().send();
            collector.logStepPass("Call getDecimals of LATToken, decimal: " + decimal.getValue());

            Uint64 balance = contract.balanceOf(new WasmAddress(credentials.getAddress())).send();
            collector.logStepPass("Call balanceOf of LATToken, balance: " + balance.getValue());

            // transfer
            WasmAddress receiver = new WasmAddress("0x299b21a1fa308abc9e983c4a202f3fdf64a5a16f");
            balance = contract.balanceOf(receiver).send();
            collector.logStepPass("Call balanceOf of LATToken, token before balance: " + balance.getValue());
            TransactionReceipt trasferTR = contract.transfer(receiver, Uint64.of(100000000)).send();
            collector.logStepPass("Send trasnsfer, hash: " + trasferTR.getTransactionHash() + " gasUsed: " + trasferTR.getGasUsed());
            collector.logStepPass("Send transfer, logs: " + trasferTR.getLogs().size());

            balance = contract.balanceOf(receiver).send();
            collector.logStepPass("Call balanceOf of LATToken, token after balance: " + balance.getValue());

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("contract_LATToken and could not call contract function");
            }else{
                collector.logStepFail("contract_LATToken failure,exception msg:" , e.getMessage());
            }
            e.printStackTrace();
        }
    }


}
