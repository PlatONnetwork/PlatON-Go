package network.platon.autotest.junit.log;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import network.platon.autotest.junit.core.LogModule;
import network.platon.autotest.junit.core.Observer;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.StepType;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.LogStepInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import network.platon.autotest.junit.rules.DriverService;
import org.apache.log4j.Logger;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;
import network.platon.autotest.utils.DateUtil;
import network.platon.autotest.utils.HttpUtil;

@SuppressWarnings("unused")
public class DatabaseLog implements Observer {

	private static Logger logger = Logger.getLogger(DatabaseLog.class);
	private String logTestSuiteId;
	private String logTestModuleId;
	private String logTestCaseId;
	private long logPictureId = 0;

	// public void setLogTestSuiteId(long logTestSuiteId) {
	// this.logTestSuiteId = logTestSuiteId;
	// }
	//
	// public long getLogTestSuiteId() {
	// return logTestSuiteId;
	// }

	// 尝试通过flysky的dblog接口来做
	@Override
	public void suiteRunStart(SuiteInfo suiteInfo) {
		Map<String, String> param = new HashMap<String, String>();
		Map<String, String> params = new HashMap<String, String>();
		String machineIp = "";
		try {
			InetAddress addr = InetAddress.getLocalHost();
			machineIp = addr.getHostAddress().toString();// 获得本机IP
		} catch (UnknownHostException e1) {
			e1.printStackTrace();
		}
		param.put("buildId", String.valueOf(suiteInfo.getBuildId()));
		param.put("buildTaskId", String.valueOf(suiteInfo.getBuildTaskId()));
		param.put("buildTestSuiteId", String.valueOf(suiteInfo.getBuildTestSuiteId()));
		param.put("submitter",suiteInfo.getSubmitter());
		param.put("submitdate",suiteInfo.getSubmitdate());
		param.put("submitnote",suiteInfo.getSubmitnote());
		// 这个要平台那边通过buildTestsuiteId找到testsuiteId，然后存进去，日志这里不做处理
		param.put("testSuiteId", "0");
		param.put("projectName", DriverService.PROPERTIES.getProperty("project"));
		param.put("testSuiteName", suiteInfo.getSuiteName());
		param.put("machineIp", machineIp);
		param.put("runner", DriverService.PROPERTIES.getProperty("runner"));
		param.put("runStatus", String.valueOf(suiteInfo.getSuiteStatus()));
		param.put("runResult", String.valueOf(suiteInfo.getSuiteResult()));
		param.put("caseNum", "0");
		param.put("passNum", "0");
		param.put("failNum", "0");
		String dbLogUrl = DriverService.PROPERTIES.getProperty("dbLogUrl");
		try {//log_suite_start.do
			params.put("cname", "admin");
			params.put("cpwd", "admin");
			//String json1 = HttpUtil.postRequest("http://localhost:8080/AutoTestPlats/demo/" + "userAction!doNotNeedSession_login.action", params, null, null);
			String json = HttpUtil.postRequest(dbLogUrl + "LogTestSuite.action", param, null, null);
			JSONObject result = JSON.parseObject(json);
			if (!result.getBoolean("result")) {
				throw new Exception(result.getString("msg"));
			} else {
				//logTestSuiteId = Long.parseLong(result.getString("msg"));
				logTestSuiteId = result.getString("msg");
				LogModule.logStepPass(StepType.DATABASE, "测试套件（" + suiteInfo.getSuiteName() + "）开始运行时数据库日志存储成功", RunResult.PASS);
			}
		} catch (Exception e) {
			LogModule.logStepFail(StepType.DATABASE, "测试套件（" + suiteInfo.getSuiteName() + "）开始运行时数据库日志存储失败", RunResult.FAIL, e.getMessage());
			// throw new RuntimeException(e.getMessage());
		}
	}

	@Override
	public void suiteRunStop(SuiteInfo suiteInfo) {
		Map<String, String> param = new HashMap<String, String>();
		param.put("runStatus", String.valueOf(suiteInfo.getSuiteStatus()));
		param.put("runResult", String.valueOf(suiteInfo.getSuiteResult()));
		param.put("moduleNum", String.valueOf(suiteInfo.getSuiteModuleNum()));
		param.put("caseNum", String.valueOf(suiteInfo.getSuiteCaseNum()));
		param.put("passNum", String.valueOf(suiteInfo.getPassCaseNum()));
		param.put("failNum", String.valueOf(suiteInfo.getFailCaseNum()));
		param.put("percent", String.valueOf(suiteInfo.getPassPercent()));
		param.put("sid", String.valueOf(logTestSuiteId));
		String dbLogUrl = DriverService.PROPERTIES.getProperty("dbLogUrl");
		try {
			String json = HttpUtil.postRequest(dbLogUrl + "LogTestSuite.action", param, null, null);
			JSONObject result = JSON.parseObject(json);
			if (!result.getBoolean("result")) {
				throw new Exception(result.getString("msg"));
			}
			LogModule.logStepPass(StepType.DATABASE, "测试套件（" + suiteInfo.getSuiteName() + "）结束运行时数据库日志存储成功", RunResult.PASS);
		} catch (Exception e) {
			LogModule.logStepFail(StepType.DATABASE, "测试套件（" + suiteInfo.getSuiteName() + "）结束运行时数据库日志存储失败", RunResult.FAIL, e.getMessage());
			// throw new RuntimeException(e.getMessage());
		}

	}

	@Override
	public void moduleRunStart(ModuleInfo moduleInfo) {
		Map<String, String> param = new HashMap<String, String>();
		//String machineIp = "";
		try {
			InetAddress addr = InetAddress.getLocalHost();
			//machineIp = addr.getHostAddress().toString();// 获得本机IP
		} catch (UnknownHostException e1) {
			// TODO Auto-generated catch block
			e1.printStackTrace();
		}
		//param.put("machineIp", machineIp);
		param.put("sid", String.valueOf(logTestSuiteId));
		param.put("moduleName", moduleInfo.getModuleName());
		param.put("runStatus", String.valueOf(moduleInfo.getModuleStatus()));
		param.put("runResult", String.valueOf(moduleInfo.getModuleResult()));
		String dbLogUrl = DriverService.PROPERTIES.getProperty("dbLogUrl");
		try {
			String json = HttpUtil.postRequest(dbLogUrl + "LogTestModule.action", param, null, null);
			JSONObject result = JSON.parseObject(json);
			if (!result.getBoolean("result")) {
				throw new Exception(result.getString("msg"));
			} else {
				//logTestModuleId = Long.parseLong(result.getString("msg"));
				logTestModuleId = result.getString("msg");
			}
			LogModule.logStepPass(StepType.DATABASE, "测试模块（" + moduleInfo.getModuleName() + "）开始运行时数据库日志存储成功", RunResult.PASS);
		} catch (Exception e) {
			LogModule.logStepFail(StepType.DATABASE, "测试模块（" + moduleInfo.getModuleName() + "）开始运行时数据库日志存储失败", RunResult.FAIL, e.getMessage());
			// throw new RuntimeException(e.getMessage());
		}
	}

	@Override
	public void moduleRunStop(ModuleInfo moduleInfo) {
		Map<String, String> param = new HashMap<String, String>();
		param.put("runStatus", String.valueOf(moduleInfo.getModuleStatus()));
		param.put("runResult", String.valueOf(moduleInfo.getModuleResult()));
		param.put("caseNum", String.valueOf(moduleInfo.getModuleCaseNum()));
		param.put("passNum", String.valueOf(moduleInfo.getPassCaseNum()));
		param.put("failNum", String.valueOf(moduleInfo.getFailCaseNum()));
		param.put("mid", String.valueOf(logTestModuleId));
		String dbLogUrl = DriverService.PROPERTIES.getProperty("dbLogUrl");
		try {
			String json = HttpUtil.postRequest(dbLogUrl + "LogTestModule.action", param, null, null);
			JSONObject result = JSON.parseObject(json);
			if (!result.getBoolean("result")) {
				throw new Exception(result.getString("msg"));
			}
			LogModule.logStepPass(StepType.DATABASE, "测试模块（" + moduleInfo.getModuleName() + "）结束运行时数据库日志存储成功", RunResult.PASS);
		} catch (Exception e) {
			LogModule.logStepFail(StepType.DATABASE, "测试模块（" + moduleInfo.getModuleName() + "）结束运行时数据库日志存储失败", RunResult.FAIL, e.getMessage());
			// throw new RuntimeException(e.getMessage());
		}

	}

	@Override
	public void caseRunStart(CaseInfo caseInfo) {
		Map<String, String> param = new HashMap<String, String>();
		//param.put("logTestSuiteId", String.valueOf(logTestSuiteId));
		param.put("mid", String.valueOf(logTestModuleId));
		param.put("caseName", caseInfo.getCaseName());
		param.put("description", caseInfo.getCaseDesc());
		param.put("runStatus", String.valueOf(caseInfo.getCaseStatus()));
		param.put("runResult", String.valueOf(caseInfo.getCaseResult()));
		String dbLogUrl = DriverService.PROPERTIES.getProperty("dbLogUrl");
		try {
			String json = HttpUtil.postRequest(dbLogUrl + "LogTestCase.action", param, null, null);
			JSONObject result = JSON.parseObject(json);
			if (!result.getBoolean("result")) {
				throw new Exception(result.getString("msg"));
			} else {
				//logTestCaseId = Long.parseLong(result.getString("msg"));
				logTestCaseId = result.getString("msg");
			}
			LogModule.logStepPass(StepType.DATABASE, "测试用例（" + caseInfo.getCaseName() + "）开始运行时数据库日志存储成功", RunResult.PASS);
		} catch (Exception e) {
			LogModule.logStepFail(StepType.DATABASE, "测试用例（" + caseInfo.getCaseName() + "）开始运行时数据库日志存储失败", RunResult.FAIL, e.getMessage());
			// throw new RuntimeException(e.getMessage());
		}

	}

	@Override
	public void caseRunStop(CaseInfo caseInfo) {
		Map<String, String> param = new HashMap<String, String>();
		param.put("runStatus", String.valueOf(caseInfo.getCaseStatus()));
		param.put("runResult", String.valueOf(caseInfo.getCaseResult()));
		param.put("rerunNum", "");
		param.put("caseid", String.valueOf(logTestCaseId));
		String dbLogUrl = DriverService.PROPERTIES.getProperty("dbLogUrl");
		try {
			String json = HttpUtil.postRequest(dbLogUrl + "LogTestCase.action", param, null, null);
			JSONObject result = JSON.parseObject(json);
			if (!result.getBoolean("result")) {
				throw new Exception(result.getString("msg"));
			}
			LogModule.logStepPass(StepType.DATABASE, "测试用例（" + caseInfo.getCaseName() + "）结束运行时数据库日志存储成功", RunResult.PASS);
		} catch (Exception e) {
			LogModule.logStepFail(StepType.DATABASE, "测试用例（" + caseInfo.getCaseName() + "）结束运行时数据库日志存储失败", RunResult.FAIL, e.getMessage());
			// throw new RuntimeException(e.getMessage());
		}
		List<LogStepInfo> logStepInfoList = caseInfo.getLogStepInfoList();
		for (LogStepInfo logStepInfo : logStepInfoList) {
			// 图片保存
			param.clear();
			String name = logStepInfo.getPicture();
			if (name != null && !name.trim().equals("")) {
				param.put("name", logStepInfo.getPictureName());
				// param.put("content", FileUtil.);

				try {
					//String json = HttpUtil.uploadFile(new File(name), dbLogUrl + "upload");
//					JSONObject result = JSON.parseObject(json);
//					if (!result.getBoolean("result")) {
//						throw new Exception(result.getString("msg"));
//					} else {
//						logPictureId = Long.parseLong(result.getString("msg"));
//					}
				} catch (Exception e) {
					LogModule.logStepFail(StepType.DATABASE, "测试用例步骤截图（" + logStepInfo.getPicture() + "）日志截图存储失败", RunResult.FAIL, e.getMessage());
					// throw new RuntimeException(e.getMessage());
				}
			} else {
				logPictureId = 0;
			}
			// 步骤保存
			param.clear();
			//param.put("logTestSuiteId", String.valueOf(logTestSuiteId));
			//param.put("logTestModuleId", String.valueOf(logTestModuleId));
			param.put("caseid", String.valueOf(logTestCaseId));
			//param.put("stepId", String.valueOf(logStepInfo.getStepId()));
			param.put("stepname", String.valueOf(logStepInfo.getStepDesc()));
			param.put("comment", String.valueOf(logStepInfo.getUrl()));
			param.put("type", String.valueOf(logStepInfo.getStepType()));
			param.put("result", String.valueOf(logStepInfo.getStepResult()));
//			param.put("logPictureId", String.valueOf(logPictureId));
			//param.put("logPictureId", String.valueOf(logStepInfo.getPictureName()));
			param.put("reason", logStepInfo.getFailReason());

			try {
				String today = DateUtil.dateToStr(new Date(), "yyyy-MM-dd");
				Date stepTime1 = DateUtil.strToDate(today + " " + logStepInfo.getStepTime(), "yyyy-MM-dd HH:mm:ss");
				SimpleDateFormat sf = new SimpleDateFormat("yyyy-MM-dd HH:mm:ss");
				String stepTime = sf.format(stepTime1);
				param.put("stepTime", stepTime);

			} catch (ParseException e1) {
				e1.printStackTrace();
			}
//			try {
//				String today = DateUtil.dateToStr(new Date(), "yyyy-MM-dd");
//				Date stepTime = DateUtil.strToDate(today + " " + logStepInfo.getStepTime(), "yyyy-MM-dd HH:mm:ss");
//				param.put("stepTime", stepTime.toGMTString());
//
//			} catch (ParseException e1) {
//				// TODO Auto-generated catch block
//				e1.printStackTrace();
//			}

			try {
				String json = HttpUtil.postRequest(dbLogUrl + "LogTestStep.action", param, null, null);
				JSONObject result = JSON.parseObject(json);
				if (!result.getBoolean("result")) {
					throw new Exception(result.getString("msg"));
				}
			} catch (Exception e) {
				LogModule.logStepFail(StepType.DATABASE, "测试用例步骤（" + logStepInfo.getStepDesc() + "）日志存储失败", RunResult.FAIL, e.getMessage());
				// throw new RuntimeException(e.getMessage());
			}
		}

	}

}
