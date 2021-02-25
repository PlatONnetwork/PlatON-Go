package network.platon.test.wasm.contract_create;

import com.platon.rlp.datatypes.Uint16;
import com.platon.rlp.datatypes.Uint64;
import com.platon.rlp.datatypes.Uint8;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.InitWithVector;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import network.platon.test.wasm.beforetest.WASMContractPrepareTest;

/**
 * @title 创建合约init包含vector测试
 * @description:
 * @author: hudenian
 * @create: 2020/02/16
 */
public class InitWithVectorTest extends WASMContractPrepareTest {
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_create创建合约init带vector",sourcePrefix = "wasm")
    public void testNewContract() {

        Uint16 age = Uint16.of(20);

        //vector中要添加的元素
        String vector1 = "vector1";
        String vector2 = "vector2";
        String vector3 = "vector3";
        String vector4 = "vector4";
        String vector5 = "vector5";

        try {
            prepare();
            InitWithVector initWithVector = InitWithVector.deploy(web3j, transactionManager, provider, chainId,age).send();
            String contractAddress = initWithVector.getContractAddress();
            String transactionHash = initWithVector.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("InitWithVector issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            collector.logStepPass("deploy gas used:" + initWithVector.getTransactionReceipt().get().getGasUsed());

            Uint8 idx = Uint8.of(0);
            Uint64 chainAge = initWithVector.get_vector(idx).send();
            collector.assertEqual(chainAge.value,age.value);

            //vctor添加第一个元素
            TransactionReceipt tx = initWithVector.vector_push_back_element(vector1).send();
            collector.logStepPass("InitWithVectorTest call vector_push_back_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //vctor添加第二、三个元素
            tx = initWithVector.vector_push_back_element(vector2).send();
            collector.logStepPass("InitWithVectorTest call vector_push_back_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());
            tx = initWithVector.vector_push_back_element(vector3).send();
            collector.logStepPass("InitWithVectorTest call vector_push_back_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看vector中元素个数
            Uint8 vectorSize = initWithVector.get_strvector_size().send();
            collector.assertEqual("3",vectorSize.value.toString());

            //获取第二个元素(下标为1)
            String chainElement = initWithVector.get_vector_element_by_position(Uint8.of("1")).send();
            collector.assertEqual(vector2,chainElement);

            //去掉最后一个元素
            tx = initWithVector.vector_pop_back_element().send();
            collector.logStepPass("InitWithVectorTest call vector_pop_back_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看vector中元素个数
            vectorSize = initWithVector.get_strvector_size().send();
            collector.assertEqual("2",vectorSize.value.toString());

            //在指定位置添加元素
            tx = initWithVector.vector_insert_element(vector4,Uint8.of("2")).send();
            collector.logStepPass("InitWithVectorTest call vector_insert_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看vector中元素个数
            vectorSize = initWithVector.get_strvector_size().send();
            collector.assertEqual("3",vectorSize.value.toString());

            //在指定位置添加元素(下标超过数组大小，则插入到最后面)
            tx = initWithVector.vector_insert_element(vector4,Uint8.of("6")).send();
            collector.logStepPass("InitWithVectorTest call vector_insert_element successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查看vector中元素个数
            vectorSize = initWithVector.get_strvector_size().send();
            collector.assertEqual("4",vectorSize.value.toString());

            //验证for循环
            tx = initWithVector.vectorFor(new String[]{vector1,vector2,vector3,vector4,vector5}).send();
            collector.logStepPass("InitWithVectorTest call vectorFor successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //验证switch语句
            tx = initWithVector.vectorCase(new String[]{vector1,vector2,vector3,vector4,vector5}).send();
            collector.logStepPass("InitWithVectorTest call vectorCase successfully.contractAddress:" + contractAddress + ", hash:" + tx.getTransactionHash());

            //查询switch调用结果
            String caseResult = initWithVector.get_vectorCase_result().send();
            collector.assertEqual("5",caseResult);

        } catch (Exception e) {
            collector.logStepFail("InitWithVectorTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
