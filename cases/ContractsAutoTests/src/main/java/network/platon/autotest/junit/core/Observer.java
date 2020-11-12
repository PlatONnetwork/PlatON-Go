package network.platon.autotest.junit.core;

import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;


/**
 * 套件、模块、用例开始与结束的接口
 * @author qcxiao
 *
 */
public interface Observer {
	public void suiteRunStart(SuiteInfo suiteInfo);
	public void moduleRunStart(ModuleInfo moduleInfo);
	public void caseRunStart(CaseInfo caseInfo);
	public void suiteRunStop(SuiteInfo suiteInfo);
	public void moduleRunStop(ModuleInfo moduleInfo);
	public void caseRunStop(CaseInfo caseInfo);
}
