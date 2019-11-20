from environment.env import TestEnvironment
from environment.node import Node
from .config import StakingConfig
from .economic import Economic
from .restricting import Restricting
from .delegate import Delegate
from .staking import Staking
from .pip import Pip
from .duplicate_sign import DuplicateSign
from typing import List


class Client:
    """
    Test client, the interface call method required for the collection test
    example:
    >>>cfg = StakingConfig("1111","test-node","http://test-node.com","I'm tester")
    >>>client = Client(env, node, cfg)
    >>>client.staking.create_staking(...)
    """

    def __init__(self, env: TestEnvironment, node: Node, cfg: StakingConfig):
        self.node = node
        self.economic = Economic(env)
        self.staking = Staking(env, node, cfg)
        self.restricting = Restricting(env, node)
        self.delegate = Delegate(env, node)
        self.pip = Pip(env, node)
        self.duplicatesign = DuplicateSign(env, node)

    @property
    def ppos(self):
        """
        use sdk ppos object
        :return:
        """
        return self.node.ppos


def get_client_obj(nodeid, client_obj_list: List[Client]) -> Client:
    """
    Get the client object according to the node id
    :param nodeid:
    :param client_obj_list:
    :return:
    """
    for client_obj in client_obj_list:
        if nodeid == client_obj.node.node_id:
            return client_obj


def get_client_obj_list(nodeid_list, client_obj_list: List[Client]) -> List[Client]:
    """
    Get the client object list according to the node id list
    :param nodeid_list:
    :param client_obj_list:
    :return:
    """
    new_client_obj_list = []
    for nodeid in nodeid_list:
        new_client_obj_list.append(get_client_obj(nodeid, client_obj_list))
    return new_client_obj_list
