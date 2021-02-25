package network.platon.test.wasm.contract_docs;

import com.platon.rlp.datatypes.Uint128;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.test.datatypes.Xuint128;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.TweetRegistry;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * @author zjsunzone
 *
 * This class exists for docs.
 */
public class ContractTweetRegistryTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_TweetAccount",sourcePrefix = "wasm")
    public void testTweetAccount() {

        try {
            // deploy contract.
            TweetRegistry contract = TweetRegistry.deploy(web3j, transactionManager, provider,chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("TweetRegistry deploy successfully. contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("TweetRegistry deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // registry
            String addr = credentials.getAddress(chainId);
            addr = "lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2";
            TransactionReceipt registrTr = contract.registry("bob", new WasmAddress(addr)).send();
            collector.logStepPass("Send registry, hash: " + registrTr.getTransactionHash() + " gasUsed: " + registrTr.getGasUsed());

            // getNumberOfAccounts
            Uint128 numberOfAccounts = contract.getNumberOfAccounts().send();
            collector.logStepPass("Call getNumberOfAccounts, result: " + numberOfAccounts);

            // getAddressOfName
            WasmAddress addressOfName = contract.getAddressOfName("bob").send();
            collector.logStepPass("Call getAddressOfName, result: " + addressOfName);
            collector.assertEqual(addressOfName.getAddress(), addr);

            // getNameOfAddress
            String nameOfAddress = contract.getNameOfAddress(new WasmAddress(addr)).send();
            collector.logStepPass("Call getNameOfAddress, result: " + nameOfAddress);
            collector.assertEqual(nameOfAddress, "bob");

            // getAddressOfId
            WasmAddress addressOfId = contract.getAddressOfId(Xuint128.ZERO).send();
            collector.logStepPass("Call getAddressOfId, result: " + addressOfId.getAddress());
            collector.assertEqual(addressOfId.getAddress(), addr);

            // unregister
            /*TransactionReceipt unregistterTr = contract.unregister().send();
            collector.logStepPass("Send unregister, txHash: " + unregistterTr.getTransactionHash() +
                    " gasUsed: " + unregistterTr.getGasUsed());
            nameOfAddress = contract.getNameOfAddress(new WasmAddress(addr)).send();
            collector.logStepPass("Call getNameOfAddress after unregister, result: " + nameOfAddress);
            collector.assertEqual(nameOfAddress, "");*/

            // adminUnregister
            TransactionReceipt adminUnregisterTr = contract.adminUnregister("bob").send();
            collector.logStepPass("Send adminUnregister, txHash: " + adminUnregisterTr.getTransactionHash() +
                    " gasUsed: " + adminUnregisterTr.getGasUsed());
            nameOfAddress = contract.getNameOfAddress(new WasmAddress(addr)).send();
            collector.logStepPass("Call getNameOfAddress after unregister, result: " + nameOfAddress);
            collector.assertEqual(nameOfAddress, "");

            // adminSetRegistrationDisable
            boolean disabled = contract.getRegistrationDisabled().send();
            collector.logStepPass("Call getRegistrationDisabled, before result " + disabled);
            TransactionReceipt disableTr = contract.adminSetRegistrationDisable(false).send();
            collector.logStepPass("Send adminDeleteRegistry, hash: " + disableTr.getTransactionHash() +
                    "gasUsed: " + disableTr.getGasUsed());
            disabled = contract.getRegistrationDisabled().send();
            collector.logStepPass("Call getRegistrationDisabled, after result " + disabled);

            // transfer to contract
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal(1), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , before balance: " + cbalance);

            // adminRetrieveDonations
            TransactionReceipt retrieveTr = contract.adminRetrieveDonations().send();
            collector.logStepPass("Send adminRetrieveDonations, hash: " + retrieveTr.getTransactionHash() +
                    "gasUsed: " + retrieveTr.getGasUsed());

            cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , after balance: " + cbalance);
            collector.assertEqual(cbalance, BigInteger.ZERO);

            //
            // transfer to contract
            t.sendFunds(contractAddress, new BigDecimal(1), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("02-Transfer to contract , before balance: " + cbalance);

            // adminDeleteRegistry
            TransactionReceipt deleteTR = contract.adminDeleteRegistry().send();
            collector.logStepPass("Send adminDeleteRegistry, hash: " + deleteTR.getTransactionHash() +
                    "gasUsed: " + deleteTR.getGasUsed());

            cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("02-Transfer to contract , after balance: " + cbalance);
            collector.assertEqual(cbalance, BigInteger.ZERO);

        } catch (Exception e) {
            collector.logStepFail("TweetAccount failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }


}
