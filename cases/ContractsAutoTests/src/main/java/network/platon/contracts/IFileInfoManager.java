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
public class IFileInfoManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_LISTBYGROUP = "listByGroup";

    public static final String FUNC_DELETEBYID = "deleteById";

    public static final String FUNC_PAGEBYGROUP = "pageByGroup";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_GETCURRENTPAGECOUNT = "getCurrentPageCount";

    public static final String FUNC_FIND = "find";

    public static final String FUNC_GENERATEFILEID = "generateFileID";

    public static final String FUNC_GETGROUPPAGECOUNT = "getGroupPageCount";

    public static final String FUNC_GETCOUNT = "getCount";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_PAGEFILES = "pageFiles";

    public static final String FUNC_GETCURRENTPAGESIZE = "getCurrentPageSize";

    public static final String FUNC_GETGROUPFILECOUNT = "getGroupFileCount";

    @Deprecated
    protected IFileInfoManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IFileInfoManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IFileInfoManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IFileInfoManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> listByGroup(String _group) {
        final Function function = new Function(FUNC_LISTBYGROUP, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_group)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> deleteById(String _fileId) {
        final Function function = new Function(
                FUNC_DELETEBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_fileId)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> pageByGroup(String _group, BigInteger _pageNo) {
        final Function function = new Function(FUNC_PAGEBYGROUP, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_group), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageNo)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> update(String _fileJson) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_fileJson)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getCurrentPageCount() {
        final Function function = new Function(FUNC_GETCURRENTPAGECOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> find(String _fileId) {
        final Function function = new Function(FUNC_FIND, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_fileId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> generateFileID(String _salt, String _groupID, String _serverId, String _filename) {
        final Function function = new Function(FUNC_GENERATEFILEID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_salt), 
                new org.web3j.abi.datatypes.Utf8String(_groupID), 
                new org.web3j.abi.datatypes.Utf8String(_serverId), 
                new org.web3j.abi.datatypes.Utf8String(_filename)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getGroupPageCount(String _group) {
        final Function function = new Function(FUNC_GETGROUPPAGECOUNT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_group)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getCount() {
        final Function function = new Function(FUNC_GETCOUNT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> insert(String _fileJson) {
        final Function function = new Function(
                FUNC_INSERT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_fileJson)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> pageFiles(BigInteger _pageNo, BigInteger _pageSize) {
        final Function function = new Function(FUNC_PAGEFILES, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_pageNo), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getCurrentPageSize() {
        final Function function = new Function(FUNC_GETCURRENTPAGESIZE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getGroupFileCount(String _group) {
        final Function function = new Function(FUNC_GETGROUPFILECOUNT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_group)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<IFileInfoManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IFileInfoManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IFileInfoManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IFileInfoManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IFileInfoManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IFileInfoManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IFileInfoManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IFileInfoManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IFileInfoManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IFileInfoManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IFileInfoManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IFileInfoManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IFileInfoManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IFileInfoManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IFileInfoManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IFileInfoManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
