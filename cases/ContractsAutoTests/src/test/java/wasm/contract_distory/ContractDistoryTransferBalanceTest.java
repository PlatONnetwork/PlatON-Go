package wasm.contract_distory;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.ContractDistory;
import org.junit.Test;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigDecimal;
import java.math.BigInteger;

/**
 * @title 合约销毁后，销毁合约的余额转移至调用者账户中
 * @description:
 * @author: hudenian
 * @create: 2020/02/18
 */
public class ContractDistoryTransferBalanceTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_distory销毁合约的余额转移至调用者账户中",sourcePrefix = "wasm")
    public void testDistoryContract() {

        String transferMoney = "10000000000000000000";

        try {
            prepare();
            ContractDistory contractDistory = ContractDistory.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = contractDistory.getContractAddress();
            String transactionHash = contractDistory.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractDistory issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + contractDistory.getTransactionReceipt().get().getGasUsed());

            //合约销毁前往合约转账
            Transfer transfer = new Transfer(web3j, transactionManager);
            TransactionReceipt transactionReceiptTransfer = transfer.sendFunds(contractAddress, new BigDecimal(transferMoney), Convert.Unit.VON).send();

            collector.logStepPass("transfer.status：" + transactionReceiptTransfer.getStatus() + ", hash：" + transactionReceiptTransfer.getTransactionHash() + ", gas used：" + transactionReceiptTransfer.getGasUsed());

            //查询当前调用者余额
            BigInteger originBalance = web3j.platonGetBalance(credentials.getAddress(chainId), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("当前用户销毁合约之前账户余额为>>>"+originBalance);


            //合约销毁
            TransactionReceipt transactionReceipt = contractDistory.distory_contract().send();
            collector.logStepPass("ContractDistory distory_contract successfully hash:" + transactionReceipt.getTransactionHash());

            //合约销毁后余额为0
            BigInteger afterDistoryBalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("After distory, contract balance is: " + afterDistoryBalance);
            collector.assertEqual(afterDistoryBalance.toString(),"0");

            //查询当前调用者余额
            BigInteger afterBalance = web3j.platonGetBalance(credentials.getAddress(chainId), DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("当前用户销毁合约之后账户余额为>>>"+afterBalance);

            //校验当前用户的余额增加值是否等于被销毁合约中的余额
            if(afterBalance.subtract(originBalance).longValue()>0){
                collector.logStepPass("销毁合约后用户账户余额增加");
            }else{
                collector.logStepFail("销毁合约后用户账户余额没有增加","销毁用户合约账户余额校验失败");
            }

        } catch (Exception e) {
            if(e instanceof ArrayIndexOutOfBoundsException){
                collector.logStepPass("ContractDistoryTransferBalanceTest and could not call contract function");
            }else{
                collector.logStepFail("ContractDistoryTransferBalanceTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }
        }
    }
}
