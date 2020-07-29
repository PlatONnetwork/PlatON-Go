package network.platon.contracts;

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
 * <p>Generated with web3j version 0.13.0.7.
 */
public class NamedCall extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610104806100206000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c8063d4b7eac3146037578063e9e3370e146087575b600080fd5b606a60048036036040811015604b57600080fd5b81019080803590602001909291908035906020019092919050505060aa565b604051808381526020018281526020019250505060405180910390f35b608d60ba565b604051808381526020018281526020019250505060405180910390f35b6000808284915091509250929050565b60008060c76001600260aa565b91509150909156fea265627a7a72315820f4e478d1119cbe4d98f7a801a121aec81ec54317d9d911cae67feeee23f491d464736f6c634300050d0032";

    public static final String FUNC_EXCHANGE = "exchange";

    public static final String FUNC_NAMECALL = "namecall";

    protected NamedCall(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected NamedCall(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> exchange(BigInteger key, BigInteger value) {
        final Function function = new Function(FUNC_EXCHANGE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(key), 
                new org.web3j.abi.datatypes.generated.Uint256(value)), 
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

    public RemoteCall<Tuple2<BigInteger, BigInteger>> namecall() {
        final Function function = new Function(FUNC_NAMECALL, 
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

    public static RemoteCall<NamedCall> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(NamedCall.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<NamedCall> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(NamedCall.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static NamedCall load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new NamedCall(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static NamedCall load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new NamedCall(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
