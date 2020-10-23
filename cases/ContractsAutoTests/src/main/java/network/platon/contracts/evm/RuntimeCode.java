package network.platon.contracts.evm;

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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class RuntimeCode extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506102f2806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063d489e0041461003b578063f5f5ba72146100be575b600080fd5b610043610141565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610083578082015181840152602081019050610068565b50505050905090810190601f1680156100b05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100c6610168565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101065780820151818401526020810190506100eb565b50505050905090810190601f1680156101335780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6040518060200161015190610193565b6020820181038252601f19601f8201166040525081565b60606040518060200161017a90610193565b6020820181038252601f19601f82011660405250905090565b61011d806101a18339019056fe6080604052348015600f57600080fd5b506004361060285760003560e01c80632096525514602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525090509056fea265627a7a72315820b2c1484bf518c2712d9261d046e522de82053b11ee9ab5497dc8dc8a5d34d62164736f6c634300050d0032a265627a7a72315820dd80b9cf03ac803e0eeda00f6397eb11109ca0382b1581c9494dfe255e0d01e564736f6c634300050d0032";

    public static final String FUNC_GETCONTRACTNAME = "getContractName";

    public static final String FUNC_RUNTIMECODEINFO = "runtimeCodeInfo";

    protected RuntimeCode(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RuntimeCode(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> getContractName() {
        final Function function = new Function(FUNC_GETCONTRACTNAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<byte[]> runtimeCodeInfo() {
        final Function function = new Function(FUNC_RUNTIMECODEINFO, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<RuntimeCode> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RuntimeCode.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RuntimeCode> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RuntimeCode.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static RuntimeCode load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RuntimeCode(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RuntimeCode load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RuntimeCode(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
