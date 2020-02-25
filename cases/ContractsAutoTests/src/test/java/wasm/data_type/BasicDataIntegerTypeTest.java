package wasm.data_type;


import com.platon.rlp.datatypes.*;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.BasicDataIntegerTypeContract;
import org.junit.Before;
import org.junit.Test;

import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试整型基本类型
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class BasicDataIntegerTypeTest extends WASMContractPrepareTest {

    private String aValue;
    private String bValue;
    private String cValue;
    private String dValue;

    @Before
    public void before() {
        aValue = driverService.param.get("aValue");
        bValue = driverService.param.get("bValue");
        cValue = driverService.param.get("cValue");
        dValue = driverService.param.get("dValue");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.basicDataTypeTest整型基本类型验证测试",sourcePrefix = "wasm")
    public void testBasicDataIntegerTypeTest() {

         //部署合约
        BasicDataIntegerTypeContract basicDataIntegerTypeContract = null;
        try {
            prepare();
            basicDataIntegerTypeContract = BasicDataIntegerTypeContract.deploy(web3j, transactionManager, provider).send();
            String contractAddress = basicDataIntegerTypeContract.getContractAddress();
            TransactionReceipt tx = basicDataIntegerTypeContract.getTransactionReceipt().get();
            collector.logStepPass("basicDataIntegerTypeContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("basicDataIntegerTypeContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：整型有符号/无符号类型
            Int8 int8Value = Int8.of(2);
            Int16 int16Value = Int16.of(3);
            Int32 int32Value = Int32.of(-100);
            Int64 int64Value = Int64.of(200);
            Uint8 uint8Value = Uint8.of(10);
            Uint16 uint16Value = Uint16.of(20);
            Uint32 uint32Value = Uint32.of(30);
            Uint64 uint64Value = Uint64.of(40);


            //int8
            TransactionReceipt  transactionReceipt = basicDataIntegerTypeContract.set_int8(int8Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证int8整型有符号/无符号类型】 successfully hash:" + transactionReceipt.getTransactionHash());
            Int8 actualInt8Value = basicDataIntegerTypeContract.get_int8().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数int8_t取值】 执行get_int8() successfully actualInt8Value:" + actualInt8Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");
            //int16
            TransactionReceipt  transactionReceipt1 = basicDataIntegerTypeContract.set_int16(int16Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证int16整型有符号/无符号类型】 successfully hash:" + transactionReceipt1.getTransactionHash());
            Int16 actualInt16Value = basicDataIntegerTypeContract.get_int16().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数int16_t取值】 执行get_int16() successfully actualInt16Value:" + actualInt16Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");
            //int32
            TransactionReceipt  transactionReceipt2 = basicDataIntegerTypeContract.set_int32(int32Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证int32整型有符号/无符号类型】 successfully hash:" + transactionReceipt2.getTransactionHash());
            Int32 actualInt32Value = basicDataIntegerTypeContract.get_int32().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数int32_t取值】 执行get_int32() successfully actualInt32Value:" + actualInt32Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");
            //int64
            TransactionReceipt  transactionReceipt3 = basicDataIntegerTypeContract.set_int64(int64Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证int64整型有符号/无符号类型】 successfully hash:" + transactionReceipt3.getTransactionHash());
            Int64 actualInt64Value = basicDataIntegerTypeContract.get_int64().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数int32_t取值】 执行get_int64() successfully actualInt64Value:" + actualInt64Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //uint8
            TransactionReceipt  transactionReceipt4 = basicDataIntegerTypeContract.set_uint8(uint8Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证uint8整型有符号/无符号类型】 successfully hash:" + transactionReceipt4.getTransactionHash());
            Uint8 actualUint8Value = basicDataIntegerTypeContract.get_uint8().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数uint8_t取值】 执行get_uint8() successfully actualUint8Value:" + actualUint8Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //uint16
            TransactionReceipt  transactionReceipt5 = basicDataIntegerTypeContract.set_uint16(uint16Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证uint16整型有符号/无符号类型】 successfully hash:" + transactionReceipt5.getTransactionHash());
            Uint16 actualUint16Value = basicDataIntegerTypeContract.get_uint16().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数uint16_t取值】 执行get_uint16() successfully actualUint16Value:" + actualUint16Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //uint32
            TransactionReceipt  transactionReceipt6 = basicDataIntegerTypeContract.set_uint32(uint32Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证uint32整型有符号/无符号类型】 successfully hash:" + transactionReceipt6.getTransactionHash());
            Uint32 actualUint32Value = basicDataIntegerTypeContract.get_uint32().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数uint32_t取值】 执行get_uint32() successfully actualUint32Value:" + actualUint32Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //uint64
            TransactionReceipt  transactionReceipt7 = basicDataIntegerTypeContract.set_uint64(uint64Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证uint64整型有符号/无符号类型】 successfully hash:" + transactionReceipt7.getTransactionHash());
            Uint64 actualUint64Value = basicDataIntegerTypeContract.get_uint64().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数uint64_t取值】 执行get_uint64() successfully actualUint64Value:" + actualUint64Value.toString());
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //u160
            TransactionReceipt  transactionReceipt8 = basicDataIntegerTypeContract.set_u160(uint64Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证u160整型有符号/无符号类型】 successfully hash:" + transactionReceipt8.getTransactionHash());
            String actualU160Value = basicDataIntegerTypeContract.get_u160().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数u160取值】 执行get_u160() successfully actualU160Value:" + actualU160Value);
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //u256
            TransactionReceipt  transactionReceipt9 = basicDataIntegerTypeContract.set_u256(uint64Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证u256整型有符号/无符号类型】 successfully hash:" + transactionReceipt9.getTransactionHash());
            String actualU256Value = basicDataIntegerTypeContract.get_u256().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数u256取值】 执行get_u160() successfully actualU256Value:" + actualU256Value);
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");

            //bigint
            TransactionReceipt  transactionReceipt10 = basicDataIntegerTypeContract.set_bigint(uint64Value).send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证bigint整型有符号/无符号类型】 successfully hash:" + transactionReceipt10.getTransactionHash());
            String actualBigIntValue = basicDataIntegerTypeContract.get_bigint().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证整数bigInt取值】 执行get_bigint() successfully actualBigIntValue:" + actualBigIntValue);
            //collector.assertEqual(actualInt8Value,a.getValue(), "checkout  execute success.");


            //5、验证：整型uint8_t取值
            //Uint8 actualUint8Value = basicDataIntegerTypeContract.getStorageUint8().send();
           // collector.logStepPass("basicDataIntegerTypeContract 【验证整数uint8取值】 执行getStorageUint8() successfully actualValue:" + actualUint8Value);
           // collector.assertEqual(actualUint8Value,d, "checkout  execute success.");
            //5、验证：无符号整数位数
          /*  TransactionReceipt  transactionReceipt1 = basicDataIntegerTypeContract.setUint().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证无符号整数位数】 successfully hash:" + transactionReceipt1.getTransactionHash());
*/

            //6、验证：有符号整数位数
           /* TransactionReceipt  transactionReceipt2 = basicDataIntegerTypeContract.setInt().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证有符号整数位数】 successfully hash:" + transactionReceipt2.getTransactionHash());
            //7、验证：大位数整型赋值
            TransactionReceipt  transactionReceipt3 = basicDataIntegerTypeContract.setBigInt().send();
            collector.logStepPass("basicDataIntegerTypeContract 【验证验证大位数整型赋值】 successfully hash:" + transactionReceipt3.getTransactionHash());
*/
        } catch (Exception e) {
            collector.logStepFail("basicDataIntegerTypeContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
