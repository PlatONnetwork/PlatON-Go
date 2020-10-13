package network.platon.autotest.junit.modules;

import java.util.Date;
import java.util.List;
import java.util.Map;

import network.platon.autotest.junit.enums.BrowserType;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.RunStatus;
import network.platon.autotest.utils.DateUtil;

public class CaseInfo {
	private SuiteInfo suiteInfo;
	private ModuleInfo moduleInfo;
	private String caseName;
	private String caseDesc;
	private Boolean caseRun;
	private RunStatus caseStatus;
	private RunResult caseResult;
	private Map<String, String> caseParams;
	private String caseLevel;
	private BrowserType caseBrowserType;
	private Date caseStartTime;
	private Date caseStopTime;
	private String casePriority;
	/**
	 * 错误用例重试次数
	 */
	private int caseRerunNum;
	private List<LogStepInfo> LogStepInfoList;
	/**
	 * 是否在测试流程执行过程中运行断言
	 */
	private boolean caseAssert;
	/**
	 * 用例运行次数
	 */
	private Integer caseRunNum;

	public RunStatus getCaseStatus() {
		return caseStatus;
	}

	public void setCaseStatus(RunStatus caseStatus) {
		this.caseStatus = caseStatus;
	}

	public RunResult getCaseResult() {
		return caseResult;
	}

	public void setCaseResult(RunResult caseResult) {
		this.caseResult = caseResult;
	}

	public SuiteInfo getSuiteInfo() {
		return suiteInfo;
	}

	public void setSuiteInfo(SuiteInfo suiteInfo) {
		this.suiteInfo = suiteInfo;
	}

	public ModuleInfo getModuleInfo() {
		return moduleInfo;
	}

	public void setModuleInfo(ModuleInfo moduleInfo) {
		this.moduleInfo = moduleInfo;
	}

	public String getCaseName() {
		return caseName;
	}

	public void setCaseName(String caseName) {
		this.caseName = caseName;
	}

	public String getCaseDesc() {
		return caseDesc;
	}

	public void setCaseDesc(String caseDesc) {
		this.caseDesc = caseDesc;
	}

	public Boolean getCaseRun() {
		return caseRun;
	}

	public void setCaseRun(Boolean caseRun) {
		this.caseRun = caseRun;
	}

	public Map<String, String> getCaseParams() {
		return caseParams;
	}

	public void setCaseParams(Map<String, String> caseParams) {
		this.caseParams = caseParams;
	}

	public String getCaseLevel() {
		return caseLevel;
	}

	public void setCaseLevel(String caseLevel) {
		this.caseLevel = caseLevel;
	}

	public void setCaseStartTime(Date caseStartTime) {
		this.caseStartTime = caseStartTime;
	}

	public Date getCaseStartTime() {
		return caseStartTime;
	}

	public void setCaseStopTime(Date caseStopTime) {
		this.caseStopTime = caseStopTime;
	}

	public Date getCaseStopTime() {
		return caseStopTime;
	}

	public String getCaseRunTimeStr() {
		return String.valueOf((caseStopTime.getTime() - caseStartTime.getTime()) / 1000);
	}

	public String getCaseStartTimeStr() {
		return DateUtil.dateToStr(caseStartTime, "HH:mm:ss");
	}

	public String getCaseStopTimeStr() {
		return DateUtil.dateToStr(caseStopTime, "HH:mm:ss");
	}

	public void setCaseBrowserType(BrowserType caseBrowserType) {
		this.caseBrowserType = caseBrowserType;
	}

	public BrowserType getCaseBrowserType() {
		return caseBrowserType;
	}

	public void setCasePriority(String casePriority) {
		this.casePriority = casePriority;
	}

	public String getCasePriority() {
		return casePriority;
	}

	public void setCaseRerunNum(int caseRerunNum) {
		this.caseRerunNum = caseRerunNum;
	}

	public int getCaseRerunNum() {
		return caseRerunNum;
	}

	public void setLogStepInfoList(List<LogStepInfo> logStepInfoList) {
		LogStepInfoList = logStepInfoList;
	}

	public List<LogStepInfo> getLogStepInfoList() {
		return LogStepInfoList;
	}

	public boolean isCaseAssert() {
		return caseAssert;
	}

	public void setCaseAssert(boolean caseAssert) {
		this.caseAssert = caseAssert;
	}

	public Integer getCaseRunNum() {
		return caseRunNum;
	}

	public void setCaseRunNum(Integer caseRunNum) {
		this.caseRunNum = caseRunNum;
	}
	
	
	// attr_accessor :case_desc
	// attr_accessor :case_run
	// attr_accessor :case_flow
	// attr_accessor :case_params
	// attr_accessor :case_dbparams
	// attr_accessor :case_level #场景级别
	// attr_accessor :case_rerun_num #重跑次数
	// attr_accessor :case_browser_level #场景浏览器级别
	// attr_accessor :case_browser #用
}
