package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class AbstractContractGSubclass extends Contract {
    private static final String BINARY = "60806040526000805534801561001457600080fd5b5060f3806100236000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c8063262a9dff14604157806335f99d4814605d578063d1eef376146088575b600080fd5b604760a4565b6040518082815260200191505060405180910390f35b608660048036036020811015607157600080fd5b810190808035906020019092919050505060aa565b005b608e60b4565b6040518082815260200191505060405180910390f35b60005481565b8060008190555050565b6000805490509056fea264697066735822122088798519f8d0f0b8a4ee7367534a6b6cdb6cf2e0c46f3fed1cf13ae5ca238d5964736f6c63430007010033";

    public static final String FUNC_AINTERAGE = "aInterAge";

    public static final String FUNC_AGE = "age";

    public static final String FUNC_SETINTERAGE = "setInterAge";

    protected AbstractContractGSubclass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AbstractContractGSubclass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> aInterAge() {
        final Function function = new Function(
                FUNC_AINTERAGE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> age() {
        final Function function = new Function(
                FUNC_AGE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setInterAge(BigInteger v) {
        final Function function = new Function(
                FUNC_SETINTERAGE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Int256(v)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<AbstractContractGSubclass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractGSubclass.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AbstractContractGSubclass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractGSubclass.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AbstractContractGSubclass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractGSubclass(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AbstractContractGSubclass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractGSubclass(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
