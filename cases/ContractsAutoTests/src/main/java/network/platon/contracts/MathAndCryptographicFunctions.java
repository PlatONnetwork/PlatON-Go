package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Bytes32;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
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
public class MathAndCryptographicFunctions extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061037e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806301c740441461006757806301f56b78146100855780635b4aa3ee14610114578063aa4e874414610132578063cc98f30e14610150578063f9b416911461016e575b600080fd5b61006f61018c565b6040518082815260200191505060405180910390f35b6100d26004803603608081101561009b57600080fd5b8101908080359060200190929190803560ff1690602001909291908035906020019092919080359060200190929190505050610206565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b61011c61027e565b6040518082815260200191505060405180910390f35b61013a6102bb565b6040518082815260200191505060405180910390f35b6101586102d0565b6040518082815260200191505060405180910390f35b61017661033d565b6040518082815260200191505060405180910390f35b6000600260405180807f41424300000000000000000000000000000000000000000000000000000000008152506003019050602060405180830381855afa1580156101db573d6000803e3d6000fd5b5050506040513d60208110156101f057600080fd5b8101908080519060200190929190505050905090565b60008060018686868660405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610266573d6000803e3d6000fd5b50505060206040510351905080915050949350505050565b600060405180807f414243000000000000000000000000000000000000000000000000000000000081525060030190506040518091039020905090565b60006003806102c657fe5b6003600209905090565b6000600360405180807f41424300000000000000000000000000000000000000000000000000000000008152506003019050602060405180830381855afa15801561031f573d6000803e3d6000fd5b5050506040515160601b6bffffffffffffffffffffffff1916905090565b600060038061034857fe5b600360020890509056fea165627a7a723058201675c7d75fa2ea48bae0088edb0613c231b2f87fb7c4da45bb387372da6911f40029";

    public static final String FUNC_CALLSHA256 = "callSha256";

    public static final String FUNC_CALLECRECOVER = "callEcrecover";

    public static final String FUNC_CALLKECCAK256 = "callKeccak256";

    public static final String FUNC_CALLMULMOD = "callMulMod";

    public static final String FUNC_CALLRIPEMD160 = "callRipemd160";

    public static final String FUNC_CALLADDMOD = "callAddMod";

    protected MathAndCryptographicFunctions(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MathAndCryptographicFunctions(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<byte[]> callSha256() {
        final Function function = new Function(FUNC_CALLSHA256, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<String> callEcrecover(byte[] hash, BigInteger v, byte[] r, byte[] s) {
        final Function function = new Function(FUNC_CALLECRECOVER, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Bytes32(hash), 
                new org.web3j.abi.datatypes.generated.Uint8(v), 
                new org.web3j.abi.datatypes.generated.Bytes32(r), 
                new org.web3j.abi.datatypes.generated.Bytes32(s)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Address>() {}));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public RemoteCall<byte[]> callKeccak256() {
        final Function function = new Function(FUNC_CALLKECCAK256, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> callMulMod() {
        final Function function = new Function(FUNC_CALLMULMOD, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<byte[]> callRipemd160() {
        final Function function = new Function(FUNC_CALLRIPEMD160, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bytes32>() {}));
        return executeRemoteCallSingleValueReturn(function, byte[].class);
    }

    public RemoteCall<BigInteger> callAddMod() {
        final Function function = new Function(FUNC_CALLADDMOD, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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
