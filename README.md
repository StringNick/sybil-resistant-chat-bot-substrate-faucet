Sybil-resistant faucet is a generic chat-bot-based faucet solution that can be used on any existing parachain (substrate-based chain, either pallets or ink! smart contracts).

## Getting Started

1. Configure environment variables that are needed for chat bot authentication, faucet wallet, etc. Copy `config.example.toml` file into `config.toml` and start setting up the toml configuration. The rationale behind different variables setup can be found in [Variable](#variables) section. Each variable, along with its setup instructions, as well as default values, is described in [Configuration](#configuration) section.

2. Once all environment variables are ready, run the development server through docker with substrate default blockchain:
```
docker-compose -up -d
```

**Note**: all environment variables need to be set correctly in order for the faucet to work 🚨

## Discord

* Go to the Discord Developer Portal (https://discord.com/developers/applications) and log in with your Discord account.
* Click on "New Application" and give your bot a name.
* Navigate to the "Bot" tab on the left sidebar and click on "Add Bot." Confirm your action when prompted.
* Under the "TOKEN" section, click on "Copy" to copy the bot token. This token is essential for your bot to authenticate itself and interact with the Discord API.
* Now, you need to add the bot to your Discord server (channel). To do this, go to the "OAuth2" tab in the developer portal, scroll down to "OAuth2 URL Generator," select the "bot" scope, and then select the permissions you want your bot to have. Once you've chosen the permissions, copy the generated URL and open it in your web browser. From there, you can add the bot to one of your Discord servers.


## Matrix

* You need to create an account. 
* Add this account to the channel where you want to use it, with permissions to read/write.
* Provide credentials to [Configuration] (#configuration)

## Configuration

To make the faucet generic, many of its parts are configurable. Configuration settings are stored in `.env` files, one per each environment. Read more about environments and their setup in [environments](#environments) section.

| Environment                | Variable         | Description                                                                                                    | Default                                               |
| -------------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------- |
| `DRIP_CAP`                 | cap              | How many tokens to send per each claim.                                                                        | `0.025`                                               |
| `DRIP_DELAY`               | delay            | How often user's can request to drip tokens (in milliseconds).                                                 | `86400000`                                            |
| `DRIP_NETWORK_DECIMALS`    | network_decimals | Decimal places used for network tokens.                                                                        | `12`                                                  |
| `REDIS_ENDPOINT`           | endpoint         | Redis instance endpoint. It's easiest to setup Redis instance at Redis Cloud, or you may run a local instance. | `None`                                                |
| `SUBSTRATE_ENDPOINT`       | endpoint         | Substrate or Ink! based blockchain endpoint. Optionally, for testing purposes.                                 | `ws://substrate:9944` (for docker-compose addressing) |
| `SUBSTRATE_SEED_OR_PHRASE` | seed_or_phrase   | Mnemonic or secret seed of faucet's wallet from which funds will be drawn.                                     |                                                       |
| `DISCORD_ENABLED`          | enabled          | Enable or disable discord chat bot.                                                                            | `false`                                               |
| `DISCORD_TOKEN`            | token            | Discord bot api token.                                                                                         |                                                       |
| `MATRIX_ENABLED`           | enabled          | Enable or disable matrix chat bot.                                                                             | `false`                                               |
| `MATRIX_DEVICE_ID`         | device_id        | Device id option in matrix sdk, more u can read in matrix section.                                             |                                                       |
| `MATRIX_HOST`              | host             | Address of matrix webspace.                                                                                    |                                                       |
| `MATRIX_USERNAME`          | username         | Username(login) of an account, that will work as chat bot.                                                     |                                                       |
| `MATRIX_PASSWORD`          | password         | Password of an account, that will work as chat bot.                                                            |                                                       |


