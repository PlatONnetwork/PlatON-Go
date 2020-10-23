package network.platon.test.evm.function.specialVariablesAndFunctions;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.BlockTransactionPropertiesFunctions;
import network.platon.utils.DataChangeUtil;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigInteger;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

/**
 * @title 验证区块和交易函数
 * @description:
 * @author: liweic
 * @create: 2019/12/31 11:01
 **/

public class BlockTransactionPropertiesFunctionsTest extends ContractPrepareTest {

    private String coinbase;
//    private String gaslimit;
    private String msgdata;
    private String gasleft;
    private String msgsig;
    private String gasprice;

    @Before
    public void before() {
        this.prepare();
        coinbase = driverService.param.get("coinbase");
//        gaslimit = driverService.param.get("gaslimit");
        msgdata = driverService.param.get("msgdata");
        gasleft = driverService.param.get("gasleft");
        msgsig = driverService.param.get("msgsig");
        gasprice = driverService.param.get("gasprice");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "function.BlockTransactionPropertiesFunctionsTest-区块和交易函数测试", sourcePrefix = "evm")
    public void BlockTransactionPropertiesfunction() {
        try {
            BlockTransactionPropertiesFunctions blockTransactionPropertiesFunctions = BlockTransactionPropertiesFunctions.deploy(web3j, transactionManager, provider, chainId).send();

            String contractAddress = blockTransactionPropertiesFunctions.getContractAddress();
            TransactionReceipt tx = blockTransactionPropertiesFunctions.getTransactionReceipt().get();
            collector.logStepPass("BlockTransactionPropertiesFunctionsTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("BlockTransactionPropertiesFunctionsTest deploy gasUsed:" + blockTransactionPropertiesFunctions.getTransactionReceipt().get().getGasUsed());

            //验证block.number函数(获取块高)
            BigInteger PlatONBlocknumber = web3j.platonBlockNumber().send().getBlockNumber();
            collector.logStepPass("web3j拿到的区块高度：" + PlatONBlocknumber);

            BigInteger blocknumber = blockTransactionPropertiesFunctions.getBlockNumber().send();
            collector.logStepPass("block.number函数返回值：" + blocknumber);
            int a = Integer.valueOf(blocknumber.toString());
            int b = Integer.valueOf(PlatONBlocknumber.toString());
            int blocknumdiff = a - b;
            List<Integer> list = new ArrayList<Integer>();
            list.add(new Integer(0));
            list.add(new Integer(1));
            list.add(new Integer(2));
            Boolean bndf = list.contains(blocknumdiff);
            collector.assertEqual(true, bndf);

            //验证blockhash(blockNumber)函数(获取区块Hash)
            String blocknumberNow = web3j.platonBlockNumber().send().getBlockNumber().toString();
            int number = Integer.valueOf(blocknumberNow).intValue()-100;
            byte[] resultB = blockTransactionPropertiesFunctions.getBlockhash(new BigInteger(String.valueOf(number))).send();
            String hexValue = DataChangeUtil.bytesToHex(resultB);
            collector.logStepPass("blockhash(blockNumber)函数返回值：" + hexValue);
            String errorhash = "0000000000000000000000000000000000000000000000000000000000000000";
            boolean isBool =  hexValue.equals(errorhash);
            collector.assertTrue(!isBool,"success");

            //验证block.coinbase函数(获取矿工地址)
            String resultC = blockTransactionPropertiesFunctions.getBlockCoinbase().send();
            collector.logStepPass("block.coinbase函数返回值：" + resultC);
            String coinaddr = "0x0000000000000000000000000000000000000000";
            boolean iscoinbase = resultC.equals(coinaddr);
            collector.assertEqual(!iscoinbase ,true);

            //验证block.difficulty(获取当前块的难度)
            BigInteger resultD = blockTransactionPropertiesFunctions.getBlockDifficulty().send();
            collector.logStepPass("block.difficulty函数返回值：" + resultD);
            collector.assertEqual("0" ,resultD.toString());

            //验证block.gaslimit(获取当前区块的gas限额)
            BigInteger resultE = blockTransactionPropertiesFunctions.getGaslimit().send();
            collector.logStepPass("block.gaslimit函数返回值：" + resultE);
            boolean gas = "0".toString().equals(resultE);
            collector.assertEqual(false ,gas);

            //验证block.timestamp(获取当前区块的UNIX时间戳)
            BigInteger resultF = blockTransactionPropertiesFunctions.getBlockTimestamp().send();
            collector.logStepPass("block.timestamp函数返回值：" + resultF);
            SimpleDateFormat sdf=new SimpleDateFormat("yyyy-MM-dd");
            String resultTime = sdf.format(new Date(Long.parseLong(String.valueOf(resultF))));
            String today = sdf.format(new Date());
            collector.assertEqual(today,resultTime);


            //验证msg.data(获取完整的calldata)
            byte[] resultG = blockTransactionPropertiesFunctions.getData().send();
            String hexvalue2 = DataChangeUtil.bytesToHex(resultG);
            collector.logStepPass("msg.data函数返回值：" + hexvalue2);
            collector.assertEqual(msgdata ,hexvalue2);

            //验证gasleft()(剩余的gas)
            BigInteger resultH = blockTransactionPropertiesFunctions.getGasleft().send();
            collector.logStepPass("gasleft函数返回值：" + resultH);
            collector.assertEqual("9223372036854754307" ,resultH.toString());

            //验证msg.sender(获取消息发送者（当前调用))
            String resultI = blockTransactionPropertiesFunctions.getSender().send();
            collector.logStepPass("msg.sender函数返回值：" + resultI);
            collector.assertEqual(walletAddress.toLowerCase() ,resultI);

            //验证msg.sig(calldata 的前 4 字节(也就是函数标识符))
            byte[] resultJ = blockTransactionPropertiesFunctions.getSig().send();
            String hexvalue3 = DataChangeUtil.bytesToHex(resultJ);
            collector.logStepPass("msg.sig函数返回值：" + hexvalue3);
            collector.assertEqual(msgsig ,hexvalue3);

            //验证msg.value(随消息发送的以太币的数量)
            TransactionReceipt transactionReceipt = blockTransactionPropertiesFunctions.getValue(new BigInteger("2")).send();
            String status = transactionReceipt.getStatus();
            collector.logStepPass("msg.value的transactionReceipt是：" + transactionReceipt);
            collector.logStepPass("msg.value的status是：" + status);
            collector.assertEqual("0x1" ,status);

            //验证now(目前区块时间戳)
            BigInteger resultK = blockTransactionPropertiesFunctions.getNow().send();
            collector.logStepPass("now函数返回值：" + resultK);
            SimpleDateFormat sdf2=new SimpleDateFormat("yyyy-MM-dd");
            String resultTime2 = sdf2.format(new Date(Long.parseLong(String.valueOf(resultF))));
            String now = sdf2.format(new Date());
            collector.assertEqual(now,resultTime2);

            //验证tx.gasprice(交易的 gas 价格)
            BigInteger resultL = blockTransactionPropertiesFunctions.getGasprice().send();
            collector.logStepPass("tx.gasprice函数返回值：" + resultL);
            collector.assertEqual("0" ,resultL.toString());

            //验证tx.origin(交易发起者(完全的调用链))
            String resultM = blockTransactionPropertiesFunctions.getOrigin().send();
            collector.logStepPass("tx.origin函数返回值：" + resultM);
            collector.assertEqual(walletAddress.toLowerCase() ,resultM);

        } catch (Exception e) {
            collector.logStepFail("BlockTransactionPropertiesFunctionsContract Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
