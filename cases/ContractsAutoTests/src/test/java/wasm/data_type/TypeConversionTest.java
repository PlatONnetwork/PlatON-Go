package wasm.data_type;

import com.platon.rlp.datatypes.*;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BasicDataTypeContract;
import network.platon.contracts.wasm.TypeConversionContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约数据类型转换
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class TypeConversionTest extends WASMContractPrepareTest {




    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.typeConversionTest测试合约数据类型转换",sourcePrefix = "wasm")
    public void testTypeConversion() {
         //部署合约
        TypeConversionContract typeConversionContract = null;
        try {
            prepare();
            typeConversionContract = TypeConversionContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = typeConversionContract.getContractAddress();
            TransactionReceipt tx = typeConversionContract.getTransactionReceipt().get();
            collector.logStepPass("typeConversionContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());

        } catch (Exception e) {
            collector.logStepFail("typeConversionContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证:隐式类型转换(混合类型的算术运算)
            Uint8 a = Uint8.of(5);
            Uint64 b = Uint64.of(10);
            Uint64  actualUint64Value = typeConversionContract.get_add(a,b).send();
            collector.logStepPass("basicDataTypeContract 【验证隐式类型转换(混合类型的算术运算)】 执行get_add() successfully actualUint64Value:" + actualUint64Value);
            collector.assertEqual(actualUint64Value,Uint64.of(15), "checkout  execute success.");
            //2、验证：隐式类型转换(不同类型的赋值操作)
            Boolean bool = true;
            Uint64 actualUint64Value1 = typeConversionContract.get_different_type_(bool).send();
            collector.logStepPass("basicDataTypeContract 【验证隐式类型转换(不同类型的赋值操作)】 执行get_different_type_() successfully actualUint64Value1:" + actualUint64Value1);
             collector.assertEqual(actualUint64Value1,Uint64.of(1), "checkout  execute success.");

            //3、、验证：函数参数传值时类型转换
            Uint32 actualValue = typeConversionContract.get_pram_type().send();
            collector.logStepPass("basicDataTypeContract 【验证函数参数传值时类型转换】 执行get_pram_type() successfully actualValue:" + actualValue);
             collector.assertEqual(actualValue,Uint32.of(1), "checkout  execute success.");

            //4、验证：隐式类型转换（函数返回值）
            Uint8 aValue = Uint8.of(10);
            Uint8 bValue = Uint8.of(10);
            Uint64 actualValue1 = typeConversionContract.get_pram_return(aValue,bValue).send();
            collector.logStepPass("basicDataTypeContract 【验证隐式类型转换（函数返回值）】 执行get_pram_return() successfully actualValue1:" + actualValue1);
             collector.assertEqual(actualValue1,Uint64.of(20), "checkout  execute success.");

            //5、验证：显示类型转换(强制类型转换)
            Uint8 actualValue2 = typeConversionContract.get_convert().send();
            collector.logStepPass("basicDataTypeContract 【验证强制类型转换】 执行get_convert() successfully actualValue2:" + actualValue2);
            collector.assertEqual(actualValue2,Uint8.of(234), "checkout  execute success.");

            //6、验证：强制类型转换(static_cast)
            Uint8 actualValue3 = typeConversionContract.get_convert_static_cast().send();
            collector.logStepPass("basicDataTypeContract 【验证强制类型转换(static_cast)】 执行get_convert_static_cast() successfully actualValue3:" + actualValue3);
            collector.assertEqual(actualValue3,Uint8.of(244), "checkout  execute success.");

            //7、验证：强制类型转换(const_cast)
            Uint8 actualValue4 = typeConversionContract.get_convert_const_cast().send();
            collector.logStepPass("basicDataTypeContract 【验证强制类型转换(const_cast)】 执行get_convert_const_cast() successfully actualValue4:" + actualValue4);
             collector.assertEqual(actualValue4,Uint8.of(10), "checkout  execute success.");



        } catch (Exception e) {
            collector.logStepFail("basicDataTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
