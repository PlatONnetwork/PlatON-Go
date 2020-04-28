package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
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
public class IRoleManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_GETROLEMODULEVERSION = "getRoleModuleVersion";

    public static final String FUNC_DELETEBYID = "deleteById";

    public static final String FUNC_ACTIONUSED = "actionUsed";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_ADDACTIONTOROLE = "addActionToRole";

    public static final String FUNC_GETROLELISTBYMODULENAME = "getRoleListByModuleName";

    public static final String FUNC_PAGEBYNAME = "pageByName";

    public static final String FUNC_FINDBYID = "findById";

    public static final String FUNC_CHECKROLEACTION = "checkRoleAction";

    public static final String FUNC_FINDBYNAME = "findByName";

    public static final String FUNC_CHECKROLEACTIONWITHKEY = "checkRoleActionWithKey";

    public static final String FUNC_GETROLEMODULENAME = "getRoleModuleName";

    public static final String FUNC_GETROLEMODULEID = "getRoleModuleId";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_PAGEBYNAMEANDMODULENAME = "pageByNameAndModuleName";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_GETROLELISTBYMODULEID = "getRoleListByModuleId";

    public static final String FUNC_ROLEEXISTS = "roleExists";

    public static final String FUNC_GETROLELISTBYCONTRACTID = "getRoleListByContractId";

    public static final String FUNC_PAGEBYNAMEANDMODULEID = "pageByNameAndModuleId";

    public static final String FUNC_GETROLEIDBYACTIONIDANDINDEX = "getRoleIdByActionIdAndIndex";

    @Deprecated
    protected IRoleManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IRoleManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IRoleManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IRoleManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> getRoleModuleVersion(String _roleId) {
        final Function function = new Function(FUNC_GETROLEMODULEVERSION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> deleteById(String _roleId) {
        final Function function = new Function(
                FUNC_DELETEBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> actionUsed(String _actionId) {
        final Function function = new Function(FUNC_ACTIONUSED, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId)), 
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

    public RemoteCall<TransactionReceipt> addActionToRole(String _roleId, String _actionId) {
        final Function function = new Function(
                FUNC_ADDACTIONTOROLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId), 
                new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getRoleListByModuleName(String _moduleName, String _moduleVersion) {
        final Function function = new Function(FUNC_GETROLELISTBYMODULENAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
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

    public RemoteCall<BigInteger> checkRoleAction(String _roleId, String _actionId) {
        final Function function = new Function(FUNC_CHECKROLEACTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId), 
                new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> findByName(String _name) {
        final Function function = new Function(FUNC_FINDBYNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> checkRoleActionWithKey(String _roleId, String _resKey, String _opKey) {
        final Function function = new Function(FUNC_CHECKROLEACTIONWITHKEY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId), 
                new org.web3j.abi.datatypes.Address(_resKey), 
                new org.web3j.abi.datatypes.Utf8String(_opKey)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getRoleModuleName(String _roleId) {
        final Function function = new Function(FUNC_GETROLEMODULENAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getRoleModuleId(String _roleId) {
        final Function function = new Function(FUNC_GETROLEMODULEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> pageByNameAndModuleName(String _moduleName, String _name, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEBYNAMEANDMODULENAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
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

    public RemoteCall<String> getRoleListByModuleId(String _moduleId) {
        final Function function = new Function(FUNC_GETROLELISTBYMODULEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> roleExists(String _roleId) {
        final Function function = new Function(FUNC_ROLEEXISTS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_roleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getRoleListByContractId(String _contract) {
        final Function function = new Function(FUNC_GETROLELISTBYCONTRACTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_contract)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> pageByNameAndModuleId(String _moduleId, String _name, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEBYNAMEANDMODULEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId), 
                new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getRoleIdByActionIdAndIndex(String _actionId, BigInteger _index) {
        final Function function = new Function(FUNC_GETROLEIDBYACTIONIDANDINDEX, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId), 
                new org.web3j.abi.datatypes.generated.Uint256(_index)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<IRoleManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IRoleManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRoleManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRoleManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IRoleManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IRoleManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRoleManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRoleManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IRoleManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRoleManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IRoleManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRoleManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IRoleManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IRoleManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IRoleManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IRoleManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
