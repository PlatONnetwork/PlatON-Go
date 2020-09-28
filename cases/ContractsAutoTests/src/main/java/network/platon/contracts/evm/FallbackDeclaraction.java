package network.platon.contracts.evm;

import com.alaya.abi.solidity.EventEncoder;
import com.alaya.abi.solidity.TypeReference;
import com.alaya.abi.solidity.datatypes.DynamicBytes;
import com.alaya.abi.solidity.datatypes.Event;
import com.alaya.abi.solidity.datatypes.Function;
import com.alaya.abi.solidity.datatypes.Type;
import com.alaya.abi.solidity.datatypes.generated.Uint256;
import com.alaya.crypto.Credentials;
import com.alaya.protocol.Web3j;
import com.alaya.protocol.core.DefaultBlockParameter;
import com.alaya.protocol.core.RemoteCall;
import com.alaya.protocol.core.methods.request.PlatonFilter;
import com.alaya.protocol.core.methods.response.Log;
import com.alaya.protocol.core.methods.response.TransactionReceipt;
import com.alaya.tx.Contract;
import com.alaya.tx.TransactionManager;
import com.alaya.tx.gas.GasProvider;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import rx.Observable;
import rx.functions.Func1;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://github.com/PlatONnetwork/client-sdk-java/releases">platon-web3j command line tools</a>,
 * or the com.alaya.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/PlatONnetwork/client-sdk-java/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.13.2.0.
 */
public class FallbackDeclaraction extends Contract {
    private static final String BINARY = "608060405234801561001057600080fd5b50610286806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630dbe671f146100ba57806342a78883146100d857806366b0bae014610106578063d46300fd14610124575b606f6000819055507f17c1956f6e992470102c5fc953bf560fda31fabee8737cf8e77bdde00eb5698d60003660405180806020018281038252848482818152602001925080828437600081840152601f19601f820116905080830192505050935050505060405180910390a1005b6100c2610142565b6040518082815260200191505060405180910390f35b610104600480360360208110156100ee57600080fd5b8101908080359060200190929190505050610148565b005b61010e6101b7565b6040518082815260200191505060405180910390f35b61012c610248565b6040518082815260200191505060405180910390f35b60005481565b7fb776d49293459725ca7d6a5abc60e389d2f3d067d4f028ba9cd790f6965998466000368360405180806020018381526020018281038252858582818152602001925080828437600081840152601f19601f82011690508083019250505094505050505060405180910390a150565b60003073ffffffffffffffffffffffffffffffffffffffff1660405180807f66756e6374696f6e4e6f744578697374282900000000000000000000000000008152506012019050600060405180830381855af49150503d8060008114610239576040519150601f19603f3d011682016040523d82523d6000602084013e61023e565b606091505b5050506001905090565b6000805490509056fea265627a7a7231582033afe4b14c98740de1cb4511cc535ad0a32e98063feb8e22347836e38b0296e164736f6c63430005110032";

    public static final String FUNC_A = "a";

    public static final String FUNC_CALLNONEXISTFUNC = "callNonExistFunc";

    public static final String FUNC_EXISTFUNC = "existFunc";

    public static final String FUNC_GETA = "getA";

    public static final Event EXISTFUNCCALLED_EVENT = new Event("ExistFuncCalled", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}, new TypeReference<Uint256>() {}));
    ;

    public static final Event FALLBACKCALLED_EVENT = new Event("FallbackCalled", 
            Arrays.<TypeReference<?>>asList(new TypeReference<DynamicBytes>() {}));
    ;

    protected FallbackDeclaraction(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider, chainId);
    }

    protected FallbackDeclaraction(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider, Long chainId) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider, chainId);
    }

    public List<ExistFuncCalledEventResponse> getExistFuncCalledEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(EXISTFUNCCALLED_EVENT, transactionReceipt);
        ArrayList<ExistFuncCalledEventResponse> responses = new ArrayList<ExistFuncCalledEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
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
                EventValuesWithLog eventValues = extractEventParametersWithLog(EXISTFUNCCALLED_EVENT, log);
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

    public List<FallbackCalledEventResponse> getFallbackCalledEvents(TransactionReceipt transactionReceipt) {
        List<EventValuesWithLog> valueList = extractEventParametersWithLog(FALLBACKCALLED_EVENT, transactionReceipt);
        ArrayList<FallbackCalledEventResponse> responses = new ArrayList<FallbackCalledEventResponse>(valueList.size());
        for (EventValuesWithLog eventValues : valueList) {
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
                EventValuesWithLog eventValues = extractEventParametersWithLog(FALLBACKCALLED_EVENT, log);
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

    public RemoteCall<BigInteger> a() {
        final Function function = new Function(FUNC_A, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
    }

    public RemoteCall<TransactionReceipt> callNonExistFunc() {
        final Function function = new Function(
                FUNC_CALLNONEXISTFUNC, 
                Arrays.<Type>asList(), 
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<TransactionReceipt> existFunc(BigInteger para) {
        final Function function = new Function(
                FUNC_EXISTFUNC, 
                Arrays.<Type>asList(new Uint256(para)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteCall<BigInteger> getA() {
        final Function function = new Function(FUNC_GETA, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Uint256>() {}));
        return executeRemoteCallSingleValueReturn(function, BigInteger.class);
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

    public static class ExistFuncCalledEventResponse {
        public Log log;

        public byte[] data;

        public BigInteger para;
    }

    public static class FallbackCalledEventResponse {
        public Log log;

        public byte[] data;
    }
}
