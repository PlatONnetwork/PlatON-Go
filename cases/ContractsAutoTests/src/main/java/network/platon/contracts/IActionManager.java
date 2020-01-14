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
public class IActionManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_DELETEBYID = "deleteById";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_FINDBYID = "findById";

    public static final String FUNC_LISTCONTRACTACTIONS = "listContractActions";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_FINDACTIONBYTYPE = "findActionByType";

    public static final String FUNC_QUERYACTIONENABLE = "queryActionEnable";

    public static final String FUNC_ACTIONEXISTS = "actionExists";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_FINDBYKEY = "findByKey";

    public static final String FUNC_GETACTIONLISTBYCONTRACTID = "getActionListByContractId";

    public static final String FUNC_GETACTIONLISTBYCONTRACTNAME = "getActionListByContractName";

    public static final String FUNC_GETACTIONLISTBYMODULEID = "getActionListByModuleId";

    public static final String FUNC_GETACTIONLISTBYMODULENAME = "getActionListByModuleName";

    public static final String FUNC_CHECKACTIONWITHKEY = "checkActionWithKey";

    @Deprecated
    protected IActionManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IActionManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IActionManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IActionManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> deleteById(String _actionId) {
        final Function function = new Function(
                FUNC_DELETEBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> update(String _actionJson) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionJson)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findById(String _actionId) {
        final Function function = new Function(FUNC_FINDBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> listContractActions(String _contractName) {
        final Function function = new Function(FUNC_LISTCONTRACTACTIONS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_contractName)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findActionByType(BigInteger _type) {
        final Function function = new Function(FUNC_FINDACTIONBYTYPE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_type)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> queryActionEnable(String _actionId, Boolean _checkOwner) {
        final Function function = new Function(FUNC_QUERYACTIONENABLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId), 
                new org.web3j.abi.datatypes.Bool(_checkOwner)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> actionExists(String _actionId) {
        final Function function = new Function(FUNC_ACTIONEXISTS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> insert(String _actionJson) {
        final Function function = new Function(
                FUNC_INSERT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionJson)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByKey(String _resKey, String _opKey) {
        final Function function = new Function(FUNC_FINDBYKEY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_resKey), 
                new org.web3j.abi.datatypes.Utf8String(_opKey)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getActionListByContractId(String _contractId) {
        final Function function = new Function(FUNC_GETACTIONLISTBYCONTRACTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_contractId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getActionListByContractName(String _moduleName, String _moduleVersion, String _contractName, String _contractVersion) {
        final Function function = new Function(FUNC_GETACTIONLISTBYCONTRACTNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Utf8String(_contractName), 
                new org.web3j.abi.datatypes.Utf8String(_contractVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getActionListByModuleId(String _moduleId) {
        final Function function = new Function(FUNC_GETACTIONLISTBYMODULEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getActionListByModuleName(String _moduleName, String _moduleVersion) {
        final Function function = new Function(FUNC_GETACTIONLISTBYMODULENAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> checkActionWithKey(String _actionId, String _contractAddr, String _opSha3Key) {
        final Function function = new Function(FUNC_CHECKACTIONWITHKEY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_actionId), 
                new org.web3j.abi.datatypes.Address(_contractAddr), 
                new org.web3j.abi.datatypes.Utf8String(_opSha3Key)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<IActionManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IActionManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IActionManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IActionManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IActionManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IActionManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IActionManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IActionManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IActionManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IActionManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IActionManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IActionManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IActionManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IActionManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IActionManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IActionManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
