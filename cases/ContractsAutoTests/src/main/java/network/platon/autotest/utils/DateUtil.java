package network.platon.autotest.utils;

import java.text.DateFormat;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.Date;
import java.util.GregorianCalendar;
import java.util.Locale;

/**
 * @author qcxiao
 * @Description: 日期操作工具
 */
public class DateUtil {
	/**
	 * 日期转换成指定格式的字符串(默认格式为yyyy-MM-dd HH:mm:ss.SSS)
	 * @param date 日期
	 * @param pattern 格式
	 * @return 字符串格式的日期
	 * @Example
	 * DateUtil.dateToStr(new GregorianCalendar().getTime(),"yyyy-MM-dd");//2013-11-15
	 */
	public static String dateToStr(Date date, String pattern) {
		return dateToStr(date, pattern, Locale.CHINA);
	}

	/**
	 * 日期转换成指定格式的字符串（默认格式为yyyy-MM-dd HH:mm:ss.SSS）
	 * @param date Date date， 日期；
	 * @param pattern String pattern，格式；
	 * @param locale Locale locale，区域语言；
	 * @return 字符串格式的日期
	 * dateToStr(new Date(), "ddMMM", Locale.ENGLISH);//29Nov
	 */
	public static String dateToStr(Date date, String pattern, Locale locale) {
		if (pattern == null) {
			pattern = "yyyy-MM-dd HH:mm:ss.SSS";
		}
		DateFormat ymdhmsFormat = new SimpleDateFormat(pattern, locale);

		return ymdhmsFormat.format(date);
	}

	/**
	 * 字符串转为Date对象（默认格式为yyyy-MM-dd HH:mm:ss.SSS）
	 * @param str 字符串
	 * @param pattern 格式
	 * @return
	 * @throws ParseException
	 * @Example 
	 * DateUtil.strToDate("2013-11-15","yyyy-MM-dd");//Fri Nov 15 00:00:00 CST 2013
	 */
	public static Date strToDate(String str, String pattern) throws ParseException {
		return strToDate(str, pattern, Locale.CHINA);
	}

	/**
	 * 指定格式的字符串转换成日期(默认格式为yyyy-MM-dd HH:mm:ss.SSS)
	 * @param str 字符串
	 * @param pattern 格式
	 * @param locale 区域语言
	 * @return
	 * @throws ParseException
	 * @Example DateUtil.strToDate("15OCT2013", "ddMMMyyyy", Locale.ENGLISH);//Tue Oct 15 00:00:00 CST 2013
	 */
	public static Date strToDate(String str, String pattern, Locale locale) throws ParseException {
		if (pattern == null) {
			pattern = "yyyy-MM-dd HH:mm:ss.SSS";
		}
		DateFormat ymdhmsFormat = new SimpleDateFormat(pattern, locale);
		return ymdhmsFormat.parse(str);
	}

	/**
	 * 得到当天的日期
	 * @return Date
	 * @Example
	 * DateUtil.getToday();//Fri Nov 15 16:41:05 CST 2013
	 */
	public static Date getToday() {
		Calendar ca = Calendar.getInstance();
		return ca.getTime();
	}

	/**
	 * 生成日期
	 * @param year
	 * @param month
	 * @param date
	 * @return Date
	 */
	public static Date mkDate(int year, int month, int date) {
		Calendar ca = Calendar.getInstance();
		ca.set(year, month - 1, date);
		SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd hh:mm:ss");
		sdf.format(ca.getTime());
		return ca.getTime();
	}

	/**
	 * get GMT Time
	 * 
	 * @param calendar
	 * @return
	 */
	public Date getGmtDate(Long time) {
		Calendar calendar = Calendar.getInstance();
		calendar.setTimeInMillis(time);
		int offset = calendar.get(Calendar.ZONE_OFFSET) / 3600000 + calendar.get(Calendar.DST_OFFSET) / 3600000;
		calendar.add(Calendar.HOUR, -offset);
		Date date = calendar.getTime();
		return date;
	}

	/**
	 * 得到指定间隔天数的日期
	 * @param interval 间隔数
	 * @param format 格式
	 * @return String 字符串格式的日期
	 * @Example
	 * <div>DateUtil.getSpecifyDate(2,"yyyy-MM-dd");//两天后的日期</div>
	 * <div>DateUtil.getSpecifyDate(-2,"yyyy-MM-dd");//两天前的日期</div>
	 */
	public static String getSpecifyDate(int interval, String format) {
		return getSpecifyDate(interval, format, Locale.CHINA);
	}

	/**
	 * 得到指定间隔天数的日期
	 * @param interval 间隔数
	 * @param format 格式
	 * @param locale 区域语言
	 * @return String 字符串格式的日期
	 * @Example
	 * <div>getSpecifyDate(2, "ddMMM", Locale.ENGLISH);//两天后的日期：01Dec</div>
	 * <div>getSpecifyDate(-2, "ddMMM", Locale.ENGLISH);//两天前的日期：27Nov</div>
	 */
	public static String getSpecifyDate(int interval, String format, Locale locale) {

		Calendar cal = new GregorianCalendar();
		cal.add(Calendar.DATE, interval);
		return dateToStr(cal.getTime(), format, locale);
	}

	/**
	 * 得到指定间隔月数的日期
	 * @param interval 间隔数
	 * @param format 间隔数
	 * @return String 字符串格式的日期
	 * @Example
	 * <div>DateUtil.getSpecifyMonth(2,"yyyy-MM-dd");//两个月后的日期</div>
	 * <div>DateUtil.getSpecifyMonth(-2,"yyyy-MM-dd");//两个月前的日期</div>
	 */
	public static String getSpecifyMonth(int interval, String format) {
		return getSpecifyMonth(interval, format, Locale.CHINA);
	}

	/**
	 * 得到指定间隔月数的日期
	 * @param interval 间隔数
	 * @param format 间隔数
	 * @return String 字符串格式的日期
	 * @param locale 区域语言
	 * @Example
	 * <div>getSpecifyMonth(2, "ddMMM", Locale.ENGLISH);//两个月后的日期：29Jan</div>
	 * <div>getSpecifyMonth(-2, "ddMMM", Locale.ENGLISH);//两个月前的日期：29Sep</div>
	 */
	public static String getSpecifyMonth(int interval, String format, Locale locale) {
		Calendar cal = new GregorianCalendar();
		cal.add(Calendar.MONTH, interval);
		return dateToStr(cal.getTime(), format, locale);
	}

	/**
	 * 得到指定间隔年数的日期
	 * @param interval 间隔数
	 * @param format 格式
	 * @return 字符串格式的日期
	 * @Example
	 * <div>getSpecifyYear(2, "ddMMMyyyy");//两年后的日期</div>
	 * <div>getSpecifyYear(-2, "ddMMMyyyy");//两年前的日期</div>
	 */
	public static String getSpecifyYear(int interval, String format) {
		return getSpecifyYear(interval, format, Locale.CHINA);
	}

	/**
	 * 得到指定间隔年数的日期
	 * @param interval 间隔数
	 * @param format 格式
	 * @param locale 区域语言
	 * @return 字符串格式的日期
	 * @Example
	 * <div>getSpecifyYear(2, "ddMMMyyyy", Locale.ENGLISH);//两年后的日期：29Nov2015</div>
	 * <div>getSpecifyYear(-2, "ddMMMyyyy", Locale.ENGLISH);//两年前的日期：29Nov2011</div>
	 */
	public static String getSpecifyYear(int interval, String format, Locale locale) {
		Calendar cal = new GregorianCalendar();
		cal.add(Calendar.YEAR, interval);
		return dateToStr(cal.getTime(), format, locale);
	}
	/**
	 * 得到指定日期间隔天数的日期
	 * @param date 指定字符串格式的日期
	 * @param interval 间隔数
	 * @param format 格式
	 * @return 字符串格式的日期
	 * @Example
	 * getSpecifyDate("2014-09-09",3,"yyyy-MM-dd");
	 */
	public static String getSpecifyDate(String date, int interval, String format) {
		return getSpecifyDate(date, interval, format, Locale.CHINA);
	}
	public static String getSpecifyDate(String date, int interval, String format, Locale locale) {

		Date d = null;
		try {
			d = strToDate(date,"yyyy-MM-dd");
		} catch (ParseException e) {
			e.printStackTrace();
		}
		Calendar cal = new GregorianCalendar();
		cal.setTime(d);
		cal.add(Calendar.DATE, interval);
		return dateToStr(cal.getTime(), format, locale);
	}

	public static void main(String[] args) {
		try {
			System.out.println(DateUtil.strToDate("15OCT2013", "ddMMMyyy", Locale.ENGLISH));
			System.out.println(DateUtil.getToday());
			System.out.println(dateToStr(new Date(), "ddMMM", Locale.ENGLISH));
			System.out.println(dateToStr(new Date(), "yyyy/dd/MM"));
			System.out.println(getSpecifyDate(2, "ddMMM", Locale.ENGLISH));
			System.out.println(getSpecifyDate(-2, "ddMMM", Locale.ENGLISH));
			System.out.println(getSpecifyMonth(2, "ddMMM", Locale.ENGLISH));
			System.out.println(getSpecifyMonth(-2, "ddMMM", Locale.ENGLISH));
			System.out.println(getSpecifyYear(2, "ddMMMyyyy", Locale.ENGLISH));
			System.out.println(getSpecifyYear(-2, "ddMMMyyyy", Locale.ENGLISH));

		} catch (ParseException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
}
