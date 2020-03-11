package wasm.complex_contract;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import com.platon.rlp.datatypes.WasmAddress;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.Bank;
import network.platon.contracts.wasm.ContractEmitEvent1;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import wasm.beforetest.WASMContractPrepareTest;
import network.platon.contracts.wasm.Guessing;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;

/**
 * @author hudenian
 *
 */
public class GuessingTest extends WASMContractPrepareTest {

    @Before
    public void before(){
        prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.GuessingTest-竞猜合约",sourcePrefix = "wasm")
    public void testGuessingContract() {

        Long blocks = 30L;//设置截止块高与当前块高为20

        try {
            //设置结束时块高
            BigInteger endBlock = web3j.platonBlockNumber().send().getBlockNumber().add(BigInteger.valueOf(blocks));

            Guessing guessing = Guessing.deploy(web3j, transactionManager, provider, Uint64.of(endBlock)).send();
            String contractAddress = guessing.getContractAddress();
            String transactionHash = guessing.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("Guessing issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("Guessing deploy successfully. gasUsed: " + guessing.getTransactionReceipt().get().getGasUsed().toString());

            //发起竞猜
            TransactionReceipt transactionReceipt = guessing.guessingWithLat("1000").send();
            collector.logStepPass("Guessing call guessingWithLat successfully hash:" + transactionReceipt.getTransactionHash());

            //对事件信息进行解析
            List<Guessing.Transfer1EventResponse> eventList = guessing.getTransfer1Events(transactionReceipt);
            String data = eventList.get(0).log.getData();
            collector.logStepPass("event arg1 value is:"+eventList.get(0).arg1);
            collector.logStepPass("topics is:"+eventList.get(0).log.getTopics().get(0).toString());

            String indexKey = guessing.getIndexKey().send().toString();
            collector.logStepPass("indexKey is:"+indexKey);

            //开奖操作
//            transactionReceipt = guessing.draw().send();
//            collector.logStepPass("Guessing call draw successfully hash:" + transactionReceipt.getTransactionHash());

            //查看合约中的余额
            String balance = guessing.getBalance().send().toString();
            collector.logStepPass("contract balance is:"+balance);


            //获取所有中奖人地址
            WasmAddress[] wasmAddresses = guessing.getWinnerAddresses().send();
            Arrays.asList(wasmAddresses).stream().forEach(addr ->{collector.logStepPass("获奖人地址:"+addr);});



        } catch (Exception e) {
            collector.logStepFail("Guessing failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

}
