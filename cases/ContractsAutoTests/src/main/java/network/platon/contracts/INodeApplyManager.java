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
public class INodeApplyManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_DELETEBYID = "deleteById";

    public static final String FUNC_AUDITING = "auditing";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_PAGEBYNAMEANDSTATUS = "pageByNameAndStatus";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_FINDBYAPPLYID = "findByApplyId";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_NODEAPPLYEXISTS = "nodeApplyExists";

    public static final String FUNC_FINDBYSTATE = "findByState";

    @Deprecated
    protected INodeApplyManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected INodeApplyManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected INodeApplyManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected INodeApplyManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> deleteById(String _nodeApplyId) {
        final Function function = new Function(
                FUNC_DELETEBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeApplyId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> auditing(String _json) {
        final Function function = new Function(
                FUNC_AUDITING, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_json)), 
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

    public RemoteCall<String> pageByNameAndStatus(BigInteger _status, String _deptName, BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEBYNAMEANDSTATUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_status), 
                new org.web3j.abi.datatypes.Utf8String(_deptName), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByApplyId(String _nodeApplyId) {
        final Function function = new Function(FUNC_FINDBYAPPLYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeApplyId)), 
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

    public RemoteCall<BigInteger> nodeApplyExists(String _nodeApplyId) {
        final Function function = new Function(FUNC_NODEAPPLYEXISTS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeApplyId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> findByState(BigInteger _state) {
        final Function function = new Function(FUNC_FINDBYSTATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_state)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<INodeApplyManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(INodeApplyManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<INodeApplyManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(INodeApplyManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<INodeApplyManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(INodeApplyManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<INodeApplyManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(INodeApplyManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static INodeApplyManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new INodeApplyManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static INodeApplyManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new INodeApplyManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static INodeApplyManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new INodeApplyManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static INodeApplyManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new INodeApplyManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
