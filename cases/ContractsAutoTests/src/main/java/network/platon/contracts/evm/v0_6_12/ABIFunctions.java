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
public class ABIFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061030a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c8063538fad8b14610046578063911a3363146100c9578063b19d51e41461014c575b600080fd5b61004e6101cf565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561008e578082015181840152602081019050610073565b50505050905090810190601f1680156100bb5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100d1610216565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101115780820151818401526020810190506100f6565b50505050905090810190601f16801561013e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61015461023e565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610194578082015181840152602081019050610179565b50505050905090810190601f1680156101c15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b606060405160200180807f31000000000000000000000000000000000000000000000000000000000000008152506001019050604051602081830303815290604052905090565b6060600160405160200180828152602001915050604051602081830303815290604052905090565b60606001604051602401808281526020019150506040516020818303038152906040527f60fe47b1000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505090509056fea264697066735822122049658cf6412324ebf2b7fe351fdcaf7caf6ea68c71a96edfb9ec4729e80760bd64736f6c634300060c0033";

    public static final String FUNC_GETENCODE = "getEncode";

    public static final String FUNC_GETENCODEPACKED = "getEncodePacked";

    public static final String FUNC_GETENCODEWITHSIGNATURE = "getEncodeWithSignature";

    protected ABIFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ABIFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getEncode() {
        final Function function = new Function(
                FUNC_GETENCODE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getEncodePacked() {
        final Function function = new Function(
                FUNC_GETENCODEPACKED, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getEncodeWithSignature() {
        final Function function = new Function(
                FUNC_GETENCODEWITHSIGNATURE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<ABIFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ABIFunctions.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ABIFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ABIFunctions.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ABIFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ABIFunctions(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ABIFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ABIFunctions(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
