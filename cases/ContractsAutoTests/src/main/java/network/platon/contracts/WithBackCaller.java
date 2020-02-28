package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
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
 * <p>Generated with web3j version 0.7.5.8-SNAPSHOT.
 */
public class WithBackCaller extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506109a3806100206000396000f3fe608060405234801561001057600080fd5b50600436106100565760003560e01c80621e257c1461005b5780630687590a1461009f57806308c2938b1461017a578063400f6a60146101fd578063de583cfa1461024b575b600080fd5b61009d6004803603602081101561007157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610269565b005b610178600480360360408110156100b557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156100f257600080fd5b82018360208201111561010457600080fd5b8035906020019184600183028401116401000000008311171561012657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506103be565b005b61018261069a565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101c25780820151818401526020810190506101a7565b50505050905090810190601f1680156101ef5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102496004803603604081101561021357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061073c565b005b6102536108c0565b6040518082815260200191505060405180910390f35b8073ffffffffffffffffffffffffffffffffffffffff166055603c604051602401808360ff1681526020018260ff168152602001925050506040516020818303038152906040527f771602f7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b60208310610351578051825260208201915060208101905060208303925061032e565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146103b3576040519150601f19603f3d011682016040523d82523d6000602084013e6103b8565b606091505b50505050565b600060608373ffffffffffffffffffffffffffffffffffffffff1683604051602401808060200180602001838103835260058152602001807f68656c6c6f000000000000000000000000000000000000000000000000000000815250602001838103825284818151815260200191508051906020019080838360005b8381101561045557808201518184015260208101905061043a565b50505050905090810190601f1680156104825780820380516001836020036101000a031916815260200191505b5093505050506040516020818303038152906040527fae49cd9c000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106105385780518252602082019150602081019050602083039250610515565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461059a576040519150601f19603f3d011682016040523d82523d6000602084013e61059f565b606091505b5091509150816105ae57600080fd5b8080602001905160208110156105c357600080fd5b81019080805160405193929190846401000000008211156105e357600080fd5b838201915060208201858111156105f957600080fd5b825186600182028301116401000000008211171561061657600080fd5b8083526020830192505050908051906020019080838360005b8381101561064a57808201518184015260208101905061062f565b50505050905090810190601f1680156106775780820380516001836020036101000a031916815260200191505b50604052505050600190805190602001906106939291906108c9565b5050505050565b606060018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107325780601f1061070757610100808354040283529160200191610732565b820191906000526020600020905b81548152906001019060200180831161071557829003601f168201915b5050505050905090565b600060608373ffffffffffffffffffffffffffffffffffffffff1683604051602401808281526020019150506040516020818303038152906040527f68875570000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b6020831061081857805182526020820191506020810190506020830392506107f5565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461087a576040519150601f19603f3d011682016040523d82523d6000602084013e61087f565b606091505b50915091508161088e57600080fd5b8080602001905160208110156108a357600080fd5b810190808051906020019092919050505060008190555050505050565b60008054905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061090a57805160ff1916838001178555610938565b82800160010185558215610938579182015b8281111561093757825182559160200191906001019061091c565b5b5090506109459190610949565b5090565b61096b91905b8082111561096757600081600090555060010161094f565b5090565b9056fea265627a7a7231582038e0e932b83c6c9b4fc6517eb5cc820c71c627b8a30afe7e2a0625ea0dd637f164736f6c634300050d0032";

    public static final String FUNC_CALLADDLTEST = "callAddlTest";

    public static final String FUNC_CALLDOUBLELTEST = "callDoublelTest";

    public static final String FUNC_CALLGETNAMETEST = "callgetNameTest";

    public static final String FUNC_GETSTRINGRESULT = "getStringResult";

    public static final String FUNC_GETUINTRESULT = "getuintResult";

    @Deprecated
    protected WithBackCaller(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected WithBackCaller(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected WithBackCaller(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected WithBackCaller(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> callAddlTest(String other) {
        final Function function = new Function(
                FUNC_CALLADDLTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(other)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callDoublelTest(String other, BigInteger a) {
        final Function function = new Function(
                FUNC_CALLDOUBLELTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(other), 
                new org.web3j.abi.datatypes.generated.Uint256(a)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callgetNameTest(String other, String name) {
        final Function function = new Function(
                FUNC_CALLGETNAMETEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(other), 
                new org.web3j.abi.datatypes.Utf8String(name)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> getStringResult() {
        final Function function = new Function(FUNC_GETSTRINGRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> getuintResult() {
        final Function function = new Function(FUNC_GETUINTRESULT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<WithBackCaller> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(WithBackCaller.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<WithBackCaller> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(WithBackCaller.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<WithBackCaller> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(WithBackCaller.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<WithBackCaller> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(WithBackCaller.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static WithBackCaller load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new WithBackCaller(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static WithBackCaller load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new WithBackCaller(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static WithBackCaller load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new WithBackCaller(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static WithBackCaller load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new WithBackCaller(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
