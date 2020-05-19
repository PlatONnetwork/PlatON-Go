package wasm.complex_contract;

import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ForeignBridge;
import network.platon.contracts.wasm.HomeBridge;
import network.platon.contracts.wasm.VIDToken;
import org.junit.Before;
import org.junit.Test;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.utils.Numeric;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * @author zjsunzone
 * <p>
 * This class is for docs.
 */
public class ContractVIDTokenTest extends WASMContractPrepareTest {

    @Before
    public void before() {
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "zjsunzone", showName = "wasm.contract_VIDToken", sourcePrefix = "wasm")
    public void testHomeBridge() {
        try {
            // deploy contract.
            VIDToken contract = VIDToken.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contract.getContractAddress();
            String transactionHash = contract.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("contract_VIDToken issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("contract_VIDToken deploy successfully. gasUsed: " + contract.getTransactionReceipt().get().getGasUsed().toString());

            // transfer to contract
            Transfer t = new Transfer(web3j, transactionManager);
            t.sendFunds(contractAddress, new BigDecimal(100), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            BigInteger cbalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("Transfer to contract , address: " + contractAddress + " cbalance: " + cbalance);

            // transfer in contract
            String to = "lax1fyeszufxwxk62p46djncj86rd553skpptsj8v6";
            BigInteger value = new BigInteger("100000");
            TransactionReceipt transferTr = contract.Transfer(to, value).send();
            collector.logStepPass("Send Transfer, hash:  " + transferTr.getTransactionHash()
                    + " gasUsed: " + transferTr.getGasUsed());

            // balance of
            BigInteger balance = contract.BalanceOf(to).send();
            collector.logStepPass("Call balanceOf, res: " + balance);
            collector.assertEqual(balance, value);

            // approve
            TransactionReceipt approveTR = contract.Approve(to, value).send();
            collector.logStepPass("Send Approve, hash:  " + approveTR.getTransactionHash()
                    + " gasUsed: " + approveTR.getGasUsed());

            // allowance
            BigInteger allowance = contract.Allowance(credentials.getAddress(chainId), to).send();
            collector.logStepPass("Call allowance, res: " + allowance);
            collector.assertEqual(allowance, value);

            // IncreaseApproval
            BigInteger increaseValue = new BigInteger("23300000");
            TransactionReceipt increaseTr = contract.IncreaseApproval(to, increaseValue).send();
            collector.logStepPass("Send IncreaseApproval, hash:  " + increaseTr.getTransactionHash()
                    + " gasUsed: " + increaseTr.getGasUsed());

            BigInteger afterIncreaseAllowance = contract.Allowance(credentials.getAddress(chainId), to).send();
            collector.logStepPass("Call Allowance after increaseApproval, res: " + afterIncreaseAllowance);
            collector.assertEqual(afterIncreaseAllowance, value.add(increaseValue));

            // DecreaseApproval
            BigInteger decreaseValue = new BigInteger("23300000");
            TransactionReceipt decreaseTr = contract.DecreaseApproval(to, decreaseValue).send();
            collector.logStepPass("Send DecreaseApproval, hash:  " + decreaseTr.getTransactionHash()
                    + " gasUsed: " + decreaseTr.getGasUsed());

            BigInteger afterDecreaseAllowance = contract.Allowance(credentials.getAddress(chainId), to).send();
            collector.logStepPass("Call Allowance after DecreaseApproval, res: " + afterDecreaseAllowance);
            collector.assertEqual(afterDecreaseAllowance, value);

            // TransferFrom
            // spender: a11859ce23effc663a9460e332ca09bd812acc390497f8dc7542b6938e13f8d7
            // 0x493301712671Ada506ba6Ca7891F436D29185821
            // to: 0x1E1ae3407377F7897470FEf31a80873B4FD75cA1
            Credentials spendCredentials = Credentials.create("a11859ce23effc663a9460e332ca09bd812acc390497f8dc7542b6938e13f8d7");
            t.sendFunds(spendCredentials.getAddress(chainId), new BigDecimal(10), Convert.Unit.LAT, provider.getGasPrice(), provider.getGasLimit()).send();
            TransactionManager spenderTM = transactionManager = new RawTransactionManager(web3j, spendCredentials, chainId);
            String to2 = "lax1rcdwxsrnwlmcjarslme34qy88d8awh9pnjpmz9";
            BigInteger valule2 = new BigInteger("10000");
            VIDToken v = VIDToken.load(contractAddress, web3j, spenderTM, provider, chainId);
            TransactionReceipt transferFromTr = v.TransferFrom(credentials.getAddress(chainId), to2, valule2).send();
            collector.logStepPass("Send TransferFrom, hash:  " + transferFromTr.getTransactionHash()
                    + " gasUsed: " + transferFromTr.getGasUsed());
            BigInteger to2Balance = contract.BalanceOf(to2).send();
            collector.logStepPass("Call balanceOf 2, res: " + to2Balance);
            collector.assertEqual(to2Balance, valule2);
            collector.logStepPass("Check TransferFrom() and Approve() success.");

            // TransferToken
            BigInteger transferTokenValue = new BigInteger("10000");
            TransactionReceipt transferTokenTR = contract.TransferToken(to, transferTokenValue).send();
            collector.logStepPass("Send TransferToken, hash:  " + transferTokenTR.getTransactionHash()
                    + " gasUsed: " + transferTokenTR.getGasUsed() + " logs:" + transferTokenTR.getLogs().size());
            VIDToken.TransferEvEventResponse transferTokenResponse = contract.getTransferEvEvents(transferTokenTR).get(0);
            collector.logStepPass("Send TransaferToken Logs: "
                    + " arg1: " + transferTokenResponse.arg1
                    + " arg2: " + transferTokenResponse.arg2
                    + " arg3: " + transferTokenResponse.arg3);

            // Burn
            BigInteger burnValue = new BigInteger("122");
            TransactionReceipt burnTr = contract.Burn(burnValue).send();
            collector.logStepPass("Send Burn, hash:  " + burnTr.getTransactionHash()
                    + " gasUsed: " + burnTr.getGasUsed() + " logs:" + burnTr.getLogs().size());
            VIDToken.BurnEvEventResponse burnResponse = contract.getBurnEvEvents(burnTr).get(0);
            collector.assertTrue(burnTr.getLogs().size() != 0);
            collector.logStepPass("Send Burn Logs: "
                    + " arg1: " + transferTokenResponse.arg1
                    + " arg2: " + transferTokenResponse.arg2);

            // Freeze
            TransactionReceipt freezeTr = contract.Freeze(credentials.getAddress(chainId), true).send();
            collector.logStepPass("Send Freeze, hash:  " + freezeTr.getTransactionHash()
                    + " gasUsed: " + freezeTr.getGasUsed() + " logs:" + freezeTr.getLogs().size());
            VIDToken.FreezeEvEventResponse freezeResponse = contract.getFreezeEvEvents(freezeTr).get(0);
            collector.assertTrue(freezeTr.getLogs().size() != 0);
            collector.logStepPass("Send Burn Logs: "
                    + " arg1: " + freezeResponse.arg1
                    + " arg2: " + freezeResponse.arg2);

            // ValidatePublisher
            TransactionReceipt publisherTr = contract.ValidatePublisher(to, true, credentials.getAddress(chainId)).send();
            collector.logStepPass("Send ValidatePublisher, hash:  " + publisherTr.getTransactionHash()
                    + " gasUsed: " + publisherTr.getGasUsed() + " logs:" + publisherTr.getLogs().size());
            VIDToken.ValidatePublisherEvEventResponse publisherResponse = contract.getValidatePublisherEvEvents(publisherTr).get(0);
            collector.assertTrue(publisherTr.getLogs().size() != 0);
            collector.logStepPass("Send ValidatePublisher Logs: "
                    + " arg1: " + publisherResponse.arg1
                    + " arg2: " + publisherResponse.arg2
                    + " arg3: " + publisherResponse.arg3);

            // ValidateWallet
            TransactionReceipt walletTr = contract.ValidateWallet(to, true, credentials.getAddress(chainId)).send();
            collector.logStepPass("Send ValidateWallet, hash:  " + walletTr.getTransactionHash()
                    + " gasUsed: " + walletTr.getGasUsed() + " logs:" + walletTr.getLogs().size());
            VIDToken.ValidateWalletEvEventResponse walletResponse = contract.getValidateWalletEvEvents(walletTr).get(0);
            collector.assertTrue(publisherTr.getLogs().size() != 0);
            collector.logStepPass("Send ValidateWallet Logs: "
                    + " arg1: " + walletResponse.arg1
                    + " arg2: " + walletResponse.arg2
                    + " arg3: " + walletResponse.arg3);

            // listFiles
            TransactionReceipt listFileTR = contract.ListFiles(BigInteger.ZERO, BigInteger.TEN).send();
            collector.logStepPass("Send ListFiles, hash:  " + listFileTR.getTransactionHash()
                    + " gasUsed: " + listFileTR.getGasUsed() + " logs:" + listFileTR.getLogs().size());

            for (int i = 0; i < listFileTR.getLogs().size(); i++) {
                VIDToken.LogEventEvEventResponse response = contract.getLogEventEvEvents(listFileTR).get(i);
                collector.logStepPass("Send ListFiles Logs: "
                        + " arg1: " + response.arg1
                        + " arg2: " + response.arg2);
            }

            // VerifyFile
            boolean verifyFileRes = contract.VerifyFile("").send();
            collector.logStepPass("Call VerifyFile, res: " + verifyFileRes);

            // set price
            TransactionReceipt setPriceTr = contract.SetPrice(new BigInteger("30000000000")).send();
            collector.logStepPass("Send SetPrice, hash:  " + setPriceTr.getTransactionHash()
                    + " gasUsed: " + setPriceTr.getGasUsed() + " logs:" + setPriceTr.getLogs().size());

            // set wallet
            TransactionReceipt setWalltTr = contract.SetWallet(credentials.getAddress(chainId)).send();
            collector.logStepPass("Send SetWallet, hash:  " + setWalltTr.getTransactionHash()
                    + " gasUsed: " + setWalltTr.getGasUsed() + " logs:" + setWalltTr.getLogs().size());

        } catch (Exception e) {
            collector.logStepFail("contract_VIDToken failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }


}
