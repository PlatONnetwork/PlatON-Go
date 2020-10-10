package network.platon.contracts.evm.v0_4_26;

import com.alaya.abi.solidity.FunctionEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tuples.generated.Tuple6;
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
public class SameNameConstructorPublicVisibility extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5060405160208061016783398101806040528101908080519060200190929190505050806000819055505061011d8061004a6000396000f300608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634c9b47a614604d575b348015604a57600080fd5b50005b348015605857600080fd5b50605f6098565b60405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390f35b60008060008060008060008060008060006301e133809450680dd2d5fcf3bc9c0000935060ff925060ff9150680dd2d5fcf3bc9c0000905060005485858585859a509a509a509a509a509a5050505050509091929394955600a165627a7a723058202b4a959a09674db3619868c043e25deae894ba277821181d7b80fcc5ba0742ce0029";

    public static final String FUNC_DISCARDLITERALSANDSUFFIXES = "discardLiteralsAndSuffixes";

    protected SameNameConstructorPublicVisibility(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected SameNameConstructorPublicVisibility(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>> discardLiteralsAndSuffixes() {
        final Function function = new Function(FUNC_DISCARDLITERALSANDSUFFIXES, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple6<BigInteger, BigInteger, BigInteger, BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (BigInteger) results.get(3).getValue(), 
                                (BigInteger) results.get(4).getValue(), 
                                (BigInteger) results.get(5).getValue());
                    }
                });
    }

    public static RemoteCall<SameNameConstructorPublicVisibility> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId, BigInteger param) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(param)));
        return deployRemoteCall(SameNameConstructorPublicVisibility.class, web3j, credentials, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static RemoteCall<SameNameConstructorPublicVisibility> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId, BigInteger param) {
        String encodedConstructor = FunctionEncoder.encodeConstructor(Arrays.<Type>asList(new Uint256(param)));
        return deployRemoteCall(SameNameConstructorPublicVisibility.class, web3j, transactionManager, contractGasProvider, BINARY, encodedConstructor, chainId);
    }

    public static SameNameConstructorPublicVisibility load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new SameNameConstructorPublicVisibility(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static SameNameConstructorPublicVisibility load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new SameNameConstructorPublicVisibility(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
