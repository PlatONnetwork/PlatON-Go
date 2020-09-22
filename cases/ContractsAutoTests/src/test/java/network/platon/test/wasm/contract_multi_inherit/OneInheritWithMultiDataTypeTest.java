package network.platon.test.wasm.contract_multi_inherit;

import com.platon.rlp.datatypes.Uint32;
import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.OneInheritWithMultiDataType;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 合约基础类单继承带有多种类型参数测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/17
 */
public class OneInheritWithMultiDataTypeTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.OneInheritWithMultiDataType合约基础类单继承带有多种类型参数测试",sourcePrefix = "wasm")
    public void testOneInhertWithMultiDataType() {

        String head = "myHead";
        String body = "myBody";
        String end = "myEnd";
        Uint32 age = Uint32.of(20);
        Uint64 money = Uint64.of(100000L);

        try {
            prepare();
            OneInheritWithMultiDataType oneInheritWithMultiDataType = OneInheritWithMultiDataType.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = oneInheritWithMultiDataType.getContractAddress();
            String transactionHash = oneInheritWithMultiDataType.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("oneInheritWithMultiDataType issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + oneInheritWithMultiDataType.getTransactionReceipt().get().getGasUsed());

            OneInheritWithMultiDataType.My_message my_message = new OneInheritWithMultiDataType.My_message();
            OneInheritWithMultiDataType.Message message = new OneInheritWithMultiDataType.Message();
            message.head = head;
            message.age = age;
            message.money = money;

            my_message.baseClass = message;
            my_message.body = body;
            my_message.end = end;

            TransactionReceipt transactionReceipt = oneInheritWithMultiDataType.add_my_message(my_message).send();
            collector.logStepPass("OneInheritWithMultiDataTypeTest call add_my_message successfully hash:" + transactionReceipt.getTransactionHash());

            //查询vector中对象数量
            Uint8 size = oneInheritWithMultiDataType.get_my_message_size().send();
            collector.logStepPass("vector中my_message 数量为："+size);

            //查询消息头信息
            Uint8 idx = Uint8.of(0);
            String chainHead = oneInheritWithMultiDataType.get_my_message_head(idx).send();
            collector.logStepPass("OneInheritWithMultiDataTypeTest call get_my_message_head successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainHead,head);

            //查询消息体信息
            String chainBody = oneInheritWithMultiDataType.get_my_message_body(idx).send();
            collector.logStepPass("OneInheritWithMultiDataTypeTest call get_my_message_body successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainBody,body);

            //查询age信息
            Uint32 chainAge = oneInheritWithMultiDataType.get_my_message_age(idx).send();
            collector.logStepPass("OneInheritWithMultiDataTypeTest call get_my_message_age successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainAge,age);

            //查询money信息
            Uint64 chainMoney = oneInheritWithMultiDataType.get_my_message_money(idx).send();
            collector.logStepPass("OneInheritWithMultiDataTypeTest call get_my_message_money successfully hash:" + transactionReceipt.getTransactionHash());
            collector.assertEqual(chainMoney,money);

        } catch (Exception e) {
            collector.logStepFail("OneInheritWithMultiDataTypeTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
