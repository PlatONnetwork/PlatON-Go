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
//            //credentials = Credentials.create(driverService.param.get("privateKeyOfTransferFrom"));
//
//
//            String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "all_addr_and_private_keys.json").toUri().getPath());
//            String jsonContent = OneselfFileUtil.readFile(filePath);
//
//            JSONArray jsonArray = JSONArray.parseArray(jsonContent);
//            for (int i = 0; i < jsonArray.size(); i++) {
//                Thread.sleep(100); // 7928
//                // transferTo = jsonArray.getJSONObject(i).getString("address");
//                try {
//                    transferTo = driverService.param.get("address");
//
//                    credentials = Credentials.create(jsonArray.getJSONObject(i).getString("private_key"));
//                    RawTransactionManager transactionManager = new RawTransactionManager(web3j, credentials, chainId);
//                    Transfer transfer = new Transfer(web3j, transactionManager);
//
//                    amount = "900000";
//                    TransactionReceipt transactionReceipt = transfer.sendFunds(transferTo, new BigDecimal(amount), Convert.Unit.LAT).send();
//
//                    BigInteger endBalance = web3j.platonGetBalance(transferTo, DefaultBlockParameterName.LATEST).send().getBalance();
//                    collector.logStepPass("transferTo:" + transferTo + ",endBalance:" + endBalance);
//                    collector.logStepPass("time:" + i + ", txHash:" + transactionReceipt.getTransactionHash());
//                } catch (Exception e) {
//
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
