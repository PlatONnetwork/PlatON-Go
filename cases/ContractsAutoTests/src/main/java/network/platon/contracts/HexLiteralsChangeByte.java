package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes1;
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
public class HexLiteralsChangeByte extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506102e5806100206000396000f3fe608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630b7f16651461005c578063420343a4146100cb578063ee4950021461012d575b600080fd5b34801561006857600080fd5b5061007161019c565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6100d36101d1565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561013957600080fd5b50610142610288565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b60008060009054906101000a90047f010000000000000000000000000000000000000000000000000000000000000002905090565b6000807f6162000000000000000000000000000000000000000000000000000000000000905060f17f0100000000000000000000000000000000000000000000000000000000000000026000806101000a81548160ff02191690837f0100000000000000000000000000000000000000000000000000000000000000900402179055506000809054906101000a90047f01000000000000000000000000000000000000000000000000000000000000000291505090565b6000809054906101000a90047f0100000000000000000000000000000000000000000000000000000000000000028156fea165627a7a723058209679a5a763e7564d9196c2523c6edfc0f5bd7cdd25d9e9149a703632202ff9010029";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_TESTCHANGE = "testChange";

    public static final String FUNC_B1 = "b1";

    protected HexLiteralsChangeByte(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected HexLiteralsChangeByte(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> getY() {
        final Function function = new Function(FUNC_GETY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<TransactionReceipt> testChange(BigInteger vonValue) {
        final Function function = new Function(
                FUNC_TESTCHANGE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function, vonValue);
    }

    public RemoteCall<byte[]> b1() {
        final Function function = new Function(FUNC_B1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public static RemoteCall<HexLiteralsChangeByte> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(HexLiteralsChangeByte.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<HexLiteralsChangeByte> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(HexLiteralsChangeByte.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static HexLiteralsChangeByte load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new HexLiteralsChangeByte(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static HexLiteralsChangeByte load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new HexLiteralsChangeByte(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
