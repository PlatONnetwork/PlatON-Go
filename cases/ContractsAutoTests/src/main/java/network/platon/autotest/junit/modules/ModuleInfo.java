package network.platon.autotest.junit.modules;

import java.util.Date;
import java.util.List;

import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.RunStatus;
import network.platon.autotest.utils.DateUtil;

public class ModuleInfo {

	private SuiteInfo suiteInfo;
	private String moduleName;
	/**
	 * 为了日志报告中模块名称为中文显示而定义的模块显示名称
	 */
	private String moduleShowName;
	private String moduleDesc;
	/**
	 * 脚本负责人
	 */
	private String moduleAuthor;
	/**
	 * 业务负责人
	 */
	private String moduleExpert;
	private Boolean moduleRun;
	private RunStatus moduleStatus;
	private RunResult moduleResult;
	private String moduleBrowser;
	private List<CaseInfo> caseInfoList;
	private int moduleCaseNum;
	private int passCaseNum;
	private int failCaseNum;
	private float passPercent;
	private Date moduleStartTime;
	private Date moduleStopTime;

	public SuiteInfo getSuiteInfo() {
		return suiteInfo;
	}

	public void setSuiteInfo(SuiteInfo suiteInfo) {
		this.suiteInfo = suiteInfo;
	}

	public RunStatus getModuleStatus() {
		return moduleStatus;
	}

	public void setModuleStatus(RunStatus moduleStatus) {
		this.moduleStatus = moduleStatus;
	}

	public String getModuleName() {
		return moduleName;
	}

	public void setModuleName(String moduleName) {
		this.moduleName = moduleName;
	}

	public RunResult getModuleResult() {
		return moduleResult;
	}

	public void setModuleResult(RunResult moduleResult) {
		this.moduleResult = moduleResult;
	}

	public String getModuleBrowser() {
		return moduleBrowser;
	}

	public void setModuleBrowser(String moduleBrowser) {
		this.moduleBrowser = moduleBrowser;
	}

	public List<CaseInfo> getCaseInfoList() {
		return caseInfoList;
	}

	public void setCaseInfoList(List<CaseInfo> caseInfoList) {
		this.caseInfoList = caseInfoList;
	}

	public int getModuleCaseNum() {
		return moduleCaseNum;
	}

	public void setModuleCaseNum(int moduleCaseNum) {
		this.moduleCaseNum = moduleCaseNum;
	}

	public int getPassCaseNum() {
		return passCaseNum;
	}

	public void setPassCaseNum(int passCaseNum) {
		this.passCaseNum = passCaseNum;
	}

	public int getFailCaseNum() {
		return failCaseNum;
	}

	public void setFailCaseNum(int failCaseNum) {
		this.failCaseNum = failCaseNum;
	}

	public float getPassPercent() {
		return passPercent;
	}

	public void setPassPercent(float passPercent) {
		this.passPercent = passPercent;
	}

	public Date getModuleStopTime() {
		return moduleStopTime;
	}

	public void setModuleStopTime(Date moduleStopTime) {
		this.moduleStopTime = moduleStopTime;
	}

	public Date getModuleStartTime() {
		return moduleStartTime;
	}

	public void setModuleStartTime(Date moduleStartTime) {
		this.moduleStartTime = moduleStartTime;
	}

	public String getModuleAuthor() {
		return moduleAuthor;
	}

	public void setModuleAuthor(String moduleAuthor) {
		this.moduleAuthor = moduleAuthor;
	}

	public String getModuleBusinessAuthor() {
		return moduleExpert;
	}

	public void setModuleBusinessAuthor(String moduleExpert) {
		this.moduleExpert = moduleExpert;
	}

	public void setModuleDesc(String moduleDesc) {
		this.moduleDesc = moduleDesc;
	}

	public String getModuleDesc() {
		return moduleDesc;
	}

	public String getModuleRunTimeStr() {
		long times = (moduleStopTime.getTime() - moduleStartTime.getTime()) / 1000;
		long minutes = times / 60;
		long second = times % 60;
		String timeStr = String.valueOf(second) + "秒";
		long hour = minutes / 60;
		long minute = minutes % 60;
		if (hour > 0) {
			timeStr = String.valueOf(hour) + "小时" + String.valueOf(minute) + "分" + timeStr;
		} else if (minute > 0) {
			timeStr = String.valueOf(minute) + "分" + timeStr;
		}
		return timeStr;
	}

	public String getModuleStartTimeStr() {
		return DateUtil.dateToStr(moduleStartTime, "HH:mm:ss");
	}

	public String getModuleStopTimeStr() {
		return DateUtil.dateToStr(moduleStopTime, "HH:mm:ss");
	}

	public void setModuleRun(Boolean moduleRun) {
		this.moduleRun = moduleRun;
	}

	public Boolean getModuleRun() {
		return moduleRun;
	}

	public String getModuleShowName() {
		return moduleShowName;
	}

	public void setModuleShowName(String moduleShowName) {
		this.moduleShowName = moduleShowName;
	}

	public String getModuleExpert() {
		return moduleExpert;
	}

	public void setModuleExpert(String moduleExpert) {
		this.moduleExpert = moduleExpert;
	}

	@Override
	public String toString() {
		return "ModuleInfo [suiteInfo=" + suiteInfo + ", moduleName=" + moduleName + ", moduleShowName=" + moduleShowName + ", moduleDesc=" + moduleDesc + ", moduleAuthor=" + moduleAuthor + ", moduleExpert=" + moduleExpert + ", moduleRun=" + moduleRun + ", moduleStatus=" + moduleStatus + ", moduleResult=" + moduleResult + ", moduleBrowser=" + moduleBrowser + ", caseInfoList=" + caseInfoList + ", moduleCaseNum=" + moduleCaseNum + ", passCaseNum=" + passCaseNum + ", failCaseNum=" + failCaseNum + ", passPercent=" + passPercent + ", moduleStartTime=" + moduleStartTime + ", moduleStopTime=" + moduleStopTime + "]";
	}
	
}
