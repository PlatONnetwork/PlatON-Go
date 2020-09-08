package wasm.complex_contract;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.MultiSigWallet;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.util.HashSet;
import java.util.Map;
import java.util.Set;

/**
 * @author hudenian
 */
public class MultiSigWalletTest extends WASMContractPrepareTest {

    @Before
    public void before() {
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.MultiSigWalletTest-MultiSigWallet", sourcePrefix = "wasm")
    public void testMultiSigWallet() {

        String address1 = "lax1w2kjkufl4g2v93xd94a0lewc75ufdr66rnzuw2";
        String address2 = "lax1fyeszufxwxk62p46djncj86rd553skpptsj8v6";
        /*String address3 = "lax1570qv7qgyd0nrsaprht5e6m285wrk3ddyeafxt";
        String address4 = "lax1tvma40k6uphdk9pz27qeltfqwxves6vjx3wjjk";*/
        String address5 = "lax1uqug0zq7rcxddndleq4ux2ft3tv6dqljphydrl";
        Set<WasmAddress> owners = new HashSet<>();
        owners.add(new WasmAddress(address1));
        owners.add(new WasmAddress(address2));
//        owners.add(new WasmAddress(address3));

        try {
            MultiSigWallet multiSigWallet = MultiSigWallet.deploy(web3j, transactionManager, provider,chainId, Uint64.of("2"), owners).send();
            String contractAddress = multiSigWallet.getContractAddress();
            String transactionHash = multiSigWallet.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MultiSigWallet issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MultiSigWallet deploy successfully. gasUsed: " + multiSigWallet.getTransactionReceipt().get().getGasUsed().toString());

            Map<WasmAddress, Boolean> isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:" + isOwner.size());

            TransactionReceipt transactionReceipt = multiSigWallet.addOwner(new WasmAddress(address5)).send();
            collector.logStepPass("MultiSigWallet call addOwner hash is:" + transactionReceipt.getTransactionHash());
            isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:" + isOwner.size());

            transactionReceipt = multiSigWallet.removeOwner(new WasmAddress(address5)).send();
            collector.logStepPass("MultiSigWallet call removeOwner hash is:" + transactionReceipt.getTransactionHash());
            isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:" + isOwner.size());

            transactionReceipt = multiSigWallet.changeThreshold(Uint64.of(2L)).send();
            collector.logStepPass("MultiSigWallet call changeThreshold hash is:" + transactionReceipt.getTransactionHash());
            isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:" + isOwner.size());


//            byte[] _data = "helloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallet".getBytes();
//            byte[] _signatures = "helloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallet".getBytes();
//
//            transactionReceipt = multiSigWallet.execute(new WasmAddress(address4),Uint64.of(10L),_data,_signatures).send();
//            collector.logStepPass("MultiSigWallet call execute hash is:"+transactionReceipt.getTransactionHash());

        } catch (Exception e) {
            collector.logStepFail("MultiSigWallet failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

}
