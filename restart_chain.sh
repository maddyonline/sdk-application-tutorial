#!/usr/bin/env bash

if [ -d "$HOME/.nsd" ] || [ -d "$HOME/.nscli" ]
then
  echo "Either $HOME/.nsd or $HOME/.nscli exits."
  echo "Please remove both these directories before continuing."
  echo "rm -rf $HOME/.ns*"
  exit 1
else
 echo "All good"
fi


MONIKER=my-mac
CHAIN_ID=my-chain
DEFAULT_PASSWORD=password12

USER_1=jack
USER_2=alice


# Initialize configuration files and genesis file
# moniker is the name of your node
nsd init ${MONIKER} --chain-id ${CHAIN_ID}


# Copy the `Address` output here and save it for later use
# [optional] add "--ledger" at the end to use a Ledger Nano S
echo ${DEFAULT_PASSWORD} | nscli keys add ${USER_1}

# Copy the `Address` output here and save it for later use
echo ${DEFAULT_PASSWORD} | nscli keys add ${USER_2}

# Add both accounts, with coins to the genesis file
nsd add-genesis-account $(nscli keys show ${USER_1} -a) 1000nametoken,100000000stake
nsd add-genesis-account $(nscli keys show ${USER_2} -a) 1000nametoken,100000000stake

# Configure your CLI to eliminate need for chain-id flag
nscli config chain-id ${CHAIN_ID}
nscli config output json
nscli config indent true
nscli config trust-node true

# Genesis transactions
echo ${DEFAULT_PASSWORD} | nsd gentx --name ${USER_1}
nsd collect-gentxs
nsd validate-genesis
nsd start
