import rlp
import json
from client_sdk_python import HTTPProvider, Web3, WebsocketProvider
from client_sdk_python.eth import Eth
from client_sdk_python.middleware import geth_poa_middleware
from hexbytes import HexBytes

# from common.key import mock_duplicate_sign
#
# result = mock_duplicate_sign(1,'a71f82ac2b15f6dbffa9f87de2d237d6bc7581d453bcaf973251b079ef9b20db','15338821af5128c21696180b7ea047905803a519c43153880b7da0c210563254',35531)
# print(result)
#
3000
interval = int((10000 / 10) / 1000)

consensus_wheel = (5 * 60) // (interval * 10 * 7)
print(consensus_wheel)
print(consensus_wheel * 70)
annual_cycle = (28 * 60) // 280
print(annual_cycle)