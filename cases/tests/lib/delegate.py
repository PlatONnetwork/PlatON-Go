from environment.env import TestEnvironment
from environment.node import Node
from .economic import Economic


class Delegate:
    """
    Used to initiate a delegate transaction,
    if you need to use the call method, please call ppos
    example:
    >>>delegate=Delegate(env, node)
    >>>delegate.ppos.getDelegateInfo(...)
    """

    def __init__(self, env: TestEnvironment, node: Node):
        self.node = node
        self.economic = Economic(env)

    @property
    def ppos(self):
        """
        use sdk ppos object
        :return:
        """
        return self.node.ppos

    def delegate(self, typ, from_address, node_id=None, amount=None, tansaction_cfg=None):
        """
        Initiate delegate
        :param typ: Amount type
        :param from_address: Initiating a delegate account
        :param node_id: The id of the delegate node
        :param amount: delegate amount
        :param tansaction_cfg:
        :return:
        """
        if node_id is None:
            node_id = self.node.node_id
        if amount is None:
            amount = self.economic.delegate_limit
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.delegate(typ, node_id, amount, pri_key, transaction_cfg=tansaction_cfg)

    def withdrew_delegate(self, staking_blocknum, from_address, node_id=None, amount=None, transaction_cfg=None):
        """
        Release delegate
        :param staking_blocknum: staking block height
        :param from_address: Initiating a delegate account
        :param node_id: The id of the delegate node
        :param amount: Release delegate amount
        :param transaction_cfg:
        :return:
        """
        if node_id is None:
            node_id = self.node.node_id
        if amount is None:
            amount = self.economic.delegate_limit
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.withdrewDelegate(staking_blocknum, node_id, amount, pri_key, transaction_cfg)

    def withdraw_delegate_reward(self, from_address, transaction_cfg=None):
        pri_key = self.economic.account.find_pri_key(from_address)
        return self.ppos.withdrawDelegateReward(pri_key, transaction_cfg)

    def get_delegate_reward_by_nodeid(self, from_address, node_ids=[]):
        reward = 0
        data = self.ppos.getDelegateReward(from_address, node_ids)
        if len(node_ids) == 1:
            return data.get('Ret')[0].get('reward')

        elif len(node_ids) == 0:
            for rewardinfo in data.get('Ret'):
                reward += rewardinfo.get('reward')
        else:
            for rewardinfo in data.get('Ret'):
                if rewardinfo.get('nodeID') in node_ids:
                    reward += rewardinfo.get('reward')
        return reward

    def get_staking_num_by_nodeid(self, from_address, node_ids=None):
        if node_ids is None:
            node_ids = []
        else:
            node_ids = [node_ids]
        data = self.ppos.getDelegateReward(from_address, node_ids)
        if data.get('Code') != 0:
            raise Exception('address {} no reward information'.format(from_address))
        for rewardinfo in data.get('Ret'):
            if node_ids[0] == rewardinfo.get('nodeID'):
                return rewardinfo.get('stakingNum')
        raise Exception('Address {} no reward inforamtion about nodeid {}'.format(
            from_address, node_ids[0]
        ))



