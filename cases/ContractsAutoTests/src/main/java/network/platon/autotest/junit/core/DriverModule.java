package network.platon.autotest.junit.core;

import java.io.File;
import java.util.Enumeration;
import java.util.Vector;
import network.platon.autotest.exception.StepException;
import network.platon.autotest.junit.enums.FileType;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.StepType;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import network.platon.autotest.junit.rules.DriverService;

/**
 * 观察者模式，用于管理套件、模块、用例的开始与结束；后续如果有新增的类似对象也可以直接往里加
 * @author qcxiao
 *
 */
public class DriverModule implements Observer {
	private String sourcesDir = DriverService.PROPERTIES.getProperty("sourcesDir", "src/test/resources/");

	private DriverModule() {
	}

	/**
	 * 单例模式
	 */
	private static DriverModule singleton;

	public static DriverModule getInstance() {
		if (singleton == null) {
			singleton = new DriverModule();
		}
		return singleton;
	}

	private Vector<Observer> observersVector = new Vector<Observer>();

	public void attach(Observer observer) {
		observersVector.addElement(observer);
	}

	public void detach(Observer observer) {
		observersVector.removeElement(observer);
	}

	public void detachAll() {
		observersVector.clear();
	}

	@SuppressWarnings("unchecked")
	public Enumeration<Observer> observers() {
		return ((Vector<Observer>) observersVector.clone()).elements();
	}

	@Override
	public void suiteRunStart(SuiteInfo suiteInfo) {
		Enumeration<Observer> enumeration = observers();
		while (enumeration.hasMoreElements()) {
			((Observer) enumeration.nextElement()).suiteRunStart(suiteInfo);
		}

	}

	@Override
	public void moduleRunStart(ModuleInfo moduleInfo) {
		Enumeration<Observer> enumeration = observers();
		while (enumeration.hasMoreElements()) {
			((Observer) enumeration.nextElement()).moduleRunStart(moduleInfo);
		}

	}

	@Override
	public void caseRunStart(CaseInfo caseInfo) {
		Enumeration<Observer> enumeration = observers();
		while (enumeration.hasMoreElements()) {
			((Observer) enumeration.nextElement()).caseRunStart(caseInfo);
		}

	}

	@Override
	public void suiteRunStop(SuiteInfo suiteInfo) {
		Enumeration<Observer> enumeration = observers();
		while (enumeration.hasMoreElements()) {
			((Observer) enumeration.nextElement()).suiteRunStop(suiteInfo);
		}
	}

	@Override
	public void moduleRunStop(ModuleInfo moduleInfo) {
		Enumeration<Observer> enumeration = observers();
		while (enumeration.hasMoreElements()) {
			((Observer) enumeration.nextElement()).moduleRunStop(moduleInfo);
		}
	}

	@Override
	public void caseRunStop(CaseInfo caseInfo) {
		Enumeration<Observer> enumeration = observers();
		while (enumeration.hasMoreElements()) {
			((Observer) enumeration.nextElement()).caseRunStop(caseInfo);
		}

	}

	/**
	 * 主要用于数据库等初始化：initialData
	 * 
	 * @param caseInfo
	 */
	public void initialData(CaseInfo caseInfo) {
		String fileName = caseInfo.getCaseParams().get("initialData");
		if (fileName == null) {
			return;
		}
		String fileType = fileName.substring(fileName.lastIndexOf(".") + 1).toUpperCase();
		fileName = findFile(fileName);
		if (fileName == null) {
			LogModule.logStepFail(StepType.DATABASE, "initialData操作失败", RunResult.FAIL, caseInfo.getCaseParams().get("initialData") + "文件不存在!");
			throw new StepException("initialData操作失败!" + caseInfo.getCaseParams().get("initialData") + "文件不存在!");
		}
		switch (FileType.valueOf(fileType)) {
		case JAVA:
			break;
		}

	}

	/**
	 * 主要用于数据库等销毁：destroyData
	 * 
	 * @param caseInfo
	 */
	public void destroyData(CaseInfo caseInfo) {
		String fileName = caseInfo.getCaseParams().get("destroyData");
		if (fileName == null) {
			return;
		}
		String fileType = fileName.substring(fileName.lastIndexOf(".") + 1).toUpperCase();
		fileName = findFile(fileName);
		if (fileName == null) {
			LogModule.logStepFail(StepType.DATABASE, "destroyData操作失败", RunResult.FAIL, caseInfo.getCaseParams().get("destroyData") + "文件不存在!");
			throw new StepException("destroyData操作失败!" + caseInfo.getCaseParams().get("destroyData") + "文件不存在!");
		}
		switch (FileType.valueOf(fileType)) {
		case JAVA:
			break;
		}

	}

	/**
	 * 根据文件名得到文件路径
	 * @param fileName
	 * @return
	 */
	private String findFile(String fileName) {
		// 根目录下
		String filePath = System.getProperty("user.dir") + "/" + sourcesDir + fileName;
		File file = new File(filePath);
		// 测试类的数据源目录下
		if (!file.exists()) {
			filePath = System.getProperty("user.dir") + "/" + sourcesDir + DriverService.DESCRIPTION.getTestClass().getSimpleName() + "/" + fileName;
			file = new File(filePath);
		}
		if (!file.exists()) {
			filePath = null;
		}
		return filePath;
	}

}
