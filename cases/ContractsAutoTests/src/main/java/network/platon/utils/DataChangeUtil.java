package network.platon.utils;


import com.platon.rlp.RLPCodec;
import com.platon.rlp.RLPList;
import org.apache.commons.lang.StringUtils;
import org.web3j.utils.Numeric;

import java.math.BigInteger;

/**
 * @title ContractsAutoTests
 * @description: 数据类型转换工具类
 * @author: hudenian
 * @create: 2019/12/26 18:17
 */
public class DataChangeUtil {

    /**
     * hex字符串转byte数组
     * @param inHex 待转换的Hex字符串
     * @return  转换后的byte数组结果
     */
    public static byte[] hexToByteArray(String inHex){
        if(inHex.startsWith("0x"))inHex = inHex.substring(2);
        int hexlen = inHex.length();
        byte[] result;
        if (hexlen % 2 == 1){
            //奇数
            hexlen++;
            result = new byte[(hexlen/2)];
            inHex="0"+inHex;
        }else {
            //偶数
            result = new byte[(hexlen/2)];
        }
        int j=0;
        for (int i = 0; i < hexlen; i+=2){
            result[j]=hexToByte(inHex.substring(i,i+2));
            j++;
        }
        return result;
    }


    /**
     * Hex字符串转byte
     * @param inHex 待转换的Hex字符串
     * @return  转换后的byte
     */
    public static byte hexToByte(String inHex){
        return (byte)Integer.parseInt(inHex,16);
    }

    /**
     * 字节转十六进制
     * @param b 需要进行转换的byte字节
     * @return  转换后的Hex字符串
     */
    public static String byteToHex(byte b){
        String hex = Integer.toHexString(b & 0xFF);
        if(hex.length() < 2){
            hex = "0" + hex;
        }
        return hex;
    }

    /**
     * 字节数组转16进制
     * @param bytes 需要转换的byte数组
     * @return  转换后的Hex字符串
     */
    public static String bytesToHex(byte[] bytes) {
        StringBuffer sb = new StringBuffer();
        for(int i = 0; i < bytes.length; i++) {
            String hex = Integer.toHexString(bytes[i] & 0xFF);
            if(hex.length() < 2){
                sb.append(0);
            }
            sb.append(hex);
        }
        return sb.toString();
    }

    public  static String subHexData(String hexStr) {
        if (StringUtils.isBlank(hexStr)) {
            throw new IllegalArgumentException("string is blank");
        }
        if (StringUtils.startsWith(hexStr, "0x")) {
            hexStr = StringUtils.substringAfter(hexStr, "0x");
        }
        byte[] addi = hexStr.getBytes();
        for (int i = 0; i < addi.length; i++) {
            if (addi[i] != 0) {
                hexStr = StringUtils.substring(hexStr, i - 1);
                break;
            }
        }
        return hexStr;
    }

    public static byte[] stringToBytes32(String string) {
        byte[] byteValue = string.getBytes();
        byte[] byteValueLen32 = new byte[32];
        System.arraycopy(byteValue, 0, byteValueLen32, 0, byteValue.length);
        return byteValueLen32;
    }

    /**
     *
     * @param string 需要转换成byte数组的字符串
     * @param n byten
     * @return
     */
    public static byte[] stringToBytesN(String string,int n) {
        byte[] byteValue = string.getBytes();
        byte[] byteValueLen = new byte[n];
        System.arraycopy(byteValue, 0, byteValueLen, 0, byteValue.length);
        return byteValueLen;
    }


    /**
     * ppos系统约data字段rlp解码处理
     * 系统合约的最外层是一个RLPList
     * RlpList里面放实际的RLP编码值
     *
     * @param hexRlp
     * @return
     */
	public static String decodeSystemContractRlp(String hexRlp, long chainId) {
		        byte[] data = Numeric.hexStringToByteArray(hexRlp);
			        RLPList rlpList = RLPCodec.decode(data,RLPList.class, chainId);
				        return RLPCodec.decode(rlpList.get(0),String.class, chainId);
					    }



    public static void main(String[] args) {
//        String hexvalue = "aaaa";
//        byte bytess = hexToByte(hexvalue);
//        System.out.println(bytess);
        test();
    }



    public static void test() {
//        byte[] bytes = new byte[10000000];
//
//        for (int i = 0; i < 10000000; i++) {
//            if (i%3 == 0) {
//                bytes[i] = 0;
//            } else {
//                bytes[i] = 1;
//            }
//        }

        byte[] bytes = hexToByteArray("aaaa");

        System.out.println("可以转换的进制范围：" + Character.MIN_RADIX + "-" + Character.MAX_RADIX);
        System.out.println("2进制："   + binary(bytes, 2));
        System.out.println("5进制："   + binary(bytes, 5));
        System.out.println("8进制："   + binary(bytes, 8));

        System.out.println("16进制："  + binary(bytes, 16));
        System.out.println("32进制："  + binary(bytes, 32));
        System.out.println("64进制："  + binary(bytes, 64));// 这个已经超出范围，超出范围后变为10进制显示

        System.exit(0);
    }

    /*
     * 将byte[]转为各种进制的字符串
     * @param bytes byte[]
     * @param radix 基数可以转换进制的范围(2-36)，从Character.MIN_RADIX到Character.MAX_RADIX，超出范围后变为10进制
     * @return 转换后的字符串
     */
    public static String binary(byte[] bytes, int radix){
        return new BigInteger(1, bytes).toString(radix);// 这里的1代表正数
    }


    public static byte[] toPrimitives(Byte[] oBytes){
        byte[] bytes = new byte[oBytes.length];

        for(int i = 0; i < oBytes.length; i++) {
            bytes[i] = oBytes[i];
        }

        return bytes;
    }

}
