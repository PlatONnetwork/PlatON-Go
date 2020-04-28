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
public class INodeInfoManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_ACTIVATEENODE = "ActivateEnode";

    public static final String FUNC_SETADMIN = "setAdmin";

    public static final String FUNC_FINDBYNODEADMIN = "findByNodeAdmin";

    public static final String FUNC_NODEINFOEXISTS = "nodeInfoExists";

    public static final String FUNC_GETREVISION = "getRevision";

    public static final String FUNC_DELETEBYID = "deleteById";

    public static final String FUNC_UPDATE = "update";

    public static final String FUNC_CHECKWRITEPERMISSION = "checkWritePermission";

    public static final String FUNC_FINDBYID = "findById";

    public static final String FUNC_FINDBYNAME = "findByName";

    public static final String FUNC_IPUSED = "IPUsed";

    public static final String FUNC_ISINWHITELIST = "isInWhiteList";

    public static final String FUNC_LISTALL = "listAll";

    public static final String FUNC_GETNODEADMIN = "getNodeAdmin";

    public static final String FUNC_INSERT = "insert";

    public static final String FUNC_ERASEADMINBYADD = "eraseAdminByAdd";

    public static final String FUNC_FINDBYPUBKEY = "findByPubkey";

    public static final String FUNC_FINDBYDEPARTMENTID = "findByDepartmentId";

    public static final String FUNC_GETENODELIST = "getEnodeList";

    @Deprecated
    protected INodeInfoManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected INodeInfoManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected INodeInfoManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected INodeInfoManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> ActivateEnode(String _pubkey) {
        final Function function = new Function(
                FUNC_ACTIVATEENODE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_pubkey)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setAdmin(String _nodeId, String _adminAddr) {
        final Function function = new Function(
                FUNC_SETADMIN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeId), 
                new org.web3j.abi.datatypes.Address(_adminAddr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByNodeAdmin(String _nodeAdmin) {
        final Function function = new Function(FUNC_FINDBYNODEADMIN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_nodeAdmin)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> nodeInfoExists(String _nodeId) {
        final Function function = new Function(FUNC_NODEINFOEXISTS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getRevision() {
        final Function function = new Function(FUNC_GETREVISION, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> deleteById(String _nodeId) {
        final Function function = new Function(
                FUNC_DELETEBYID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeId)), 
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

    public RemoteCall<BigInteger> checkWritePermission(String _addr, String _nodeInfoId) {
        final Function function = new Function(FUNC_CHECKWRITEPERMISSION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_addr), 
                new org.web3j.abi.datatypes.Utf8String(_nodeInfoId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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

    public RemoteCall<BigInteger> IPUsed(String _ip) {
        final Function function = new Function(FUNC_IPUSED, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_ip)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> isInWhiteList(String _commonName, String _ip) {
        final Function function = new Function(FUNC_ISINWHITELIST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_commonName), 
                new org.web3j.abi.datatypes.Utf8String(_ip)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> listAll() {
        final Function function = new Function(FUNC_LISTALL, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getNodeAdmin(String _nodeId) {
        final Function function = new Function(FUNC_GETNODEADMIN, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_nodeId)), 
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

    public RemoteCall<TransactionReceipt> eraseAdminByAdd(String _userAddr) {
        final Function function = new Function(
                FUNC_ERASEADMINBYADD, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_userAddr)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> findByPubkey(String _pubkey) {
        final Function function = new Function(FUNC_FINDBYPUBKEY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_pubkey)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> findByDepartmentId(String _departmentId) {
        final Function function = new Function(FUNC_FINDBYDEPARTMENTID, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_departmentId)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<String> getEnodeList() {
        final Function function = new Function(FUNC_GETENODELIST, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<INodeInfoManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(INodeInfoManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<INodeInfoManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(INodeInfoManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<INodeInfoManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(INodeInfoManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<INodeInfoManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(INodeInfoManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static INodeInfoManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new INodeInfoManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static INodeInfoManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new INodeInfoManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static INodeInfoManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new INodeInfoManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static INodeInfoManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new INodeInfoManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
