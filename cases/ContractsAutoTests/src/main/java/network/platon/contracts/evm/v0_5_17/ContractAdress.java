package network.platon.contracts.evm.v0_5_17;

import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class ContractAdress extends Contract {
    private static final String BINARY = "6080604052348015600f57600080fd5b50604780601d6000396000f3fe6080604052348015600f57600080fd5b5000fea265627a7a723158201e2b8da001f5bfa8d514df4902c1fbcf2e91f58798e8a03d0c9fec5a0b9e4c9364736f6c63430005110032";

    protected ContractAdress(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ContractAdress(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<ContractAdress> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ContractAdress.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ContractAdress> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ContractAdress.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ContractAdress load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ContractAdress(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ContractAdress load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ContractAdress(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
