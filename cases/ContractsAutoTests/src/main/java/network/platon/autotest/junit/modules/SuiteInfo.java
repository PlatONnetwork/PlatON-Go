package network.platon.autotest.junit.modules;

import java.util.Date;
import java.util.List;

import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.RunStatus;
import network.platon.autotest.utils.DateUtil;

/**
 * @author qcxiao
 *
 */
public class SuiteInfo {
	private String projectName;
	private String suiteName;
	private RunStatus suiteStatus;
	private RunResult suiteResult;
	private String suiteBrowser;
	private List<ModuleInfo> moduleInfoList;
	private String suiteAuthor;
	private String suiteExpert;
	private int suiteModuleNum;
	private int suiteCaseNum;
	private int passCaseNum;
	private int failCaseNum;
	private float passPercent;
	private Date suiteStartTime;
	private Date suiteStopTime;
	private Boolean reportMerged;
	private long buildId;
	private long buildTaskId;
	private long buildTestSuiteId;
	private String submitter;//结合Californium平台使用的提交人
	private String submitdate;//结合Californium平台使用的提交时间
	private String submitnote;//结合Californium平台使用的提交信息

	public String getSuiteName() {
		return suiteName;
	}

	public void setSuiteName(String suiteName) {
		this.suiteName = suiteName;
	}

	public RunStatus getSuiteStatus() {
		return suiteStatus;
	}

	public void setSuiteStatus(RunStatus suiteStatus) {
		this.suiteStatus = suiteStatus;
	}

	public RunResult getSuiteResult() {
		return suiteResult;
	}

	public void setSuiteResult(RunResult suiteResult) {
		this.suiteResult = suiteResult;
	}

	public String getSuiteBrowser() {
		return suiteBrowser;
	}

	public void setSuiteBrowser(String suiteBrowser) {
		this.suiteBrowser = suiteBrowser;
	}

	public List<ModuleInfo> getModuleInfoList() {
		return moduleInfoList;
	}

	public void setModuleInfoList(List<ModuleInfo> moduleInfoList) {
		this.moduleInfoList = moduleInfoList;
	}

	public String getSuiteAuthor() {
		return suiteAuthor;
	}

	public void setSuiteAuthor(String suiteAuthor) {
		this.suiteAuthor = suiteAuthor;
	}

	public String getSuiteExpert() {
		return suiteExpert;
	}

	public void setSuiteExpert(String suiteExpert) {
		this.suiteExpert = suiteExpert;
	}

	public int getSuiteCaseNum() {
		return suiteCaseNum;
	}

	public void setSuiteCaseNum(int suiteCaseNum) {
		this.suiteCaseNum = suiteCaseNum;
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

	public Date getSuiteStartTime() {
		return suiteStartTime;
	}

	public void setSuiteStartTime(Date suiteStartTime) {
		this.suiteStartTime = suiteStartTime;
	}

	public Date getSuiteStopTime() {
		return suiteStopTime;
	}

	public void setSuiteStopTime(Date suiteStopTime) {
		this.suiteStopTime = suiteStopTime;
	}

	public String getSuiteRunTimeStr() {
		long times = (suiteStopTime.getTime() - suiteStartTime.getTime()) / 1000;
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

	public String getSuiteStartTimeStr() {
		return DateUtil.dateToStr(suiteStartTime, "HH:mm:ss");
	}

	public String getSuiteStopTimeStr() {
		return DateUtil.dateToStr(suiteStopTime, "HH:mm:ss");
	}

	public void setSuiteModuleNum(int suiteModuleNum) {
		this.suiteModuleNum = suiteModuleNum;
	}

	public int getSuiteModuleNum() {
		return suiteModuleNum;
	}

	public void setReportMerged(Boolean reportMerged) {
		this.reportMerged = reportMerged;
	}

	public Boolean getReportMerged() {
		return reportMerged;
	}

	public void setBuildId(long buildId) {
		this.buildId = buildId;
	}

	public long getBuildId() {
		return buildId;
	}

	public void setBuildTaskId(long buildTaskId) {
		this.buildTaskId = buildTaskId;
	}

	public long getBuildTaskId() {
		return buildTaskId;
	}

	public void setBuildTestSuiteId(long buildTestSuiteId) {
		this.buildTestSuiteId = buildTestSuiteId;
	}

	public long getBuildTestSuiteId() {
		return buildTestSuiteId;
	}

	public String getSubmitter() {
		return submitter;
	}

	public void setSubmitter(String submitter) {
		this.submitter = submitter;
	}

	public String getSubmitdate() {
		return submitdate;
	}

	public void setSubmitdate(String submitdate) {
		this.submitdate = submitdate;
	}

	public String getSubmitnote() {
		return submitnote;
	}

	public void setSubmitnote(String submitnote) {
		this.submitnote = submitnote;
	}

	public String getProjectName() {
		return projectName;
	}

	public void setProjectName(String projectName) {
		this.projectName = projectName;
	}

}
