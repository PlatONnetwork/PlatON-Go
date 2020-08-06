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

import java.util.List;

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
            LATToken contract = LATToken.deploy(web3j, transactionManager, provider,chainId,
                    Uint64.of("100000000000000"),
                    "LTT Token",
                    Uint8.of(6),
                    "LTT").send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_LATToken deploy successfully. contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_LATToken deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            WasmAddress sender = contract.getSender().send();
            collector.logStepPass("sender address is: " + sender.toString());

            // the info of contract
            String name = contract.getName().send();
            collector.logStepPass("Call getName of LATToken, name: " + name);
            String symbol = contract.getSymbol().send();
            collector.logStepPass("Call getSymbol of LATToken, symbol: " + symbol);
            Uint64 totalSUpply = contract.getTotalSupply().send();
            collector.logStepPass("Call totalSUpply of LATToken, totalSUpply: " + totalSUpply.getValue());
            Uint8 decimal = contract.getDecimals().send();
            collector.logStepPass("Call getDecimals of LATToken, decimal: " + decimal.getValue());

            Uint64 balance = contract.balanceOf(sender).send();
            collector.logStepPass("Call balanceOf of LATToken, balance: " + balance.getValue());

            // transfer
            WasmAddress receiver = new WasmAddress("lax19xdjrg06xz9te85c839zqtelmaj2tgt047fh5q");
            balance = contract.balanceOf(receiver).send();
            collector.logStepPass("Call balanceOf of LATToken, token before balance: " + balance.getValue());
            TransactionReceipt trasferTR = contract.transfer(receiver, Uint64.of(100000000)).send();
            collector.logStepPass("Send trasnsfer, hash: " + trasferTR.getTransactionHash() + " gasUsed: " + trasferTR.getGasUsed());

            // parse logs
            List<LATToken.TransferEventResponse> responses = contract.getTransferEvents(trasferTR);

            collector.logStepPass("Send transfer, logs: " + trasferTR.getLogs().size()
                    + " from: " + responses.get(0).topic1 + " to: " + responses.get(0).topic2 + " value: "
                    + responses.get(0).arg1);



            balance = contract.balanceOf(receiver).send();
            collector.logStepPass("Call balanceOf of LATToken, token after balance: " + balance.getValue());
            collector.assertEqual(balance.getValue().longValue(), Long.valueOf(100000000));

        } catch (Exception e) {
            collector.logStepFail("contract_LATToken failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
