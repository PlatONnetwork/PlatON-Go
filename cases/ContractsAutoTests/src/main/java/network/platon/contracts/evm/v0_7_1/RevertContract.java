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
public class RevertContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550610171806100606000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806314fef9361461003b578063de29278914610087575b600080fd5b6100716004803603604081101561005157600080fd5b8101908080359060200190929190803590602001909291905050506100a5565b6040518082815260200191505060405180910390f35b61008f610131565b6040518082815260200191505060405180910390f35b60008183101561011d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260168152602001807f6669727374206c657373207468616e207365636f6e640000000000000000000081525060200191505060405180910390fd5b818303600181905550600154905092915050565b600060015490509056fea2646970667358221220ba771aebc074f16f3fb259177d1d3f115df8311c6720673957d300a9ae3c78c464736f6c63430007010033";

    public static final String FUNC_GETRESULT = "getResult";

    public static final String FUNC_TOSENDERAMOUNT = "toSenderAmount";

    protected RevertContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected RevertContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static RemoteCall<RevertContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RevertContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<RevertContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(RevertContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public RemoteCall<TransactionReceipt> getResult() {
        final Function function = new Function(
                FUNC_GETRESULT, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> toSenderAmount(BigInteger first, BigInteger second) {
        final Function function = new Function(
                FUNC_TOSENDERAMOUNT, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Uint256(first), 
                new com.alaya.abi.solidity.datatypes.generated.Uint256(second)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RevertContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new RevertContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static RevertContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new RevertContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
