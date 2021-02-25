package network.platon.test.wasm.contract_create;

import com.platon.rlp.datatypes.Int8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitOverloadWithString;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约包含字符串操作
 * @description:
 * @author: hudenian
 * @create: 2020/02/29
 */
public class InitOverloadWithStringTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约包含字符串操作",sourcePrefix = "wasm")
    public void testStringOpt() {

        String name = "how are you";
        String who = "lily";
        try {
            prepare();
            InitOverloadWithString initOverloadWithString = InitOverloadWithString.deploy(web3j, transactionManager, provider, chainId, name).send();
            String contractAddress = initOverloadWithString.getContractAddress();
            String transactionHash = initOverloadWithString.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitOverloadWithString issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initOverloadWithString.getTransactionReceipt().get().getGasUsed());

            //获取字符串大小
            String chainStrLength = initOverloadWithString.string_length().send().value.toString();
            collector.assertEqual(String.valueOf(name.length()),chainStrLength);

            //字符串连接
            String chainSpliceStr = initOverloadWithString.string_splice(who).send().toString();
            collector.assertEqual(name+who,chainSpliceStr);

            //字符串比较
            Int8 result = initOverloadWithString.string_compare("abc","efg").send();
            collector.assertEqual("-1",result.toString());

            //字符串查找
            Int8 location = Int8.of(initOverloadWithString.string_find("how").send().value);
            collector.logStepPass("how location is:"+location);

            //字符串倒置
            TransactionReceipt tx = initOverloadWithString.string_reverse(name).send();
            String reverseStr = initOverloadWithString.get_string().send();
            collector.logStepPass("reverse Strirng is:"+reverseStr);


        } catch (Exception e) {
            collector.logStepFail("InitOverloadWithStringTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
