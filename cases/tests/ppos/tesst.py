import math
from web3 import Web3
from common.connect import connect_web3
from common.key import mock_duplicate_sign

w3 = connect_web3('192.168.112.172:9789')
info = mock_duplicate_sign(1, 'a71f82ac2b15f6dbffa9f87de2d237d6bc7581d453bcaf973251b079ef9b20db', '15338821af5128c21696180b7ea047905803a519c43153880b7da0c210563254', 50850)
print(info)