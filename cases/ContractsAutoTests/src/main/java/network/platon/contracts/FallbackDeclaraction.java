package network.platon.contracts;

import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import org.web3j.abi.EventEncoder;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.DynamicBytes;
import org.web3j.abi.datatypes.Event;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.generated.Uint256;
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
public class FallbackDeclaraction extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b506102d4806100206000396000f3fe608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630dbe671f146100dd57806342a788831461010857806366b0bae014610143578063d46300fd1461016e575b34801561006e57600080fd5b50606f6000819055507f17c1956f6e992470102c5fc953bf560fda31fabee8737cf8e77bdde00eb5698d60003660405180806020018281038252848482818152602001925080828437600081840152601f19601f820116905080830192505050935050505060405180910390a1005b3480156100e957600080fd5b506100f2610199565b6040518082815260200191505060405180910390f35b34801561011457600080fd5b506101416004803603602081101561012b57600080fd5b810190808035906020019092919050505061019f565b005b34801561014f57600080fd5b5061015861020e565b6040518082815260200191505060405180910390f35b34801561017a57600080fd5b5061018361029f565b6040518082815260200191505060405180910390f35b60005481565b7fb776d49293459725ca7d6a5abc60e389d2f3d067d4f028ba9cd790f6965998466000368360405180806020018381526020018281038252858582818152602001925080828437600081840152601f19601f82011690508083019250505094505050505060405180910390a150565b60003073ffffffffffffffffffffffffffffffffffffffff1660405180807f66756e6374696f6e4e6f744578697374282900000000000000000000000000008152506012019050600060405180830381855af49150503d8060008114610290576040519150601f19603f3d011682016040523d82523d6000602084013e610295565b606091505b5050506001905090565b6000805490509056fea165627a7a723058209f81750ddd4501de806fb9df6e19a11c41119c4412c9221167c4efd2f4093ef70029";

    public static final String FUNC_A = "a";

    public static final String FUNC_EXISTFUNC = "existFunc";

    public static final String FUNC_CALLNONEXISTFUNC = "callNonExistFunc";

    public static final String FUNC_GETA = "getA";

    public static final Event FALLBACKCALLED_EVENT = new Event("FallbackCalled", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
    ;

    public static final Event EXISTFUNCCALLED_EVENT = new Event("ExistFuncCalled", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}, new TypeReference<Uint256>() {}));
    ;

    protected FallbackDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected FallbackDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> existFunc(BigInteger para) {
        final Function function = new Function(
                FUNC_EXISTFUNC, 
                Arrays.<Type>asList(new org.web3j.abi.datatypes.generated.Uint256(para)), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> callNonExistFunc() {
        final Function function = new Function(
                FUNC_CALLNONEXISTFUNC, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getA() {
        final Function function = new Function(FUNC_GETA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public List<FallbackCalledEventResponse> getFallbackCalledEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(FALLBACKCALLED_EVENT, transactionReceipt);
        ArrayList<FallbackCalledEventResponse> responses = new ArrayList<FallbackCalledEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            FallbackCalledEventResponse typedResponse = new FallbackCalledEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.data = (byte[]) eventValues.getNonIndexedValues().get(0).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<FallbackCalledEventResponse> fallbackCalledEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, FallbackCalledEventResponse>() {
            @Override
            public FallbackCalledEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(FALLBACKCALLED_EVENT, log);
                FallbackCalledEventResponse typedResponse = new FallbackCalledEventResponse();
                typedResponse.log = log;
                typedResponse.data = (byte[]) eventValues.getNonIndexedValues().get(0).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<FallbackCalledEventResponse> fallbackCalledEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(FALLBACKCALLED_EVENT));
        return fallbackCalledEventObservable(filter);
    }

    public List<ExistFuncCalledEventResponse> getExistFuncCalledEvents(TransactionReceipt transactionReceipt) {
        List<Contract.EventValuesWithLog> valueList = extractEventParametersWithLog(EXISTFUNCCALLED_EVENT, transactionReceipt);
        ArrayList<ExistFuncCalledEventResponse> responses = new ArrayList<ExistFuncCalledEventResponse>(valueList.size());
        for (Contract.EventValuesWithLog eventValues : valueList) {
            ExistFuncCalledEventResponse typedResponse = new ExistFuncCalledEventResponse();
            typedResponse.log = eventValues.getLog();
            typedResponse.data = (byte[]) eventValues.getNonIndexedValues().get(0).getValue();
            typedResponse.para = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
            responses.add(typedResponse);
        }
        return responses;
    }

    public Observable<ExistFuncCalledEventResponse> existFuncCalledEventObservable(PlatonFilter filter) {
        return web3j.platonLogObservable(filter).map(new Func1<Log, ExistFuncCalledEventResponse>() {
            @Override
            public ExistFuncCalledEventResponse call(Log log) {
                Contract.EventValuesWithLog eventValues = extractEventParametersWithLog(EXISTFUNCCALLED_EVENT, log);
                ExistFuncCalledEventResponse typedResponse = new ExistFuncCalledEventResponse();
                typedResponse.log = log;
                typedResponse.data = (byte[]) eventValues.getNonIndexedValues().get(0).getValue();
                typedResponse.para = (BigInteger) eventValues.getNonIndexedValues().get(1).getValue();
                return typedResponse;
            }
        });
    }

    public Observable<ExistFuncCalledEventResponse> existFuncCalledEventObservable(DefaultBlockParameter startBlock, DefaultBlockParameter endBlock) {
        PlatonFilter filter = new PlatonFilter(startBlock, endBlock, getContractAddress());
        filter.addSingleTopic(EventEncoder.encode(EXISTFUNCCALLED_EVENT));
        return existFuncCalledEventObservable(filter);
    }

    public static RemoteCall<FallbackDeclaraction> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FallbackDeclaraction.class, web3j, credentials, contractGasProvider, BINARY,  "", chainId);
    }

    public static RemoteCall<FallbackDeclaraction> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return deployRemoteCall(FallbackDeclaraction.class, web3j, transactionManager, contractGasProvider, BINARY,  "", chainId);
    }

    public static FallbackDeclaraction load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        return new FallbackDeclaraction(contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    public static FallbackDeclaraction load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        return new FallbackDeclaraction(contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public static class FallbackCalledEventResponse {
        public Log log;

        public byte[] data;
    }

    public static class ExistFuncCalledEventResponse {
        public Log log;

        public byte[] data;

        public BigInteger para;
    }
}
