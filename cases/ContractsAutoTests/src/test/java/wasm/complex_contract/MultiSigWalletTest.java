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
 *
 */
public class MultiSigWalletTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.MultiSigWalletTest-MultiSigWallet",sourcePrefix = "wasm")
    public void testMultiSigWallet() {

        String address1 = "0xBE0af016941Acaf08Bf5f4ad185155Df6B7388ce";
        String address2 = "0x08b2320Ef2482f0a5ad9411CCB1a748BcE7c2823";
        String address3 = "0xa79e067808235f31C3A11dd74Ceb6A3d1c3B45ad";
        String address4 = "0x5B37dABeDAE06Edb142257819FAD207199986992";
        String address5 = "0x90CfBCf969F35a0E03b9d4B0FC59e83ff05A81Cd";
        Set<WasmAddress> owners =new HashSet<>();
        owners.add(new WasmAddress(address1));
        owners.add(new WasmAddress(address2));
        owners.add(new WasmAddress(address3));

        try {
            MultiSigWallet multiSigWallet = MultiSigWallet.deploy(web3j, transactionManager, provider,Uint64.of("2"),owners).send();
            String contractAddress = multiSigWallet.getContractAddress();
            String transactionHash = multiSigWallet.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("MultiSigWallet issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("MultiSigWallet deploy successfully. gasUsed: " + multiSigWallet.getTransactionReceipt().get().getGasUsed().toString());

            Map<WasmAddress,Boolean> isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:"+ isOwner.size());

            TransactionReceipt transactionReceipt = multiSigWallet.addOwner(new WasmAddress(address5)).send();
            collector.logStepPass("MultiSigWallet call addOwner hash is:"+transactionReceipt.getTransactionHash());
            isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:"+ isOwner.size());

            transactionReceipt = multiSigWallet.removeOwner(new WasmAddress(address5)).send();
            collector.logStepPass("MultiSigWallet call removeOwner hash is:"+transactionReceipt.getTransactionHash());
            isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:"+ isOwner.size());

            transactionReceipt = multiSigWallet.changeThreshold(Uint64.of(2L)).send();
            collector.logStepPass("MultiSigWallet call changeThreshold hash is:"+transactionReceipt.getTransactionHash());
            isOwner = multiSigWallet.getIsOwner().send();
            collector.logStepPass("isOwner map size is:"+ isOwner.size());


//            byte[] _data = "helloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallet".getBytes();
//            byte[] _signatures = "helloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallethelloMultiSigWalletTesttestMultiSigWalletMultiSigWalletTesttestMultiSigWallet".getBytes();
//
//            transactionReceipt = multiSigWallet.execute(new WasmAddress(address4),Uint64.of(10L),_data,_signatures).send();
//            collector.logStepPass("MultiSigWallet call execute hash is:"+transactionReceipt.getTransactionHash());

        } catch (Exception e) {
            collector.logStepFail("MultiSigWallet failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
