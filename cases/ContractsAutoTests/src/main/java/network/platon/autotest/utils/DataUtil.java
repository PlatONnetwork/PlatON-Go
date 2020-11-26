package network.platon.autotest.utils;

import java.text.NumberFormat;

/**
 * 数据处理工具类
 * 
 * @author qcxiao
 * 
 */
public class DataUtil {
	/**
	 * 四舍五入，并保留多少位小数
	 * @param v 参数为double类型
	 * @param scale
	 * @return 返回值为String类型
	 */
	public static String round(double value, Integer max) {
        return roundObject(value,max);
	}
	public static String roundObject(double value, Integer max) {
		NumberFormat nf = NumberFormat.getInstance();
        nf.setMaximumFractionDigits(max);
        nf.setMinimumFractionDigits(max);
        nf.setGroupingUsed(false);
        return nf.format(value);
	}
	/**
	 * 四舍五入，并保留多少位小数
	 * @param v 参数为String类型
	 * @param scale
	 * @return 返回值为String类型
	 */
	public static String round(String value, Integer max) {
        return roundObject(Double.parseDouble(value),max);
	}
	/**
	 * 四舍五入，并保留多少位小数
	 * @param v 参数为double类型
	 * @param scale
	 * @return 返回值为Double类型
	 */
	public static Double round(double value, int max) {
        return Double.parseDouble(roundObject(value,max));
	}
	/**
	 * 四舍五入，并保留多少位小数
	 * @param v 参数为String类型
	 * @param scale
	 * @return 返回值为Double类型
	 */
	public static Double round(String value, int max) {
        return Double.parseDouble(roundObject(Double.parseDouble(value),max));
	}
	public static void main(String[] args) {
		System.out.println(round(12.0036,Integer.valueOf(3)));
	}
}
