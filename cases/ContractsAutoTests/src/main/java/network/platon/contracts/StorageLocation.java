package network.platon.contracts;

import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
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
public class StorageLocation extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610524806100206000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063066cfad114610051578063fcbc6ad714610105575b600080fd5b34801561005d57600080fd5b5061008a6004803603810190808035906020019082018035906020019190919293919293905050506101e7565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100ca5780820151818401526020810190506100af565b50505050905090810190601f1680156100f75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561011157600080fd5b5061016c600480360381019080803590602001908201803590602001908080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050919291929050505061029e565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ac578082015181840152602081019050610191565b50505050905090810190601f1680156101d95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60608282600091906101fa9291906103d3565b5060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102915780601f1061026657610100808354040283529160200191610291565b820191906000526020600020905b81548152906001019060200180831161027457829003601f168201915b5050505050905092915050565b606081600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002090805190602001906102f3929190610453565b50600160003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103c75780601f1061039c576101008083540402835291602001916103c7565b820191906000526020600020905b8154815290600101906020018083116103aa57829003601f168201915b50505050509050919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061041457803560ff1916838001178555610442565b82800160010185558215610442579182015b82811115610441578235825591602001919060010190610426565b5b50905061044f91906104d3565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061049457805160ff19168380011785556104c2565b828001600101855582156104c2579182015b828111156104c15782518255916020019190600101906104a6565b5b5090506104cf91906104d3565b5090565b6104f591905b808211156104f15760008160009055506001016104d9565b5090565b905600a165627a7a72305820c294cf6133f7facad482cbb918c4fb806bd4b06bc5a03ae7596776718a1d32720029";

    public static final String FUNC_TRANSFER = "transfer";

    public static final String FUNC_STORAGELOCALTIONCHECK = "storageLocaltionCheck";

    protected StorageLocation(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected StorageLocation(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> transfer(byte[] _data) {
        final Function function = new Function(FUNC_TRANSFER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicBytes(_data)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> storageLocaltionCheck(byte[] _data) {
        final Function function = new Function(FUNC_STORAGELOCALTIONCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.DynamicBytes(_data)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<StorageLocation> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(StorageLocation.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<StorageLocation> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(StorageLocation.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static StorageLocation load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new StorageLocation(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static StorageLocation load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new StorageLocation(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
