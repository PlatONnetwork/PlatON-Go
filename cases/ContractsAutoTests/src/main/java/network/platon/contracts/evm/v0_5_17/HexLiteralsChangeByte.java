package network.platon.contracts.evm.v0_5_17;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Bytes1;
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
 * <p>Generated with web3j version 0.13.2.1.
 */
public class HexLiteralsChangeByte extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061022f806100206000396000f3fe6080604052600436106100345760003560e01c80630b7f166514610039578063420343a4146100a8578063ee4950021461010a575b600080fd5b34801561004557600080fd5b5061004e610179565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b6100b061018f565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b34801561011657600080fd5b5061011f6101e8565b60405180827effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff19167effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916815260200191505060405180910390f35b60008060009054906101000a900460f81b905090565b6000807f6162000000000000000000000000000000000000000000000000000000000000905060f160f81b6000806101000a81548160ff021916908360f81c02179055506000809054906101000a900460f81b91505090565b6000809054906101000a900460f81b8156fea265627a7a723158209764a2de70097315ecd3b947f08b64a94f63e56b573639ae9049969e75f2ef6864736f6c63430005110032";

    public static final String FUNC_B1 = "b1";

    public static final String FUNC_GETY = "getY";

    public static final String FUNC_TESTCHANGE = "testChange";

    protected HexLiteralsChangeByte(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected HexLiteralsChangeByte(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> b1() {
        final Function function = new Function(FUNC_B1, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes1>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
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
