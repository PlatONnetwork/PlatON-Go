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
public class IRoleFilterManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_FINDCONTRACTBYMODTEXT = "findContractByModText";

    public static final String FUNC_CHANGEMODULEOWNER = "changeModuleOwner";

    public static final String FUNC_FINDBYMODULETEXT = "findByModuleText";

    public static final String FUNC_ADDACTION = "addAction";

    public static final String FUNC_QRYMODULEDETAIL = "qryModuleDetail";

    public static final String FUNC_DELMODULE = "delModule";

    public static final String FUNC_LISTCONTRACTBYMODULENAME = "listContractByModuleName";

    public static final String FUNC_ADDROLE = "addRole";

    public static final String FUNC_FINDBYNAME = "findByName";

    public static final String FUNC_AUTHORIZEPROCESSOR = "authorizeProcessor";

    public static final String FUNC_QRYMODULES = "qryModules";

    public static final String FUNC_GETMODULECOUNT = "getModuleCount";

    public static final String FUNC_LISTCONTRACTBYMODTEXTANDCTTNAME = "listContractByModTextAndCttName";

    public static final String FUNC_UPDMODULE = "updModule";

    public static final String FUNC_MODULEISEXIST = "moduleIsExist";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_ADDACTIONTOROLE = "addActionToRole";

    public static final String FUNC_ADDMODULE = "addModule";

    public static final String FUNC_LISTCONTRACTBYMODULEID = "listContractByModuleId";

    public static final String FUNC_ADDCONTRACT = "addContract";

    public static final String FUNC_ADDMENU = "addMenu";

    @Deprecated
    protected IRoleFilterManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IRoleFilterManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IRoleFilterManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IRoleFilterManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> findContractByModText(String _moduleText) {
        final Function function = new Function(FUNC_FINDCONTRACTBYMODTEXT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleText)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> changeModuleOwner(String _moduleName, String _moduleVersion, String _newOwner) {
        final Function function = new Function(
                FUNC_CHANGEMODULEOWNER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Address(_newOwner)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByModuleText(String _moduleText) {
        final Function function = new Function(FUNC_FINDBYMODULETEXT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleText)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> addAction(String _json) {
        final Function function = new Function(
                FUNC_ADDACTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> qryModuleDetail(String _moduleName, String _moduleVersion) {
        final Function function = new Function(FUNC_QRYMODULEDETAIL, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> delModule(String _moduleId) {
        final Function function = new Function(
                FUNC_DELMODULE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> listContractByModuleName(String _moduleName, String _moduleVersion) {
        final Function function = new Function(FUNC_LISTCONTRACTBYMODULENAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> addRole(String _json) {
        final Function function = new Function(
                FUNC_ADDROLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByName(String _name) {
        final Function function = new Function(FUNC_FINDBYNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> authorizeProcessor(String _from, String _to, String _funcHash, String _extraData) {
        final Function function = new Function(FUNC_AUTHORIZEPROCESSOR, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_from), 
                new org.web3j.abi.datatypes.Address(_to), 
                new org.web3j.abi.datatypes.Utf8String(_funcHash), 
                new org.web3j.abi.datatypes.Utf8String(_extraData)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> qryModules() {
        final Function function = new Function(FUNC_QRYMODULES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getModuleCount() {
        final Function function = new Function(FUNC_GETMODULECOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> listContractByModTextAndCttName(String _moduleText, String _cttName, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_LISTCONTRACTBYMODTEXTANDCTTNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleText), 
                new org.web3j.abi.datatypes.Utf8String(_cttName), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> qryModuleDetail(String _moduleId) {
        final Function function = new Function(FUNC_QRYMODULEDETAIL, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> updModule(String _json) {
        final Function function = new Function(
                FUNC_UPDMODULE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> moduleIsExist(String _moduleId) {
        final Function function = new Function(FUNC_MODULEISEXIST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> addActionToRole(String _moduleId, String _roleId, String _actionId) {
        final Function function = new Function(
                FUNC_ADDACTIONTOROLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId), 
                new org.web3j.abi.datatypes.Utf8String(_roleId), 
                new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> addModule(String _json) {
        final Function function = new Function(
                FUNC_ADDMODULE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> listContractByModuleId(String _moduleId) {
        final Function function = new Function(FUNC_LISTCONTRACTBYMODULEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> addContract(String _json) {
        final Function function = new Function(
                FUNC_ADDCONTRACT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> addMenu(String _json) {
        final Function function = new Function(
                FUNC_ADDMENU, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<IRoleFilterManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IRoleFilterManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRoleFilterManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRoleFilterManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IRoleFilterManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IRoleFilterManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRoleFilterManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRoleFilterManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IRoleFilterManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRoleFilterManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IRoleFilterManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRoleFilterManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IRoleFilterManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IRoleFilterManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IRoleFilterManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IRoleFilterManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
