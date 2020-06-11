//package evm.regreetest;
//
//import com.alibaba.fastjson.JSONArray;
//import network.platon.autotest.junit.annotations.DataSource;
//import network.platon.autotest.junit.enums.DataSourceType;
//import network.platon.autotest.junit.rules.AssertCollector;
//import network.platon.autotest.junit.rules.DriverService;
//import network.platon.autotest.utils.FileUtil;
//import network.platon.utils.OneselfFileUtil;
//import org.junit.Before;
//import org.junit.Rule;
//import org.junit.Test;
//import org.web3j.crypto.Credentials;
//import org.web3j.protocol.Web3j;
//import org.web3j.protocol.core.DefaultBlockParameterName;
//import org.web3j.protocol.core.methods.response.TransactionReceipt;
//import org.web3j.protocol.http.HttpService;
//import org.web3j.tx.RawTransactionManager;
//import org.web3j.tx.Transfer;
//import org.web3j.utils.Convert;
//
//import java.math.BigDecimal;
//import java.math.BigInteger;
//import java.nio.file.Paths;
//
//
///**
// * @title PlatON普通有回执转账交易
// * @description: 步骤：账户A向账户B转账amount，预期1：账户A的余额减少amount，预期2：账户B的余额增加amount
// * @author: qcxiao
// * @create: 2019/12/16 11:03
// **/
//public class PlatONTransfer2Test {
//
//    @Rule
//    public DriverService driverService = new DriverService();
//
//    @Rule
//    public AssertCollector collector = new AssertCollector();
//
//    // 底层链ID
//    private long chainId;
//    // 发行代币的地址
//    private String transferFrom;
//    // 接收代币的地址
//    private String transferTo;
//    private Web3j web3j;
//    // 转账的金额
//    private String amount;
//
//    @Before
//    public void before() {
//        chainId = Integer.valueOf(driverService.param.get("chainId"));
//        //transferFrom = driverService.param.get("transferFrom");
//        transferFrom = driverService.param.get("address");
//        transferTo = driverService.param.get("transferTo");
//        amount = driverService.param.get("amount");
//    }
//
//    @Test
//    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
//            author = "qcxiao", showName = "evm.PlatonTransferTest-普通有回执转账交易", sourcePrefix = "evm")
//    public void testTransfer() {
//        Credentials credentials = null;
//        BigInteger nonce = null;
//        try {
//            web3j = Web3j.build(new HttpService(driverService.param.get("nodeUrl")));
//            credentials = Credentials.create(driverService.param.get("privateKey"));
//
//
//            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "lax_bech32_all_addr_and_private_keys.json").toUri().getPath());
//            String jsonContent = OneselfFileUtil.readFile(filePath);
//
//            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
//            for (int i = 0; i < jsonArray.size(); i++) {
//                Thread.sleep(1000); // 1680
//                transferTo = jsonArray.getJSONObject(i).getString("address");
//                try {
//                    RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//                    Transfer transfer = new Transfer(web3j, transactionManager);
//
//                    amount = "100000";
//                    //TransactionReceipt transactionReceipt = transfer.sendFunds(transferTo, new BigDecimal(amount), Convert.Unit.LAT).send();
//
//                    transfer.sendFunds(transferTo, new BigDecimal(amount), Convert.Unit.LAT).sendAsync();
//
//                    //BigInteger endBalance = web3j.platonGetBalance(transferTo, DefaultBlockParameterName.LATEST).send().getBalance();
//                    //collector.logStepPass("transferTo:" + transferTo + ",endBalance:" + endBalance);
//                    collector.logStepPass("time:" + i);
//                } catch (Exception e) {
//                    e.printStackTrace();
//                }
//
//            }
//
//        } catch (Exception e) {
//            collector.logStepFail("transfer fail.", e.toString());
//            e.printStackTrace();
//        }
//
//    }
//
//}
