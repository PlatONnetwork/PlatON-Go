package network.platon.autotest.junit.core;

import java.io.File;
import java.io.IOException;
import java.lang.reflect.Method;
import java.text.NumberFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.HashSet;
import java.util.Iterator;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.Set;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import lombok.extern.slf4j.Slf4j;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.LogStepInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import network.platon.autotest.junit.rules.DriverService;
import org.dom4j.DocumentException;
import org.junit.Ignore;
import org.junit.runner.Description;
import com.alibaba.fastjson.JSON;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.BrowserType;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.RunStatus;
import network.platon.autotest.junit.enums.StepType;
import network.platon.autotest.junit.enums.TestSuiteType;
import network.platon.autotest.utils.ClassUtil;
import network.platon.autotest.utils.DateUtil;
import network.platon.autotest.utils.ExcelUtil;
import network.platon.autotest.utils.FileUtil;
import network.platon.autotest.utils.XmlUtil;

/**
 * @Title: SuiteInitial.java
 * @Package network.platon.autotest.junit.modules
 * @Description: TODO(用一句话描述该文件做什么)
 * @author qcxiao
 * @date 2013-12-12 下午04:49:47
 */
@Slf4j
public class SuiteObserver implements Observer {
	/**
	 * test.properties里的配置信息
	 */
	private Properties properties = DriverService.PROPERTIES;
	/**
	 * 初始化日志路径，如果没有就默认为C:/autosky_log/
	 */
	private String logDir = DriverService.PROPERTIES.getProperty("logDir", "C:/autosky_log/");
	/**
	 * 初始化配置文件路径，如果没有配置就默认为src/test/resources/
	 */
	private String sourcesDir = DriverService.PROPERTIES.getProperty("sourcesDir", "src/test/resources/");
	/**
	 * 计划要执行的套件信息
	 */
	private SuiteInfo planedSuiteInfo = new SuiteInfo();
	/**
	 * website.properties资源文件中的键值匹配对信息
	 */
	public static Map<String, Map<String, String>> websiteMap = DriverService.PROPERTIES_MAP;
	/**
	 * 存放日志步骤信息
	 */
	public List<LogStepInfo> logStepInfoList = new ArrayList<LogStepInfo>();
	/**
	 * 收集异常
	 */
	List<Throwable> errors = new ArrayList<Throwable>();

	@Override
	public void suiteRunStart(SuiteInfo suiteInfo) {
		suiteInfo.setSuiteResult(RunResult.RUNNING);
		suiteInfo.setSuiteStatus(RunStatus.RUNNING);
		Boolean suiteMerged = DriverService.SUITE_MERGED;
		suiteInfo.setReportMerged(suiteMerged);
		if (suiteMerged) {
			// mvn test -Dreport.merged=true
			suiteInfo.setSuiteName(properties.getProperty("project"));
		} else {
			// 单个类执行或者 mvn test 或者 mvn test -Dreport.merged=false
			suiteInfo.setSuiteName(DriverService.DESCRIPTION.getTestClass().getSimpleName());
		}
		String runId = DriverService.PROPERTIES.getProperty("runId");
		if (runId == null || runId.trim().equals("") || "$runId".equals(runId)) {
			runId = "buildId:0:buildTaskId:0:buildTestSuiteId:0";
		}
		suiteInfo.setBuildId(Long.parseLong(runId.split(":")[1]));
		suiteInfo.setBuildTaskId(Long.parseLong(runId.split(":")[3]));
		suiteInfo.setBuildTestSuiteId(Long.parseLong(runId.split(":")[5]));
		
		String submitInfo = DriverService.PROPERTIES.getProperty("submitInfo");
		if (submitInfo == null || submitInfo.trim().equals("") || "$submitInfo".equals(submitInfo)) {
			submitInfo = "submitter:0:submitdate:0:submitnote:0";
		}
		suiteInfo.setSubmitter(submitInfo.split(":")[1]);
		suiteInfo.setSubmitdate(submitInfo.split(":")[3]);
		suiteInfo.setSubmitnote(submitInfo.split(":")[5]);
		
		// 后续加上负责人的信息
		suiteInfo.setSuiteStartTime(new Date());
		List<ModuleInfo> moduleInfoList = new ArrayList<ModuleInfo>();
		
		suiteInfo.setModuleInfoList(moduleInfoList);
		
		// 后续加一个suite集合
		addSuiteModules(suiteInfo);
		
		// 后续走多个class的话，需要再这里更新
		// LogModule.suiteInfo = suiteInfo;
		/*
		 * 获取webSite.properties里的公共信息形成键值对信息
		 */
		websiteMap.clear();
		websiteMap = getWebsiteMap();// 各个网站信息集合
		/**
		 * 判断是否maven执行,决定编码方式，主要是兼容windows窗口控制台输出的gbk编码，后续做改进
		 * 获取栈信息中是否含有maven的信息，如果是就给编码赋值为GBK
		 */
		StackTraceElement stack[] = Thread.currentThread().getStackTrace();
		for (StackTraceElement ste : stack) {
			if (ste.getClassName().contains("maven")) {
				DriverService.ENCODE = "GBK";
				break;
			}
		}
		logDir = logDir + suiteInfo.getSuiteName() + "/";
		String buildNumber = System.getProperty("build.number");
		if (buildNumber != null) {
			//如果是jenkins执行时，buildNumber会在以前的基础上自动加1然后发给框架
			logDir += buildNumber + "/";
		} else {
			logDir += DateUtil.dateToStr(new Date(), "yyyy-MM-dd_HH-mm-ss") + "/";
		}
		LogModule.SUITE_PATH = logDir;
		File directory = new File(logDir);
		if (!directory.exists()) {
			directory.mkdirs();
		}
		System.out.println("测试套件（" + suiteInfo.getSuiteName() + "）执行开始。");
	}

	/**
	 * 两种情况： 1、合并报告：suiteName就为Project属性值；执行过程中会去收集target\test-
	 * classes目录下的所有符合条件类与符合条件方法 2、不合并报告：suiteName就为本次运行的类名；将收集本次运行类里的符合条件方法
	 * 
	 * @param suiteInfo
	 * @return
	 */
	public List<ModuleInfo> addSuiteModules(SuiteInfo suiteInfo) {
		List<ModuleInfo> moduleInfoList = new ArrayList<ModuleInfo>();
		Set<Class<?>> classes = new HashSet<Class<?>>();
		if (suiteInfo.getReportMerged()) {
			try {
				/**
				 * 收集当前目录下的所有类
				 */
				String filePath = this.getClass().getResource("/").getPath();
				//System.out.println("filePath: " + filePath);
				if(System.getProperties().getProperty("os.name").contains("Windows")){
					filePath = this.getClass().getResource("/").getPath().replaceFirst("/", "").replace("/", "\\").replace("%20", " ");
				}
				System.out.println(filePath);
				classes = ClassUtil.getClasses(filePath);
			} catch (IOException e) {
				e.printStackTrace();
			}
		} else {
			classes.add(DriverService.DESCRIPTION.getTestClass());
		}
		for (Class<?> cls : classes) {
			// 不执行忽略的类
			if (cls.getAnnotation(Ignore.class) != null) {
				continue;
			}
			List<ModuleInfo> moduleNameList = getSuiteModuleNameList(cls);
			for (ModuleInfo moduleInfo : moduleNameList) {
				moduleInfo.setSuiteInfo(suiteInfo);
//				moduleInfo.setModuleName(cls.getSimpleName() + "." + moduleInfo.getModuleName());
				moduleInfo.setModuleName(cls.getName() + "." + moduleInfo.getModuleName());
				moduleInfo.setModuleRun(true);
				moduleInfo.setModuleStatus(RunStatus.WAITING);
				moduleInfo.setModuleResult(RunResult.WAITING);
				moduleInfoList.add(moduleInfo);
			}
		}
		List<ModuleInfo> planedModuleInfoList = mergePlanModuleInfo(moduleInfoList);
		suiteInfo.setModuleInfoList(planedModuleInfoList);
		return moduleInfoList;
	}

	@Override
	public void moduleRunStart(ModuleInfo moduleInfo) {
		// 由于此处无法遍历suiteInfo中的module信息，只能根据方法名来做
//		moduleInfo.setModuleName(DriverService.DESCRIPTION.getTestClass().getSimpleName() + "." + DriverService.DESCRIPTION.getMethodName());
		moduleInfo.setModuleName(DriverService.DESCRIPTION.getTestClass().getName() + "." + DriverService.DESCRIPTION.getMethodName());
		moduleInfo.setModuleStatus(RunStatus.RUNNING);
		moduleInfo.setModuleResult(RunResult.RUNNING);
		moduleInfo.setModuleRun(true);
		moduleInfo.setModuleStartTime(new Date());
		addModuleCases(moduleInfo);// addModuleCases(moduleInfo, //
		updateSuiteInfo(LogModule.SUITE_INFO, moduleInfo); // base, des);
		System.out.println("测试模块（" + moduleInfo.getModuleName() + "）执行开始。");
	}

	@Override
	public void caseRunStart(CaseInfo caseInfo) {
		LogModule.logStepInfoList.clear();
		caseInfo.setCaseResult(RunResult.RUNNING);
		caseInfo.setCaseStatus(RunStatus.RUNNING);
		caseInfo.setCaseStartTime(new Date());
		updateModuleInfo(LogModule.MODULE_INFO, caseInfo);
		System.out.println("测试用例（" + caseInfo.getCaseName() + "）执行开始。");
	}

	@Override
	public void suiteRunStop(SuiteInfo suiteInfo) {
		suiteInfo.setSuiteResult(RunResult.PASS);
		suiteInfo.setSuiteStatus(RunStatus.COMPLETED);
		for (ModuleInfo moduleInfo : suiteInfo.getModuleInfoList()) {
			if (moduleInfo.getModuleStatus() != RunStatus.COMPLETED) {
				suiteInfo.setSuiteStatus(RunStatus.RUNNING);
				suiteInfo.setSuiteResult(RunResult.RUNNING);
				log.info("-----------" + moduleInfo.getModuleName());
				break;
			}
		}
		if (suiteInfo.getSuiteStatus() == RunStatus.COMPLETED) {
			for (ModuleInfo moduleInfo : suiteInfo.getModuleInfoList()) {
				if (moduleInfo.getModuleResult() != RunResult.PASS) {//缺陷修改(2015-01-14)，原来是：moduleInfo.getModuleResult() != RunResult.FAIL
					suiteInfo.setSuiteResult(RunResult.FAIL);
					break;
				}
			}

			System.out.println("测试套件（" + suiteInfo.getSuiteName() + "）执行结束。");
			suiteInfo.setSuiteStopTime(new Date());
			NumberFormat numberFormat = NumberFormat.getInstance();
			numberFormat.setMaximumFractionDigits(2);
			float percent = 100;
			if (suiteInfo.getSuiteCaseNum() > 0) {
				percent = Float.parseFloat(numberFormat.format((float) suiteInfo.getPassCaseNum() / (float) suiteInfo.getSuiteCaseNum() * 100));
			}
			suiteInfo.setPassPercent(percent);
		}
	}

	@Override
	public void moduleRunStop(ModuleInfo moduleInfo) {
		moduleInfo.setModuleStatus(RunStatus.COMPLETED);
		moduleInfo.setModuleResult(RunResult.PASS);
		moduleInfo.setModuleStopTime(new Date());
		NumberFormat numberFormat = NumberFormat.getInstance();
		numberFormat.setMaximumFractionDigits(2);
		float percent = 100;
		if (moduleInfo.getModuleCaseNum() > 0) {
			percent = Float.parseFloat(numberFormat.format((float) moduleInfo.getPassCaseNum() / (float) moduleInfo.getModuleCaseNum() * 100));
		}
		moduleInfo.setPassPercent(percent);
		for (CaseInfo caseinfo : moduleInfo.getCaseInfoList()) {
			if (caseinfo.getCaseResult() == RunResult.FAIL) {
				moduleInfo.setModuleResult(RunResult.FAIL);
				break;
			}
		}
		updateSuiteInfo(LogModule.SUITE_INFO, moduleInfo);
		System.out.println("测试模块（" + moduleInfo.getModuleName() + "）执行结束。");
	}

	@Override
	public void caseRunStop(CaseInfo caseInfo) {
		// 判断logstep中有fail没
		caseInfo.setCaseResult(RunResult.PASS);
		caseInfo.setCaseStopTime(new Date());
		caseInfo.setLogStepInfoList(LogModule.logStepInfoList);
		for (LogStepInfo logStepInfo : LogModule.logStepInfoList) {
			if (logStepInfo.getStepResult() == RunResult.FAIL) {
				caseInfo.setCaseResult(RunResult.FAIL);
				break;
			}
		}
		String logLevel = properties.getProperty("logLevel", "error");
		// 当用例运行成功时，只保留检验日志，删除控件操作等其他日志
		if (logLevel.toLowerCase().trim().equals("error") && caseInfo.getCaseResult() == RunResult.PASS) {

			List<LogStepInfo> neededLogStepInfoList = new ArrayList<LogStepInfo>();
			int id = 0;
			for (LogStepInfo logStepInfo : LogModule.logStepInfoList) {
				// 只保留断言的日志
				if (logStepInfo.getStepType() == StepType.ASSERT) {
					id++;
					logStepInfo.setStepId(id);
					neededLogStepInfoList.add(logStepInfo);
				}
			}
			LogModule.logStepInfoList = neededLogStepInfoList;
		}
		caseInfo.setLogStepInfoList(LogModule.logStepInfoList);
		caseInfo.setCaseStatus(RunStatus.COMPLETED);
		System.out.println("测试用例（" + caseInfo.getCaseName() + "）执行结束。");
		caseInfo.setCaseRerunNum(caseInfo.getCaseRerunNum() - 1);
		caseInfo.setCaseRunNum(caseInfo.getCaseRunNum() - 1);
		if (caseInfo.getCaseResult().equals(RunResult.PASS) || caseInfo.getCaseRerunNum() < 0) {
			updateModuleInfo(LogModule.MODULE_INFO, caseInfo);
		}
	}

	/**
	 * 获取webSite.properties里的公共信息
	 * @return
	 */
	@SuppressWarnings("rawtypes")
	private Map<String, Map<String, String>> getWebsiteMap() {
		Map<String, String> websiteInfoMap = new HashMap<String, String>();
		Map<String, Map<String, String>> map = new HashMap<String, Map<String, String>>();
		Properties properties = new Properties();
		if (websiteMap.isEmpty()) {
			try {
				properties = FileUtil.getProperties("/website.properties");
			} catch (Exception e) {
				return map;
			}
			Iterator iter = properties.entrySet().iterator();
			while (iter.hasNext()) {
				Map.Entry entry = (Map.Entry) iter.next();
				String websiteKey = (String) entry.getKey();
				if (websiteKey != null && !websiteKey.trim().equals("") && websiteKey.contains(".")) {
					String websiteName = websiteKey.split("\\.")[0].trim();
					if (websiteMap.get(websiteName) == null) {
						websiteMap.put(websiteName, new HashMap<String, String>());
					}
					String websiteInfoKey = websiteKey.split("\\.")[1].trim();
					String websiteInfoValue = ((String) entry.getValue()).trim();
					websiteInfoMap = websiteMap.get(websiteName);
					websiteInfoMap.put(websiteInfoKey, websiteInfoValue);
					map.put(websiteName, websiteInfoMap);
				}
			}
			websiteMap = map;
		}
		return map;
	}

	@SuppressWarnings({ "unused", "rawtypes" })
	private Map<String, Map<String, String>> getPropertiesMap() {
		Map<String, String> websiteInfoMap = new HashMap<String, String>();
		Map<String, Map<String, String>> map = new HashMap<String, Map<String, String>>();
		Properties properties = new Properties();
		if (websiteMap.isEmpty()) {
			try {
				properties = FileUtil.getProperties("/webSite.properties");
			} catch (Exception e) {
				return map;
			}
			Iterator iter = properties.entrySet().iterator();
			while (iter.hasNext()) {
				Map.Entry entry = (Map.Entry) iter.next();
				String websiteKey = (String) entry.getKey();
				if (websiteKey != null && !websiteKey.trim().equals("") && websiteKey.contains(".")) {
					String[] aa = websiteKey.split("\\.");
					String websiteName = websiteKey.split("\\.")[0].trim();
					if (websiteMap.get(websiteName) == null) {
						websiteMap.put(websiteName, new HashMap<String, String>());
					}
					String websiteInfoKey = websiteKey.split("\\.")[1].trim();
					String websiteInfoValue = ((String) entry.getValue()).trim();
					websiteInfoMap = websiteMap.get(websiteName);
					websiteInfoMap.put(websiteInfoKey, websiteInfoValue);

					map.put(websiteName, websiteInfoMap);
				}
			}
			websiteMap = map;
		}
		return map;
	}

	/**
	 * 合并计划与模块，找设置与模块信息的交集做为新的测试计划
	 * 通过设定的测试计划(testPlan)生成
	 * @param moduleInfoList
	 * @return
	 */
	private List<ModuleInfo> mergePlanModuleInfo(List<ModuleInfo> moduleInfoList) {
		List<ModuleInfo> planedModuleInfoList = new ArrayList<ModuleInfo>();
		/**
		 * 获取test.properties中plan类型，如果没有就默认为create
		 */
		String testPlan = properties.getProperty("plan", "create");
		if (testPlan.trim().toLowerCase().equals("error")) {
			System.out.println("提示：本次计划执行的是plan.error.xml中的用例.");
			try{
				planedSuiteInfo = XmlUtil.parsePlanXml(sourcesDir + "plan.error.xml");
			}catch(DocumentException e){
				planedSuiteInfo = new SuiteInfo();
			}
			if (planedSuiteInfo.getModuleInfoList() == null || planedSuiteInfo.getModuleInfoList().size() == 0) {
				System.out.println("提示：plan.error.xml中无用例数据.");
				planedSuiteInfo.setModuleInfoList(planedModuleInfoList);
				return planedModuleInfoList;
			}
		} else if (testPlan.trim().toLowerCase().equals("all")) {
			System.out.println("提示：本次计划执行的是plan.xml中的用例.");
			try{
				planedSuiteInfo = XmlUtil.parsePlanXml(sourcesDir + "plan.xml");
			}catch(DocumentException e){
				planedSuiteInfo = new SuiteInfo();
			}
			if (planedSuiteInfo.getModuleInfoList() == null || planedSuiteInfo.getModuleInfoList().size() == 0) {
				System.out.println("提示：由于您未手动增加plan.xml,默认执行所有.");
				planedSuiteInfo.setModuleInfoList(planedModuleInfoList);
				return moduleInfoList;
			}
		} else if (testPlan.contains("{")) {
			System.out.println("testPlan:" + testPlan);
			SuiteInfo suite = (SuiteInfo) JSON.parseObject(testPlan, SuiteInfo.class);
			/**
			 * 处理JSON中只传类名，且不传方法名的情况，如下：
			 * 单个类：mvn test -Dplan={\"suiteName\":\"AutoSky\",\"moduleInfoList\":[{\"moduleName\":\"TestJunit\"}]} -DsuiteMerged=true
			 * 多个类：mvn test -Dplan={\"suiteName\":\"AutoSky\",\"moduleInfoList\":[{\"moduleName\":\"TestJunit\"},{\"moduleName\":\"TestJunit2\"}]} -DsuiteMerged=true
			 */
			boolean flag = false;
			List<ModuleInfo> tempModuleList = new ArrayList<ModuleInfo>();
			for(ModuleInfo module : moduleInfoList){
				for(ModuleInfo jModule : suite.getModuleInfoList()){
					if(module.getModuleName().split("\\.")[0].equals(jModule.getModuleName())){
						tempModuleList.add(module);
						flag = true;
					}
					// 处理部分传了类名加方法名的情况，即"TestJunit.testExists"，此情况直接写入临时集合
					if(module.getModuleName().equals(jModule.getModuleName())){
						tempModuleList.add(jModule);
					}
				}
			}
			if(flag){
				suite.getModuleInfoList().clear();
				suite.getModuleInfoList().addAll(tempModuleList);
			}
			
			/*for(ModuleInfo module : moduleInfoList){
			System.out.println("module.getModuleName():" + module.getModuleName());
			}
			for(ModuleInfo jModule : suite.getModuleInfoList()){
				System.out.println("jModule.getModuleName():" + jModule.getModuleName());
			}*/
			
			
			for (int i = 0; i < suite.getModuleInfoList().size(); i++) {
				for(int j = 0; j < moduleInfoList.size(); j++){
					if (moduleInfoList.get(j).getModuleName().equals(suite.getModuleInfoList().get(i).getModuleName())) {
						suite.getModuleInfoList().get(i).setModuleShowName(moduleInfoList.get(j).getModuleShowName());
						suite.getModuleInfoList().get(i).setModuleAuthor(moduleInfoList.get(j).getModuleAuthor());
						suite.getModuleInfoList().get(i).setModuleExpert(moduleInfoList.get(j).getModuleExpert());
						break;
					}
				}
			}
			planedSuiteInfo.setModuleInfoList(suite.getModuleInfoList());
			return suite.getModuleInfoList();
		} else {
			planedSuiteInfo.setModuleInfoList(moduleInfoList);
			return moduleInfoList;
		}
		// 没有plan.xml（plan.error.xml）或者模块为空
		/**
		 * 遍历plan.xml或plan.error.xml转换成的planedSuiteInfo里的模块名
		 */
		for (ModuleInfo planedModuleInfo : planedSuiteInfo.getModuleInfoList()) {
			if (planedModuleInfo.getModuleRun()) {
				for (ModuleInfo moduleInfo : moduleInfoList) {
					if (moduleInfo.getModuleName().equals(planedModuleInfo.getModuleName())) {
						planedModuleInfoList.add(moduleInfo);
						break;
					}
				}
			}
		}
		return planedModuleInfoList;
	}

	/**
	 * 将用例信息加入到模块中
	 * 
	 * @param moduleInfo
	 * @return
	 */
	public List<CaseInfo> addModuleCases(ModuleInfo moduleInfo) {
		List<Map<String, String>> datas = getDatas(DriverService.DESCRIPTION);
		List<CaseInfo> caseInfoList = new ArrayList<CaseInfo>();
		if (datas.size() == 0) {
			Map<String, String> data = new HashMap<String, String>();
			data.put("caseName", "用例一");
			data.put("caseDescription", "用例一描述信息");
			datas.add(data);
		}
		for (Map<String, String> data : datas) {
			String website = properties.getProperty("website");
			if (website != null && !websiteMap.isEmpty()) {
				Map<String, String> websiteInfo = websiteMap.get(website);
				if (websiteInfo != null && !websiteInfo.isEmpty())
					data.putAll(websiteInfo);
			}
			/**
			 * 多浏览器运行时，会把一个用例一分二或为三，形成不同的用例(即只有用例名(用例名为原始用例名_浏览器类型)与浏览器类型不一样)
			 */
			List<Map<String, String>> browserTypeCases = new ArrayList<Map<String, String>>();
			if (properties.getProperty("testSuiteType") != null && !properties.getProperty("testSuiteType").trim().equals("") && TestSuiteType.valueOf(properties.getProperty("testSuiteType")) != TestSuiteType.WEB_UI) {
				data.remove("browserType");
				//System.out.println("test.properties中设置了非WEB_UI，因此移除了browserType属性。");
				browserTypeCases.add(data);
			} else {
				browserTypeCases = getBrowserTypeCases(data);
			}
			for (Map<String, String> browserTypeCase : browserTypeCases) {
				CaseInfo caseInfo = new CaseInfo();
				String caseName = data.get("caseName");
				if (caseName != null) {
					Pattern pa = Pattern.compile("\\s*|\t|\r|\n");
					Matcher m = pa.matcher(caseName);
					//去掉空格等特殊字符
					caseName = m.replaceAll("");
				} else {
					caseName = "无用例名称";
				}
				if(properties.getProperty("testSuiteType") == null || properties.getProperty("testSuiteType").trim().equals("") || TestSuiteType.valueOf(properties.getProperty("testSuiteType")) == TestSuiteType.WEB_UI){
					//caseName = caseName + "_" + browserTypeCase.get("browserType");
					caseInfo.setCaseBrowserType(BrowserType.valueOf(browserTypeCase.get("browserType")));
				}
				caseInfo.setModuleInfo(moduleInfo);
				caseInfo.setSuiteInfo(moduleInfo.getSuiteInfo());
				caseInfo.setCaseName(caseName);
				caseInfo.setCaseAssert(data.get("caseAssert")==null?true:data.get("caseAssert").equals("Y")?true:false);
				caseInfo.setCaseDesc(data.get("caseDescription"));
				caseInfo.setCaseRun(data.get("caseRun") == null || "Y".equals(data.get("caseRun").toUpperCase().trim()));
				caseInfo.setCasePriority(data.get("casePriority"));
				String caseRerunNum = (properties.getProperty("caseRerunNum") == null || properties.getProperty("caseRerunNum").trim().equals("")) ? "0" : properties.getProperty("caseRerunNum");
				caseInfo.setCaseRerunNum(Integer.parseInt(caseRerunNum));
				String caseRunNum = (properties.getProperty("caseRunNum") == null || properties.getProperty("caseRunNum").trim().equals("")) ? "0" : properties.getProperty("caseRunNum");
				caseInfo.setCaseRunNum(Integer.parseInt(caseRunNum));
				caseInfo.setCaseResult(RunResult.WAITING);
				caseInfo.setCaseStatus(RunStatus.WAITING);
				caseInfo.setCaseParams(data);
				caseInfo.setModuleInfo(moduleInfo);
				caseInfoList.add(caseInfo);
			}
		}

		moduleInfo.setCaseInfoList(caseInfoList);
		List<CaseInfo> planedCaseInfoList = mergePlanCaseInfo(moduleInfo);
		moduleInfo.setCaseInfoList(planedCaseInfoList);
		return caseInfoList;
	}

	/**
	 * 合并计划与用例信息
	 * 
	 * @param moduleInfo
	 * @return
	 */
	private List<CaseInfo> mergePlanCaseInfo(ModuleInfo moduleInfo) {
		List<CaseInfo> planedCaseInfoList = new ArrayList<CaseInfo>();
		ModuleInfo planedModuleInfo = null;
		if (planedSuiteInfo.getModuleInfoList().size() == 0)
			return moduleInfo.getCaseInfoList();
		for (ModuleInfo module : planedSuiteInfo.getModuleInfoList()) {

			if (moduleInfo.getModuleName().equals(module.getModuleName())) {
				planedModuleInfo = module;
				break;
			}
		} // 没有plan.xml（plan.error.xml）或者模块为空
		if (planedModuleInfo == null) {
			return planedCaseInfoList;
		}
		if (planedModuleInfo.getCaseInfoList() == null || planedModuleInfo.getCaseInfoList().size() == 0) {
			return moduleInfo.getCaseInfoList();
		}
		for (CaseInfo planedCaseInfo : planedModuleInfo.getCaseInfoList()) {
			if (planedCaseInfo.getCaseRun() == null || planedCaseInfo.getCaseRun()) {
				for (CaseInfo caseInfo : moduleInfo.getCaseInfoList()) {
					if (caseInfo.getCaseName().toUpperCase().equals(planedCaseInfo.getCaseName().toUpperCase())) {
						planedCaseInfoList.add(caseInfo);
						break;
					}
				}
			}
		}
		return planedCaseInfoList;
	}

	private List<Map<String, String>> getBrowserTypeCases(Map<String, String> data) {
		List<Map<String, String>> browserTypeCases = new ArrayList<Map<String, String>>();
		if (data.get("browserType") != null) {
			String[] types = data.get("browserType").split(",");
			for (String type : types) {
				Map<String, String> browserTypeCase = new HashMap<String, String>();
				browserTypeCase.putAll(data);
				if (type.trim().toUpperCase().equals("IE")) {
					browserTypeCase.put("browserType", "IE");
				} else if (type.trim().toUpperCase().equals("FF") || type.toUpperCase().equals("FIREFOX")) {
					browserTypeCase.put("browserType", "FIREFOX");
				} else if (type.trim().toUpperCase().equals("CHROME")) {
					browserTypeCase.put("browserType", "CHROME");
				} else if (type.trim().equals("")) {
					browserTypeCase.put("browserType", "IE");
				} else {
					System.err.println("设置的浏览器类型为:" + type + "，格式有问题！请参照：单个浏览器如IE；多个浏览器用英文逗号分隔，如IE,FF,CHROME。");
					throw new RuntimeException("数据池中的browserType格式有问题，请改正！\n设置的浏览器类型为:" + type + "，格式有问题！请参照：单个浏览器如IE；多个浏览器用英文逗号分隔，如IE,FF,CHROME。");
				}
				browserTypeCases.add(browserTypeCase);
			}
		} else {
			data.put("browserType", "IE");
			browserTypeCases.add(data);
		}
		return browserTypeCases;
	}

	/**
	 * 将moduleInfo中的字段状态更新到suiteInfo对象中，并为LogModule.SUITE_INFO重新赋值
	 * 
	 * @param suiteInfo
	 * @param currentModuleInfo
	 */
	public void updateSuiteInfo(SuiteInfo suiteInfo, ModuleInfo currentModuleInfo) {
		List<ModuleInfo> moduleInfoList = suiteInfo.getModuleInfoList();
		for (int i = 0; i < moduleInfoList.size(); i++) {
			if (moduleInfoList.get(i).getModuleName().equals(currentModuleInfo.getModuleName())) {
				if (currentModuleInfo.getModuleStatus() == RunStatus.COMPLETED) {
					suiteInfo.setPassCaseNum(suiteInfo.getPassCaseNum() + currentModuleInfo.getPassCaseNum());
					suiteInfo.setFailCaseNum(suiteInfo.getFailCaseNum() + currentModuleInfo.getFailCaseNum());
					suiteInfo.setSuiteCaseNum(suiteInfo.getSuiteCaseNum() + currentModuleInfo.getModuleCaseNum());
					suiteInfo.setSuiteModuleNum(suiteInfo.getSuiteModuleNum() + 1);
				}
				currentModuleInfo.setModuleShowName(moduleInfoList.get(i).getModuleShowName());
				currentModuleInfo.setModuleAuthor(moduleInfoList.get(i).getModuleAuthor());
				currentModuleInfo.setModuleExpert(moduleInfoList.get(i).getModuleExpert());
				moduleInfoList.set(i, currentModuleInfo);
				break;
			}
		}
		suiteInfo.setModuleInfoList(moduleInfoList);
		LogModule.SUITE_INFO = suiteInfo;
	}

	/**
	 * 更新模块中的用例执行结果与用例数
	 * 
	 * @param moduleInfo
	 * @param currentCaseInfo
	 */
	public void updateModuleInfo(ModuleInfo moduleInfo, CaseInfo currentCaseInfo) {
		List<CaseInfo> caseInfoList = moduleInfo.getCaseInfoList();
		for (int i = 0; i < caseInfoList.size(); i++) {
			if (caseInfoList.get(i).getCaseName().equals(currentCaseInfo.getCaseName())) {
				if (currentCaseInfo.getCaseResult() == RunResult.PASS) {
					moduleInfo.setPassCaseNum(moduleInfo.getPassCaseNum() + 1);
					moduleInfo.setModuleCaseNum(moduleInfo.getModuleCaseNum() + 1);
				} else if (currentCaseInfo.getCaseResult() == RunResult.FAIL) {
					moduleInfo.setFailCaseNum(moduleInfo.getFailCaseNum() + 1);
					moduleInfo.setModuleCaseNum(moduleInfo.getModuleCaseNum() + 1);
				}
				caseInfoList.set(i, currentCaseInfo);
				break;
			}
		}
		moduleInfo.setCaseInfoList(caseInfoList);
		LogModule.MODULE_INFO = moduleInfo;
	}

	/**
	 * 获取Test Class里面有@Test的Method列表
	 * 
	 * @param Clazz
	 * @return
	 */
	private List<ModuleInfo> getSuiteModuleNameList(Class<?> Clazz) {
		List<ModuleInfo> moduleNameList = new ArrayList<ModuleInfo>();
		Method[] methods = Clazz.getMethods();
		// Map<String, String> map = new LinkedHashMap<String, String>();
		for (Method method : methods) {
			ModuleInfo moduleInfo = new ModuleInfo();
			String methodName = method.getName();
			DataSource dataSource = method.getAnnotation(DataSource.class);
			if (method.isAnnotationPresent(org.junit.Test.class) && method.getAnnotation(Ignore.class) == null) {
				moduleInfo.setModuleName(methodName);
				if(null != dataSource){
					moduleInfo.setModuleShowName(dataSource.showName());
					moduleInfo.setModuleAuthor(dataSource.author());
					moduleInfo.setModuleExpert(dataSource.expert());
				}
				moduleNameList.add(moduleInfo);
			}
		}
		return moduleNameList;
	}

	public List<Map<String, String>> getDatas(Description des) {
		List<Map<String, String>> datas = new ArrayList<Map<String, String>>();
		if (des.getAnnotation(DataSource.class) == null) {
			return datas;
		}
		DataSourceType type = ((DataSource) des.getAnnotation(DataSource.class)).type();
		if (type == null) {
			return datas;
		}
		String fileName = ((DataSource) des.getAnnotation(DataSource.class)).file().toString().trim();
		String filePath = sourcesDir + des.getTestClass().getSimpleName() + "/" + fileName;
		String sourcePrefix = ((DataSource)des.getAnnotation(DataSource.class)).sourcePrefix();
		if (!sourcePrefix.trim().equals("")) {
			filePath = sourcesDir  + "/" +  sourcePrefix.trim()  + "/" +  des.getTestClass().getSimpleName() + "/" + fileName;
		}
		switch (type) {
		case EXCEL:
			String sheetName = ((DataSource) des.getAnnotation(DataSource.class)).sheetName().toString().trim();
			ExcelUtil excelUtil = new ExcelUtil();
			datas = excelUtil.excelDatas(filePath, sheetName);
			break;
		case XML:
			datas = XmlUtil.xmlDatas(filePath);
			break;
		case CSV:
			break;
		default:
		}
		return datas;
	}

	/**
	 * 获取Test Class里面有@Test的Method列表
	 * @param Clazz
	 * @return
	 */
	@SuppressWarnings("unused")
	private static Map<String, String> getTestMethods(Class<?> Clazz) {
		Method[] methods = Clazz.getMethods();
		Map<String, String> map = new LinkedHashMap<String, String>();
		int i = 0;
		for (Method method : methods) {
			String methodName = method.getName();
			if (method.isAnnotationPresent(org.junit.BeforeClass.class) || method.isAnnotationPresent(org.junit.AfterClass.class)) {
				map.put(method.getAnnotations()[0].annotationType().getSimpleName(), methodName);
			} else if (method.isAnnotationPresent(org.junit.Test.class)) {
				map.put(String.valueOf(i), methodName);
				i++;
			}
		}
		return map;
	}

}
