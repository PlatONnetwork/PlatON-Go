package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
public class RevertHandle extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610157806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063852da1631461003b578063f76051e714610069575b600080fd5b6100676004803603602081101561005157600080fd5b8101908080359060200190929190505050610097565b005b6100956004803603602081101561007f57600080fd5b8101908080359060200190929190505050610111565b005b600a81111561010e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f636865636b20636174636820657863657074696f6e000000000000000000000081525060200191505060405180910390fd5b50565b600a81111561011f57600080fd5b5056fea265627a7a72315820b6e7ab798b03b222aa2af17f558166bf2e414728a9959f715a3818482d4a673264736f6c634300050d0032";

    public static final String FUNC_REVERTCHECK = "revertCheck";

    public static final String FUNC_REVERTREASONCHECK = "revertReasonCheck";

    protected RevertHandle(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RevertHandle(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> revertCheck(BigInteger param) {
        final Function function = new Function(
                FUNC_REVERTCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> revertReasonCheck(BigInteger param) {
        final Function function = new Function(
                FUNC_REVERTREASONCHECK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(param)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<RevertHandle> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RevertHandle.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RevertHandle> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RevertHandle.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static RevertHandle load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RevertHandle(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RevertHandle load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RevertHandle(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
