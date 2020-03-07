package network.platon.utils;

import lombok.extern.slf4j.Slf4j;
import network.platon.autotest.utils.FileUtil;
import org.web3j.rlp.RlpEncoder;
import org.web3j.rlp.RlpList;
import org.web3j.rlp.RlpString;
import org.web3j.rlp.RlpType;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Paths;
import java.util.List;

/**
 * @title ContractsAutoTests
 * @description: RLP编解码及生成合约升级参数工具类
 * @author: hudenian
 * @create: 2020/2/12 18:03
 */
@Slf4j
public class RlpUtil {

    /**
     *
     * @param wasmFile 根据wasm文件
     * @param args  init参数数组
     * @return
     */
    public static Byte[] loadInitArg(String wasmFile,List<String> args) {

        File file = new File(wasmFile);

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


            byte[] initAndParamsRlp = null;

            if(args.size() == 0 ){
                initAndParamsRlp = RlpEncoder.encode(RlpString.create("init"));
            }else{
                RlpType[] values = new RlpType[args.size()+1];
                values[0] = RlpString.create("init");
                for(int i=0;i<args.size();i++){
                    values[i+1] = RlpString.create(args.get(i));
                }
                initAndParamsRlp = RlpEncoder.encode(new RlpList(values));
            }


            //merge two arrays
            byte[] bt3 = new byte[bufferFinish.length+initAndParamsRlp.length];
            System.arraycopy(bufferFinish, 0, bt3, 0, bufferFinish.length);
            System.arraycopy(initAndParamsRlp, 0, bt3, bufferFinish.length, initAndParamsRlp.length);

            Byte[] bodyRlp = toObjects(RlpEncoder.encode(RlpString.create(bt3)));


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
            log.error("load wasm file fail:{}", e.getMessage());
        } catch (IOException e) {
            log.error("load wasm file fail:{}", e.getMessage());
            e.printStackTrace();
        }
//        System.out.println("utl>>>"+DataChangeUtil.bytesToHex(DataChangeUtil.toPrimitives(finalByte)));
        return finalByte;
    }


    /**
     * transfer between byte[] and Byte[]
     * @param bytesPrim
     * @return
     */
    public static  Byte[] toObjects(byte[] bytesPrim) {
        Byte[] bytes = new Byte[bytesPrim.length];

        int i = 0;
        for (byte b : bytesPrim) bytes[i++] = b;

        return bytes;
    }
}
