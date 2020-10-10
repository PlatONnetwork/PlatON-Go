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
public class CreationCode extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610314806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063ade003e81461003b578063f5f5ba72146100be575b600080fd5b610043610141565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610083578082015181840152602081019050610068565b50505050905090810190601f1680156100b05780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6100c6610168565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101065780820151818401526020810190506100eb565b50505050905090810190601f1680156101335780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6040518060200161015190610193565b6020820181038252601f19601f8201166040525081565b60606040518060200161017a90610193565b6020820181038252601f19601f82011660405250905090565b61013e806101a18339019056fe608060405234801561001057600080fd5b5061011e806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80632096525514602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600581526020017f68656c6c6f00000000000000000000000000000000000000000000000000000081525090509056fea2646970667358221220846d2d25a17e63808df82aac08490cf1c808bf05153d59600a435bd040bb8b1564736f6c63430007010033a26469706673582212205237accec072f89c86ce6a9ac90e49d250b6c257ffe5a2c371b248686caaf4fb64736f6c63430007010033";

    public static final String FUNC_CREATIONCODEINFO = "creationCodeInfo";

    public static final String FUNC_GETCONTRACTNAME = "getContractName";

    protected CreationCode(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected CreationCode(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<TransactionReceipt> creationCodeInfo() {
        final Function function = new Function(
                FUNC_CREATIONCODEINFO, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getContractName() {
        final Function function = new Function(
                FUNC_GETCONTRACTNAME, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<CreationCode> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CreationCode.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<CreationCode> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(CreationCode.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static CreationCode load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new CreationCode(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static CreationCode load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new CreationCode(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
