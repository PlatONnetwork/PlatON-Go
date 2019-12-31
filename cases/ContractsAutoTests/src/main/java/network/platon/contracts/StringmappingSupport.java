package network.platon.contracts;

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
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.7.5.0.
 */
public class StringmappingSupport extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610769806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063215e59a314610051578063d4d7306b1461005b578063e78855a8146101ad578063f55aa68e146102e1575b600080fd5b610059610364565b005b6101ab6004803603604081101561007157600080fd5b810190808035906020019064010000000081111561008e57600080fd5b8201836020820111156100a057600080fd5b803590602001918460018302840111640100000000831117156100c257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561012557600080fd5b82018360208201111561013757600080fd5b8035906020019184600183028401116401000000008311171561015957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610409565b005b610266600480360360208110156101c357600080fd5b81019080803590602001906401000000008111156101e057600080fd5b8201836020820111156101f257600080fd5b8035906020019184600183028401116401000000008311171561021457600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061048b565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156102a657808201518184015260208101905061028b565b50505050905090810190601f1680156102d35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102e9610596565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561032957808201518184015260208101905061030e565b50505050905090810190601f1680156103565780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000686c697975656368756e60b81b90506040518060400160405280600881526020017f687564656e69616e000000000000000000000000000000000000000000000000815250600160008376ffffffffffffffffffffffffffffffffffffffffffffff191676ffffffffffffffffffffffffffffffffffffffffffffff19168152602001908152602001600020908051906020019061040592919061068f565b5050565b806000836040518082805190602001908083835b60208310610440578051825260208201915060208101905060208303925061041d565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020908051906020019061048692919061068f565b505050565b60606000826040518082805190602001908083835b602083106104c357805182526020820191506020810190506020830392506104a0565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561058a5780601f1061055f5761010080835404028352916020019161058a565b820191906000526020600020905b81548152906001019060200180831161056d57829003601f168201915b50505050509050919050565b60606000686c697975656368756e60b81b9050600160008276ffffffffffffffffffffffffffffffffffffffffffffff191676ffffffffffffffffffffffffffffffffffffffffffffff191681526020019081526020016000208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156106845780601f1061065957610100808354040283529160200191610684565b820191906000526020600020905b81548152906001019060200180831161066757829003601f168201915b505050505091505090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106d057805160ff19168380011785556106fe565b828001600101855582156106fe579182015b828111156106fd5782518255916020019190600101906106e2565b5b50905061070b919061070f565b5090565b61073191905b8082111561072d576000816000905550600101610715565b5090565b9056fea265627a7a7231582050d6ac405abc0e66e59825881f35e8b257b2a58887fa2b0e5677d5dcf166d6ae64736f6c634300050d0032";

    public static final String FUNC_GETBYTE32MAPVALUE = "getByte32mapValue";

    public static final String FUNC_GETSTRINGMAPVALUE = "getStringmapValue";

    public static final String FUNC_SETBYTE32MAPVALUE = "setByte32mapValue";

    public static final String FUNC_SETSTRINGMAPVALUE = "setStringmapValue";

    @Deprecated
    protected StringmappingSupport(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected StringmappingSupport(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    @Deprecated
    protected StringmappingSupport(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(BINARY, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected StringmappingSupport(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public RemoteCall<TransactionReceipt> getByte32mapValue() {
        final Function function = new Function(
                FUNC_GETBYTE32MAPVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> getStringmapValue(String _key) {
        final Function function = new Function(
                FUNC_GETSTRINGMAPVALUE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_key)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setByte32mapValue() {
        final Function function = new Function(
                FUNC_SETBYTE32MAPVALUE, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setStringmapValue(String _key, String _value) {
        final Function function = new Function(
                FUNC_SETSTRINGMAPVALUE, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(_key), 
                new org.web3j.abi.datatypes.Utf8String(_value)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<StringmappingSupport> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return deployRemoteCall(StringmappingSupport.class, web3j, credentials, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<StringmappingSupport> deploy(Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(StringmappingSupport.class, web3j, credentials, gasPrice, gasLimit, BINARY, "");
    }

    public static RemoteCall<StringmappingSupport> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return deployRemoteCall(StringmappingSupport.class, web3j, transactionManager, contractGasProvider, BINARY, "");
    }

    @Deprecated
    public static RemoteCall<StringmappingSupport> deploy(Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(StringmappingSupport.class, web3j, transactionManager, gasPrice, gasLimit, BINARY, "");
    }

    @Deprecated
    public static StringmappingSupport load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new StringmappingSupport(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    @Deprecated
    public static StringmappingSupport load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new StringmappingSupport(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    public static StringmappingSupport load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new StringmappingSupport(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static StringmappingSupport load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new StringmappingSupport(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
