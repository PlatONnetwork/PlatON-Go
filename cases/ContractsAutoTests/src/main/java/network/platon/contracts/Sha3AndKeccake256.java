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
public class Sha3AndKeccake256 extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061024f806100206000396000f3fe60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063598a717b146100515780638f9e1a5d1461007c575b600080fd5b34801561005d57600080fd5b50610066610144565b6040518082815260200191505060405180910390f35b34801561008857600080fd5b506101426004803603602081101561009f57600080fd5b81019080803590602001906401000000008111156100bc57600080fd5b8201836020820111156100ce57600080fd5b803590602001918460018302840111640100000000831117156100f057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061014d565b005b60008054905090565b80816040516020018083805190602001908083835b6020831015156101875780518252602082019150602081019050602083039250610162565b6001836020036101000a03801982511681845116808217855250505050505090500182805190602001908083835b6020831015156101da57805182526020820191506020810190506020830392506101b5565b6001836020036101000a0380198251168184511680821785525050505050509050019250505060405160208183030381529060405280519060200120600190046000819055505056fea165627a7a7230582051570af14839d43ac105f9f30d052ce82780e6b0339cc211886ce8532668d26b0029";

    public static final String FUNC_GETKECCAK256VALUE = "getKeccak256Value";

    public static final String FUNC_KECCAK = "keccak";

    protected Sha3AndKeccake256(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected Sha3AndKeccake256(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getKeccak256Value() {
        final Function function = new Function(FUNC_GETKECCAK256VALUE, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> keccak(String sha256value) {
        final Function function = new Function(
                FUNC_KECCAK, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(sha256value)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<Sha3AndKeccake256> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Sha3AndKeccake256.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<Sha3AndKeccake256> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(Sha3AndKeccake256.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static Sha3AndKeccake256 load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new Sha3AndKeccake256(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static Sha3AndKeccake256 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new Sha3AndKeccake256(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
