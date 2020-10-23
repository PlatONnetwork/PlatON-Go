package network.platon.test.evm.complexcontracts;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.Guessing;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * @title 竞猜合约测试
 * @description:
 * @author: hudenian
 * @create: 2020/03/04 16:42
 **/

public class GuessingTest extends ContractPrepareTest {
    private String endBlock = "12823"; //设置竞猜截止块高

    @Before
    public void before() {
        this.prepare();
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "function.GuessingTest-竞猜合约测试", sourcePrefix = "evm")
    public void guessingTest() {

        try {

            Guessing guessing = Guessing.deploy(web3j, transactionManager, provider, chainId, new BigInteger(endBlock)).send();

            String contractAddress = guessing.getContractAddress();
            TransactionReceipt tx = guessing.getTransactionReceipt().get();
            collector.logStepPass("GuessingTest deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("GuessingTest deploy gasUsed:" + guessing.getTransactionReceipt().get().getGasUsed());

            //查询截止块高
            endBlock = guessing.endBlock().send().toString();
            collector.logStepPass("查询截止块高为：" + endBlock);

            //查询合约余额
            String contractBalance = guessing.getBalanceOf().send().toString();
            collector.assertEqual("0", contractBalance);


            //查询总奖池
            String balance = guessing.balance().send().toString();
            collector.logStepPass("发起竞猜前奖池总金额为：" + balance);

            //发起转账(触发竞猜操作)
            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt transactionReceipt = transfer.sendFunds(contractAddress, new BigDecimal("1000"), Convert.Unit.LAT, new BigInteger("1000000000"), new BigInteger("4712388")).send();
            collector.logStepPass("gas used>>>>>>>" + transactionReceipt.getGasUsed().toString());

            //查询合约余额
            contractBalance = guessing.getBalanceOf().send().toString();

            balance = guessing.balance().send().toString();
            collector.logStepPass("发起第一次竞猜前奖池总金额为：" + balance);

            //发起竞猜(客户端发起的单位是von与 lat差10^18次方)
            tx = guessing.guessingWithLat(new BigInteger("1000000000000000000000")).send();
            collector.logStepPass(tx.getLogs().toString());
            collector.logStepPass(" gas used>>>>>>>" + transactionReceipt.getGasUsed());

            //自增序列下标
            String indexKey = guessing.indexKey().send().toString();
            collector.logStepPass("自增序列下标为：" + indexKey);

            //查询总奖池
            balance = guessing.balance().send().toString();
            collector.logStepPass("发起第二次竞猜后奖池总金额为：" + balance);

            contractBalance = guessing.getBalanceOf().send().toString();

            //自增序列下标
            indexKey = guessing.indexKey().send().toString();
            collector.logStepPass("自增序列下标为：" + indexKey);

            //查询当前调用者余额
            BigInteger originBalance = web3j.platonGetBalance(credentials.getAddress(chainId), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("开奖前用户账户余额为>>>" + originBalance);

            //查询当前块高
            long currentBlockNumber = new Long(web3j.platonBlockNumber().send().getBlockNumber().toString()).intValue();

            //入参为截止块高的hash，真实环境需要从浏览器获取，此处为了测试方便直接从合约取当前块高hash
            byte[] blockhash = guessing.generateBlockHash(new BigInteger(String.valueOf(currentBlockNumber - 20))).send();


            //开奖操作
            tx = guessing.draw(blockhash).send();
            collector.logStepPass(tx.getLogs().toString());


            //查询当前调用者余额
            BigInteger afterBalance = web3j.platonGetBalance(credentials.getAddress(chainId), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("开奖后用户账户余额为>>>" + afterBalance);


        } catch (Exception e) {
            collector.logStepFail("GuessingTest Calling Method fail.", e.toString());
            e.printStackTrace();
        }
    }
}
