package network.platon.contracts.evm;

import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tuples.generated.Tuple2;
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
public class ReferenceDataTypeArrayContract extends Contract {
    private static final String BINARY = "60806040526040518060a00160405280600160ff168152602001600260ff168152602001600360ff168152602001600460ff168152602001600560ff16815250600090600561004f92919061030b565b506040518060c001604052806040518060400160405280600181526020017f310000000000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600181526020017f320000000000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600181526020017f330000000000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600181526020017f340000000000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600181526020017f350000000000000000000000000000000000000000000000000000000000000081525081526020016040518060400160405280600181526020017f360000000000000000000000000000000000000000000000000000000000000081525081525060059060066101cb929190610350565b5060056040519080825280602002602001820160405280156101fc5781602001602082028038833980820191505090505b50600690805190602001906102129291906103b0565b506040518060c001604052806040518060400160405280600060ff168152602001600060ff1681525081526020016040518060400160405280600060ff168152602001600160ff1681525081526020016040518060400160405280600060ff168152602001600260ff1681525081526020016040518060400160405280600160ff168152602001600060ff1681525081526020016040518060400160405280600160ff168152602001600160ff1681525081526020016040518060400160405280600160ff168152602001600260ff1681525081525060079060066102f89291906103fd565b5034801561030557600080fd5b50610610565b826005810192821561033f579160200282015b8281111561033e578251829060ff1690559160200191906001019061031e565b5b50905061034c9190610458565b5090565b82805482825590600052602060002090810192821561039f579160200282015b8281111561039e57825182908051906020019061038e92919061047d565b5091602001919060010190610370565b5b5090506103ac91906104fd565b5090565b8280548282559060005260206000209081019282156103ec579160200282015b828111156103eb5782518255916020019190600101906103d0565b5b5090506103f99190610458565b5090565b828054828255906000526020600020908101928215610447579160200282015b8281111561044657825182906002610436929190610529565b509160200191906001019061041d565b5b509050610454919061057b565b5090565b61047a91905b8082111561047657600081600090555060010161045e565b5090565b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106104be57805160ff19168380011785556104ec565b828001600101855582156104ec579182015b828111156104eb5782518255916020019190600101906104d0565b5b5090506104f99190610458565b5090565b61052691905b80821115610522576000818161051991906105a7565b50600101610503565b5090565b90565b82805482825590600052602060002090810192821561056a579160200282015b82811115610569578251829060ff16905591602001919060010190610549565b5b5090506105779190610458565b5090565b6105a491905b808211156105a0576000818161059791906105ef565b50600101610581565b5090565b90565b50805460018160011615610100020316600290046000825580601f106105cd57506105ec565b601f0160209004906000526020600020908101906105eb9190610458565b5b50565b508054600082559060005260206000209081019061060d9190610458565b50565b6103cc8061061f6000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c80630849cc99146100675780630dca60821461008557806354c73338146100bd57806357933804146100c7578063ab35ec6314610182578063c3d1f404146101c4575b600080fd5b61006f6101e9565b6040518082815260200191505060405180910390f35b6100bb6004803603604081101561009b57600080fd5b8101908080359060200190929190803590602001909291905050506101f6565b005b6100c561020d565b005b610180600480360360208110156100dd57600080fd5b81019080803590602001906401000000008111156100fa57600080fd5b82018360208201111561010c57600080fd5b8035906020019184600183028401116401000000008311171561012e57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610243565b005b6101ae6004803603602081101561019857600080fd5b8101908080359060200190929190505050610285565b6040518082815260200191505060405180910390f35b6101cc61029c565b604051808381526020018281526020019250505060405180910390f35b6000600580549050905090565b806000836005811061020457fe5b01819055505050565b6064600760018154811061021d57fe5b9060005260206000200160008154811061023357fe5b9060005260206000200181905550565b60058190806001815401808255809150509060018203906000526020600020016000909192909190915090805190602001906102809291906102f2565b505050565b600080826005811061029357fe5b01549050919050565b60008060076001815481106102ad57fe5b906000526020600020016000815481106102c357fe5b906000526020600020015460076000815481106102dc57fe5b9060005260206000200180549050915091509091565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061033357805160ff1916838001178555610361565b82800160010185558215610361579182015b82811115610360578251825591602001919060010190610345565b5b50905061036e9190610372565b5090565b61039491905b80821115610390576000816000905550600101610378565b5090565b9056fea265627a7a723158204186c4211035ee13053e71c41b6bd991dd22b28f596a1f428965db61acbece7364736f6c634300050d0032";

    public static final String FUNC_GETARRAY = "getArray";

    public static final String FUNC_GETARRAYLENGTH = "getArrayLength";

    public static final String FUNC_GETMULTIARRAY = "getMultiArray";

    public static final String FUNC_SETARRAY = "setArray";

    public static final String FUNC_SETARRAYPUSH = "setArrayPush";

    public static final String FUNC_SETMULTIARRAY = "setMultiArray";

    protected ReferenceDataTypeArrayContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected ReferenceDataTypeArrayContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> getArray(BigInteger index) {
        final Function function = new Function(FUNC_GETARRAY, 
                Arrays.<Type>asList(new Uint256(index)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<BigInteger> getArrayLength() {
        final Function function = new Function(FUNC_GETARRAYLENGTH, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<Tuple2<BigInteger, BigInteger>> getMultiArray() {
        final Function function = new Function(FUNC_GETMULTIARRAY, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}, new TypeReference<Uint256>() {}));
        return new RemoteCall<Tuple2<BigInteger, BigInteger>>(
                new Callable<Tuple2<BigInteger, BigInteger>>() {
                    @Override
                    public Tuple2<BigInteger, BigInteger> call() throws Exception {
                        List<Type> results = executeCallMultipleValueReturn(function);
                        return new Tuple2<BigInteger, BigInteger>(
                                (BigInteger) results.get(0).getValue(), 
                                (BigInteger) results.get(1).getValue());
                    }
                });
    }

    public RemoteCall<TransactionReceipt> setArray(BigInteger index, BigInteger value) {
        final Function function = new Function(
                FUNC_SETARRAY, 
                Arrays.<Type>asList(new Uint256(index),
                new Uint256(value)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setArrayPush(String x) {
        final Function function = new Function(
                FUNC_SETARRAYPUSH, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(x)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> setMultiArray() {
        final Function function = new Function(
                FUNC_SETMULTIARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<ReferenceDataTypeArrayContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeArrayContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<ReferenceDataTypeArrayContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(ReferenceDataTypeArrayContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static ReferenceDataTypeArrayContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeArrayContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static ReferenceDataTypeArrayContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new ReferenceDataTypeArrayContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }
}
