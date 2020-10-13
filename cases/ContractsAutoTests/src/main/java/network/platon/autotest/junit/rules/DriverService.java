package network.platon.autotest.junit.rules;

import java.util.*;

import lombok.extern.slf4j.Slf4j;
import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import org.junit.rules.TestRule;
import org.junit.runner.Description;
import org.junit.runners.model.MultipleFailureException;
import org.junit.runners.model.Statement;
import network.platon.autotest.exception.StepException;
import network.platon.autotest.junit.core.DriverModule;
import network.platon.autotest.junit.core.LogModule;
import network.platon.autotest.junit.core.PlanObserver;
import network.platon.autotest.junit.core.SuiteObserver;
import network.platon.autotest.junit.enums.RunResult;
import network.platon.autotest.junit.enums.RunStatus;
import network.platon.autotest.junit.enums.StepType;
import network.platon.autotest.junit.log.DatabaseLog;
import network.platon.autotest.junit.log.LocalFileLog;
import network.platon.autotest.utils.FileUtil;

/**
 * 框架测试驱动
 *
 * @author qcxiao
 */
@Slf4j
public class DriverService implements TestRule {
    public static String PLAN_VM_CONTENT = "<?xml version=\"1.0\" encoding=\"$!encode\"?>\n<suite name=\"$!suiteInfo.suiteName\">\n#foreach ($!moduleInfo in $!suiteInfo.moduleInfoList)\n<module  run=\"$!moduleInfo.moduleRun\" name=\"$!moduleInfo.moduleName\">\n#foreach ($caseInfo in $moduleInfo.caseInfoList)\n<case  run=\"$!caseInfo.caseRun\" name=\"$!caseInfo.caseName\"> \n</case>\n#end\n </module>\n#end\n</suite>";
    /**
     * 系统参数和test.properties文件中的值
     */
    public static Properties PROPERTIES = FileUtil.getProperties();
    public static String ENCODE = PROPERTIES.getProperty("encode") == null ? "utf-8" : PROPERTIES.getProperty("encode");
    /**
     * 资源文件中的键值匹配对信息,即对应properties文件中的值
     */
    public static Map<String, Map<String, String>> PROPERTIES_MAP = new HashMap<>();
    public static Description DESCRIPTION = null;
    public static Statement STATEMENT = null;
    /**
     * 是否合并所有类的测试报告
     * System.getProperty获取的是JVM运行时的参数，因此需要在maven执行时对其进行赋值
     */
    public static Boolean SUITE_MERGED = Boolean.valueOf(System.getProperty("suiteMerged") == null ? "false" : System.getProperty("suiteMerged").toLowerCase());

    /**
     * 获取用例在数据池中对应的参数信息
     *
     * @Example driverService.param;
     */
    public Map<String, String> param = new HashMap<>();
    /**
     * 配置文件参数
     *
     * @Example <div>driverService.propertiesMap;</div>
     * <div>driverService.propertiesMap;</div>
     */
    public Map<String, Map<String, String>> propertiesMap = new HashMap<>();
    private List<Throwable> errors = new ArrayList<>();
    // 观察者模式
    private DriverModule driverModule = DriverModule.getInstance();

    @Override
    public Statement apply(final Statement base, final Description des) {
        return new Statement() {
            @Override
            public void evaluate() throws Throwable {
                DESCRIPTION = des;
                STATEMENT = base;
                /**
                 * 两种情况： 1、合并报告：suiteName就为Project属性值；判断是否初始进入本项目，如果是则执行，
                 * 执行过程中会去收集target\test-classes目录下的所有符合条件类与符合条件方法
                 * 2、不合并报告：suiteName就为上次运行的类名
                 * ；判断是否初始进入本项目，如果是则执行，如果当前执行类与上次执行类不一致时也执行；
                 */
                if (LogModule.SUITE_INFO.getSuiteName() == null || !(SUITE_MERGED || des.getTestClass().getSimpleName().equals(LogModule.SUITE_INFO.getSuiteName()))) {
                    // 后续可以传入suite集合的上一层
                    Properties props = System.getProperties();
                    if (System.getProperty("website") == null || System.getProperty("website").trim().equals("")) {
                        props.remove("website");
                    }
                    if (System.getProperty("casePriority") == null || System.getProperty("casePriority").trim().equals("")) {
                        props.remove("casePriority");
                    }
                    System.setProperties(props);
                    PROPERTIES.putAll(System.getProperties());
                    driverModule.detachAll();
                    driverModule.attach(new SuiteObserver());
                    driverModule.attach(new PlanObserver());
                    if (PROPERTIES.getProperty("logType") != null) {
                        if (PROPERTIES.getProperty("logType").toUpperCase().equals("ALL")) {
                            driverModule.attach(new DatabaseLog());
                            driverModule.attach(new LocalFileLog());
                        } else if (PROPERTIES.getProperty("logType").toUpperCase().equals("DATABASE")) {
                            driverModule.attach(new DatabaseLog());
                        } else {
                            driverModule.attach(new LocalFileLog());
                        }
                    } else {
                        driverModule.attach(new LocalFileLog());
                    }
                    LogModule.SUITE_INFO = new SuiteInfo();
                    driverModule.suiteRunStart(LogModule.SUITE_INFO);
                }
                for (ModuleInfo planedModuleInfo : LogModule.SUITE_INFO.getModuleInfoList()) {
                    if (planedModuleInfo.getModuleName().equals(des.getTestClass().getSimpleName() + "." + des.getMethodName())) {
                        LogModule.MODULE_INFO = new ModuleInfo();
                        driverModule.moduleRunStart(LogModule.MODULE_INFO);
                        List<CaseInfo> caseInfoList = LogModule.MODULE_INFO.getCaseInfoList();
                        List<CaseInfo> caseInfoListUnrun = new ArrayList<CaseInfo>();
                        // 清除module信息，以及新增新module信息
                        for (int i = 0; i < caseInfoList.size(); i++) {
                            CaseInfo caseInfo = caseInfoList.get(i);
                            // 判断用例是否设置了要执行且用例中的用例级别是否为空或配置文件即PROPERTIES的用例级别是否或配置文件中的用例级别是否包含用例里的设置的用例级别
                            //if (caseInfo.getCaseRun() && (caseInfo.getCasePriority() == null || PROPERTIES.getProperty("casePriority") == null || PROPERTIES.getProperty("casePriority").trim().equals("") || PROPERTIES.getProperty("casePriority").toLowerCase().contains(caseInfo.getCasePriority().toLowerCase()))) {
                            if (!caseInfo.getCaseRun()) {
                                System.out.println("测试用例（" + caseInfo.getCaseName() + "）被设置成不执行。");
                                caseInfoListUnrun.add(caseInfo);
                                continue;
                            }
                            if ((null == PROPERTIES.getProperty("casePriority")
                                    || "".equals(PROPERTIES.getProperty("casePriority").trim()))
                                    || ((null != caseInfo.getCasePriority() && (null != PROPERTIES.getProperty("casePriority")
                                        || "".equals(PROPERTIES.getProperty("casePriority").trim())))
                                        && Arrays.asList(PROPERTIES.getProperty("casePriority").toLowerCase().split(",")).contains(caseInfo.getCasePriority().toLowerCase()))) {

                                LogModule.CASE_INFO = caseInfo;
                                driverModule.caseRunStart(LogModule.CASE_INFO);
                                param = caseInfo.getCaseParams();
                                propertiesMap = DriverService.PROPERTIES_MAP;
                                assertConfig(param);
                                try {
                                    driverModule.initialData(LogModule.CASE_INFO);
                                    base.evaluate();
                                } catch (StepException e) {
                                    String error = "测试用例（ " + param.get("caseName") + " ）执行失败! \n";
                                    errors.add(new Throwable(error, e));
                                } catch (RuntimeException e) {
                                    String error = "测试用例（ " + param.get("caseName") + " ）执行失败! \n";
                                    LogModule.logStepFail(StepType.CUSTOM, errorMessage(e) + "语句执行错误", RunResult.FAIL, "出错原因为" + e.getClass().getSimpleName() + ":" + e.getMessage() + "!");
                                    errors.add(new Throwable(error, e));
                                } finally {
                                    // 后期如有需求，才对数据销毁做处理
                                    driverModule.caseRunStop(LogModule.CASE_INFO);
                                    // 重跑失败
                                    if (caseInfo.getCaseResult().equals(RunResult.FAIL) && caseInfo.getCaseRerunNum() >= 0) {
                                        i--;
                                    }
                                    if (caseInfo.getCaseRunNum() > 0) {
                                        i--;
                                    }
                                }
                            } else {
                                System.out.println("测试用例（" + caseInfo.getCaseName() + "）被设置成不执行。");
                                caseInfoListUnrun.add(caseInfo);
                            }
                        }
                        List<CaseInfo> caseInfoListRuned = LogModule.MODULE_INFO.getCaseInfoList();
                        for (CaseInfo caseInfo : caseInfoListUnrun) {
                            caseInfoListRuned.remove(caseInfo);
                        }
                        LogModule.MODULE_INFO.setCaseInfoList(caseInfoListRuned);
                        driverModule.moduleRunStop(LogModule.MODULE_INFO);
                        /*
                         * 判断套件中的模块是否全部执行完成
                         */
                        Boolean suiteRunCompleted = true;
                        for (ModuleInfo moduleInfo : LogModule.SUITE_INFO.getModuleInfoList()) {
                            if (moduleInfo.getModuleStatus() != RunStatus.COMPLETED) {
                                suiteRunCompleted = false;
                                log.info("-----------" + moduleInfo.getModuleName());
                                break;
                            }
                        }
                        if (suiteRunCompleted) {
                            driverModule.suiteRunStop(LogModule.SUITE_INFO);
                            MultipleFailureException.assertEmpty(errors);
                            break;
                        }
                    }
                }
            }

        };
    }

    private String errorMessage(RuntimeException e) {
        String message = "";
        for (StackTraceElement ee : e.getStackTrace()) {
            message = message + ee.toString() + "\n";
            if (ee.toString().contains("network.platon.")) {
                break;
            }
        }
        return message;
    }

    /**
     * 配置是否执行断言
     *
     * @param param
     */
    private void assertConfig(Map<String, String> param) {
        // test.properties与website.properties都没有配置caseAssert默认为true
        //PROPERTIES表示test.properties中的内容，param表示website.properties中的内容
        //只要param中的caseAssert为N，此用例一定不执行
        if (param.get("caseAssert") != null && param.get("caseAssert").trim().equals("N")) {
            PROPERTIES.put("caseAssert", "false");
            return;
        }
        //param中的caseAssert不为Y，PROPERTIES中的caseAssert为N或false，此用例也不执行
        else if (param.get("caseAssert") != null
                && !param.get("caseAssert").trim().equals("Y")
                && PROPERTIES.get("caseAssert") != null
                && (PROPERTIES.get("caseAssert").toString().trim().equals("N")
                || PROPERTIES.get("caseAssert").toString().trim().equals("false"))) {
            PROPERTIES.put("caseAssert", "false");
            return;
        } else {
            PROPERTIES.put("caseAssert", "true");
        }
    }
}