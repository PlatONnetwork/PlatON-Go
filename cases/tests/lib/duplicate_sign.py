from environment.env import TestEnvironment
from environment.node import Node
from .economic import Economic


class DuplicateSign:
    def __init__(self, env: TestEnvironment, node: Node):
        self.node = node
        self.economic = Economic(env)

    def reportDuplicateSign(self, typ, data, from_address, transaction_cfg=None):
        """
        Report duplicate sign
        :param typ: Represents duplicate sign type, 1:prepareBlock, 2: prepareVote, 3:viewChange
        :param data: Json value of single evidence, format reference RPC interface Evidences
        :param from_address: address for transaction
        :param transaction_cfg:
        :return:
        """
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.reportDuplicateSign(typ, data, pri_key, transaction_cfg)

    @property
    def ppos(self):
        return self.node.ppos
