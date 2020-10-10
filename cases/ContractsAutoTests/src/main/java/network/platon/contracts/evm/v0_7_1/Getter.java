package network.platon.contracts.evm.v0_7_1;

import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
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
public class Getter extends Contract {
    private static final String BINARY = "6080604052600a60005534801561001557600080fd5b50610149806100256000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806326121ff01461003b57806373d4a13a14610060575b600080fd5b61004361007e565b604051808381526020018281526020019250505060405180910390f35b61006861010d565b6040518082815260200191505060405180910390f35b6000806000543073ffffffffffffffffffffffffffffffffffffffff166373d4a13a6040518163ffffffff1660e01b815260040160206040518083038186803b1580156100ca57600080fd5b505afa1580156100de573d6000803e3d6000fd5b505050506040513d60208110156100f457600080fd5b8101908080519060200190929190505050915091509091565b6000548156fea2646970667358221220130eaa544f0358226ceaa823fbc69c08f60b7c83de60c8cb3c4de46019f0b5c164736f6c63430007010033";

    public static final String FUNC_DATA = "data";

    public static final String FUNC_F = "f";

    protected Getter(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Getter(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> data() {
        final Function function = new Function(
                FUNC_DATA, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> f() {
        final Function function = new Function(
                FUNC_F, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Getter> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Getter.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Getter> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Getter.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Getter load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Getter(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Getter load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Getter(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
