package network.platon.contracts.evm.v0_6_12;

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
public class InheritContractSubclass extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060e68061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c8063430fe9c1146041578063825b71d914605d578063887c5838146079575b600080fd5b60476095565b6040518082815260200191505060405180910390f35b6063609e565b6040518082815260200191505060405180910390f35b607f60a7565b6040518082815260200191505060405180910390f35b60006002905090565b60006004905090565b6000600390509056fea2646970667358221220e1e9ea957ec116e9956913e482fdf55ed11a20514296f57ea9da377e1182b2cb64736f6c634300060c0033";

    public static final String FUNC_GETDATATHREE = "getDataThree";

    public static final String FUNC_GETDATE = "getDate";

    public static final String FUNC_GETSUBDATA = "getSubData";

    protected InheritContractSubclass(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected InheritContractSubclass(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getDataThree() {
        final Function function = new Function(
                FUNC_GETDATATHREE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getDate() {
        final Function function = new Function(
                FUNC_GETDATE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getSubData() {
        final Function function = new Function(
                FUNC_GETSUBDATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<InheritContractSubclass> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractSubclass.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<InheritContractSubclass> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(InheritContractSubclass.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static InheritContractSubclass load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractSubclass(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static InheritContractSubclass load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new InheritContractSubclass(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
