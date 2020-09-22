package network.platon.test.wasm.data_type;

import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ReferenceDataTypeComplexContract;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 测试合约引用类型复杂合约
 * @description:
 * @author: qudong
 * @create: 2020/02/07
 */
public class ReferenceDataTypeComplexContractTest extends WASMContractPrepareTest {

    private String sIdStr;
    private String nameStr;
    private String ageStr;

    @Before
    public void before() {
        sIdStr = driverService.param.get("sIdStr");
        nameStr = driverService.param.get("nameStr");
        ageStr = driverService.param.get("ageStr");
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qudong", showName = "wasm.referenceDataTypeComplexContractTest合约引用类型复杂合约",sourcePrefix = "wasm")
    public void testReferenceDataTypeComplexContract() {

         //部署合约
        ReferenceDataTypeComplexContract referenceDataTypeComplexContract = null;
        try {
            prepare();
            referenceDataTypeComplexContract = ReferenceDataTypeComplexContract.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = referenceDataTypeComplexContract.getContractAddress();
            TransactionReceipt tx = referenceDataTypeComplexContract.getTransactionReceipt().get();
            collector.logStepPass("referenceDataTypeComplexContract issued successfully.contractAddress:" + contractAddress
                                  + ", hash:" + tx.getTransactionHash() + ",deploy gas used:" + tx.getGasUsed());
            collector.logStepPass("deployFinishCurrentBlockNumber:" + tx.getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeComplexContract deploy fail.", e.toString());
            e.printStackTrace();
        }
        //调用合约方法
        try {
            //1、验证：学生赋值基本信息
            Uint64 sId = Uint64.of(sIdStr);
            String name = nameStr;
            Uint8 age = Uint8.of(ageStr);
            Boolean sex = false;
            TransactionReceipt  transactionReceipt = referenceDataTypeComplexContract.set_student_info(sId,name,age,sex).send();
            collector.logStepPass("referenceDataTypeComplexContract 【验证学生赋值基本信息】 successfully hash:" + transactionReceipt.getTransactionHash());
            //学生信息取值
            Uint64 actualValueSid = referenceDataTypeComplexContract.get_student_id().send();
            collector.logStepPass("referenceDataTypeComplexContract 【取值学生编号】 执行get_student_id() successfully actualValueSid:" + actualValueSid);
           // collector.assertEqual(actualValue,expectIndexValue, "checkout  execute success.");
            String actualValueName = referenceDataTypeComplexContract.get_name().send();
            collector.logStepPass("referenceDataTypeComplexContract 【取值学生姓名】 执行get_name() successfully actualValueName:" + actualValueName);
            // collector.assertEqual(actualValue,expectIndexValue, "checkout  execute success.");
            Uint8 actualValueAge = referenceDataTypeComplexContract.get_age().send();
            collector.logStepPass("referenceDataTypeComplexContract 【取值学生年龄】 执行get_age() successfully actualValueAge:" + actualValueAge);
            // collector.assertEqual(actualValue,expectIndexValue, "checkout  execute success.");
            Boolean actualValueSex = referenceDataTypeComplexContract.get_sex().send();
            collector.logStepPass("referenceDataTypeComplexContract 【取值学生性别】 执行get_sex() successfully actualValueSex:" + actualValueSex);
            // collector.assertEqual(actualValue,expectIndexValue, "checkout  execute success.");

          /*  String[] arrayCourse = {"English","Mathematics"};
            TransactionReceipt  transactionReceipt1 = referenceDataTypeComplexContract.set_array_course(arrayCourse).send();
            collector.logStepPass("referenceDataTypeComplexContract 【验证数组赋值课程】 successfully hash:" + transactionReceipt1.getTransactionHash());
            *///2、验证：数组课程取值
            String[] courseArr = referenceDataTypeComplexContract.get_array_course().send();
            collector.logStepPass("referenceDataTypeComplexContract 【验证数组课程取值】 执行get_array_course() successfully courseArr.length:" + courseArr.length);
            // collector.assertEqual(actualValue,expectIndexValue, "checkout  execute success.");



        } catch (Exception e) {
            collector.logStepFail("referenceDataTypeComplexContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }

    }
}
