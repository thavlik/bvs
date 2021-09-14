import os
from blockfrost import BlockFrostApi, ApiError, ApiUrls

with open(os.path.expanduser('~/.blockfrost/key')) as f:
    project_id = f.read().strip()

print(f'Using project_id "{project_id}"')

api = BlockFrostApi(
    project_id=project_id,
    # optional: pass base_url or export BLOCKFROST_API_URL to use testnet, defaults to ApiUrls.mainnet.value
    #base_url=ApiUrls.testnet.value,
)


def get_tx_info(api, txid):
    pass
    
try:
    health = api.health()
    print(health.is_healthy)  # True

    txid = 'd628f645544bc68739a1f589f7d3a4d961b62ab1f74db70eb5cc05ddf5ece869'
    get_tx_info(api, txid)    

    account_rewards = api.account_rewards(
        stake_address='stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7',
        count=20,
    )
    print(account_rewards[0].epoch)  # prints 221
    print(len(account_rewards))  # prints 20

    account_rewards = api.account_rewards(
        stake_address='stake1ux3g2c9dx2nhhehyrezyxpkstartcqmu9hk63qgfkccw5rqttygt7',
        count=20,
        gather_pages=True,
    )
    print(account_rewards[0].epoch)  # prints 221
    print(len(account_rewards))  # prints 57

    address = api.address(
        address='addr1qxqs59lphg8g6qndelq8xwqn60ag3aeyfcp33c2kdp46a09re5df3pzwwmyq946axfcejy5n4x0y99wqpgtp2gd0k09qsgy6pz')
    print(address.type)  # prints 'shelley'
    for amount in address.amount:
        print(amount.unit)  # prints 'lovelace'

except ApiError as e:
    print(e)