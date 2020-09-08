package network.platon.test.wasm.contract_multi_inherit;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.OneInherit;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约基础类单继承测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/11
 */
public class OneInheritTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.OneInherit合约基础类单继承测试",sourcePrefix = "wasm")
    public void testOneInhert() {

        String head = "myHead";
        String body = "myBody";
        String end = "myEnd";

        try {
            prepare();
            OneInherit oneInherit = OneInherit.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = oneInherit.getContractAddress();
            String transactionHash = oneInherit.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("OneInherit issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + oneInherit.getTransactionReceipt().get().getGasUsed());

            OneInherit.My_message my_message = new OneInherit.My_message();
            OneInherit.Message message = new OneInherit.Message();
            message.head = head;

            my_message.baseClass = message;
            my_message.body = body;
            my_message.end = end;

            TransactionReceipt transactionReceipt = oneInherit.add_my_message(my_message).send();
            collector.logStepPass("OneInheritTest call add_my_message successfully hash:" + transactionReceipt.getTransactionHash());

            //查询vector中对象数量
            Uint8 size = oneInherit.get_my_message_size().send();
            collector.logStepPass("vector中my_message 数量为："+size);

            //查询消息头信息
            Uint8 idx = Uint8.of(0);
            String chainHead = oneInherit.get_my_message_head(idx).send();
            collector.logStepPass("OneInheritTest call get_my_message_head successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainHead,head);

            //查询消息体信息
            String chainBody = oneInherit.get_my_message_body(idx).send();
            collector.logStepPass("OneInheritTest call get_my_message_body successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainBody,body);


        } catch (Exception e) {
            collector.logStepFail("OneInheritTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
