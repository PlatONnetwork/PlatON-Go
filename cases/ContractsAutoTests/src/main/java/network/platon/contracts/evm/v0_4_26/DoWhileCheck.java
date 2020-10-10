package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tuples.generated.Tuple2;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.List;
import java.util.concurrent.Callable;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class DoWhileCheck extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610101806100206000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063db60eb9e146044575b600080fd5b348015604f57600080fd5b5060566073565b604051808381526020018281526020019250505060405180910390f35b6000806000806000809250600091505b828060010193505081806001019250506014821115609f5760b7565b600a83111560ab576083565b601e83101560b7576083565b601583141560c457600c90505b600e905082819450945050505090915600a165627a7a723058201e424c84649f9e69730ffbc213d5ba9d0674a612e78d15dd11a60d329fb133290029";

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
