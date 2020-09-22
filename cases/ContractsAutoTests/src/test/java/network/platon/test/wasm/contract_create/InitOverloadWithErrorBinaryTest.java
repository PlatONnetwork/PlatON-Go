package network.platon.test.wasm.contract_create;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitOverload;
import org.junit.Test;
import org.web3j.protocol.exceptions.TransactionException;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 使用异常编译文件部署合约
 * @description:
 * @author: hudenian
 * @create: 2020/02/19
 */
public class InitOverloadWithErrorBinaryTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create使用异常编译文件部署合约",sourcePrefix = "wasm")
    public void testNewContractWithErrorBinary() {

        String name = "hudenian";
        try {
            prepare();
            //截掉binary前面5位数字
            InitOverload.BINARY = InitOverload.BINARY.substring(5);

            InitOverload initOverload = InitOverload.deploy(web3j, transactionManager, provider, chainId).send();
            String contractAddress = initOverload.getContractAddress();
            String transactionHash = initOverload.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitOverload issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

        } catch (Exception e) {
            if(e instanceof TransactionException){
                collector.logStepPass("合约binary文件错误，预期部署不成功");
            }else{
                collector.logStepFail("InitOverloadTest failure,exception msg:" , e.getMessage());
                e.printStackTrace();
            }

        }
    }
}
