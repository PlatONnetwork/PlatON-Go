package wasm.data_type;

import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.FunctionTemplateContract;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约类模板std :: function
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class FunctionTemplateTest extends WASMContractPrepareTest {


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.autoTypeTest测试合约类模板std :: function",sourcePrefix = "wasm")
    public void testFunctionTemplate() {

         //部署合约
        FunctionTemplateContract functionTemplateContract = null;
        try {
            prepare();
            functionTemplateContract = FunctionTemplateContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = functionTemplateContract.getContractAddress();
            TransactionReceipt tx = functionTemplateContract.getTransactionReceipt().get();
            collector.logStepPass("functionTemplateContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("functionTemplateContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证:调用lambda表达式
            Uint8  actualValue = functionTemplateContract.get_lambda_function().send();
            collector.logStepPass("functionTemplateContract 【验证调用lambda表达式】 执行get_lambda_function() successfully actualValue:" + actualValue);
            collector.assertEqual(actualValue,Uint8.of(6), "checkout  execute success.");

            //2、验证:调用普通函数
            Uint8  actualValue1 = functionTemplateContract.get_normal_function().send();
            collector.logStepPass("functionTemplateContract 【验证调用普通函数】 执行get_normal_function() successfully actualValue1:" + actualValue1);
            collector.assertEqual(actualValue1,Uint8.of(6), "checkout  execute success.");

            //3、验证:调用类静态成员函数
            Uint8  actualValue2 = functionTemplateContract.get_class_static_function().send();
            collector.logStepPass("functionTemplateContract 【验证调用类静态成员函数】 执行get_class_static_function() successfully actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,Uint8.of(6), "checkout  execute success.");
        } catch (Exception e) {
            collector.logStepFail("functionTemplateContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
