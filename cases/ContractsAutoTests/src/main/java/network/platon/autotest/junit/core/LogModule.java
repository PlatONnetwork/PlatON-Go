package network.platon.autotest.junit.core;

import java.io.UnsupportedEncodingException;
import java.net.URLEncoder;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.StepType;
import network.platon.autotest.junit.enums.TestSuiteType;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.LogStepInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import org.openqa.selenium.NoAlertPresentException;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.phantomjs.PhantomJSDriver;

import network.platon.autotest.utils.DateUtil;
import network.platon.autotest.utils.ScreenShotUtil;

/**
 * 存放日志
 * @author qcxiao
 *
 */
public class LogModule {
	public static SuiteInfo SUITE_INFO = new SuiteInfo();
	public static ModuleInfo MODULE_INFO = new ModuleInfo();
	public static CaseInfo CASE_INFO = new CaseInfo();
	public static String SUITE_PATH = "";
	public static List<LogStepInfo> logStepInfoList = new ArrayList<LogStepInfo>();

	public static List<LogStepInfo> onLogStep(LogStepInfo logStepInfo) {
		logStepInfo.setStepId(logStepInfoList.size() + 1);
		logStepInfoList.add(logStepInfo);
		return logStepInfoList;
	}

	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, RunResult stepResult, String failReason, TestSuiteType testSuiteType) {
		switch (testSuiteType) {
		case WEB_UI:
			return logStepFail(stepType, stepDesc, stepResult, failReason, TestSuiteType.IOS);
		default:
			return logStepFail(stepType, stepDesc, stepResult, failReason);
		}
	}
	/**
	 * 日志中输入期望值与实际值的assertEqual断言
	 */
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, String actual, String expect, RunResult stepResult, String failReason, TestSuiteType testSuiteType) {
		switch (testSuiteType) {
		case WEB_UI:
			return logStepFail(stepType, stepDesc, actual, expect, stepResult, failReason, TestSuiteType.IOS);
		default:
			return logStepFail(stepType, stepDesc, actual, expect, stepResult, failReason);
		}
	}
	
	

	/**
	 * 通用步骤失败日志
	 */
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, RunResult stepResult, String failReason) {
		LogStepInfo logStepInfo = new LogStepInfo();
		logStepInfo.setStepType(stepType);
		logStepInfo.setStepDesc(stepDesc.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setStepResult(stepResult);
		logStepInfo.setFailReason(failReason);
		String url = "";
		logStepInfo.setUrl(url);
		String stepTime = DateUtil.dateToStr(new Date(), "HH:mm:ss");
		logStepInfo.setStepTime(stepTime);
		String failType = "";
		logStepInfo.setFailType(failType);
		System.err.println(stepDesc + "\n" + failReason);
		return onLogStep(logStepInfo);
	}
	
	/**
	 * 通用步骤失败日志
	 * 日志中输入期望值与实际值的assertEqual断言
	 */
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, String actual, String expect, RunResult stepResult, String failReason) {
		LogStepInfo logStepInfo = new LogStepInfo();
		logStepInfo.setStepType(stepType);
		//logStepInfo.setStepDesc(stepDesc.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setStepDesc(stepDesc);
		logStepInfo.setActual(actual.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setExpect(expect.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setStepResult(stepResult);
		logStepInfo.setFailReason(failReason);
		String url = "";
		logStepInfo.setUrl(url);
		String stepTime = DateUtil.dateToStr(new Date(), "HH:mm:ss");
		logStepInfo.setStepTime(stepTime);
		String failType = "";
		logStepInfo.setFailType(failType);
		//System.err.println(stepDesc + "\n" + failReason);
		return onLogStep(logStepInfo);
	}

	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, RunResult stepResult, String failReason, String hWnd) {
		return logStepFail(stepType, stepDesc, stepResult, failReason, hWnd, "");
	}

	/**
	 * WINDOWS_UI步骤失败日志，暂时没对hWnd做处理，直接全屏截图
	 */
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, RunResult stepResult, String failReason, String hWnd, String failType) {
		LogStepInfo logStepInfo = new LogStepInfo();
		List<LogStepInfo> logStepInfoList = logStepFail(stepType, stepDesc, stepResult, failReason);
		logStepInfo = logStepInfoList.remove(logStepInfoList.size() - 1);
		String url = "";
		logStepInfo.setUrl(url);
		logStepInfo.setFailType(failType);
		return onLogStep(logStepInfo);
	}

	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, RunResult stepResult, String failReason, WebDriver driver) {
		return logStepFail(stepType, stepDesc, stepResult, failReason, driver, "");
	}
	
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, String actual, String expect, RunResult stepResult, String failReason, WebDriver driver) {
		return logStepFail(stepType, stepDesc, actual, expect, stepResult, failReason, driver, "");
	}
	

	/**
	 * WEB_UI步骤失败日志
	 */
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, RunResult stepResult, String failReason, WebDriver driver, String failType) {
		LogStepInfo logStepInfo = new LogStepInfo();
		List<LogStepInfo> logStepInfoList = logStepFail(stepType, stepDesc, stepResult, failReason);
		logStepInfo = logStepInfoList.remove(logStepInfoList.size() - 1);
		String url = "";
		String picture = "";
		if (driver instanceof PhantomJSDriver) {
			url = driver.getCurrentUrl();
			//picture = ScreenShotUtil.screenShotByDriver(driver, SUITE_PATH);
		} else {
//			try {
//				if (isAlert(driver)) {
//					url = driver.switchTo().alert().getText();
//					picture = ScreenShotUtil.screenShotByDesktop(SUITE_PATH);
//
//				} else {
//					url = driver.getCurrentUrl();
//					picture = ScreenShotUtil.screenShotByDriver(driver, SUITE_PATH);
//				}
//			} catch (Exception e) {
//				if (picture.equals("")) {
//					try {
//						picture = ScreenShotUtil.screenShotByDesktop(SUITE_PATH);
//					} catch (Exception ex) {
//						ex.printStackTrace();
//					}
//
//				}
//			}
		}
		logStepInfo.setUrl(url);
		logStepInfo.setPicture(picture);
		if (!picture.equals("")) {
			String[] pictureAddress = picture.split("screenshot");
			String pictureName = pictureAddress[pictureAddress.length - 1].substring(1);
			logStepInfo.setPictureName(pictureName);
			String pictureRelative = "../screenshot/" + pictureName;
			logStepInfo.setPictureRelative(pictureRelative);
		}
		logStepInfo.setFailType(failType);
		return onLogStep(logStepInfo);
	}
	
	/**
	 * WEB_UI失败日志
	 * 日志中输入期望值与实际值的assertEqual断言
	 */
	public static List<LogStepInfo> logStepFail(StepType stepType, String stepDesc, String actual, String expect, RunResult stepResult, String failReason, WebDriver driver, String failType) {
		LogStepInfo logStepInfo = new LogStepInfo();
		List<LogStepInfo> logStepInfoList = logStepFail(stepType, stepDesc, actual, expect, stepResult, failReason);
		logStepInfo = logStepInfoList.remove(logStepInfoList.size() - 1);
		String url = "";
		String picture = "";
		try {
			if (isAlert(driver)) {
				url = driver.switchTo().alert().getText();
				picture = ScreenShotUtil.screenShotByDesktop(SUITE_PATH);

			} else {
				url = driver.getCurrentUrl();
				picture = ScreenShotUtil.screenShotByDriver(driver, SUITE_PATH);
			}
		} catch (Exception e) {
			if (picture.equals("")) {
				try {
					picture = ScreenShotUtil.screenShotByDesktop(SUITE_PATH);
				} catch (Exception ex) {
					ex.printStackTrace();
				}

			}
		}
		logStepInfo.setUrl(url);
		logStepInfo.setStepDesc(stepDesc);
		logStepInfo.setPicture(picture);
		if (!picture.equals("")) {
			String[] pictureAddress = picture.split("screenshot");
			String pictureName = pictureAddress[pictureAddress.length - 1].substring(1);
			logStepInfo.setPictureName(pictureName);
			String pictureRelative = "../screenshot/" + pictureName;
			logStepInfo.setPictureRelative(pictureRelative);
		}
		logStepInfo.setFailType(failType);
		return onLogStep(logStepInfo);
	}

	public static List<LogStepInfo> logStepPass(StepType stepType, String stepDesc, RunResult stepResult, TestSuiteType testSuiteType) {
		switch (testSuiteType) {
		case WEB_UI:
			return logStepPass(stepType, stepDesc, RunResult.PASS, TestSuiteType.WEB_UI);
		default:
			return logStepPass(stepType, stepDesc, RunResult.PASS);
		}

	}

	/**
	 * 通用步骤成功日志
	 */
	public static List<LogStepInfo> logStepPass(StepType stepType, String stepDesc, RunResult stepResult) {
		LogStepInfo logStepInfo = new LogStepInfo();
		logStepInfo.setStepType(stepType);
		logStepInfo.setStepDesc(stepDesc.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setStepResult(stepResult);
		String url = "";
		logStepInfo.setUrl(url);
		String stepTime = DateUtil.dateToStr(new Date(), "HH:mm:ss");
		logStepInfo.setStepTime(stepTime);
		try {
			Thread.sleep(100);
		} catch (InterruptedException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		System.out.println(stepDesc);
		return onLogStep(logStepInfo);
	}
	/**
	 * 通用步骤成功日志
	 * 日志中输入期望值与实际值
	 */
	public static List<LogStepInfo> logStepPass(StepType stepType, String stepDesc, String actual, String expect, RunResult stepResult) {
		LogStepInfo logStepInfo = new LogStepInfo();
		logStepInfo.setStepType(stepType);
		logStepInfo.setStepDesc(stepDesc);
		logStepInfo.setActual(actual.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setExpect(expect.replace("<","&lt;").replace(">","&gt;"));
		logStepInfo.setStepResult(stepResult);
		String url = "";
		logStepInfo.setUrl(url);
		String stepTime = DateUtil.dateToStr(new Date(), "HH:mm:ss");
		logStepInfo.setStepTime(stepTime);
		try {
			Thread.sleep(100);
		} catch (InterruptedException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		System.out.println(stepDesc);
		return onLogStep(logStepInfo);
	}
	
	/**
	 * WINDOWS_UI步骤成功日志
	 */
	public static List<LogStepInfo> logStepPass(StepType stepType, String stepDesc, RunResult stepResult, String hWnd) {
		LogStepInfo logStepInfo = new LogStepInfo();
		List<LogStepInfo> logStepInfoList = logStepPass(stepType, stepDesc, stepResult);
		logStepInfo = logStepInfoList.remove(logStepInfoList.size() - 1);
		String url = "";
		logStepInfo.setUrl(url);
		return onLogStep(logStepInfo);
	}

	/**
	 * WEB_UI步骤成功日志
	 */
	public static List<LogStepInfo> logStepPass(StepType stepType, String stepDesc, RunResult stepResult, WebDriver driver) {
		LogStepInfo logStepInfo = new LogStepInfo();
		List<LogStepInfo> logStepInfoList = logStepPass(stepType, stepDesc, stepResult);
		logStepInfo = logStepInfoList.remove(logStepInfoList.size() - 1);
		String url = "";
		try {
			if (isAlert(driver)) {
				url = "alert";
			} else {
				url = driver.getCurrentUrl();
			}
		} catch (Exception e) {
		}
		logStepInfo.setUrl(url);
		return onLogStep(logStepInfo);
	}

	/**
	 * 判断弹出框是否存在
	 */
	public static boolean isAlert(WebDriver driver) {
		try {
			driver.switchTo().alert();
			return true;
		} catch (NoAlertPresentException e) {
			return false;
		}
	}

	public static String encode(String str, String encoding) {
		try {
			return URLEncoder.encode(str, encoding);
		} catch (UnsupportedEncodingException e) {
			return str;
		}
	}

}
