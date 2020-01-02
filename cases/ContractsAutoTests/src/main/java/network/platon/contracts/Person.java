package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
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
 * <p>Generated with web3j version 0.7.5.0.
 */
public class Person extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b5b60b4600181905550601460008190555033600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040805190810160405280600681526020017f71637869616f0000000000000000000000000000000000000000000000000000815250600390805190602001906100ac9291906100b3565b505b610158565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100f457805160ff1916838001178555610122565b82800160010185558215610122579182015b82811115610121578251825591602001919060010190610106565b5b50905061012f9190610133565b5090565b61015591905b80821115610151576000816000905550600101610139565b5090565b90565b6104b3806101676000396000f3006060604052361561008c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100915780630ef267431461012057806326121ff014610149578063262a9dff1461015e57806341c0e1b514610187578063741a39441461019c5780638da5cb5b146101bf578063d5dcf12714610214575b600080fd5b341561009c57600080fd5b6100a4610237565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100e55780820151818401525b6020810190506100c9565b50505050905090810190601f1680156101125780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561012b57600080fd5b6101336102e0565b6040518082815260200191505060405180910390f35b341561015457600080fd5b61015c6102eb565b005b341561016957600080fd5b6101716102f8565b6040518082815260200191505060405180910390f35b341561019257600080fd5b61019a610302565b005b34156101a757600080fd5b6101bd6004808035906020019091905050610396565b005b34156101ca57600080fd5b6101d26103a1565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b341561021f57600080fd5b61023560048080359060200190919050506103cc565b005b61023f610473565b60038054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102d55780601f106102aa576101008083540402835291602001916102d5565b820191906000526020600020905b8154815290600101906020018083116102b857829003601f168201915b505050505090505b90565b600060015490505b90565b6102f560036103d7565b5b565b6000805490505b90565b3373ffffffffffffffffffffffffffffffffffffffff16600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141561039357600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b5b565b806001819055505b50565b6000600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690505b90565b806000819055505b50565b7f5100000000000000000000000000000000000000000000000000000000000000816000815460018160011615610100020316600290048110151561041857fe5b8154600116156104375790600052602060002090602091828204019190065b601f036101000a81548160ff021916907f0100000000000000000000000000000000000000000000000000000000000000840402179055505b50565b6020604051908101604052806000815250905600a165627a7a72305820031578f357f7f9eabfb7d7f49272c303cf53ac11ea08ff49535df030381597b60029";

    public static final String FUNC_NAME = "name";

    public static final String FUNC_HEIGHT = "height";

    public static final String FUNC_F = "f";

    public static final String FUNC_AGE = "age";

    public static final String FUNC_KILL = "kill";

    public static final String FUNC_SETHEIGHT = "setHeight";

    public static final String FUNC_OWNER = "owner";

    public static final String FUNC_SETAGE = "setAge";

    @Deprecated
    protected Person(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected Person(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected Person(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected Person(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<String> name() {
        final Function function = new Function(FUNC_NAME, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<BigInteger> height() {
        final Function function = new Function(FUNC_HEIGHT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> f() {
        final Function function = new Function(
                FUNC_F, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> age() {
        final Function function = new Function(FUNC_AGE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public void kill() {
        throw new RuntimeException("cannot call constant function with void return type");
    }

    public RemoteCall<TransactionReceipt> setHeight(BigInteger height) {
        final Function function = new Function(
                FUNC_SETHEIGHT, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(height)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<String> owner() {
        final Function function = new Function(FUNC_OWNER, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<TransactionReceipt> setAge(BigInteger age) {
        final Function function = new Function(
                FUNC_SETAGE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(age)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Person> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(Person.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    public static RemoteCall<Person> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(Person.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<Person> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(Person.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<Person> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(Person.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static Person load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new Person(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static Person load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new Person(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static Person load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new Person(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static Person load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new Person(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
