package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.Utf8String;
import com.alaya.abi.solidity.datatypes.generated.Int256;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class AbstractContractSon extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061022f806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806305a83329146100465780637839360014610064578063ed55abd4146100e7575b600080fd5b61004e61016a565b6040518082815260200191505060405180910390f35b61006c610178565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100ac578082015181840152602081019050610091565b50505050905090810190601f1680156100d95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100ef6101b7565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561012f578082015181840152602081019050610114565b50505050905090810190601f16801561015c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600080601490508091505090565b6060806040518060400160405280600a81526020017f6661746865724e616d650000000000000000000000000000000000000000000081525090505090565b6060806040518060400160405280600781526020017f736f6e4e616d65000000000000000000000000000000000000000000000000008152509050809150509056fea2646970667358221220409111651484062d29643dcd0fe79342c68320db5eacf7f32763205c1f62ee1964736f6c63430007010033";

    public static final String FUNC_FATHERAGE = "fatherAge";

    public static final String FUNC_FATHERNAME = "fatherName";

    public static final String FUNC_SONNAME = "sonName";

    protected AbstractContractSon(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected AbstractContractSon(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> fatherAge() {
        final Function function = new Function(FUNC_FATHERAGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Int256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> fatherName() {
        final Function function = new Function(
                FUNC_FATHERNAME, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> sonName() {
        final Function function = new Function(FUNC_SONNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<AbstractContractSon> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractSon.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<AbstractContractSon> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(AbstractContractSon.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static AbstractContractSon load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractSon(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static AbstractContractSon load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new AbstractContractSon(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
