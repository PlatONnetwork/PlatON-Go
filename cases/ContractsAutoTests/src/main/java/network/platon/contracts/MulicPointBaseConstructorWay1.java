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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class MulicPointBaseConstructorWay1 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50600a80600081905550506101528061002a6000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f1461005157806382ab890a1461007c575b600080fd5b34801561005d57600080fd5b506100666100fe565b6040518082815260200191505060405180910390f35b34801561008857600080fd5b506100b56004803603602081101561009f57600080fd5b8101908080359060200190929190505050610104565b604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390f35b60005481565b600080826000808282540192505081905550336000548191509150915091509156fea165627a7a72305820746ed30553cc2d5a071894b25a4ade6e0b31c4a32ab9b3ba3959c232988ad79f0029";

    public static final String FUNC_A = "a";

    public static final String FUNC_UPDATE = "update";

    @Deprecated
    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected MulicPointBaseConstructorWay1(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
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

    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<MulicPointBaseConstructorWay1> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(MulicPointBaseConstructorWay1.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static MulicPointBaseConstructorWay1 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new MulicPointBaseConstructorWay1(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
