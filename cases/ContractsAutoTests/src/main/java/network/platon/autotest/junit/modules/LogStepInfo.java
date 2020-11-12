package network.platon.autotest.junit.modules;

import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.StepType;

public class LogStepInfo {
	private SuiteInfo suiteInfo;
	private ModuleInfo moduleInfo;
	private String caseName;
	private int stepId;
	private StepType stepType;
	private String stepDesc;
	private String actual = "";
	private String expect = "";
	private RunResult stepResult;
	private String failReason;
	private String picture;
	private String pictureName;
	private String pictureRelative;
	private String stepTime;
	private String url;
	private String failType;

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

	public String getStepDesc() {
		return stepDesc;
	}

	public void setStepDesc(String stepDesc) {
		this.stepDesc = stepDesc;
	}

	public RunResult getStepResult() {
		return stepResult;
	}

	public void setStepResult(RunResult stepResult) {
		this.stepResult = stepResult;
	}

	public String getFailReason() {
		return failReason;
	}

	public void setFailReason(String failReason) {
		this.failReason = failReason;
	}

	public String getPicture() {
		return picture;
	}

	public void setPicture(String picture) {
		this.picture = picture;
	}

	public String getStepTime() {
		return stepTime;
	}

	public void setStepTime(String stepTime) {
		this.stepTime = stepTime;
	}

	public StepType getStepType() {
		return stepType;
	}

	public void setStepType(StepType stepType) {
		this.stepType = stepType;
	}

	public String getUrl() {
		return url;
	}

	public void setUrl(String url) {
		this.url = url;
	}

	public void setFailType(String failType) {
		this.failType = failType;
	}

	public String getFailType() {
		return failType;
	}

	public void setStepId(int stepId) {
		this.stepId = stepId;
	}

	public int getStepId() {
		return stepId;
	}

	public void setPictureRelative(String pictureRelative) {
		this.pictureRelative = pictureRelative;
	}

	public String getPictureRelative() {
		return pictureRelative;
	}

	public void setPictureName(String pictureName) {
		this.pictureName = pictureName;
	}

	public String getPictureName() {
		return pictureName;
	}

	public String getActual() {
		return actual;
	}

	public void setActual(String actual) {
		this.actual = actual;
	}

	public String getExpect() {
		return expect;
	}

	public void setExpect(String expect) {
		this.expect = expect;
	}

}
