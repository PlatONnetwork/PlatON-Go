package wasm.contract_migrate;

import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.utils.FileUtil;
import network.platon.contracts.wasm.ContractMigrate_v1;
import org.junit.Test;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.rlp.RlpEncoder;
import org.web3j.rlp.RlpString;
import org.web3j.tx.Transfer;
import org.web3j.utils.Convert;
import org.web3j.protocol.core.DefaultBlockParameterName;

import wasm.beforetest.WASMContractPrepareTest;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.math.BigDecimal;
import java.math.BigInteger;
import java.nio.file.Paths;

/**
 * @title contract migrate
 * @description:
 * @author: yuanwenjun
 * @create: 2020/02/12
 */
public class ContractMigrateBalanceTest extends WASMContractPrepareTest {

    //the file name of migrate contract
    private String wasmFile = "ContractMigrate_v1.bin";

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "Sheet1",
            author = "yuanwenjun", showName = "wasm.contract_migrate",sourcePrefix = "wasm")
    public void testMigrateContractBalance() {

        Byte[] init_arg = null;
        Long transfer_value = 100000L;
        Long gas_value = 200000L;
        BigInteger origin_contract_value = BigInteger.valueOf(10000);

        try {
            prepare();
            ContractMigrate_v1 contractMigratev1 = ContractMigrate_v1.deploy(web3j, transactionManager, provider).send();
            String contractAddress = contractMigratev1.getContractAddress();
            String transactionHash = contractMigratev1.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("ContractMigrateV1 issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            
            Transfer transfer = new Transfer(web3j, transactionManager);
            transfer.sendFunds(contractAddress, new BigDecimal(origin_contract_value), Convert.Unit.VON).send();
            BigInteger originBalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("origin contract balance is: " + originBalance);
            
            init_arg = loadInitArg();
            TransactionReceipt transactionReceipt = contractMigratev1.migrate_contract(init_arg,transfer_value,gas_value).send();
            collector.logStepPass("Contract Migrate V1  successfully hash:" + transactionReceipt.getTransactionHash());
            
            BigInteger originAfterMigrateBalance = web3j.platonGetBalance(contractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("After migrate, origin contract balance is: " + originAfterMigrateBalance);
            collector.assertEqual(originAfterMigrateBalance, BigInteger.valueOf(0), "checkout origin contract balance");
            
            String newContractAddress = contractMigratev1.getPlaton_event1_transferEvents(transactionReceipt).get(0).arg1;
            collector.logStepPass("new Contract Address is:"+newContractAddress);
            BigInteger newMigrateBalance = web3j.platonGetBalance(newContractAddress, DefaultBlockParameterName.LATEST).send().getBalance();
            collector.logStepPass("new contract balance is: " + newMigrateBalance);
            collector.assertEqual(newMigrateBalance, origin_contract_value.add(BigInteger.valueOf(transfer_value)), "checkout new contract balance");

        } catch (Exception e) {
            collector.logStepFail("ContractDistoryTest failure,exception msg:" , e.getMessage());
            e.printStackTrace();
        }
    }

    /**
     * generate the migrate int_arg parameter
     * init_arg magic number +  RLP(code, RLP("init", init_paras...))
     * @return
     */
    private Byte[] loadInitArg() {
        //create file object wasmFile
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
            // make sure get all data
            if (offset != buffer.length) {
                throw new IOException("Could not completely read file "
                        + file.getName());
            }
            fi.close();
            byte[] bufferFinish = buffer;//wasm

            byte[] initAndParamsRlp = RlpEncoder.encode(RlpString.create("init"));

            //merge two arrays
            byte[] bt3 = new byte[bufferFinish.length+initAndParamsRlp.length];
            System.arraycopy(bufferFinish, 0, bt3, 0, bufferFinish.length);
            System.arraycopy(initAndParamsRlp, 0, bt3, bufferFinish.length, initAndParamsRlp.length);

            Byte[] bodyRlp = toObjects(RlpEncoder.encode(RlpString.create(bt3)));


            //magic number is 0x0061736d
//            String magicNumber = "0x0061736d";
//            byte[] magicNumberRlp = RlpEncoder.encode(RlpString.create(magicNumber));
            Byte[] magicNumberRlp = new Byte[4];
            //0x,00,61,73,6d
            magicNumberRlp[0] = 0x00;
            magicNumberRlp[1] = 0x61;
            magicNumberRlp[2] = 0x73;
            magicNumberRlp[3] = 0x6d;


            //pass new contract code
            finalByte = new Byte[magicNumberRlp.length+bodyRlp.length];
            System.arraycopy(magicNumberRlp,0,finalByte,0,magicNumberRlp.length);
            System.arraycopy(bodyRlp,0,finalByte,magicNumberRlp.length,bodyRlp.length);
        } catch (FileNotFoundException e) {
            e.printStackTrace();
            collector.logStepFail("load wasm file fail:",e.getMessage());
        } catch (IOException e) {
        	collector.logStepFail("load wasm file fail:",e.getMessage());
            e.printStackTrace();
        }
//        System.out.println(DataChangeUtil.bytesToHex(DataChangeUtil.toPrimitives(finalByte)));
        return finalByte;
    }


    /**
     * transfer between byte[] and Byte[]
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
