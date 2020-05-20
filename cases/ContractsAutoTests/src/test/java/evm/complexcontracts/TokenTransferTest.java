package evm.complexcontracts;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.junit.rules.AssertCollector;
import network.platon.autotest.junit.rules.DriverService;
import network.platon.contracts.HumanStandardToken;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.math.BigInteger;


/**
 * @title 代币转移
 * @description:
 * @author: qcxiao
 * @create: 2019/12/16 13:39
 **/
public class TokenTransferTest {

    @Rule
    public DriverService driverService = new DriverService();

    @Rule
    public AssertCollector collector = new AssertCollector();

    // 底层链ID
    private long chainId;
    // 每次转移的代币数量
    private String transferAmount;
    // 发行代币的总额
    private String ownerAmount;
    // 接收代币的地址
    private final static String transferTo = "lax1354ckckjla0869lernuzrjh7arslu3vypxek6h";
    // 代币名称
    private String tokenName;
    // 发行代币的地址
    private String transferFrom;


    @Before
    public void before() {
        chainId = Integer.valueOf(driverService.param.get("chainId"));
        ownerAmount = driverService.param.get("ownerAmount");
        transferAmount = driverService.param.get("transferAmount");
        tokenName = driverService.param.get("tokenName");
        transferFrom = driverService.param.get("address");
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "qcxiao", showName = "complexcontracts.TokenTransferTest-代币转移", sourcePrefix = "evm")
    public void testTransfer() {
        Web3j web3j = null;
        Credentials credentials = null;
        try {
            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
            credentials = Credentials.create(driverService.param.get("privateKey"));
            collector.logStepPass("currentBlockNumber:" + web3j.platonBlockNumber().send().getBlockNumber());
        } catch (Exception e) {
            collector.logStepFail("The node is unable to connect", e.toString());
            e.printStackTrace();
        }


        ContractGasProvider provider = new ContractGasProvider(new BigInteger("5000000000000"), new BigInteger("3000000"));
        RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);

        try {
            HumanStandardToken token = HumanStandardToken.deploy(web3j, transactionManager, provider, chainId,
                    new BigInteger(ownerAmount), tokenName, BigInteger.valueOf(18), "USDT").send();

            String contractAddress = token.getContractAddress();
            TransactionReceipt tx = token.getTransactionReceipt().get();

            collector.logStepPass("Token issued successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash() +
                    ", tokenName:" + token.name().send() + ", symbol:" + token.symbol().send());
            collector.logStepPass("deploy gas used:" + token.getTransactionReceipt().get().getGasUsed());
            collector.assertEqual(tokenName, token.name().send(), "checkout tokenName");

            collector.logStepPass("5次循环调用...");
            for (int i = 1; i < 6; i++) {
                TransactionReceipt transactionReceipt = HumanStandardToken.load(contractAddress, web3j, transactionManager, provider, chainId)
                        .transfer(transferTo, new BigInteger(transferAmount)).send();
                BigInteger toBalance = token.balanceOf(transferTo).send();
                BigInteger fromBalance = token.balanceOf(transferFrom).send();
                collector.logStepPass("Token transfer successful.transactionHash:" + transactionReceipt.getTransactionHash());
                collector.logStepPass("currentTransferTime:" + i + ", currentBlockNumber:" + transactionReceipt.getBlockNumber());
                collector.logStepPass("transferToBalance:" + toBalance + ", transferFromBalance:" + fromBalance);
                // 累计转移的数量
                BigInteger amount = new BigInteger(transferAmount).multiply(BigInteger.valueOf(i));
                // 判断代币接收地址的余额是否正确
                collector.assertEqual(toBalance, amount, "checkout every time transferTo balance.");
                // 判断代币转出地址余额是否正确
                collector.assertEqual(fromBalance, (new BigInteger(ownerAmount)).subtract(amount), "checkout every time transferFrom balance.");
            }
        } catch (Exception e) {
            e.printStackTrace();
        }

    }

}
