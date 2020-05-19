package wasm.complex_contract;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.GuessingWasm;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;

/**
 * @author hudenian
 *
 */
public class GuessingWasmTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.GuessingWasmTest-竞猜合约",sourcePrefix = "wasm")
    public void testGuessing() {

        Long blocks = 30L;//设置截止块高与当前块高为20
        try {
            //设置结束时块高
            BigInteger endBlock = web3j.platonBlockNumber().send().getBlockNumber().add(BigInteger.valueOf(blocks));

            GuessingWasm guessing = GuessingWasm.deploy(web3j, transactionManager, provider,chainId, Uint64.of(endBlock)).send();
            String contractAddress = guessing.getContractAddress();
            String transactionHash = guessing.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("GuessingWasm issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("GuessingWasm deploy successfully. gasUsed: " + guessing.getTransactionReceipt().get().getGasUsed().toString());

            //发起竞猜 至少5LAT
            TransactionReceipt transactionReceipt = guessing.guessingWithLat(new BigInteger("5000000000000000000")).send();
            collector.logStepPass("GuessingWasm call guessingWithLat successfully hash:" + transactionReceipt.getTransactionHash());

            //对事件信息进行解析
            List<GuessingWasm.Transfer1EventResponse> eventList = guessing.getTransfer1Events(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.logStepPass("event arg1 value is:"+eventList.get(0).arg1);
            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());

            String indexKey = guessing.getIndexKey().send().toString();
            collector.logStepPass("indexKey is:"+indexKey);

            //开奖操作
            transactionReceipt = guessing.draw().send();
            collector.logStepPass("GuessingWasm call draw successfully hash:" + transactionReceipt.getTransactionHash());

            //查看合约中的余额
            String balance = guessing.getBalance().send().toString();
            collector.logStepPass("contract balance is:"+balance);


            //获取所有中奖人地址
            WasmAddress[] wasmAddresses = guessing.getWinnerAddresses().send();
            for (WasmAddress wasmAddresse : wasmAddresses) {
                collector.logStepPass("获奖人地址:" + wasmAddresse);
            }
            //Arrays.asList(wasmAddresses).stream().forEach(addr ->{collector.logStepPass("获奖人地址:"+addr);});

        } catch (Exception e) {
            collector.logStepFail("Guessing wasm failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }

    }

}
