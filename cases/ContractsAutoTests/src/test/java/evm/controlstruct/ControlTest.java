package evm.controlstruct;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.Control;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;


/**
 * @title   控制结构
 *          1. if...else
 *          2. do...while
 *          3. for循环
 *          4. for循环包含break
 *          5. for循环包含continue
 *          6. for循环包含return
 *          7. 三目运算符
 * @description:
 * @author: hudenian
 * @create: 2019/12/30
 */
public class ControlTest extends ContractPrepareTest {

    //需要判断年龄的是否大于20的数值
    private String ifControlValue;

    //需要判断年龄的是否大于20的数值
    private String threeControlControlValue;

    @Before
    public void before() {
        this.prepare();
        ifControlValue = driverService.param.get("age");
        threeControlControlValue = driverService.param.get("age");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "ControlTest-控制结构测试", sourcePrefix = "evm")
    public void controlStructCheck() {
        try {

            Control control = Control.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = control.getContractAddress();
            TransactionReceipt tx = control.getTransactionReceipt().get();

            collector.logStepPass("ControlTest 常用控制结构功能测试");
            collector.logStepPass("deploy gas used:" + control.getTransactionReceipt().get().getGasUsed());

            //1.if控制结构验证 ifControlValue
            TransactionReceipt transactionReceipt = control.ifControl(new BigInteger(ifControlValue)).send();
            collector.logStepPass("ControlTest ifControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //if控制结构验证结果值
            String chainIfControlValue = control.getIfControlResult().send();
            collector.logStepPass( "if控制结构测试获取链上的结果是:" + chainIfControlValue);
            collector.assertEqual(ifControlCheck(ifControlValue),chainIfControlValue);

            //2.doWhile控制结构
            transactionReceipt = control.doWhileControl().send();
            collector.logStepPass("ControlTest doWhileControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //doWhile控制结构执行结果值
            String chaindoWhileValue = control.getdoWhileResult().send().toString();
            collector.logStepPass( "doWhile控制结构执行结果是:" + chaindoWhileValue);
            collector.assertEqual("45",chaindoWhileValue);//需要调整

            //3.forControl控制结构
            transactionReceipt = control.forControl().send();
            collector.logStepPass("ControlTest forControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //forControl控制结构执行结果值
            String chainGetForControlValue = control.getForControlResult().send().toString();
            collector.logStepPass( "forControl控制结构执行结果是:" + chainGetForControlValue);
            collector.assertEqual("45",chainGetForControlValue);//需要调整

            //4. for循环包含break控制结构
            transactionReceipt = control.forBreakControl().send();
            collector.logStepPass("ControlTest forBreakControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //for循环包含break控制结构执行结果值
            String chainForBreakControlValue = control.getForBreakControlResult().send().toString();
            collector.logStepPass( "for循环包含break控制结构执行结果是:" + chainForBreakControlValue);
            collector.assertEqual("1",chainForBreakControlValue);//需要调整

            //5. for循环包含continue控制结构
            transactionReceipt = control.forContinueControl().send();
            collector.logStepPass("ControlTest forContinueControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            // for循环包含continue控制结构执行结果值
            String chainForContinueControlValue = control.getForContinueControlResult().send().toString();
            collector.logStepPass( " for循环包含continue控制结构执行结果是:" + chainForContinueControlValue);
            collector.assertEqual("25",chainForContinueControlValue);//需要调整

            //6. for循环包含return控制结构
            transactionReceipt = control.forReturnControl().send();
            collector.logStepPass("ControlTest forReturnControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //for循环包含return执行结果值
            String chainForReturnControlValue = control.getForReturnControlResult().send().toString();
            collector.logStepPass( "for循环包含return执行结果是:" + chainForReturnControlValue);
            collector.assertEqual("10",chainForReturnControlValue);//需要调整

            //7.三目运算符控制结构
            transactionReceipt = control.forThreeControlControl(new BigInteger(threeControlControlValue)).send();
            collector.logStepPass("ControlTest forThreeControlControl successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass( "currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //三目运算符控制结构执行结果值
            String chainThreeControlControlValue = control.getForThreeControlControlResult().send().toString();
            collector.logStepPass( "三目运算符控制结构执行结果是:" + chainThreeControlControlValue);
            collector.assertEqual(Integer.valueOf(threeControlControlValue).intValue()> 20?"less than 20":"more than 20",chainThreeControlControlValue);//需要调整



        } catch (Exception e) {
            collector.logStepFail("ControlTest controlStructCheck process fail.", e.toString());
            e.printStackTrace();
        }
    }

    public static String ifControlCheck(String ifControlValue){
        int age = Integer.valueOf(ifControlValue);
        String ifControlResult = "";
        if(age < 20){
            ifControlResult = "you are a young man";
        }else if (age < 60){
            ifControlResult = "you are a middle man";
        }else {
            ifControlResult = "you are a old man";
        }
        return ifControlResult;
    }

}
