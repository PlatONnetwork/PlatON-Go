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
public class MathAndCryptographicFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061035a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806301c740441461006757806301f56b78146100855780635b4aa3ee146100fe578063aa4e87441461011c578063cc98f30e1461013a578063f9b4169114610158575b600080fd5b61006f610176565b6040518082815260200191505060405180910390f35b6100d26004803603608081101561009b57600080fd5b8101908080359060200190929190803560ff16906020019092919080359060200190929190803590602001909291905050506101f0565b604051808273ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b610106610265565b6040518082815260200191505060405180910390f35b61012461028d565b6040518082815260200191505060405180910390f35b6101426102a2565b6040518082815260200191505060405180910390f35b61016061030f565b6040518082815260200191505060405180910390f35b6000600260405180807f41424300000000000000000000000000000000000000000000000000000000008152506003019050602060405180830381855afa1580156101c5573d6000803e3d6000fd5b5050506040513d60208110156101da57600080fd5b8101908080519060200190929190505050905090565b60008060018686868660405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561024d573d6000803e3d6000fd5b50505060206040510351905080915050949350505050565b60007fe1629b9dda060bb30c7908346f6af189c16773fa148d3366701fbaa35d54f3c8905090565b600060038061029857fe5b6003600209905090565b6000600360405180807f41424300000000000000000000000000000000000000000000000000000000008152506003019050602060405180830381855afa1580156102f1573d6000803e3d6000fd5b5050506040515160601b6bffffffffffffffffffffffff1916905090565b600060038061031a57fe5b600360020890509056fea2646970667358221220dec0fd0b0c96b2f39870776d6a9e183d65792718b5221c52ab3d4de28c4b92b564736f6c63430007010033";

    public static final String FUNC_CALLADDMOD = "callAddMod";

    public static final String FUNC_CALLECRECOVER = "callEcrecover";

    public static final String FUNC_CALLKECCAK256 = "callKeccak256";

    public static final String FUNC_CALLMULMOD = "callMulMod";

    public static final String FUNC_CALLRIPEMD160 = "callRipemd160";

    public static final String FUNC_CALLSHA256 = "callSha256";

    protected MathAndCryptographicFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MathAndCryptographicFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> callAddMod() {
        final Function function = new Function(
                FUNC_CALLADDMOD, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callEcrecover(byte[] hash, BigInteger v, byte[] r, byte[] s) {
        final Function function = new Function(
                FUNC_CALLECRECOVER, 
                Arrays.<Type>asList(new com.alaya.abi.solidity.datatypes.generated.Bytes32(hash), 
                new com.alaya.abi.solidity.datatypes.generated.Uint8(v), 
                new com.alaya.abi.solidity.datatypes.generated.Bytes32(r), 
                new com.alaya.abi.solidity.datatypes.generated.Bytes32(s)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callKeccak256() {
        final Function function = new Function(
                FUNC_CALLKECCAK256, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callMulMod() {
        final Function function = new Function(
                FUNC_CALLMULMOD, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callRipemd160() {
        final Function function = new Function(
                FUNC_CALLRIPEMD160, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callSha256() {
        final Function function = new Function(
                FUNC_CALLSHA256, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<MathAndCryptographicFunctions> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MathAndCryptographicFunctions.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<MathAndCryptographicFunctions> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MathAndCryptographicFunctions.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static MathAndCryptographicFunctions load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MathAndCryptographicFunctions(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MathAndCryptographicFunctions load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MathAndCryptographicFunctions(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
