package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
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
public class IRegisterManager extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_REGISTER = "register";

    public static final String FUNC_FINDMODULENAMEBYADDRESS = "findModuleNameByAddress";

    public static final String FUNC_IFCONTRACTREGIST = "IfContractRegist";

    public static final String FUNC_UNREGISTER = "unRegister";

    public static final String FUNC_IFMODULEREGIST = "IfModuleRegist";

    public static final String FUNC_GETMODULEADDRESS = "getModuleAddress";

    public static final String FUNC_TRANSFERCONTRACT = "transferContract";

    public static final String FUNC_CHANGEMODULEREGISTEROWNER = "changeModuleRegisterOwner";

    public static final String FUNC_FINDCONTRACTVERSIONBYADDRESS = "findContractVersionByAddress";

    public static final String FUNC_FINDRESNAMEBYADDRESS = "findResNameByAddress";

    public static final String FUNC_CHANGECONTRACTREGISTEROWNER = "changeContractRegisterOwner";

    public static final String FUNC_GETCONTRACTADDRESS = "getContractAddress";

    public static final String FUNC_FINDMODULEVERSIONBYADDRESS = "findModuleVersionByAddress";

    public static final String FUNC_GETREGISTEREDCONTRACT = "getRegisteredContract";

    @Deprecated
    protected IRegisterManager(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected IRegisterManager(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected IRegisterManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected IRegisterManager(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> register(String _moduleName, String _moduleVersion, String _contractName, String _contractVersion) {
        final Function function = new Function(
                FUNC_REGISTER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Utf8String(_contractName), 
                new org.web3j.abi.datatypes.Utf8String(_contractVersion)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> findModuleNameByAddress(String _addr) {
        final Function function = new Function(FUNC_FINDMODULENAMEBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Boolean> IfContractRegist(String _contractAddr) {
        final Function function = new Function(FUNC_IFCONTRACTREGIST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_contractAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> unRegister() {
        final Function function = new Function(
                FUNC_UNREGISTER, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Boolean> IfModuleRegist(String _moduleName, String _moduleVersion) {
        final Function function = new Function(FUNC_IFMODULEREGIST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<String> getModuleAddress(String _moduleName, String _moduleVersion) {
        final Function function = new Function(FUNC_GETMODULEADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> transferContract(String _fromModuleNameAndVersion, String _fromNameAndVersion, String _toModuleNameAndVersion, String _toNameAndVersion, String _signString) {
        final Function function = new Function(
                FUNC_TRANSFERCONTRACT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_fromModuleNameAndVersion), 
                new org.web3j.abi.datatypes.Utf8String(_fromNameAndVersion), 
                new org.web3j.abi.datatypes.Utf8String(_toModuleNameAndVersion), 
                new org.web3j.abi.datatypes.Utf8String(_toNameAndVersion), 
                new org.web3j.abi.datatypes.Utf8String(_signString)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<Boolean> IfContractRegist(String _moduleName, String _moduleVersion, String _contractName, String _contractVersion) {
        final Function function = new Function(FUNC_IFCONTRACTREGIST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Utf8String(_contractName), 
                new org.web3j.abi.datatypes.Utf8String(_contractVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public RemoteCall<TransactionReceipt> changeModuleRegisterOwner(String _moduleName, String _moduleVersion, String _newOwner) {
        final Function function = new Function(
                FUNC_CHANGEMODULEREGISTEROWNER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Address(_newOwner)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> register(String _moduleName, String _moduleVersion) {
        final Function function = new Function(
                FUNC_REGISTER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> findContractVersionByAddress(String _addr) {
        final Function function = new Function(FUNC_FINDCONTRACTVERSIONBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> findResNameByAddress(String _addr) {
        final Function function = new Function(FUNC_FINDRESNAMEBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> changeContractRegisterOwner(String _moduleName, String _moduleVersion, String _contractName, String _contractVersion, String _newOwner) {
        final Function function = new Function(
                FUNC_CHANGECONTRACTREGISTEROWNER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Utf8String(_contractName), 
                new org.web3j.abi.datatypes.Utf8String(_contractVersion), 
                new org.web3j.abi.datatypes.Address(_newOwner)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getContractAddress(String _moduleName, String _moduleVersion, String _contractName, String _contractVersion) {
        final Function function = new Function(FUNC_GETCONTRACTADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_moduleName), 
                new org.web3j.abi.datatypes.Utf8String(_moduleVersion), 
                new org.web3j.abi.datatypes.Utf8String(_contractName), 
                new org.web3j.abi.datatypes.Utf8String(_contractVersion)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> findModuleVersionByAddress(String _addr) {
        final Function function = new Function(FUNC_FINDMODULEVERSIONBYADDRESS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_addr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<String> getRegisteredContract(BigInteger _pageNum, BigInteger _pageSize) {
        final Function function = new Function(FUNC_GETREGISTEREDCONTRACT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(_pageNum), 
                new org.web3j.abi.datatypes.generated.Uint256(_pageSize)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<Boolean> IfModuleRegist(String _moduleAddr) {
        final Function function = new Function(FUNC_IFMODULEREGIST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(_moduleAddr)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }

    public static RemoteCall<IRegisterManager> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(IRegisterManager.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRegisterManager> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRegisterManager.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<IRegisterManager> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(IRegisterManager.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<IRegisterManager> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(IRegisterManager.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static IRegisterManager load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRegisterManager(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static IRegisterManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new IRegisterManager(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static IRegisterManager load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new IRegisterManager(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static IRegisterManager load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new IRegisterManager(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
