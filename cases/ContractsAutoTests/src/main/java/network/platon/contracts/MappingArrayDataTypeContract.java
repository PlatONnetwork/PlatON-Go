package network.platon.contracts;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple4;
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
public class MappingArrayDataTypeContract extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b5061034e806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80631175c0491461003b5780638e7f79f6146100a0575b600080fd5b61009e600480360360a081101561005157600080fd5b81019080803560ff169060200190929190803560ff169060200190929190803560ff169060200190929190803560ff169060200190929190803560ff169060200190929190505050610112565b005b6100cf600480360360208110156100b657600080fd5b81019080803560ff169060200190929190505050610137565b604051808560ff1660ff1681526020018460ff1660ff1681526020018360ff1660ff1681526020018260ff1660ff16815260200194505050505060405180910390f35b61011f600086868661021b565b505061012e600286848461021b565b50505050505050565b6000806000806000806002811061014a57fe5b0160008660ff1660ff16815260200190815260200160002060009054906101000a900460ff16600060016002811061017e57fe5b0160008760ff1660ff16815260200190815260200160002060009054906101000a900460ff1660026000600281106101b257fe5b0160008860ff1660ff16815260200190815260200160002060009054906101000a900460ff1660026001600281106101e657fe5b0160008960ff1660ff16815260200190815260200160002060009054906101000a900460ff1693509350935093509193509193565b60008060008660006002811061022d57fe5b0160008760ff1660ff16815260200190815260200160002060009054906101000a900460ff16905060008760016002811061026457fe5b0160008860ff1660ff16815260200190815260200160002060009054906101000a900460ff169050858860006002811061029a57fe5b0160008960ff1660ff16815260200190815260200160002060006101000a81548160ff021916908360ff16021790555084886001600281106102d857fe5b0160008960ff1660ff16815260200190815260200160002060006101000a81548160ff021916908360ff16021790555081819350935050509450949250505056fea265627a7a72315820b57c0dc78beac52b2791610471df8d1d03e6bb0944a2d95d07386a92f45f940764736f6c634300050d0032";

    public static final String FUNC_GETVALUEBYKEY = "getValueByKey";

    public static final String FUNC_SET = "set";

    protected MappingArrayDataTypeContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected MappingArrayDataTypeContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>> getValueByKey(BigInteger key) {
        final Function function = new Function(FUNC_GETVALUEBYKEY, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint8(key)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>() {}, new TypeReference<Uint8>() {}, new TypeReference<Uint8>() {}, new TypeReference<Uint8>() {}));
        return new RemoteCall<Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>>(
                new Callable<Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>>() {
                    @Override
                    public Tuple4<BigInteger, BigInteger, BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple4<BigInteger, BigInteger, BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue(), 
                                (BigInteger) results.get(2).getValue(), 
                                (BigInteger) results.get(3).getValue());
                    }
                });
    }

    public RemoteCall<TransactionReceipt> set(BigInteger key, BigInteger value_a1, BigInteger value_a2, BigInteger value_b1, BigInteger value_b2) {
        final Function function = new Function(
                FUNC_SET, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint8(key), 
                new org.web3j.abi.datatypes.generated.Uint8(value_a1), 
                new org.web3j.abi.datatypes.generated.Uint8(value_a2), 
                new org.web3j.abi.datatypes.generated.Uint8(value_b1), 
                new org.web3j.abi.datatypes.generated.Uint8(value_b2)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<MappingArrayDataTypeContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MappingArrayDataTypeContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<MappingArrayDataTypeContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(MappingArrayDataTypeContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static MappingArrayDataTypeContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new MappingArrayDataTypeContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static MappingArrayDataTypeContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new MappingArrayDataTypeContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
