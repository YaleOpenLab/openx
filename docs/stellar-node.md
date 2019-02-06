# Setting up a Stellar node and a Horizon server

The docs for setting up a stellar node are vague and spread out. Either this process is so easy that nobody seems to have needed help or no one's tried to setup a validator on their own. In either case, this document serves as a guide for people who wish to run a stellar full node.

Get a server: Stellar's full blockchain size is approximately 20GB, so a standard 100GB server should be good to start off with. *ONLY RUN UBUNTU 16.04* since this is the only version stellar tests its releases on. The author can confirm that this does not work on debian, so it would be safest to run Ubuntu on the server instance.

Run [stellar-quickstart](https://github.com/stellar/packages#debug-symbols):

```
wget -qO - https://apt.stellar.org/SDF.asc | sudo apt-key add -
echo "deb https://apt.stellar.org/public stable/" | sudo tee -a /etc/apt/sources.list.d/SDF.list
sudo apt-get update && sudo apt-get install stellar-quickstart
```

This would create a new user `stellar` on your machine and you should be able to run the installed packages as the user. A natural step would be to do `sudo su stellar` at this point to test out stellar-core and other packages.

Verify installation by running `stellar-core -h` to see that it was properly installed.

To get attach a screen with the terminal (to run `stellar-core` in), run `script /dev/null`.

If you try to run horizon at this point, you should notice that psql spits out some nasty errors on its own. If you do want to reinstall postgresql, run

```
deb http://apt.postgresql.org/pub/repos/apt/ YOUR_DEBIAN_VERSION_HERE-pgdg main
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
```

to ensure that your copy of postgresql is not corrupted by some random influence. if you have existing data on postgresql, this step is to be proceeded ahead with caution. Another alternative is to use sqlite3 in place of postgresql and change config in the appropriate places. This guide assumes that we will be using postgresql since its more suited for production environments than sqlite3 even though the latter has better performance. In addition, under note from the stellar dev team, if one plans to run just a validator, sqlite3 should function just fine.

Open psql: `psql` and create a new database to hold the horizon database (horizon is the query interface for stellar)

```
> psql create database horizon
> psql grant all privileges on database mydb to stellar;
```

After this step, the database would be created, but you need to reset your password for this user "stellar" on postgresql. To do that, run

```
> psql ALTER USER stellar PASSWORD 'password';
```

Restart postgresql for the changes to take effect:

```
/etc/init.d/postgresql restart
```

Before re-starting stellar-core with out desired changes (stellar-core would already be running as a basic validator in the background once `stellar-quickstart` returns true), it helps to understand the different roles that stellar has with respect to the concept of "full nodes":

1. Basic Validator - Somewhat similar to a light client in ethereum, one can configure the number of blocks beyond which the client does not verify the blocks given to it. For eg, say the current blockheight is 2001156 and `CATCHUP_RECENT` is set to `1024`. This would mean that the client only downloads the last 1024 blocks and does not validate the rest. This is the default configuration under which the node would be running once `quickstart` returns true.

2. Full Validator - This is similar to a full node in bitcoin which stores the past tx history of all blocks and proceeds to verify them. Sometimes it may be difficult to sync up a full node depending on where you are since stellar full nodes are pretty less in number and more centralized location wise. Setting `NODE_IS_VALIDATOR` and `CATCHUP_COMPLETE` to `true` runs a full node.

3. Watchers - Watches activity on the blockchain but does not take part in consensus

4. Archiver nodes - Records and stores activity on the blockchain, but doesn't take part in consensus

An important component of running `stellar-core` is making sure that you have the right config file. The following config file can be used to run a VALIDATING node on TESTNET. Mainnet configuration is similar except for validators and this section will be updated once the platform transitions to it.

## Config

```
# The port which stellar-core uses
HTTP_PORT=11626
# Don't set this to true except for cases when you're on localhost
PUBLIC_HTTP_PORT=false
# Default value, don't change
NETWORK_PASSPHRASE="Test SDF Network ; September 2015"
# There are only three validators sun by the Stellar Foundation, have them here for testnet.
KNOWN_PEERS=[
"core-testnet1.stellar.org",
"core-testnet2.stellar.org",
"core-testnet3.stellar.org"]
# Paths which you may want to customize
LOG_FILE_PATH="/var/log/stellar/stellar-core.log"
BUCKET_DIR_PATH="/var/lib/stellar/buckets"
DATABASE="postgresql://dbname=stellar user=stellar"
# Set this to true ONLY for VALIDATING nodes
NODE_IS_VALIDATOR=true
CATCHUP_COMPLETE=true
# Set this option to whatever value you would like to go back in history and download blocks from if you're running a BASIC validator
#CATCHUP_RECENT=1549278755
# Accept unsafe quorums
UNSAFE_QUORUM=true
FAILURE_SAFETY=1
# Set the threshold for accepting a block proposed by a specific quorum
[QUORUM_SET]
THRESHOLD_PERCENT=51 # rounded up -> 2 nodes out of 3
VALIDATORS=[
"GDKXE2OZMJIPOSLNA6N6F2BVCI3O777I2OOC4BV7VOYUEHYX7RTRYA7Y  sdf1",
"GCUCJTIYXSOXKBSNFGNFWW5MUQ54HKRPGJUTQFJ5RQXZXNOLNXYDHRAP  sdf2",
"GC2V2EFSXN6SQTWVYA5EPJPBWWIMSD2XQNKUOHGEKB535AQE2I6IXV2Z  sdf3"]
# History of the stellar testnet, dunno why this is here and didn't remove it, so maybe one could try removing it and seeing if it works.
[HISTORY.h1]
get="curl -sf http://s3-eu-west-1.amazonaws.com/history.stellar.org/prd/core-testnet/core_testnet_001/{0} -o {1}"
[HISTORY.h2]
get="curl -sf http://s3-eu-west-1.amazonaws.com/history.stellar.org/prd/core-testnet/core_testnet_002/{0} -o {1}"
[HISTORY.h3]
get="curl -sf http://s3-eu-west-1.amazonaws.com/history.stellar.org/prd/core-testnet/core_testnet_003/{0} -o {1}"
```

Store the config in a file named `stellar-core.cfg` for stellar-core to load the config from. Now stellar-core should be ready to start. ` screen -S stellar stellar-core run` starts a `screen` with name `stellar` that we can attach to later if required. `ctrl+a+d` to detach from the screen.

Now wait for `stellar-core` to finish syncing. In the meantime, we can setup `stellar-horizon` to be ready to run once `stellar-core` is done syncing. First, we need to load the schema unto the database:

```
stellar-horizon --stellar-core-db-url="postgresql://stellar:password@localhost/stellar" --stellar-core-url="http://localhost:11626" --db-url="postgresql://stellar:password@localhost/horizon" db init
```

(replace relevant parameters with your credentials)
After initialising the schema, we need to wait for `stellar-core` to finish syncing. To check on this, run `curl localhost:11626/info` and check whether the "status" field returns "Synced!". If it does, we are ready to start horizon, which should be as simple as running:

```
stellar-horizon --stellar-core-db-url="postgresql://stellar:password@localhost/stellar" --stellar-core-url="http://localhost:11626" --db-url="postgresql://stellar:password@localhost/horizon" --port 8080 db backfill 2060000
```

Replace 2060000 with a number that is higher than the current latest block number. This is done to index blocks that have already been downloaded by `stellar-core`. Once this is done with, run:

```
screen -S horizon stellar-horizon --stellar-core-db-url="postgresql://stellar:password@localhost/stellar" --stellar-core-url="http://localhost:11626" --db-url="postgresql://stellar:password@localhost/horizon" --port 8080 --ingest=true
```

to start horizon on port 8080 (change if required). The `ingest` field exists to provide a constant stream of data from `stellar-core` to `stellar-horizon` for continuous indexing. One can now query the given horizon server like the default horizon server run by the stellar foundation. Please do note that you'd have to open the relevant ports in the firewall in order for one to be able to receive / send traffic. The main port that stellar uses is 11625, 11626 is local to the environment that stellar-core is being run in and the `stellar-horizon` port can be configured to run on any port that we want it to.

The stellar node takes ~20 hours to completely sync testnet from scratch and takes about 20GB of space on an Ubuntu 16.04LTS machine.

This document will be updated in the future once the platform transitions to mainnet from testnet.

## References

- https://github.com/stellar/docker-stellar-core-horizon/blob/master/testnet/core/etc/stellar-core.cfg
- https://github.com/stellar/docs/blob/master/other/stellar-core-validator-example.cfg
- https://www.stellar.org/developers/stellar-core/software/admin.html#why-run-a-node
- https://github.com/stellar/stellar-core/blob/master/docs/stellar-core_testnet.cfg
- https://github.com/stellar/packages
- https://www.stellar.org/developers/stellar-core/software/testnet.html
- https://www.postgresql.org/docs/9.2/libpq-connect.html#AEN38419
- https://www.stellar.org/developers/horizon/reference/admin.html
- https://www.stellar.org/developers/stellar-core/software/admin.html
