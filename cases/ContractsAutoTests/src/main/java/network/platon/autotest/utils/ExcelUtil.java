package network.platon.autotest.utils;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.Iterator;
import java.util.List;
import java.util.Map;
import org.apache.poi.hssf.usermodel.HSSFCell;
import org.apache.poi.hssf.usermodel.HSSFDateUtil;
import org.apache.poi.hssf.usermodel.HSSFWorkbook;
import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;
import org.junit.Test;

import com.csvreader.CsvReader;

/**
 * Excel的操作工具类
 * 包括对2003、2007的Excel表进行操作，还能够对CSV的表格进行操作
 * @author qcxiao
 *
 */
public class ExcelUtil {

	/** 总行数 */

	private int totalRows = 0;

	/** 总列数 */

	private int totalCells = 0;

	/** 错误信息 */

	private String errorInfo;
	/** 表单序号 */

	@SuppressWarnings("unused")
	private int sheetIndex = 0;

	/** 构造方法 */

	public ExcelUtil() {

	}

	/**
	 * 是否为excel2003格式
	 * 
	 * @param filePath
	 * @return
	 */
	private static boolean isExcel2003(String filePath) {

		return filePath.matches("^.+\\.(?i)(xls)$");

	}

	/**
	 * 是否为excel2007格式
	 * 
	 * @param filePath
	 * @return
	 */
	private static boolean isExcel2007(String filePath) {

		return filePath.matches("^.+\\.(?i)(xlsx)$");

	}

	/**
	 * 总行数
	 * 
	 * @return
	 */
	@SuppressWarnings("unused")
	private int getTotalRows() {

		return totalRows;

	}

	/**
	 * 总列数
	 * 
	 * @return
	 */
	private int getTotalCells() {

		return totalCells;

	}

	/**
	 * 错误信息
	 * 
	 * @return
	 */
	@SuppressWarnings("unused")
	private String getErrorInfo() {
		return errorInfo;

	}

	/**
	 * 
	 * @描述：验证excel文件
	 * 
	 * @参数：@param filePath　文件完整路径
	 * 
	 * @参数：@return
	 * 
	 * @返回值：boolean
	 */

	public boolean validateExcel(String filePath) {
		/** 检查文件名是否为空或者是否是Excel格式的文件 */

		if (filePath == null || !(isExcel2003(filePath) || isExcel2007(filePath))) {

			errorInfo = "文件名不是excel格式";

			return false;

		}

		/** 检查文件是否存在 */

		File file = new File(filePath);
		if (file == null || !file.exists()) {

			errorInfo = "文件不存在";
			return false;
		}
		return true;

	}

	@Test
	public void testss(){
		read("C:\\Users\\qcxiao\\Desktop\\applyeterm\\src\\test\\resources\\AutoeTermTest\\bookPnr.xls","自动预订PNR");
	}
	/**
	 * 
	 * @描述：根据文件名读取excel文件
	 * 
	 * @参数：@param filePath 文件完整路径
	 * 
	 * @参数：@return
	 * 
	 * @返回值：List
	 */

	public List<List<String>> read(String filePath, String sheetName) {

		List<List<String>> dataLst = new ArrayList<List<String>>();

		InputStream is = null;

		/** 验证文件是否合法 */

		if (!validateExcel(filePath)) {
			System.out.println(errorInfo);
			return null;

		}

		/** 判断文件的类型，是2003还是2007 */

		boolean isExcel2003 = true;

		if (isExcel2007(filePath)) {

			isExcel2003 = false;

		}

		/** 调用本类提供的根据流读取的方法 */
		try {
			File file = new File(filePath);

			is = new FileInputStream(file);

			/** 根据版本选择创建Workbook的方式 */

			Workbook wb = null;

			if (isExcel2003) {
				wb = new HSSFWorkbook(is);
			} else {
				wb = new XSSFWorkbook(is);
			}
			dataLst = read(wb, sheetName);

			is.close();

		} catch (IOException e) {

			e.printStackTrace();

		} finally {

			if (is != null) {

				try {

					is.close();

				} catch (IOException e) {

					is = null;

					e.printStackTrace();

				}

			}

		}

		/** 返回最后读取的结果 */

		return dataLst;

	}

	/**
	 * 
	 * @描述：根据流读取Excel文件
	 * 
	 * @参数：@param inputStream
	 * 
	 * @参数：@param isExcel2003
	 * 
	 * @参数：@return
	 * 
	 * @返回值：List
	 */

	public List<List<String>> read(InputStream inputStream, boolean isExcel2003) {

		List<List<String>> dataLst = null;

		try {

			/** 根据版本选择创建Workbook的方式 */

			Workbook wb = null;

			if (isExcel2003) {
				wb = new HSSFWorkbook(inputStream);
			} else {
				wb = new XSSFWorkbook(inputStream);
			}
			dataLst = read(wb);

		} catch (IOException e) {

			e.printStackTrace();

		}

		return dataLst;

	}

	/**
	 * 
	 * @描述：读取数据
	 * 
	 * @参数：@param Workbook
	 * 
	 * @参数：@return
	 * 
	 * @返回值：List<List<String>>
	 */

	private List<List<String>> read(Workbook wb, int sheetIndex) {
		List<List<String>> dataLst = new ArrayList<List<String>>();

		/** 得到指定的sheet */

		// int index =0;
		// wb.getSheetIndex(sheetName);

		Sheet sheet = wb.getSheetAt(sheetIndex);

		/** 得到Excel的行数 */

		this.totalRows = sheet.getPhysicalNumberOfRows();

		/** 得到Excel的列数 */

		if (this.totalRows >= 1 && sheet.getRow(0) != null) {

			this.totalCells = sheet.getRow(0).getPhysicalNumberOfCells();

		}

		/** 循环Excel的行 */

		for (int r = 0; r < this.totalRows; r++) {

			Row row = sheet.getRow(r);

			if (row == null) {

				continue;

			}

			List<String> rowLst = new ArrayList<String>();

			/** 循环Excel的列 */

			for (int c = 0; c < this.getTotalCells(); c++) {

				Cell cell = row.getCell(c);
				// String key=row.getCell(0).getStringCellValue();
				String cellValue = "";

				if (null != cell) {

					// cellValue=cell.getStringCellValue();
					// 以下是判断数据的类型
					switch (cell.getCellType()) {
					case HSSFCell.CELL_TYPE_NUMERIC: // 数字
						// cellValue = cell.getNumericCellValue() + "";
						if (HSSFDateUtil.isCellDateFormatted(cell)) {
							// 如果是Date类型则，取得该Cell的Date值
							Date date = cell.getDateCellValue();
							// 把Date转换成本地格式的字符串
							cellValue = DateUtil.dateToStr(date, "yyyy-MM-dd HH:mm:ss").toString();
							System.out.println(cellValue);
						}
						// 如果是纯数字
						else {
							// 取得当前Cell的数值
							Integer num = new Integer((int) cell.getNumericCellValue());
							cellValue = String.valueOf(num);
						}

						break;

					case HSSFCell.CELL_TYPE_STRING: // 字符串
						cellValue = cell.getStringCellValue().trim();
						break;

					case HSSFCell.CELL_TYPE_BOOLEAN: // Boolean
						cellValue = cell.getBooleanCellValue() + "";
						break;

					case HSSFCell.CELL_TYPE_FORMULA: // 公式
						cellValue = cell.getCellFormula() + "";
						break;

					case HSSFCell.CELL_TYPE_BLANK: // 空值
						cellValue = "";
						break;

					case HSSFCell.CELL_TYPE_ERROR: // 故障
						cellValue = "非法字符";
						break;

					default:
						cellValue = "未知类型";
						break;
					}
				}

				// System.out.print(cellValue+"  ");
				rowLst.add(cellValue);

			}
			// System.out.println();
			/** 保存第r行的第c列 */

			dataLst.add(rowLst);

		}

		return dataLst;
	}

	private List<List<String>> read(Workbook wb, String sheetName) {

		int sheetIndex = 0;
		try {
			sheetIndex = wb.getSheetIndex(sheetName);
		} catch (Exception e) {
			// 抛异常取第一个sheet
		}
		// 默认取第一个
		if (sheetIndex < 0) {
			sheetIndex = 0;
		}
		return read(wb, sheetIndex);
	}

	private List<List<String>> read(Workbook wb) {
		return read(wb, 0);
	}

	/**
	 * 将excel解析后的集合封装成Map形式
	 * 
	 * @param list
	 * @return
	 */
	@SuppressWarnings("unused")
	public static List<Map<String, String>> reflectMapList(List<List<String>> list) {
		ExcelUtil poi = new ExcelUtil();
		List<Map<String, String>> mlist = new ArrayList<Map<String, String>>();

		Map<String, String> map = new HashMap<String, String>();
		if (list != null) {

			for (int i = 1; i < list.size(); i++) {
				map = new HashMap<String, String>();
				List<String> cellList = list.get(i);

				for (int j = 0; j < cellList.size(); j++) {
					map.put(list.get(0).get(j), cellList.get(j));
				}
				mlist.add(map);
			}

		}

		return mlist;
	}

	public List<Map<String, String>> excelDatas(String filePath, String sheetName) {
		List<List<String>> lists = read(filePath, sheetName);
		// 对集合进行重新组装 Map<字段,值>
		List<Map<String, String>> datas = ExcelUtil.reflectMapList(lists);
		return datas;
	}

	/**
	 * 读取CSV的方法
	 * @param file
	 * @return
	 */
	public static List<String[]> importCsv(String file) {
		List<String[]> list = new ArrayList<String[]>();
		CsvReader reader = null;
		try {
			// 初始化CsvReader并指定列分隔符和字符编码
			reader = new CsvReader(file, ',', Charset.forName("GBK"));
			while (reader.readRecord()) {
				// 读取每行数据以数组形式返回
				String[] str = reader.getValues();
				if (str != null && str.length > 0) {
					// if (str[0] != null && !"".equals(str[0].trim())) {
					// list.add(str);
					// }
					if (str[0] != null) {
						list.add(str);
					}
				}
			}
		} catch (FileNotFoundException e) {
			// log.error("Error reading csv file.", e);
		} catch (IOException e) {
			// log.error("", e);
		} finally {
			if (reader != null)
				// 关闭CsvReader
				reader.close();
		}
		return list;
	}
	@SuppressWarnings({ "unchecked", "rawtypes" })
	public static void main(String[] args) throws Exception {
		List<String[]> listt = importCsv("C:\\Users\\qcxiao\\Documents\\Tencent Files\\280887262\\FileRecv\\GP_20141210_CSI_CAACSC.CSV");
		for(int i =0; i < listt.size(); i++){
			String [] strRows = listt.get(i);
			for(int j = 0; j < strRows.length;j++){
				System.out.println(strRows[j]);
			}
		}
		ExcelUtil poi = new ExcelUtil();
		// 获取解析后的集合
		List<List<String>> lists = poi.read("C:\\Users\\qcxiao\\Desktop\\applyeterm\\src\\test\\resources\\AutoeTermTest\\bookPnr.xls", "自动预订PNR");
		System.out.println(lists.size());
		// 对集合进行重新组装 Map<字段,值>
		List<Map<String, String>> list = ExcelUtil.reflectMapList(lists);
		// 调用工具类，组成对象集合
		// List<TravelerInfo> ts=Tool.reflectObj(TravelerInfo.class, list);
		// //遍历
		// for (TravelerInfo t : ts) {
		// System.out.println(t.getAirline_code()+" | "+
		// t.getFlight_num()+" | "+t.getSto()+" | "+t.getNationality()+"| ………………");
		// }
		int i = 1;
		for (Map<String, String> map : list) {
			System.out.println("行数" + i);
			Iterator iter = map.entrySet().iterator();
			while (iter.hasNext()) {
				Map.Entry<String, String> entry = (Map.Entry<String, String>) iter.next();
				System.out.println(entry.getValue());
			}
			i++;
		}
	}
}
