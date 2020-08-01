package network.platon.contracts;

import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
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
public class TypeName extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101ec806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806306fdde031461003b578063f5f5ba72146100be575b600080fd5b610043610141565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610083578082015181840152602081019050610068565b50505050905090810190601f1680156100b05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100c661017a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101065780820151818401526020810190506100eb565b50505050905090810190601f1680156101335780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6040518060400160405260088152806020017f547970654e616d6500000000000000000000000000000000000000000000000081525081565b60606040518060400160405260088152806020017f547970654e616d6500000000000000000000000000000000000000000000000081525090509056fea265627a7a72315820806e5021d57326e3fbd5a8e67455ed01740e2f68f1208813beadfb024074c0b264736f6c634300050d0032";

    public static final String FUNC_GETCONTRACTNAME = "getContractName";

    public static final String FUNC_NAME = "name";

    protected TypeName(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected TypeName(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getContractName() {
        final Function function = new Function(
                FUNC_GETCONTRACTNAME, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> name() {
        final Function function = new Function(FUNC_NAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public static RemoteCall<TypeName> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeName.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<TypeName> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeName.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static TypeName load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new TypeName(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static TypeName load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new TypeName(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
