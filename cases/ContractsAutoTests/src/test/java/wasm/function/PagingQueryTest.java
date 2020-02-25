package wasm.function;

import com.platon.rlp.datatypes.Uint64;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.wasm.PagingQuery;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import wasm.beforetest.WASMContractPrepareTest;

import java.math.BigInteger;

/**
 * vector分页查询
 * @create: 2020/02/20
 * @author liweic
 */

public class PagingQueryTest extends WASMContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "liweic", showName = "wasm.PagingQuery验证vector分页查询",sourcePrefix = "wasm")
    public void Pagingquery() {

        try {
            prepare();
            PagingQuery pagingquery = PagingQuery.deploy(web3j, transactionManager, provider).send();
            String contractAddress = pagingquery.getContractAddress();
            String transactionHash = pagingquery.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("PagingQuery issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            TransactionReceipt str1 = pagingquery.insertVectorValue("a").send();
            TransactionReceipt str2 = pagingquery.insertVectorValue("b").send();
            TransactionReceipt str3 = pagingquery.insertVectorValue("c").send();
            TransactionReceipt str4 = pagingquery.insertVectorValue("d").send();
            TransactionReceipt str5 = pagingquery.insertVectorValue("e").send();

            Uint64 vecsize = pagingquery.getVectorSize().send();
            collector.logStepPass("vector长度:" + vecsize.value);
            collector.assertEqual(vecsize.value, new BigInteger("5"));

            String result = pagingquery.getPagingQuery(Uint64.of("2"), Uint64.of("3")).send();
            collector.logStepPass("PagingQuery结果为:" + result);
            collector.assertEqual(result, "{\"PageTotal\":2,\"Data\":[d,e]}");

        } catch (Exception e) {
            collector.logStepFail("Pagingquery failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }
}
