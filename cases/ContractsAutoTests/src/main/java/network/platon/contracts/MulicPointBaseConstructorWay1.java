package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class MulicPointBaseConstructorWay1 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600a80600081905550506101528061002a6000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f1461005157806382ab890a1461007c575b600080fd5b34801561005d57600080fd5b506100666100fe565b6040518082815260200191505060405180910390f35b34801561008857600080fd5b506100b56004803603602081101561009f57600080fd5b8101908080359060200190929190505050610104565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b60005481565b600080826000808282540192505081905550336000548191509150915091509156fea165627a7a723058206c3a5812d90b08f0206032e044c7eccd9d69d92fe2212425cfd17a42dbf1d5400029";

    public static final String FUNC_A = "a";

    public static final String FUNC_UPDATE = "update";

    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> update(BigInteger amount) {
        final Function function = new Function(
                FUNC_UPDATE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(amount)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
