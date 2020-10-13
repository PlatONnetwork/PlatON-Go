package network.platon.autotest.utils;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;

import network.platon.autotest.junit.modules.CaseInfo;
import network.platon.autotest.junit.modules.ModuleInfo;
import network.platon.autotest.junit.modules.SuiteInfo;
import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.dom4j.io.SAXReader;

public class XmlUtil {
	
	public static String GetResAttributeXML(String Attribute, String resultList)
			throws UnsupportedEncodingException, DocumentException {
		System.out.println(resultList);
		String res = "";
		try {
			Document doc = DocumentHelper.parseText(resultList);
			Element root = doc.getRootElement();
			res = getElementAttribute(root, Attribute);
		} catch (DocumentException e) {
			e.printStackTrace();
		}
		return res;
	}

	// 递归迭代获取xml的节点；
	@SuppressWarnings("rawtypes")
	public static String getElementAttribute(Element elem, String Attribute) {
		String res = "";
		res = elem.attributeValue(Attribute);
		System.out.println(elem.getName());
		if (res == null || res.equals("")) {
			Iterator it = elem.elementIterator();
			while (it.hasNext()) {
				Element e = (Element) it.next();
				System.out.println(e.getName());
				res = getElementAttribute(e, Attribute);
				if (res != null && !res.equals("")) {
					break;
				}
			}
		}
		return res;
	}

	/**
	 * @param infoXML XML字符串信息
	 * @param attribute 属性
	 * @return 属性value值
	 */
	public static List<String> getAttributeValue(String infoXML,String attribute) {
		Document document;
		List<String> listAttributeValue = new ArrayList<String>();
		try {
			document = DocumentHelper.parseText(infoXML);
			Element root = document.getRootElement();
			listAttributeValue = getAttributeValue2(root,attribute);
			
		} catch (DocumentException e1) {
			e1.printStackTrace();
		}
		return listAttributeValue;
	}

	@SuppressWarnings("rawtypes")
	private static List<String> getAttributeValue2(Element element, String attribute) {
		List<String> value = new ArrayList<String>();
		value.add(element.attributeValue(attribute));
		System.out.println(element.getName());
		if(value.get(value.size()-1)==null||value==null||"null".equals(value.get(value.size()-1))){
			 for(Iterator it = element.elementIterator();it.hasNext();){
				 Element ele = (Element)it.next();
				 System.out.println(ele.getName());
				 value = getAttributeValue2(ele,attribute);
				 if(value.get(value.size()-1)!=null&&value!=null&&!"null".equals(value.get(value.size()-1))){
					 continue;
				 }
			 }
		}
		
		return value;
	}

	@SuppressWarnings("rawtypes")
	public static List<Map<String, String>> xmlDatas(String filePath) {
		List<Map<String, String>> datas = new ArrayList<Map<String, String>>();
		// Map<String, String> data = new HashMap<String, String>();
		SAXReader reader = new SAXReader();
		Document doc;
		try {
			doc = reader.read(filePath);
			Element root = doc.getRootElement();
			// System.out.println(root.getName());
			for (Iterator it = root.elementIterator(); it.hasNext();) {
				Element element = (Element) it.next();
				// System.out.println(element.getName());
				Map<String, String> data = new HashMap<String, String>();
				for (Iterator itt = element.elementIterator(); itt.hasNext();) {
					// System.out.println(element.attributeValue("age"));
					Element el = (Element) itt.next();
					data.put(el.getName(), el.getText());
				}
				datas.add(data);
			}
		} catch (DocumentException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		return datas;
	}

	/**
	 * 获取src/test/resources下的plan.xml或者plan.error.xml中的数据信息并转换为测试套件SuiteInfo对象
	 * @param filePath
	 * @return
	 * @throws DocumentException 
	 */
	@SuppressWarnings("rawtypes")
	public static SuiteInfo parsePlanXml(String filePath) throws DocumentException {
		SuiteInfo suiteInfo = new SuiteInfo();
		//List<Map<String, String>> datas = new ArrayList<Map<String, String>>();
		SAXReader reader = new SAXReader();
		Document doc;
		//File file = new File(filePath);
		try {
			doc = reader.read(filePath);
			Element root = doc.getRootElement();
			suiteInfo.setSuiteName(root.attributeValue("name"));
			List<ModuleInfo> moduleInfoList = new ArrayList<ModuleInfo>();
			for (Iterator it = root.elementIterator(); it.hasNext();) {
				Element element = (Element) it.next();
				ModuleInfo moduleInfo = new ModuleInfo();
				moduleInfo.setModuleName(element.attributeValue("name"));
				moduleInfo.setModuleRun(Boolean.valueOf(element.attributeValue("run")));

				List<CaseInfo> caseInfoList = new ArrayList<CaseInfo>();
				for (Iterator itt = element.elementIterator(); itt.hasNext();) {
					// System.out.println(element.attributeValue("age"));
					Element el = (Element) itt.next();
					CaseInfo caseInfo = new CaseInfo();
					caseInfo.setCaseName(el.attributeValue("name"));
					caseInfo.setCaseRun(Boolean.valueOf(el.attributeValue("run")));
					caseInfoList.add(caseInfo);
				}
				moduleInfo.setCaseInfoList(caseInfoList);
				moduleInfoList.add(moduleInfo);
			}
			suiteInfo.setModuleInfoList(moduleInfoList);
		} catch (DocumentException e) {
			throw new DocumentException(e);
		}
		return suiteInfo;
	}

	public static void main(String[] args) {
		List<Map<String, String>> a = xmlDatas("D:/workspace/autosky/src/test/resources/plan.xml");
		System.out.println(a);
	}
}
