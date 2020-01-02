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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class Caller0513 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610749806100206000396000f3fe608060405234801561001057600080fd5b50600436106100405760003560e01c80621e257c146100455780633528d8fe14610089578063dfbf80b114610146575b600080fd5b6100876004803603602081101561005b57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610203565b005b6100cb6004803603602081101561009f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610358565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561010b5780820151818401526020810190506100f0565b50505050905090810190601f1680156101385780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101886004803603602081101561015c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506105b2565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101c85780820151818401526020810190506101ad565b50505050905090810190601f1680156101f55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b8073ffffffffffffffffffffffffffffffffffffffff166055603c604051602401808360ff1681526020018260ff168152602001925050506040516020818303038152906040527f771602f7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106102eb57805182526020820191506020810190506020830392506102c8565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461034d576040519150601f19603f3d011682016040523d82523d6000602084013e610352565b606091505b50505050565b6060600060608373ffffffffffffffffffffffffffffffffffffffff166040516024018080602001828103825260058152602001807f68656c6c6f0000000000000000000000000000000000000000000000000000008152506020019150506040516020818303038152906040527f6932cf81000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106104675780518252602082019150602081019050602083039250610444565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146104c9576040519150601f19603f3d011682016040523d82523d6000602084013e6104ce565b606091505b509150915060608180602001905160208110156104ea57600080fd5b810190808051604051939291908464010000000082111561050a57600080fd5b8382019150602082018581111561052057600080fd5b825186600182028301116401000000008211171561053d57600080fd5b8083526020830192505050908051906020019080838360005b83811015610571578082015181840152602081019050610556565b50505050905090810190601f16801561059e5780820380516001836020036101000a031916815260200191505b506040525050509050809350505050919050565b6060600060608373ffffffffffffffffffffffffffffffffffffffff166070604051602401808260ff1681526020019150506040516020818303038152906040527feee97206000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040518082805190602001908083835b602083106106945780518252602082019150602081019050602083039250610671565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146106f6576040519150601f19603f3d011682016040523d82523d6000602084013e6106fb565b606091505b50915091508161070a57600080fd5b809250505091905056fea265627a7a723158202e1b99bd708d8d5511e4c2750c859c3d3bdbec6966e0bb52d5c38b268e2b7ebd64736f6c634300050d0032";

    public static final String FUNC_CALLADDLTEST = "callAddlTest";

    public static final String FUNC_CALLDOUBLELTEST = "callDoublelTest";

    public static final String FUNC_CALLGETNAMETEST = "callgetNameTest";

    @Deprecated
    protected Caller0513(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected Caller0513(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected Caller0513(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected Caller0513(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> callAddlTest(String other) {
        final Function function = new Function(
                FUNC_CALLADDLTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(other)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callDoublelTest(String other) {
        final Function function = new Function(
                FUNC_CALLDOUBLELTEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(other)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callgetNameTest(String other) {
        final Function function = new Function(
                FUNC_CALLGETNAMETEST, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Address(other)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Caller0513> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(Caller0513.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<Caller0513> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(Caller0513.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<Caller0513> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(Caller0513.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<Caller0513> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(Caller0513.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static Caller0513 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new Caller0513(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static Caller0513 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new Caller0513(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static Caller0513 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new Caller0513(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static Caller0513 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new Caller0513(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
