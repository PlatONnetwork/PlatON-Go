package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicArray;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.abi.datatypes.generated.Int256;
import org.web3j.abi.datatypes.generated.StaticArray2;
import org.web3j.abi.datatypes.generated.Uint256;
import org.web3j.abi.datatypes.generated.Uint8;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.DefaultBlockParameter;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.request.PlatonFilter;
import org.web3j.protocol.core.methods.response.Log;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.GasProvider;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.0.7.
 */
public class EventIndexedContract extends Contract {
    private static final String BINARY = "60806040526040518060a00160405280600060ff168152602001600160ff168152602001600260ff168152602001600360ff168152602001600460ff16815250600090600561004f929190610177565b506040518060c001604052806040518060400160405280600060ff168152602001600060ff1681525081526020016040518060400160405280600060ff168152602001600160ff1681525081526020016040518060400160405280600060ff168152602001600260ff1681525081526020016040518060400160405280600060ff168152602001600360ff1681525081526020016040518060400160405280600060ff168152602001600460ff1681525081526020016040518060400160405280600060ff168152602001600560ff1681525081525060019060066101359291906101c9565b506040518060800160405280605881526020016107766058913960029080519060200190610164929190610227565b5034801561017157600080fd5b5061034b565b8280548282559060005260206000209081019282156101b8579160200282015b828111156101b7578251829060ff16905591602001919060010190610197565b5b5090506101c591906102a7565b5090565b828054828255906000526020600020906002028101928215610216579160200282015b82811115610215578251829060026102059291906102cc565b50916020019190600201906101ec565b5b5090506102239190610311565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061026857805160ff1916838001178555610296565b82800160010185558215610296579182015b8281111561029557825182559160200191906001019061027a565b5b5090506102a391906102a7565b5090565b6102c991905b808211156102c55760008160009055506001016102ad565b5090565b90565b8260028101928215610300579160200282015b828111156102ff578251829060ff169055916020019190600101906102df565b5b50905061030d91906102a7565b5090565b61033a91905b80821115610336576000818161032d919061033d565b50600201610317565b5090565b90565b506000815560010160009055565b61041c8061035a6000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063bb0a5bd91161005b578063bb0a5bd9146100a0578063bbd847af146100aa578063db06ee8d146100b4578063fcf3fbb8146100e25761007d565b80633ac559931461008257806343ae41d81461008c578063b05dfcf814610096575b600080fd5b61008a6100ec565b005b610094610184565b005b61009e6101bf565b005b6100a861022e565b005b6100b2610243565b005b6100e0600480360360208110156100ca57600080fd5b81019080803590602001909291905050506102d1565b005b6100ea610301565b005b600260405180828054600181600116156101000203166002900480156101495780601f10610127576101008083540402835291820191610149565b820191906000526020600020905b815481529060010190602001808311610135575b505091505060405180910390207f617cf8a4400dd7963ed519ebe655a16e8da1282bb8fea36a21f634af912f54ab60405160405180910390a2565b6000600381111561019157fe5b7fde7a62815e0b38238b6211179d7d98017a99227a90823b0f44227e81dd3ad9c260405160405180910390a2565b60006040518082805480156101f357602002820191906000526020600020905b8154815260200190600101908083116101df575b505091505060405180910390207f38a323fa24260bbb8b86f61cd1d8c1900024088af6d08eda9e2d793da33c1b5860405160405180910390a2565b600460036002600160405160405180910390a4565b600160405180828054801561029657602002820191906000526020600020905b8160028015610287576020028201915b815481526020019060010190808311610273575b50509060020190808311610263575b505091505060405180910390207f406715adbc90cbc793dcd5707190ad1390229b2a75cf5b5ca228b518ae52de9a60405160405180910390a2565b807f834f28319998410e0230b46f29e58f7a86d825bc8505a917dae170f66d70bcdd60405160405180910390a250565b6002604051808280546001816001161561010002031660029004801561035e5780601f1061033c57610100808354040283529182019161035e565b820191906000526020600020905b81548152906001019060200180831161034a575b505091505060405180910390206000600381111561037857fe5b60006040518082805480156103ac57602002820191906000526020600020905b815481526020019060010190808311610398575b505091505060405180910390207f36b0026af468b8a488c2bd2a23a731b236d139a5211611fb763844bcb87f1abb60405160405180910390a456fea265627a7a7231582043ff9cc8771ba30c75621b31f6baef8e1910036f303b0bd2b1f77791d46ebada64736f6c634300050d003231323334353637383930303937383635343332313132333435363738393030393837363534333231313233343536373839303039373634333534363636363633323432343434343434343434343735383331353436383536";

    public static final String FUNC_TESTANONYMOUSINDEXED = "testAnonymousIndexed";

    public static final String FUNC_TESTCOMPLEX = "testComplex";

    public static final String FUNC_TESTENUM = "testEnum";

    public static final String FUNC_TESTMINUS = "testMinus";

    public static final String FUNC_TESTONEDIMENSIONALARRAY = "testOneDimensionalArray";

    public static final String FUNC_TESTSTR = "testStr";

    public static final String FUNC_TESTTWODIMENSIONALARRAY = "testTwoDimensionalArray";

    public static final Event ANONYMOUSINDEXEDNUM_EVENT = new Event("AnonymousIndexedNum", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>(true) {}, new TypeReference<Uint256>(true) {}, new TypeReference<Uint256>(true) {}, new TypeReference<Uint256>(true) {}));
    ;

    public static final Event COMPLEXINDEXEDEVENT_EVENT = new Event("ComplexIndexedEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Uint256>>(true) {}, new TypeReference<Uint8>(true) {}, new TypeReference<Utf8String>(true) {}));
    ;

    public static final Event ENUMEVENT_EVENT = new Event("EnumEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Uint8>(true) {}));
    ;

    public static final Event MINUSEVENT_EVENT = new Event("MinusEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Int256>(true) {}));
    ;

    public static final Event ONEDIMENSIONALARRAYEVENT_EVENT = new Event("OneDimensionalArrayEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<Uint256>>(true) {}));
    ;

    public static final Event STRINGEVENT_EVENT = new Event("StringEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>(true) {}));
    ;

    public static final Event TWODIMENSIONALARRAYEVENT_EVENT = new Event("TwoDimensionalArrayEvent", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicArray<StaticArray2<Uint256>>>(true) {}));
    ;

    protected EventIndexedContract(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected EventIndexedContract(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public List<AnonymousIndexedNumEventResponse> getAnonymousIndexedNumEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(ANONYMOUSINDEXEDNUM_EVENT, transactionReceipt);
        ArrayList<AnonymousIndexedNumEventResponse> responses = new ArrayList<AnonymousIndexedNumEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            AnonymousIndexedNumEventResponse typedResponse = new AnonymousIndexedNumEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.u1 = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
            typedResponse.u2 = (BigInteger) eventValues.getIndexedValues().get(1).getValue();
            typedResponse.u3 = (BigInteger) eventValues.getIndexedValues().get(2).getValue();
            typedResponse.u4 = (BigInteger) eventValues.getIndexedValues().get(3).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<AnonymousIndexedNumEventResponse> anonymousIndexedNumEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, AnonymousIndexedNumEventResponse>() {
            @Override
            public AnonymousIndexedNumEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(ANONYMOUSINDEXEDNUM_EVENT, log);
                AnonymousIndexedNumEventResponse typedResponse = new AnonymousIndexedNumEventResponse();
                typedResponse.log = log;
                typedResponse.u1 = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
                typedResponse.u2 = (BigInteger) eventValues.getIndexedValues().get(1).getValue();
                typedResponse.u3 = (BigInteger) eventValues.getIndexedValues().get(2).getValue();
                typedResponse.u4 = (BigInteger) eventValues.getIndexedValues().get(3).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<AnonymousIndexedNumEventResponse> anonymousIndexedNumEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ANONYMOUSINDEXEDNUM_EVENT));
        return anonymousIndexedNumEventObservable(filter);
    }

    public List<ComplexIndexedEventEventResponse> getComplexIndexedEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(COMPLEXINDEXEDEVENT_EVENT, transactionReceipt);
        ArrayList<ComplexIndexedEventEventResponse> responses = new ArrayList<ComplexIndexedEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            ComplexIndexedEventEventResponse typedResponse = new ComplexIndexedEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.array = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            typedResponse.choice = (BigInteger) eventValues.getIndexedValues().get(1).getValue();
            typedResponse.str = (byte[]) eventValues.getIndexedValues().get(2).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<ComplexIndexedEventEventResponse> complexIndexedEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, ComplexIndexedEventEventResponse>() {
            @Override
            public ComplexIndexedEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(COMPLEXINDEXEDEVENT_EVENT, log);
                ComplexIndexedEventEventResponse typedResponse = new ComplexIndexedEventEventResponse();
                typedResponse.log = log;
                typedResponse.array = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                typedResponse.choice = (BigInteger) eventValues.getIndexedValues().get(1).getValue();
                typedResponse.str = (byte[]) eventValues.getIndexedValues().get(2).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<ComplexIndexedEventEventResponse> complexIndexedEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(COMPLEXINDEXEDEVENT_EVENT));
        return complexIndexedEventEventObservable(filter);
    }

    public List<EnumEventEventResponse> getEnumEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(ENUMEVENT_EVENT, transactionReceipt);
        ArrayList<EnumEventEventResponse> responses = new ArrayList<EnumEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            EnumEventEventResponse typedResponse = new EnumEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.choices = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<EnumEventEventResponse> enumEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, EnumEventEventResponse>() {
            @Override
            public EnumEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(ENUMEVENT_EVENT, log);
                EnumEventEventResponse typedResponse = new EnumEventEventResponse();
                typedResponse.log = log;
                typedResponse.choices = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<EnumEventEventResponse> enumEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ENUMEVENT_EVENT));
        return enumEventEventObservable(filter);
    }

    public List<MinusEventEventResponse> getMinusEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(MINUSEVENT_EVENT, transactionReceipt);
        ArrayList<MinusEventEventResponse> responses = new ArrayList<MinusEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            MinusEventEventResponse typedResponse = new MinusEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.minus = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<MinusEventEventResponse> minusEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, MinusEventEventResponse>() {
            @Override
            public MinusEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(MINUSEVENT_EVENT, log);
                MinusEventEventResponse typedResponse = new MinusEventEventResponse();
                typedResponse.log = log;
                typedResponse.minus = (BigInteger) eventValues.getIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<MinusEventEventResponse> minusEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(MINUSEVENT_EVENT));
        return minusEventEventObservable(filter);
    }

    public List<OneDimensionalArrayEventEventResponse> getOneDimensionalArrayEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(ONEDIMENSIONALARRAYEVENT_EVENT, transactionReceipt);
        ArrayList<OneDimensionalArrayEventEventResponse> responses = new ArrayList<OneDimensionalArrayEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            OneDimensionalArrayEventEventResponse typedResponse = new OneDimensionalArrayEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.array = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<OneDimensionalArrayEventEventResponse> oneDimensionalArrayEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, OneDimensionalArrayEventEventResponse>() {
            @Override
            public OneDimensionalArrayEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(ONEDIMENSIONALARRAYEVENT_EVENT, log);
                OneDimensionalArrayEventEventResponse typedResponse = new OneDimensionalArrayEventEventResponse();
                typedResponse.log = log;
                typedResponse.array = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<OneDimensionalArrayEventEventResponse> oneDimensionalArrayEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(ONEDIMENSIONALARRAYEVENT_EVENT));
        return oneDimensionalArrayEventEventObservable(filter);
    }

    public List<StringEventEventResponse> getStringEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(STRINGEVENT_EVENT, transactionReceipt);
        ArrayList<StringEventEventResponse> responses = new ArrayList<StringEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            StringEventEventResponse typedResponse = new StringEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.str = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<StringEventEventResponse> stringEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, StringEventEventResponse>() {
            @Override
            public StringEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(STRINGEVENT_EVENT, log);
                StringEventEventResponse typedResponse = new StringEventEventResponse();
                typedResponse.log = log;
                typedResponse.str = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<StringEventEventResponse> stringEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(STRINGEVENT_EVENT));
        return stringEventEventObservable(filter);
    }

    public List<TwoDimensionalArrayEventEventResponse> getTwoDimensionalArrayEventEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(TWODIMENSIONALARRAYEVENT_EVENT, transactionReceipt);
        ArrayList<TwoDimensionalArrayEventEventResponse> responses = new ArrayList<TwoDimensionalArrayEventEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            TwoDimensionalArrayEventEventResponse typedResponse = new TwoDimensionalArrayEventEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.array = (byte[]) eventValues.getIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<TwoDimensionalArrayEventEventResponse> twoDimensionalArrayEventEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, TwoDimensionalArrayEventEventResponse>() {
            @Override
            public TwoDimensionalArrayEventEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(TWODIMENSIONALARRAYEVENT_EVENT, log);
                TwoDimensionalArrayEventEventResponse typedResponse = new TwoDimensionalArrayEventEventResponse();
                typedResponse.log = log;
                typedResponse.array = (byte[]) eventValues.getIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<TwoDimensionalArrayEventEventResponse> twoDimensionalArrayEventEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(TWODIMENSIONALARRAYEVENT_EVENT));
        return twoDimensionalArrayEventEventObservable(filter);
    }

    public RemoteCall<TransactionReceipt> testAnonymousIndexed() {
        final Function function = new Function(
                FUNC_TESTANONYMOUSINDEXED, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testComplex() {
        final Function function = new Function(
                FUNC_TESTCOMPLEX, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testEnum() {
        final Function function = new Function(
                FUNC_TESTENUM, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testMinus(BigInteger value) {
        final Function function = new Function(
                FUNC_TESTMINUS, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Int256(value)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testOneDimensionalArray() {
        final Function function = new Function(
                FUNC_TESTONEDIMENSIONALARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testStr() {
        final Function function = new Function(
                FUNC_TESTSTR, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> testTwoDimensionalArray() {
        final Function function = new Function(
                FUNC_TESTTWODIMENSIONALARRAY, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public static RemoteCall<EventIndexedContract> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(EventIndexedContract.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<EventIndexedContract> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(EventIndexedContract.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static EventIndexedContract load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new EventIndexedContract(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static EventIndexedContract load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new EventIndexedContract(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class AnonymousIndexedNumEventResponse {
        public Log log;

        public BigInteger u1;

        public BigInteger u2;

        public BigInteger u3;

        public BigInteger u4;
    }

    public static class ComplexIndexedEventEventResponse {
        public Log log;

        public byte[] array;

        public BigInteger choice;

        public byte[] str;
    }

    public static class EnumEventEventResponse {
        public Log log;

        public BigInteger choices;
    }

    public static class MinusEventEventResponse {
        public Log log;

        public BigInteger minus;
    }

    public static class OneDimensionalArrayEventEventResponse {
        public Log log;

        public byte[] array;
    }

    public static class StringEventEventResponse {
        public Log log;

        public byte[] str;
    }

    public static class TwoDimensionalArrayEventEventResponse {
        public Log log;

        public byte[] array;
    }
}
