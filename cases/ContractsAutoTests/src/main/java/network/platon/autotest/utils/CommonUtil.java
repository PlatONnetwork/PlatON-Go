package network.platon.autotest.utils;

import java.util.Calendar;
import java.util.Random;

/**
 * @Title: CommonUtil.java
 * @Package network.platon.autotest.utils
 * @Description: TODO(用一句话描述该文件做什么)
 * @author qcxiao
 * @date 2013-9-24 下午11:10:31
 */
public class CommonUtil {

	/**
	 * 指定时间内等待条件为真
	 * 
	 * @param condition
	 *            条件
	 * @param s
	 *            等待时间，秒为单位
	 * @return 等待是否成功
	 */
	private static Boolean waitUntil(Boolean condition, int s) {
		Calendar now = Calendar.getInstance();
		Calendar end = (Calendar) now.clone();
		end.add(Calendar.SECOND, s);
		while (!condition && now.before(end)) {
			try {
				Thread.sleep(500);
			} catch (InterruptedException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
			now = Calendar.getInstance();
		}
		return condition;
	}

	/**
	 * 指定时间内等待条件为真，默认60s
	 * 
	 * @param condition
	 *            条件
	 * @return 等待是否成功
	 */
	@SuppressWarnings("unused")
	private static Boolean waitUntil(Boolean condition) {
		return waitUntil(condition, 60);
	}

	public static void main(String args[]) {
		System.out.println(Calendar.getInstance().toString());
		if (waitUntil(1 == 0, 5)) {
			System.out.println("true");
		} else {
			System.out.println("false");
		}
		System.out.println(Calendar.getInstance().toString());
	}

	/**
	 * 生成指定长度的随机数，包含字母与数字
	 * @param length
	 * @return 随机数与字符
	 * @Example CommonUtil.getRandomCharacterAndNumber(8);
	 */
	public static String getRandomCharacterAndNumber(int length) {
		String val = "";

		Random random = new Random();
		for (int i = 0; i < length; i++) {
			String charOrNum = random.nextInt(2) % 2 == 0 ? "char" : "num"; // 输出字母还是数字

			if ("char".equalsIgnoreCase(charOrNum)) // 字符串
			{
				int choice = random.nextInt(2) % 2 == 0 ? 65 : 97; // 取得大写字母还是小写字母
				val += (char) (choice + random.nextInt(26));
			} else if ("num".equalsIgnoreCase(charOrNum)) // 数字
			{
				val += String.valueOf(random.nextInt(10));
			}
		}

		return val;
	}

	/**
	 * 生成指定长度的字母随机数
	 * @param length
	 * @return 随机字符
	 * @Example
	 * CommonUtil.getRandomCharacter(8);
	 */
	public static String getRandomCharacter(int length) {
		String val = "";

		Random random = new Random();
		for (int i = 0; i < length; i++) {
			int choice = random.nextInt(2) % 2 == 0 ? 65 : 97; // 取得大写字母还是小写字母
			val += (char) (choice + random.nextInt(26));
		}

		return val;
	}

	/**
	 * 生成指定长度的数字随机数
	 * @param length 
	 * @return 随机数
	 * @Example 
	 * CommonUtil.getRandomNumber(8);
	 */
	public static String getRandomNumber(int length) {
		String val = "";
		Random random = new Random();
		for (int i = 0; i < length; i++) {
			val += String.valueOf(random.nextInt(10));
		}

		return val;
	}
	public static Integer getRandomForMax(int value){
		Random random = new Random();
		int val = random.nextInt(value);
		return val == 0 ? 1 : val;
	}
}
