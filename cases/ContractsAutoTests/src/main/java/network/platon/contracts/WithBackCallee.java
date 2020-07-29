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
public class WithBackCallee extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506103eb806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063688755701461003b578063ae49cd9c1461007d575b600080fd5b6100676004803603602081101561005157600080fd5b8101908080359060200190929190505050610248565b6040518082815260200191505060405180910390f35b6101cd6004803603604081101561009357600080fd5b81019080803590602001906401000000008111156100b057600080fd5b8201836020820111156100c257600080fd5b803590602001918460018302840111640100000000831117156100e457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561014757600080fd5b82018360208201111561015957600080fd5b8035906020019184600183028401116401000000008311171561017b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610259565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561020d5780820151818401526020810190506101f2565b50505050905090810190601f16801561023a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600080828301905080915050919050565b6060610265838361026d565b905092915050565b6060808390506060839050606081518351016040519080825280601f01601f1916602001820160405280156102b15781602001600182028038833980820191505090505b5090506060819050600080905060008090505b8551811015610332578581815181106102d957fe5b602001015160f81c60f81b8383806001019450815181106102f657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535080806001019150506102c4565b5060008090505b84518110156103a75784818151811061034e57fe5b602001015160f81c60f81b83838060010194508151811061036b57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508080600101915050610339565b5082955050505050509291505056fea265627a7a7231582035f2615bd6b3d826aa402dca1ec736d876c20b2147ee88da7e0a5433c9ebc7eb64736f6c634300050d0032";

    public static final String FUNC_GETDOUBLE = "getDouble";

    public static final String FUNC_GETNAME = "getName";

    protected WithBackCallee(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected WithBackCallee(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getDouble(BigInteger a) {
        final Function function = new Function(
                FUNC_GETDOUBLE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(a)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getName(String option, String name) {
        final Function function = new Function(
                FUNC_GETNAME, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(option), 
                new org.web3j.abi.datatypes.Utf8String(name)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<WithBackCallee> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(WithBackCallee.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<WithBackCallee> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(WithBackCallee.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static WithBackCallee load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new WithBackCallee(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static WithBackCallee load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new WithBackCallee(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
