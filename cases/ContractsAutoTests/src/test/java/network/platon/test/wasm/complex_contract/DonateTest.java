package network.platon.test.wasm.complex_contract;

import com.platon.rlp.datatypes.Uint128;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Donate;
//import network.platon.contracts.wasm.VRF;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.utils.Convert;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.List;
import java.util.Map;

/**
 * @author denglonghui
 *
 */
public class DonateTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "denglonghui", showName = "wasm.DonateTest-DonateTest",sourcePrefix = "wasm")
    public void testDonate() {

        BigInteger initialVonValue = BigInteger.valueOf(100L);
        WasmAddress _charity = new WasmAddress(credentials.getAddress(chainId));
        System.out.println(_charity);
        Uint128 _openingTime = Uint128.of((System.currentTimeMillis() - 1 * 60 * 60 * 1000));
        Uint128 _closingTime = Uint128.of((System.currentTimeMillis() + 24 * 60 * 60 * 1000));
        Uint128 _minVonAmount = Uint128.of(1);
        Uint128 _maxVonAmount = Uint128.of(Convert.toVon(new BigDecimal(10000), Convert.Unit.LAT).toBigInteger());
        Uint128 _maxNumDonors = Uint128.of(100000);

        try {
            System.out.println(_openingTime);
            System.out.println(_closingTime);
            Donate donate = Donate.deploy(web3j, transactionManager, provider, chainId, _charity, _openingTime, _closingTime, _minVonAmount, _maxVonAmount, _maxNumDonors).send();
            String contractAddress = donate.getContractAddress();
            String transactionHash = donate.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("Donate issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("Donate deploy successfully. gasUsed: " + donate.getTransactionReceipt().get().getGasUsed().toString());

            // 加入白名单
            donate.addToWhitelist(new WasmAddress(credentials.getAddress(chainId))).send();
            Map whitelist = donate.getWhitelist().send();
            System.err.println("Whitelist >>>> " + whitelist);

            // 捐赠，只有加入白名单才能捐赠
            TransactionReceipt receipt = donate.donate(new WasmAddress(credentials.getAddress(chainId)),BigInteger.valueOf(1000)).send();
            List<Donate.DonatedEventResponse> list = donate.getDonatedEvents(receipt);
            collector.logStepPass(">>>>>>>>>>> " + list.get(0).topic);
            collector.logStepPass("捐赠金额>>>>>>>>>> " + list.get(0).arg1);

            collector.logStepPass("After donate, balance >>>> " + web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance());

            // 捐赠名单
            WasmAddress[] wasms = donate.getDonors().send();
            collector.logStepPass(" >>>>>>>>>>>>> " + wasms.length);
        } catch (Exception e) {
            collector.logStepFail("DonateTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
