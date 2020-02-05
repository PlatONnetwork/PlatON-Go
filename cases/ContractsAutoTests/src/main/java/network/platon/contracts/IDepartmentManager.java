package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class IDepartmentManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_SETADMIN = "setAdmin";

    public static final String FUNC_GETREVISION = "getRevision";

    public static final String FUNC_DEPARTMENTEMPTY = "departmentEmpty";

    public static final String FUNC_DELETEBYID = "deleteById";

    public static final String FUNC_DEPARTMENTEXISTSBYCN = "departmentExistsByCN";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_FINDBYPARENTID = "findByParentId";

    public static final String FUNC_CHECKWRITEPERMISSION = "checkWritePermission";

    public static final String FUNC_PAGEBYNAME = "pageByName";

    public static final String FUNC_FINDBYID = "findById";

    public static final String FUNC_FINDBYNAME = "findByName";

    public static final String FUNC_ERASEADMINBYADDRESS = "eraseAdminByAddress";

    public static final String FUNC_GETCHILDIDBYINDEX = "getChildIdByIndex";

    public static final String FUNC_DEPARTMENTEXISTS = "departmentExists";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_PAGEBYNAMEANDSTATUS = "pageByNameAndStatus";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_SETDEPARTMENTSTATUS = "setDepartmentStatus";

    public static final String FUNC_GETADMIN = "getAdmin";

    @Deprecated
    protected IDepartmentManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IDepartmentManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IDepartmentManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IDepartmentManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> setAdmin(String _departmentId, String _adminAddr) {
        final Function function = new Function(
                FUNC_SETADMIN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId), 
                new org.web3j.abi.datatypes.Address(_adminAddr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getRevision() {
        final Function function = new Function(FUNC_GETREVISION, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> departmentEmpty(String _departmentId) {
        final Function function = new Function(FUNC_DEPARTMENTEMPTY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> deleteById(String _departmentId) {
        final Function function = new Function(
                FUNC_DELETEBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> departmentExistsByCN(String _commonName) {
        final Function function = new Function(FUNC_DEPARTMENTEXISTSBYCN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_commonName)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> update(String _json) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByParentId(String _parentId) {
        final Function function = new Function(FUNC_FINDBYPARENTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_parentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> checkWritePermission(String _addr, String _departmentId) {
        final Function function = new Function(FUNC_CHECKWRITEPERMISSION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_addr), 
                new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> pageByName(String _name, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEBYNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findById(String _id) {
        final Function function = new Function(FUNC_FINDBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_id)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByName(String _name) {
        final Function function = new Function(FUNC_FINDBYNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> eraseAdminByAddress(String _userAddr) {
        final Function function = new Function(
                FUNC_ERASEADMINBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getChildIdByIndex(String _departmentId, BigInteger _index) {
        final Function function = new Function(FUNC_GETCHILDIDBYINDEX, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId), 
                new org.web3j.abi.datatypes.generated.Uint256(_index)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> departmentExists(String _departmentId) {
        final Function function = new Function(FUNC_DEPARTMENTEXISTS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> pageByNameAndStatus(String _parentId, BigInteger _status, String _name, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEBYNAMEANDSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_parentId), 
                new org.web3j.abi.datatypes.generated.Uint256(_status), 
                new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> insert(String _json) {
        final Function function = new Function(
                FUNC_INSERT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setDepartmentStatus(String _departmentId, BigInteger _status) {
        final Function function = new Function(
                FUNC_SETDEPARTMENTSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId), 
                new org.web3j.abi.datatypes.generated.Uint256(_status)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getAdmin(String _departmentId) {
        final Function function = new Function(FUNC_GETADMIN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<IDepartmentManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IDepartmentManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IDepartmentManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IDepartmentManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IDepartmentManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IDepartmentManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IDepartmentManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IDepartmentManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IDepartmentManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IDepartmentManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IDepartmentManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IDepartmentManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IDepartmentManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IDepartmentManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IDepartmentManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IDepartmentManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
