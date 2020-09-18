package network.platon.contracts.evm;

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
 * <p>Generated with web3j version 0.13.1.5.
 */
public class DoWhileLogicAnd99Style extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506101d6806100206000396000f3fe608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806349de2f961461006757806361cc1193146100925780637cf7eab0146100cd578063c1633e2214610108575b600080fd5b34801561007357600080fd5b5061007c610133565b6040518082815260200191505060405180910390f35b34801561009e57600080fd5b506100cb600480360360208110156100b557600080fd5b810190808035906020019092919050505061013c565b005b3480156100d957600080fd5b50610106600480360360208110156100f057600080fd5b8101908080359060200190929190505050610176565b005b34801561011457600080fd5b5061011d6101a0565b6040518082815260200191505060405180910390f35b60008054905090565b6000600a8201905060006009830190505b6001830192508083111561016057610161565b5b818310151561014d5782600181905550505050565b60008090505b8181101561019c578060005401600081905550808060010191505061017c565b5050565b600060015490509056fea165627a7a72305820581e094b4ab7d2207fe3b13e52014b86e10064ef3ee2f3a53267e70e59fad3820029";

    public static final String FUNC_GETFORSUM = "getForSum";

    public static final String FUNC_DOWHILE = "dowhile";

    public static final String FUNC_FORSUM = "forsum";

    public static final String FUNC_GETDOWHILESUM = "getDoWhileSum";

    protected DoWhileLogicAnd99Style(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected DoWhileLogicAnd99Style(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getForSum() {
        final Function function = new Function(FUNC_GETFORSUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> dowhile(BigInteger x) {
        final Function function = new Function(
                FUNC_DOWHILE, 
                Arrays.<Type>asList(new Uint256(x)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> forsum(BigInteger x) {
        final Function function = new Function(
                FUNC_FORSUM, 
                Arrays.<Type>asList(new Uint256(x)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getDoWhileSum() {
        final Function function = new Function(FUNC_GETDOWHILESUM, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<DoWhileLogicAnd99Style> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DoWhileLogicAnd99Style.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<DoWhileLogicAnd99Style> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(DoWhileLogicAnd99Style.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static DoWhileLogicAnd99Style load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new DoWhileLogicAnd99Style(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static DoWhileLogicAnd99Style load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new DoWhileLogicAnd99Style(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
