package wasm.data_type;

import com.platon.rlp.datatypes.Uint32;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BindFunctionContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约bind函数
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class BindFunctionTest extends WASMContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.bindFunctionTest测试合约bind函数",sourcePrefix = "wasm")
    public void testBindFunction() {

         //部署合约
        BindFunctionContract bindFunctionContract = null;
        try {
            prepare();
            bindFunctionContract = BindFunctionContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = bindFunctionContract.getContractAddress();
            TransactionReceipt tx = bindFunctionContract.getTransactionReceipt().get();
            collector.logStepPass("bindFunctionContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("bindFunctionContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：bind绑定普通函数
            Uint8  actualIntValue = bindFunctionContract.get_bind_function().send();
            collector.logStepPass("bindFunctionContract 【验证bind绑定普通函数】 执行get_bind_function() successfully actualIntValue:" + actualIntValue);
            collector.assertEqual(actualIntValue,Uint8.of(6), "checkout  execute success.");

            //2、验证:bind绑定类的成员函数
            Uint32 actualIntValue2 = bindFunctionContract.get_bind_class_function().send();
            collector.logStepPass("bindFunctionContract 【验证bind绑定类的成员函数(指针形式调用成员函数)】 执行get_bind_class_function() successfully actualIntValue2:" + actualIntValue2);
            collector.assertEqual(actualIntValue2,Uint32.of(6), "checkout  execute success.");

            //3、验证:bind绑定类的成员函数(对象形式调用成员函数)
            Uint32 actualIntValue3 = bindFunctionContract.get_bind_class_function_one().send();
            collector.logStepPass("bindFunctionContract 【验证bind绑定类的成员函数(对象形式调用成员函数)】 执行get_bind_class_function_one() successfully actualIntValue3:" + actualIntValue3);
            collector.assertEqual(actualIntValue3,Uint32.of(6), "checkout  execute success.");

            //4、验证:bind绑定类静态成员函数
            Uint32 actualIntValue4 = bindFunctionContract.get_bind_static_function().send();
            collector.logStepPass("bindFunctionContract 【验证bind绑定类静态成员函数】 执行get_bind_static_function() successfully actualIntValue4:" + actualIntValue4);
            collector.assertEqual(actualIntValue4,Uint32.of(6), "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("bindFunctionContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
