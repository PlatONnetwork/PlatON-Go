package wasm.data_type;

import com.platon.rlp.datatypes.Int32;
import com.platon.rlp.datatypes.Int64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.AutoTypeContract;
import network.platon.contracts.wasm.BasicDataTypeContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约anto关键字
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class AutoTypeTest extends WASMContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.autoTypeTest测试合约anto关键字",sourcePrefix = "wasm")
    public void testAutoType() {

         //部署合约
        AutoTypeContract autoTypeContract = null;
        try {
            prepare();
            autoTypeContract = AutoTypeContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = autoTypeContract.getContractAddress();
            TransactionReceipt tx = autoTypeContract.getTransactionReceipt().get();
            collector.logStepPass("autoTypeContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("autoTypeContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证:自动匹配int类型
            Int32  actualIntValue = autoTypeContract.get_anto_int().send();
            collector.logStepPass("autoTypeContract 【验证自动匹配int类型】 执行get_anto_int() successfully actualIntValue:" + actualIntValue);
            collector.assertEqual(actualIntValue,Int32.of(5), "checkout  execute success.");
            //2、验证:自动匹配int类型负数值
            Int32  actualIntValue2 = autoTypeContract.get_anto_int32().send();
            collector.logStepPass("autoTypeContract 【验证自动匹配int类型负数值】 执行get_anto_int32() successfully actualIntValue2:" + actualIntValue2);
            collector.assertEqual(actualIntValue,Int32.of(5), "checkout  execute success.");


            //2、验证:自动匹配double类型
          /*  Double  actualDoubleValue = autoTypeContract.get_anto_double().send();
            collector.logStepPass("autoTypeContract 【验证自动匹配double类型】 get_anto_double() successfully actualDoubleValue:" + actualDoubleValue);
            collector.assertEqual(actualDoubleValue,Double.parseDouble("1.0"), "checkout  execute success.");*/
            //3、验证:自动匹配多个值类型
            Int32  actualmultipleValue = autoTypeContract.get_anto_multiple().send();
            collector.logStepPass("autoTypeContract 【验证自动匹配多个值类型】 get_anto_multiple() successfully actualmultipleValue:" + actualmultipleValue);
            collector.assertEqual(actualmultipleValue,Int32.of(30), "checkout  execute success.");
            //4、验证:自动匹配uint8类型
            Uint8  actualUint8Value = autoTypeContract.get_anto_uint8_t().send();
            collector.logStepPass("autoTypeContract 【验证自动匹配uint8类型】 get_anto_uint8_t() successfully actualUint8Value:" + actualUint8Value);
            collector.assertEqual(actualUint8Value,Uint8.of(10), "checkout  execute success.");

            //5、验证:自动匹配表达式
        /*    Double  actualValue = autoTypeContract.get_anto_express().send();
            collector.logStepPass("autoTypeContract 【验证自动匹配表达式】 执行get_anto_express() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,Double.parseDouble("15.32"), "checkout  execute success.");*/
            //6、验证:迭代器中应用auto
            Uint8  iteratorCount = autoTypeContract.get_anto_iterator().send();
            collector.logStepPass("autoTypeContract 【验证迭代器中应用auto】 执行get_anto_iterator() successfully iteratorCount:" + iteratorCount);
            collector.assertEqual(iteratorCount,Uint8.of(3), "checkout  execute success.");

        } catch (Exception e) {
            collector.logStepFail("autoTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
