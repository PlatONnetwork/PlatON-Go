package wasm.contract_migrate;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.contracts.wasm.ContractMigrate_v1;
import org.junit.Test;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.rlp.RlpEncoder;
import org.web3j.rlp.RlpString;
import org.web3j.tx.gas.GasProvider;
import wasm.beforetest.WASMContractPrepareTest;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Paths;

/**
 * @title 合约升级
 * @description:
 * @author: hudenian
 * @create: 2020/02/10
 */
public class ContractMigrateV1Test extends WASMContractPrepareTest {

    //需要升级的合约
    private String wasmFile = "ContractMigrate_v1.bin";

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "hudenian", showName = "wasm.contract_migrate合约升级",sourcePrefix = "wasm")
    public void testMigrateContract() {

        Byte[] init_arg = null;
        Long transfer_value = 100000L;
        Long gas_value = 200000L;
        String name = "hello";

        try {
            prepare();
            ContractMigrate_v1 contractMigratev1 = ContractMigrate_v1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contractMigratev1.getContractAddress();
            String transactionHash = contractMigratev1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractMigrateV1 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //设置值
            transactionHash = contractMigratev1.set_string(name).send().getTransactionHash();
            collector.logStepPass("ContractMigrateV1 set_string successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);

            //查询结果
            String chainName = contractMigratev1.get_string().send();
            collector.assertEqual(chainName,name);


            /**
             * 加载需要升级的合约
             * init_arg 参数为  magic number +  RLP(code, RLP("init", init_paras...))
             * transfer_value 为转到新合约地址的金额，gas_value 为预估消耗的 gas
             */
            init_arg = loadInitArg();
//            System.out.println(Arrays.toString(init_arg));

            //合约升级
            TransactionReceipt transactionReceipt = contractMigratev1.migrate_contract(init_arg,transfer_value,gas_value).send();
            collector.logStepPass("Contract Migrate V1  successfully hash:" + transactionReceipt.getTransactionHash());

            //获取升级后的合约地址(需要通过事件获取)
            String newContractAddress = contractMigratev1.getPlaton_event1_transferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);

            //调用升级后的合约
            //FIXME 等bug修复后放开
//            ContractMigrate_v1 new_contractMigrate_v1 = ContractMigrate_v1.load(newContractAddress,web3j,credentials,provider);
//            String newContractChainName = new_contractMigrate_v1.get_string().send();
//            collector.assertContains(newContractChainName,name);

        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    /**
     * 拼装需要加载的合约信息
     * init_arg 参数为  magic number +  RLP(code, RLP("init", init_paras...))
     * @return
     */
    private Byte[] loadInitArg() {
        //创建一个文件对象 wasmFile
        String filePath = FileUtil.pathOptimization(Paths.get("src", "test", "resources", "contracts", "wasm", "contract_migrate").toUri().getPath());
        File file = new File(filePath+File.separator+wasmFile);

        long fileSize = file.length();
        if (fileSize > Integer.MAX_VALUE) {
            System.out.println("file too big...");
        }

        FileInputStream fi = null;
        Byte[] finalByte = null;
        try {
            fi = new FileInputStream(file);
            byte[] buffer = new byte[(int) fileSize];
            int offset = 0;
            int numRead = 0;
            while (offset < buffer.length
                    && (numRead = fi.read(buffer, offset, buffer.length - offset)) >= 0) {
                offset += numRead;
            }
            // 确保所有数据均被读取
            if (offset != buffer.length) {
                throw new IOException("Could not completely read file "
                        + file.getName());
            }
            fi.close();
            byte[] bufferFinish = buffer;//wasm

            byte[] initAndParamsRlp = RlpEncoder.encode(RlpString.create("init"));

            //将两个数组合并
            byte[] bt3 = new byte[bufferFinish.length+initAndParamsRlp.length];
            System.arraycopy(bufferFinish, 0, bt3, 0, bufferFinish.length);
            System.arraycopy(initAndParamsRlp, 0, bt3, bufferFinish.length, initAndParamsRlp.length);

            Byte[] bodyRlp = toObjects(RlpEncoder.encode(RlpString.create(bt3)));


            //magic number为固定值0x0061736d
//            String magicNumber = "0x0061736d";
//            byte[] magicNumberRlp = RlpEncoder.encode(RlpString.create(magicNumber));
            Byte[] magicNumberRlp = new Byte[4];
            //0x,00,61,73,6d
            magicNumberRlp[0] = 0x00;
            magicNumberRlp[1] = 0x61;
            magicNumberRlp[2] = 0x73;
            magicNumberRlp[3] = 0x6d;


            //需要传入进行升级合约中的代码
            finalByte = new Byte[magicNumberRlp.length+bodyRlp.length];
            System.arraycopy(magicNumberRlp,0,finalByte,0,magicNumberRlp.length);
            System.arraycopy(bodyRlp,0,finalByte,magicNumberRlp.length,bodyRlp.length);
        } catch (FileNotFoundException e) {
            e.printStackTrace();
            collector.logStepFail("加载wasm二进制文件失败，失败原因:",e.getMessage());
        } catch (IOException e) {
            e.printStackTrace();
        }
//        System.out.println(DataChangeUtil.bytesToHex(DataChangeUtil.toPrimitives(finalByte)));
        return finalByte;
    }


    /**
     * byte[] 与 Byte[]之间转换
     * @param bytesPrim
     * @return
     */
    private  Byte[] toObjects(byte[] bytesPrim) {
        Byte[] bytes = new Byte[bytesPrim.length];

        int i = 0;
        for (byte b : bytesPrim) bytes[i++] = b; // Autoboxing

        return bytes;
    }
}
