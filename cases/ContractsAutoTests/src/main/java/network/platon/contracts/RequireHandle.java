package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class RequireHandle extends Contract {
    private static final String BINARY = "6080604052734b0897b0513fdc7c541b6d9d7e929c4e5364d2db600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555034801561006557600080fd5b50610585806100756000396000f3fe6080604052600436106100705760003560e01c80635995caa71161004e5780635995caa714610185578063afcd320e1461019c578063ce602ba3146101d7578063e08302331461021257610070565b80632e230b19146100e1578063414f180e1461011c578063431d70d714610157575b34801561007c57600080fd5b5060007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc600a9081150290604051600060405180830381858888f193505050501580156100dd573d6000803e3d6000fd5b5050005b3480156100ed57600080fd5b5061011a6004803603602081101561010457600080fd5b810190808035906020019092919050505061024d565b005b34801561012857600080fd5b506101556004803603602081101561013f57600080fd5b81019080803590602001909291905050506102d3565b005b6101836004803603602081101561016d57600080fd5b8101908080359060200190929190505050610380565b005b34801561019157600080fd5b5061019a6103e4565b005b3480156101a857600080fd5b506101d5600480360360208110156101bf57600080fd5b810190808035906020019092919050505061044e565b005b3480156101e357600080fd5b50610210600480360360208110156101fa57600080fd5b810190808035906020019092919050505061045e565b005b34801561021e57600080fd5b5061024b6004803603602081101561023557600080fd5b81019080803590602001909291905050506104ca565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634d431097826040518263ffffffff1660e01b8152600401600060405180830381600088803b1580156102b757600080fd5b5087f11580156102cb573d6000803e3d6000fd5b505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663370158ea600a83906040518363ffffffff1660e01b81526004016020604051808303818589803b15801561033f57600080fd5b5088f1158015610353573d6000803e3d6000fd5b5050505050506040513d602081101561036b57600080fd5b81019080805190602001909291905050505050565b60007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f193505050501580156103df573d6000803e3d6000fd5b505050565b6040516103f09061052e565b604051809103906000f08015801561040c573d6000803e3d6000fd5b506000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550565b600a811061045b57600080fd5b50565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f193505050501580156104c6573d6000803e3d6000fd5b5050565b60007314723a09acff6d2a60dcdf7aa4aff308fddc160c90508073ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015610529573d6000803e3d6000fd5b505050565b60168061053b8339019056fe6080604052348015600f57600080fd5b50600080fdfea265627a7a723158205f536f857b4a931b9eb2b3b0ebf72ec38ba57fb3e23facd15634f89cb4a39fbd64736f6c634300050d0032";

    public static final String FUNC_FUNCTIONCALLEXCEPTION = "functionCallException";

    public static final String FUNC_NEWCONTRACTEXCEPTION = "newContractException";

    public static final String FUNC_NONPAYABLERECEIVEETHEXCEPTION = "nonPayableReceiveEthException";

    public static final String FUNC_OUTFUNCTIONCALLEXCEPTION = "outFunctionCallException";

    public static final String FUNC_PARAMEXCEPTION = "paramException";

    public static final String FUNC_PUBLICGETTERRECEIVEETHEXCEPTION = "publicGetterReceiveEthException";

    public static final String FUNC_TRANSFERCALLEXCEPTION = "transferCallException";

    protected RequireHandle(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RequireHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<RequireHandle> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RequireHandle.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RequireHandle> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RequireHandle.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> functionCallException(BigInteger param) {
        final Function function = new Function(
                FUNC_FUNCTIONCALLEXCEPTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> newContractException() {
        final Function function = new Function(
                FUNC_NEWCONTRACTEXCEPTION, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> nonPayableReceiveEthException(BigInteger count) {
        final Function function = new Function(
                FUNC_NONPAYABLERECEIVEETHEXCEPTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(count)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> outFunctionCallException(BigInteger count) {
        final Function function = new Function(
                FUNC_OUTFUNCTIONCALLEXCEPTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(count)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> paramException(BigInteger param) {
        final Function function = new Function(
                FUNC_PARAMEXCEPTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> publicGetterReceiveEthException(BigInteger count) {
        final Function function = new Function(
                FUNC_PUBLICGETTERRECEIVEETHEXCEPTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(count)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> transferCallException(BigInteger count, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TRANSFERCALLEXCEPTION, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(count)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public static RequireHandle load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RequireHandle(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RequireHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RequireHandle(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
