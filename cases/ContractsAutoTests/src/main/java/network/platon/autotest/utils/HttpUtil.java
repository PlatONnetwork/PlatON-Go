package network.platon.autotest.utils;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.InputStreamReader;
import java.net.URL;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;

import org.apache.commons.httpclient.Credentials;
import org.apache.commons.httpclient.HttpClient;
import org.apache.commons.httpclient.HttpMethod;
import org.apache.commons.httpclient.MultiThreadedHttpConnectionManager;
import org.apache.commons.httpclient.NameValuePair;
import org.apache.commons.httpclient.UsernamePasswordCredentials;
import org.apache.commons.httpclient.auth.AuthScope;
import org.apache.commons.httpclient.methods.GetMethod;
import org.apache.commons.httpclient.methods.PostMethod;
import org.apache.commons.httpclient.methods.StringRequestEntity;
import org.apache.commons.httpclient.methods.multipart.FilePart;
import org.apache.commons.httpclient.methods.multipart.MultipartRequestEntity;
import org.apache.commons.httpclient.methods.multipart.Part;
import org.apache.commons.httpclient.params.HttpConnectionManagerParams;
import org.apache.commons.httpclient.params.HttpMethodParams;
import org.apache.http.HttpEntity;
import org.apache.http.HttpResponse;
import org.apache.http.HttpStatus;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.entity.StringEntity;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.impl.conn.tsccm.ThreadSafeClientConnManager;
import org.apache.http.util.EntityUtils;
import org.json.JSONObject;

@SuppressWarnings("deprecation")
public class HttpUtil {
	private static int HTTPCLIENT_CONNECT_TIMEOUT = 10000;
	private static int HTTPCLIENT_RESPONSE_TIMEOUT = 30000;
	static DefaultHttpClient httpClient = new DefaultHttpClient(new ThreadSafeClientConnManager());
	
	public static JSONObject httpPost(String url, String path) throws Exception, FileNotFoundException {
		HttpPost method = new HttpPost(url);
		// String jsonParam =
		// "{\"requestObject\": {\"departureCode\": \"PEK\", \"arrivalCode\": \"SHA\", \"departureDate\": \"2015-10-10\", \"timeRange\": 5, \"returnDate\": \"2015-10-20\",\"numOfPerson\": 1, \"travelPolicyVO\": {\"corpCode\": \"String类型\", \"agentId\": 0, \"corpId\": 0,\"updateTime\": \"\", \"discount\": 0, \"flightTime\": 0,\"serviceCodeId\": 0,\"serviceCodeName\": \"String类型\",\"lowestPrice\": 0,\"travelPolicyId\": 0,\"timePartHour\": 0,\"forwardOrderTime\": 0,\"exceAirport\": \"String类型\",\"hotelBigCity\": 0,\"hotelSmallCity\": 0,\"travelPolicyDesc\": \"String类型\",\"highestCabin\": \"String类型\",\"international\": \"String类型\",\"airBookRemarkEn\": \"String类型\",\"hotelBookRemarkCn\": \"String类型\",\"hotelBookRemarkEn\": \"String类型\",\"secondTierCity\": 0,\"thirdTierCity\": 0,\"corpfareHotelFirst\": \"String类型\",\"lowestPriceType\": \"0\",\"areaLowestType\": \"String类型\",\"corpfareAirline\": \"String类型\",\"corpfareAirlineFirst\": \"String类型\",\"cabinChoose\": \"String类型\",\"floatAmount\": 0,\"floatRate\": 0,\"latestReservedTime\": 0,\"canUpgrade\": \"String类型\",\"upgradeHours\": \"String类型\",\"upgradeMiles\": \"String类型\",\"upgradeCabin\": \"String类型\",\"tripConstraint\": \"String类型\",\"firstSelectionAirline\": \"String类型\",\"subSelectionAirline\": \"String类型\", \"firstSelectionAirlinePriceDiff\": 0,\"subSelectionAirlinePriceDiff\": 0,\"starRate\": \"String类型\",\"carrentPriceMax\": \"String类型\",\"carrentPriceMin\": \"String类型\",\"carrentRemarkChn\": \"String类型\",\"carrentRemarkEng\": \"String类型\",\"pricePercentage\": \"String类型\",\"isDiscountPreferred\": \"String类型\",\"fourthTierCity\": 0,\"fifthTierCity\": 0,\"smallCityCodes\": \"String类型\", \"discountDetail\": 0, \"name\": \"String类型\"}, \"expectTime\": \"09:00\"}}";
		StringEntity entity = new StringEntity(path, "utf-8");
		entity.setContentEncoding("UTF-8");
		entity.setContentType("application/json");
		method.setEntity(entity);
		JSONObject response = null;
		try {
			HttpResponse res = httpClient.execute(method);
			if (res.getStatusLine().getStatusCode() == HttpStatus.SC_OK) {
				HttpEntity entitys = res.getEntity();
				response = new JSONObject(new org.json.JSONTokener(new InputStreamReader(entitys.getContent(),"UTF-8")));
				String result = EntityUtils.toString(res.getEntity());// 返回json格式：
				System.out.println(result);
				// response = JSONObject.fromObject(result);
			}
		} catch (Exception e) {
			throw new RuntimeException(e);
		}
		return response;
	}

	public static String getResult(String url, String userName, String password) throws Exception {
		HttpClient client = new HttpClient(new MultiThreadedHttpConnectionManager());
		HttpConnectionManagerParams managerParams = client.getHttpConnectionManager().getParams();

		// 设置连接超时时间(单位毫秒)
		managerParams.setConnectionTimeout(HTTPCLIENT_CONNECT_TIMEOUT);

		// 设置读数据超时时间(单位毫秒)
		managerParams.setSoTimeout(HTTPCLIENT_RESPONSE_TIMEOUT);

		client.getParams().setParameter(HttpMethodParams.HTTP_CONTENT_CHARSET, "utf-8");

		setAuth(userName, password, url, client);
		HttpMethod method = new GetMethod(url);
		byte[] resultBytes;

		try {
			int statusCode = client.executeMethod(method);
			if (!(statusCode == 200 || statusCode == 302))
				throw (new Exception(statusCode + " Error"));
			resultBytes = method.getResponseBody();
			return new String(resultBytes);
		} catch (Exception e) {
			// logger.error("404 Error", e);
			throw new Exception(e);
		} finally {
			method.releaseConnection();
		}
	}

	public static String postRequest(String url, String userName, String password) throws Exception {
		HttpClient client = new HttpClient(new MultiThreadedHttpConnectionManager());
		HttpConnectionManagerParams managerParams = client.getHttpConnectionManager().getParams();

		// 设置连接超时时间(单位毫秒)
		managerParams.setConnectionTimeout(HTTPCLIENT_CONNECT_TIMEOUT);

		// 设置读数据超时时间(单位毫秒)
		managerParams.setSoTimeout(HTTPCLIENT_RESPONSE_TIMEOUT);

		client.getParams().setParameter(HttpMethodParams.HTTP_CONTENT_CHARSET, "utf-8");
		setAuth(userName, password, url, client);

		HttpMethod method = new PostMethod(url);
		byte[] resultBytes;

		try {
			int statusCode = client.executeMethod(method);
			if (!(statusCode == 200 || statusCode == 302))
				throw (new Exception(statusCode + " Error"));
			resultBytes = method.getResponseBody();
			return new String(resultBytes);
		} catch (Exception e) {
			// logger.error("404 Error", e);
			throw new Exception(e);
		} finally {
			method.releaseConnection();
		}
	}

	public static String postRequest(String url, Map<String, String> params, String userName, String password) throws Exception {
		HttpClient client = new HttpClient(new MultiThreadedHttpConnectionManager());
		HttpConnectionManagerParams managerParams = client.getHttpConnectionManager().getParams();

		// 设置连接超时时间(单位毫秒)
		managerParams.setConnectionTimeout(HTTPCLIENT_CONNECT_TIMEOUT);

		// 设置读数据超时时间(单位毫秒)
		managerParams.setSoTimeout(HTTPCLIENT_RESPONSE_TIMEOUT);

		client.getParams().setParameter(HttpMethodParams.HTTP_CONTENT_CHARSET, "utf-8");

		setAuth(userName, password, url, client);

		PostMethod method = new PostMethod(url);

		Set<Entry<String, String>> esParams = params.entrySet();
		NameValuePair[] httpParams = new NameValuePair[esParams.size()];

		int i = 0;
		for (Entry<String, String> esParam : esParams) {
			httpParams[i] = new NameValuePair(esParam.getKey(), esParam.getValue());
			i++;
		}
		method.setRequestBody(httpParams);
		byte[] resultBytes;

		try {
			int statusCode = client.executeMethod(method);
			if (!(statusCode == 200 || statusCode == 302))
				throw (new Exception(statusCode + " Error"));
			resultBytes = method.getResponseBody();
			return new String(resultBytes);
		} catch (Exception e) {
			// logger.error("404 Error", e);
			throw new Exception(e);
		} finally {
			method.releaseConnection();
		}
	}

	public static String postXml(String url, String value, String userName, String password) throws Exception {
		HttpClient client = new HttpClient(new MultiThreadedHttpConnectionManager());
		HttpConnectionManagerParams managerParams = client.getHttpConnectionManager().getParams();

		// 设置连接超时时间(单位毫秒)
		managerParams.setConnectionTimeout(HTTPCLIENT_CONNECT_TIMEOUT);

		// 设置读数据超时时间(单位毫秒)
		managerParams.setSoTimeout(HTTPCLIENT_RESPONSE_TIMEOUT);

		client.getParams().setParameter(HttpMethodParams.HTTP_CONTENT_CHARSET, "utf-8");

		setAuth(userName, password, url, client);

		PostMethod method = new PostMethod(url);
		method.setRequestEntity(new StringRequestEntity(value, "text/xml", "UTF-8"));
		byte[] resultBytes;

		try {
			int statusCode = client.executeMethod(method);
			if (!(statusCode == 200 || statusCode == 302))
				throw (new Exception(statusCode + " Error"));
			resultBytes = method.getResponseBody();
			return new String(resultBytes);
		} catch (Exception e) {
			// logger.error("404 Error", e);
			throw new Exception(e);
		} finally {
			method.releaseConnection();
		}
	}

	private static void setAuth(String username, String password, String url, HttpClient httpClient) throws Exception {
		if (username != null) {
			httpClient.getParams().setAuthenticationPreemptive(true);
			Credentials defaultcreds = new UsernamePasswordCredentials(username, password);
			//Credentials defaultcreds = new UsernamePasswordCredentials("demo/userAction!doNotNeedSession_login.action?cname=admin&", "cpwd=admin");
			URL tUrl = new URL(url);
			httpClient.getState().setCredentials(new AuthScope(tUrl.getHost(), tUrl.getPort()), defaultcreds);
		}
	}

	public static String uploadFile(File file, String url) throws Exception {
		if (!file.exists()) {
			return null;
		}
		PostMethod method = new PostMethod(url);
		// try {
		// FilePart：用来上传文件的类
		FilePart fp = new FilePart("fileData", file);
		Part[] parts = { fp };

		// 对于MIME类型的请求，httpclient建议全用MulitPartRequestEntity进行包装
		MultipartRequestEntity mre = new MultipartRequestEntity(parts, method.getParams());
		method.setRequestEntity(mre);
		HttpClient client = new HttpClient();
		client.getHttpConnectionManager().getParams().setConnectionTimeout(50000);// 设置连接时间

		byte[] resultBytes;

		try {
			int statusCode = client.executeMethod(method);
			if (!(statusCode == 200 || statusCode == 302))
				throw (new Exception(statusCode + " Error"));
			resultBytes = method.getResponseBody();
			return new String(resultBytes);
		} catch (Exception e) {
			// logger.error("404 Error", e);
			throw new Exception(e);
		} finally {
			method.releaseConnection();
		}
		// int status = client.executeMethod(postMethod);
		// if (status == HttpStatus.SC_OK) {
		// System.out.println(postMethod.getResponseBodyAsString());
		// } else {
		// System.out.println("fail");
		// }
		// } catch (Exception e) {
		// e.printStackTrace();
		// } finally {
		// // 释放连接
		// postMethod.releaseConnection();
		// }
	}

	/**
	 * @param args
	 * @throws Exception
	 */
	public static void main(String[] args) throws Exception {
		//http://localhost:8080/AutoTestPlat/logtest/logTestPictureAction
		uploadFile(new File("C:\\autosky_log\\TestJunit\\2014-04-22_16-16-33\\screenshot\\20140422-161703.png"), "http://localhost:8080/AutoTestPlat/logtest/upload");
	}
}
