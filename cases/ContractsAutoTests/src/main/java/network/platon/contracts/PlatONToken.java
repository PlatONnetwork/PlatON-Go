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
public class PlatONToken extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610189806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c8063249bb3731461006757806375efc40d14610085578063c2412676146100a3578063c951fdf6146100ad578063d87698ae146100cb578063eecb9ce9146100e9575b600080fd5b61006f610107565b6040518082815260200191505060405180910390f35b61008d610119565b6040518082815260200191505060405180910390f35b6100ab610128565b005b6100b5610132565b6040518082815260200191505060405180910390f35b6100d361013d565b6040518082815260200191505060405180910390f35b6100f1610143565b6040518082815260200191505060405180910390f35b6000670de0b6b3a76400008101905090565b600064e8d4a510008101905090565b6001600081905550565b600060018101905090565b60005481565b600066038d7ea4c68000810190509056fea265627a7a723158205561fc564849e1eb852075ff9a9d623bca01178493678aab877edf0f3dd686ad64736f6c634300050d0032";

    public static final String FUNC_PFINNEY = "Pfinney";

    public static final String FUNC_PLAT = "Plat";

    public static final String FUNC_PSZABO = "Pszabo";

    public static final String FUNC_PVON = "Pvon";

    public static final String FUNC_TOKEN = "Token";

    public static final String FUNC_PLATONTOKEN = "platontoken";

    protected PlatONToken(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected PlatONToken(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> Pfinney() {
        final Function function = new Function(FUNC_PFINNEY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> Plat() {
        final Function function = new Function(FUNC_PLAT, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> Pszabo() {
        final Function function = new Function(FUNC_PSZABO, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> Pvon() {
        final Function function = new Function(FUNC_PVON, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> Token() {
        final Function function = new Function(
                FUNC_TOKEN, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> platontoken() {
        final Function function = new Function(FUNC_PLATONTOKEN, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<PlatONToken> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatONToken.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<PlatONToken> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(PlatONToken.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static PlatONToken load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new PlatONToken(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static PlatONToken load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new PlatONToken(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
