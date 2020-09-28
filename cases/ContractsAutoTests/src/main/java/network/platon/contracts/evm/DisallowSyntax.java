package network.platon.contracts.evm;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class DisallowSyntax extends Contract {
    private static final String BINARY = "";

    public static final String FUNC_GETMSGVALUE = "getMsgValue";

    public static final String FUNC_GETVALUE = "getValue";

    public static final String FUNC_METHOD = "method";

    public static final String FUNC_MULVALUE2 = "mulvalue2";

    public static final String FUNC_TESRETURN = "tesReturn";

    public static final String FUNC_TESTBLOCK = "testBlock";

    protected DisallowSyntax(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DisallowSyntax(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> getMsgValue(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_GETMSGVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<BigInteger> getValue(String _to, BigInteger _value) {
        final Function function = new Function(FUNC_GETVALUE, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(_to), 
                new Uint256(_value)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> method() {
        final Function function = new Function(
                FUNC_METHOD, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> mulvalue2(BigInteger a, BigInteger b) {
        final Function function = new Function(
                FUNC_MULVALUE2, 
                Arrays.<Type>asList(new Uint256(a),
                new Uint256(b)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> tesReturn(BigInteger _id, BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TESRETURN, 
                Arrays.<Type>asList(new Uint256(_id)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<TransactionReceipt> testBlock(String _to, BigInteger _value) {
        final Function function = new Function(
                FUNC_TESTBLOCK, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.Address(_to), 
                new Uint256(_value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<DisallowSyntax> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowSyntax.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DisallowSyntax> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DisallowSyntax.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DisallowSyntax load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DisallowSyntax(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DisallowSyntax load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DisallowSyntax(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
