package evm.csdccontracts;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.OrderDao;
import org.junit.Before;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.util.Date;


/**
 * @title 结算质押申请合约验证测试
 * @description:
 * @author: hudenian
 * @create: 2020/1/9
 */
public class OrderDaoTest extends ContractPrepareTest {


    //模拟简单的业务数据
    private String secApply = "2-businessNo1-bizId1-4-5-6-7-8-9-10-11-12-13-14-15-16-17-18-19-20-21-22-23-24-25-26-27-28-29-30-31-32-33-34-35-36-37-38-39-40";

    private String bizId = "2";

    @Before
    public void before() {
        this.prepare();
    }


    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "OrderDaoTest-结算复杂合约测试验证", sourcePrefix = "evm")
    public void orderDaoTest() {
        try {
            //调用者合约地址
            OrderDao orderDao = OrderDao.deploy(web3j, transactionManager, provider, chainId).send();
            String callerContractAddress = orderDao.getContractAddress();
            TransactionReceipt tx = orderDao.getTransactionReceipt().get();
            collector.logStepPass("OrderDao deploy successfully.contractAddress:" + callerContractAddress + ", hash:" + tx.getTransactionHash());
            collector.logStepPass("deploy gas used:" + orderDao.getTransactionReceipt().get().getGasUsed());

            Date start = new Date();


            tx = orderDao.insert_SecPledgeApply(secApply).send();
            //插入业务数据
            collector.logStepPass("OrderDao callDoublelTest successfully hash:" + tx.getTransactionHash());

            //根据业务id查询业务数据
           String business_id =  orderDao.select_SecPledgeApply_byId(bizId).send().toString();
           collector.logStepPass("bizId:"+bizId+"对应的business_no为："+business_id);

            Date end = new Date();
            collector.logStepPass("插入到调用一共耗时："+(end.getTime()-start.getTime()));


        } catch (Exception e) {
            collector.logStepFail("OrderDaoTest process fail.", e.toString());
            e.printStackTrace();
        }
    }

}
