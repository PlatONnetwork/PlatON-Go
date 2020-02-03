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
public class IRegisterApplyManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_UPDATEAUTOAUDITSWITCH = "updateAutoAuditSwitch";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_FINDBYID = "findById";

    public static final String FUNC_AUDIT = "audit";

    public static final String FUNC_GETAUTOAUDITSWITCH = "getAutoAuditSwitch";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_FINDBYUUID = "findByUuid";

    public static final String FUNC_LISTBYCONDITION = "listByCondition";

    @Deprecated
    protected IRegisterApplyManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IRegisterApplyManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IRegisterApplyManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IRegisterApplyManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> updateAutoAuditSwitch(BigInteger code) {
        final Function function = new Function(
                FUNC_UPDATEAUTOAUDITSWITCH, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(code)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> update(String _json) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findById(String _applyId) {
        final Function function = new Function(FUNC_FINDBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_applyId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> audit(String _json) {
        final Function function = new Function(
                FUNC_AUDIT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getAutoAuditSwitch() {
        final Function function = new Function(FUNC_GETAUTOAUDITSWITCH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> insert(String _json) {
        final Function function = new Function(
                FUNC_INSERT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByUuid(String _uuid) {
        final Function function = new Function(FUNC_FINDBYUUID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_uuid)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> listByCondition(String _name, String _mobile, BigInteger _certType, BigInteger _pageSize, BigInteger _pageNo, String _auditStatus, BigInteger _accountStatus) {
        final Function function = new Function(FUNC_LISTBYCONDITION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.Utf8String(_mobile), 
                new org.web3j.abi.datatypes.generated.Uint256(_certType), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNo), 
                new org.web3j.abi.datatypes.Utf8String(_auditStatus), 
                new org.web3j.abi.datatypes.generated.Uint256(_accountStatus)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<IRegisterApplyManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IRegisterApplyManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRegisterApplyManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRegisterApplyManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IRegisterApplyManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IRegisterApplyManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRegisterApplyManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRegisterApplyManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IRegisterApplyManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRegisterApplyManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IRegisterApplyManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRegisterApplyManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IRegisterApplyManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IRegisterApplyManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IRegisterApplyManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IRegisterApplyManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
