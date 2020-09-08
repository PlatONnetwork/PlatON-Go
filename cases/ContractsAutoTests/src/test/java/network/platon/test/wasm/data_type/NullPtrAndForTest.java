package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint32;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.NullPtrAndForContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约Nullptr 和 序列for循环
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class NullPtrAndForTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.nullPtrAndForTest合约Nullptr和序列for循环",sourcePrefix = "wasm")
    public void testNullPtrAndFor() {

         //部署合约
        NullPtrAndForContract nullPtrAndForContract = null;
        try {
            prepare();
            nullPtrAndForContract = NullPtrAndForContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = nullPtrAndForContract.getContractAddress();
            TransactionReceipt tx = nullPtrAndForContract.getTransactionReceipt().get();
            collector.logStepPass("nullPtrAndForContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("nullPtrAndForContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证:nullprt赋值
            boolean  actualValue = nullPtrAndForContract.get_nullptr().send();
            collector.logStepPass("nullPtrAndForContract 【验证nullprt赋值】 执行get_nullptr() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,true, "checkout  execute success.");
            //2、验证:nullprt赋值不同类型
            boolean  actualValue2 = nullPtrAndForContract.get_nullptr_one().send();
            collector.logStepPass("nullPtrAndForContract 【验证nullprt赋值不同类型】 执行get_nullptr_one() successfully actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,false, "checkout  execute success.");
            //3、验证:验证NULL在C++中就是0；Nullptr表示为空指针
            String  actualValue3 = nullPtrAndForContract.set_nullptr_overload().send();
            collector.logStepPass("nullPtrAndForContract 【验证NULL在C++中就是0；Nullptr表示为空指针】 set_nullptr_overload() successfully actualValue3:" + actualValue3);
            collector.assertEqual(actualValue3,"is nullptr", "checkout  execute success.");
            //4、验证:foreach遍历map容器
            String  actualValue4 = nullPtrAndForContract.get_foreach_map().send();
            collector.logStepPass("nullPtrAndForContract 【验证foreach遍历map容器】 get_foreach_map() successfully actualValue4:" + actualValue4);
            collector.assertEqual(actualValue4,"four,one,three,two,", "checkout  execute success.");
            //5、验证:foreach遍历array容器
            Uint32 actualValue5 = nullPtrAndForContract.get_foreach_array().send();
            collector.logStepPass("nullPtrAndForContract 【验证foreach遍历array容器】 get_foreach_array() successfully actualValue5:" + actualValue5);
            collector.assertEqual(actualValue5,Uint32.of(45), "checkout  execute success.");


        } catch (Exception e) {
            collector.logStepFail("autoTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
