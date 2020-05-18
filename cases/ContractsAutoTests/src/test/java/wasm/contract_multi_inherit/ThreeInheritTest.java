package wasm.contract_multi_inherit;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ThreeInherit;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约基础类三次继承测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/16
 */
public class ThreeInheritTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.TwoInheritT合约基础类三次继承测试",sourcePrefix = "wasm")
    public void testThreeInhert() {

        String head = "myHead";
        String body = "myBody";
        String end = "myEnd";
        String from = "myFrom";
        String to = "myTo";
        String level = "myLevel";
        String desc = "myDesc";

        try {
            prepare();
            ThreeInherit threeInherit = ThreeInherit.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = threeInherit.getContractAddress();
            String transactionHash = threeInherit.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ThreeInherit issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + threeInherit.getTransactionReceipt().get().getGasUsed());

            //基类
            ThreeInherit.Message message = new ThreeInherit.Message();
            message.head = head;

            //子类
            ThreeInherit.My_message my_message = new ThreeInherit.My_message();
            my_message.baseClass = message;
            my_message.body = body;
            my_message.end = end;

            //孙子类
            ThreeInherit.Sub_my_message sub_my_message = new ThreeInherit.Sub_my_message();
            sub_my_message.from = from;
            sub_my_message.to = to;
            sub_my_message.baseClass = my_message;

            //曾孙类
            ThreeInherit.Greate_sub_my_message greate_sub_my_message= new ThreeInherit.Greate_sub_my_message();
            greate_sub_my_message.level = level;
            greate_sub_my_message.desc = desc;
            greate_sub_my_message.baseClass = sub_my_message;


            TransactionReceipt transactionReceipt = threeInherit.add_greate_sub_my_message(greate_sub_my_message).send();
            collector.logStepPass("ThreeInheritTest call add_my_message successfully hash:" + transactionReceipt.getTransactionHash());

            //查询vector中对象数量
            Uint8 size = threeInherit.get_greate_sub_my_message_size().send();
            collector.logStepPass("vector中sub_my_message 数量为："+size);

            //查询消息头信息
            Uint8 idx = Uint8.of(0);
            String chainHead = threeInherit.get_greate_sub_my_message_head(idx).send();
            collector.logStepPass("ThreeInheritTest call get_sub_my_message_head successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainHead,head);

            //查询曾孙desc
            String chainDesc = threeInherit.get_greate_sub_my_message_desc(idx).send();
            collector.logStepPass("ThreeInheritTest call get_sub_my_message_from successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainDesc,desc);


        } catch (Exception e) {
            collector.logStepFail("ThreeInheritTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
