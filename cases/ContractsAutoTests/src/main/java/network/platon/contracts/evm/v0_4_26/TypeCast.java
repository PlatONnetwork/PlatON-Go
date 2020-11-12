package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes4;
import com.alaya.abi.solidity.datatypes.generated.Uint16;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tuples.generated.Tuple3;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class TypeCast extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101de806100206000396000f300608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806312e9b56d14610046575b600080fd5b34801561005257600080fd5b5061005b610103565b604051808461ffff1661ffff168152602001837bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19168152602001935050505060405180910390f35b600080600080600080600060127e01000000000000000000000000000000000000000000000000000000000000029350837e01000000000000000000000000000000000000000000000000000000000000900492506104d27c01000000000000000000000000000000000000000000000000000000000291506112347c0100000000000000000000000000000000000000000000000000000000029050828282965096509650505050509091925600a165627a7a7230582005e4c0f51008c902039dde9182acbebf7e1b055a4ef9496f9f1d523e20f9d93b0029";

    public static final String FUNC_TYPECAST = "typeCast";

    protected TypeCast(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected TypeCast(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple3<BigInteger, byte[], byte[]>> typeCast() {
        final Function function = new Function(FUNC_TYPECAST, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint16>() {}, new TypeReference<Bytes4>() {}, new TypeReference<Bytes4>() {}));
        return new RemoteCall<Tuple3<BigInteger, byte[], byte[]>>(
                new Callable<Tuple3<BigInteger, byte[], byte[]>>() {
                    @Override
                    public Tuple3<BigInteger, byte[], byte[]> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple3<BigInteger, byte[], byte[]>(
                                (BigInteger) results.get(0).getValue(), 
                                (byte[]) results.get(1).getValue(), 
                                (byte[]) results.get(2).getValue());
                    }
                });
    }

    public static RemoteCall<TypeCast> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeCast.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<TypeCast> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(TypeCast.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static TypeCast load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new TypeCast(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static TypeCast load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new TypeCast(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
