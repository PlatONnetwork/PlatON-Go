package network.platon.test.evm.complexcontracts.erc20_v_0_4_12;

import network.platon.test.evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.evm.ERC200412Token;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.math.BigDecimal;
import java.math.BigInteger;


/**
 * @title ERC200412Token测试
 * @description:
 * @author: albedo
 * @create: 2019/12/28
 */
public class ERC200412TokenTest extends ContractPrepareTest {

    //供应份额
    private String initialSupply;

    //代币名称
    private String tokenName;

    //代币简称
    private String tokenSymbol;

    //转出账号
    private String to;

    //转出金额
    private String value;

    //设置approveAddress可以创建交易者名义花费的代币数
    private String approveAddress;

    //设置approveAddress可以创建交易者名义花费的代币数approveValue
    private String approveValue;

    //创建者销毁的代币数
    private String burnValue;


    @Before
    public void before() {
        this.prepare();
        initialSupply = driverService.param.get("initialSupply");
        tokenName = driverService.param.get("tokenName");
        tokenSymbol = driverService.param.get("tokenSymbol");
        to = driverService.param.get("to");
        value = driverService.param.get("value");
        approveAddress = driverService.param.get("approveAddress");
        approveValue = driverService.param.get("approveValue");
        burnValue = driverService.param.get("burnValue");

    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "albedo", showName = "ERC200412TokenTest-测试0.4.12版本ERC20", sourcePrefix = "evm")
    public void erc20Test() {
        try {

            ERC200412Token eRC200412Token = ERC200412Token.deploy(web3j, transactionManager, provider, chainId, new BigInteger(initialSupply), tokenName, tokenSymbol).send();

            String contractAddress = eRC200412Token.getContractAddress();
            TransactionReceipt tx = eRC200412Token.getTransactionReceipt().get();

            collector.logStepPass("ERC200412Token deploy successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + eRC200412Token.getTransactionReceipt().get().getGasUsed());

            //获取代币名称
            String chainTokenName = eRC200412Token.name().send().toString();
            collector.logStepPass("开始获取代币名称" + chainTokenName);
            collector.assertEqual(tokenName, chainTokenName);

            //获取代币简称
            String chainTokenSymbol = eRC200412Token.symbol().send().toString();
            collector.logStepPass("开始获取代币简称" + chainTokenSymbol);
            collector.assertEqual(tokenSymbol, chainTokenSymbol);

            //获取代币总发行量
            String chainInitialSupply = eRC200412Token.totalSupply().send().toString();
            collector.logStepPass("开始获取代币发行量" + chainInitialSupply);
            collector.assertEqual(initialSupply + "000000000000000000", chainInitialSupply);

            //给to地址转账value值
            TransactionReceipt transactionReceipt = eRC200412Token.transfer(to, new BigInteger(value)).send();
            collector.logStepPass("发行账户向" + to + "转账" + value);

            collector.logStepPass("ERC200412Token transfer to " + to + " successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //查询to地址的余额
            String to_balance = eRC200412Token.balanceOf(to).send().toString();
            collector.logStepPass(to + "账户余额为：" + to_balance);
            collector.assertEqual(to_balance, value);

            //查询from地址的余额
            String from_balance = eRC200412Token.balanceOf(walletAddress).send().toString();
            collector.logStepPass("转账后发行账户" + walletAddress + "余额为：" + from_balance);
            collector.assertEqual(from_balance, new BigDecimal(chainInitialSupply).subtract(new BigDecimal(value)).toString());

            //创建者设置approveAddress可以创建交易者名义花费的代币数为approveValue
            transactionReceipt = eRC200412Token.approve(approveAddress, new BigInteger(approveValue)).send();
            collector.logStepPass("发行账户允许" + approveAddress + "账户以自己的名义花费的代币数为：" + approveValue);

            collector.logStepPass("ERC200412Token transfer to " + to + " successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());

            //查询approveAddress可以从我的地址转出我代币数
            String chainApproveValue = eRC200412Token.allowance(walletAddress, approveAddress).send().toString();
            collector.logStepPass("发行账户允许" + approveAddress + "账户以自己的名义花费的代币数为：" + chainApproveValue);
            collector.assertEqual(approveValue, chainApproveValue);

            //创建者销毁自己账户的指定的代币数
            transactionReceipt = eRC200412Token.burn(new BigInteger(burnValue)).send();
            collector.logStepPass("发行账销毁的代币数：" + burnValue);

            collector.logStepPass("ERC200412Token burn " + burnValue + " successful.transactionHash:" + transactionReceipt.getTransactionHash());
            collector.logStepPass("currentBlockNumber:" + transactionReceipt.getBlockNumber());


            //查询from地址的余额
            String from_balance_afterBurn = eRC200412Token.balanceOf(walletAddress).send().toString();
            collector.logStepPass("发行都销毁" + burnValue + "后账户余额为" + from_balance_afterBurn);

            //销毁指定代币后获取总发行量
            String chainInitialSupply_after_burn = eRC200412Token.totalSupply().send().toString();
            collector.logStepPass("销毁指定代币后获取总发行量：" + chainInitialSupply_after_burn);


        } catch (Exception e) {
            collector.logStepFail("erc20Test failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

}
