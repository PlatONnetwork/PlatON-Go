package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class LibDb extends Contract {
    private static final String BINARY = "6060604052341561000f57600080fd5b5b6104d38061001f6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063878487271461003e575b600080fd5b610114600480803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509190505061012a565b6040518082815260200191505060405180910390f35b6000610134610493565b6000806000606060405190810160405280602d81526020017f5b3639643938643661303463343162343630356161636237626432663734626581526020017f655d5b3037777269746564625d0000000000000000000000000000000000000081525093506101b56101a6896002610217565b8561024e90919063ffffffff16565b93506101d46101c5886004610217565b8561024e90919063ffffffff16565b93506101f36101e4876009610217565b8561024e90919063ffffffff16565b9350835191506020840192508183209050806001900494505b505050509392505050565b61021f610493565b610245836102378486516102b890919063ffffffff16565b61024e90919063ffffffff16565b90505b92915050565b610256610493565b6000806000845186510160405180591061026d5750595b908082528060200260200182016040525b50935060208601925060208501915060208401905061029f81848851610446565b6102ae86518201838751610446565b5b50505092915050565b6102c0610493565b60008060008092508591505b60008211156102f157600a828115156102e157fe5b04915082806001019350506102cc565b848310156102fd578492505b8260405180591061030b5750595b908082528060200260200182016040525b5093506001830390505b60008611156103bd576030600a8781151561033d57fe5b06017f010000000000000000000000000000000000000000000000000000000000000002848260000b81518110151561037257fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535080600190039050600a868115156103b557fe5b049550610326565b5b60008160000b12151561043c5760307f010000000000000000000000000000000000000000000000000000000000000002848260000b81518110151561040057fe5b9060200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a905350806001900390506103be565b5b50505092915050565b60005b60208210151561046f57825184526020840193506020830192505b602082039150610449565b6001826020036101000a039050801983511681855116818117865250505b50505050565b6020604051908101604052806000815250905600a165627a7a723058205f39b921c967a0a0d68f24a2dc531cc83d4088edfb1351521975c749b1a0dc7c0029";

    public static final String FUNC_WRITEDB = "writedb";

    @Deprecated
    protected LibDb(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected LibDb(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected LibDb(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected LibDb(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<BigInteger> writedb(String _name, String _key, String _value) {
        final Function function = new Function(FUNC_WRITEDB, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_name), 
                new org.web3j.abi.datatypes.Utf8String(_key), 
                new org.web3j.abi.datatypes.Utf8String(_value)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public static RemoteCall<LibDb> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(LibDb.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<LibDb> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(LibDb.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<LibDb> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(LibDb.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<LibDb> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(LibDb.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static LibDb load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new LibDb(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static LibDb load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new LibDb(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static LibDb load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new LibDb(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static LibDb load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new LibDb(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
