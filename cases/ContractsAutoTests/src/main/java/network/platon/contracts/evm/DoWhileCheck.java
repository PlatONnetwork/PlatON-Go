package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tuples.generated.Tuple2;
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
public class DoWhileCheck extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610101806100206000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063db60eb9e146044575b600080fd5b348015604f57600080fd5b5060566073565b604051808381526020018281526020019250505060405180910390f35b6000806000806000809250600091505b828060010193505081806001019250506014821115609f5760b7565b600a83111560ab576083565b601e83101560b7576083565b601583141560c457600c90505b600e905082819450945050505090915600a165627a7a723058206c4d3705f8740f6c667ac94b2189df0850071480b43a6b8e009803dd608456940029";

    public static final String FUNC_DOWHILECHECK = "doWhileCheck";

    protected DoWhileCheck(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DoWhileCheck(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> doWhileCheck() {
        final Function function = new Function(FUNC_DOWHILECHECK, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public static RemoteCall<DoWhileCheck> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DoWhileCheck.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DoWhileCheck> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DoWhileCheck.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DoWhileCheck load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DoWhileCheck(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DoWhileCheck load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DoWhileCheck(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
