package network.platon.test.wasm.contract_multi_inherit;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.TwoInherit;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约基础类双继承测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/11
 */
public class TwoInheritTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.TwoInheritT合约基础类双继承测试",sourcePrefix = "wasm")
    public void testTwoInhert() {

        String head = "myHead";
        String body = "myBody";
        String end = "myEnd";
        String from = "myFrom";
        String to = "myTo";

        try {
            prepare();
            TwoInherit twoInherit = TwoInherit.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = twoInherit.getContractAddress();
            String transactionHash = twoInherit.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("twoInherit issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + twoInherit.getTransactionReceipt().get().getGasUsed());

            //基类
            TwoInherit.Message message = new TwoInherit.Message();
            message.head = head;

            //子类
            TwoInherit.My_message my_message = new TwoInherit.My_message();
            my_message.baseClass = message;
            my_message.body = body;
            my_message.end = end;

            //孙子类
            TwoInherit.Sub_my_message sub_my_message = new TwoInherit.Sub_my_message();
            sub_my_message.from = from;
            sub_my_message.to = to;
            sub_my_message.baseClass = my_message;


            TransactionReceipt transactionReceipt = twoInherit.add_sub_my_message(sub_my_message).send();
            collector.logStepPass("TwoInheritTest call add_my_message successfully hash:" + transactionReceipt.getTransactionHash());

            //查询vector中对象数量
            Uint8 size = twoInherit.get_sub_my_message_size().send();
            collector.logStepPass("vector中sub_my_message 数量为："+size);

            //查询消息头信息
            Uint8 idx = Uint8.of(0);
            String chainHead = twoInherit.get_sub_my_message_head(idx).send();
            collector.logStepPass("TwoInheritTest call get_sub_my_message_head successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainHead,head);

            //查询消息from
            String chainFrom = twoInherit.get_sub_my_message_from(idx).send();
            collector.logStepPass("TwoInheritTest call get_sub_my_message_from successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainFrom,from);


        } catch (Exception e) {
            collector.logStepFail("TwoInheritTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
