# TFM-UPC Pharma

This is a project simulating a pharma industry use case where blockchain can be applied to fasten the medical supply process and prevent medicament falsification.

## Set the network

In order to set up the network, run the following commands in order to set and run the network.

```
cd TFM-UPC/pharma

./startPharma.sh
```

## Init the ledgers and prepare the CLI to run the commands

```
source ./setEnvVars.sh

./initLedgers.sh
```

## Login to the channels

In order to log into, we have set some predefined users.

To log in into the Pharmacy - Laboratory network set on channel 1, we have to run the following command. This command will create 3 inputs on the shell so the user can input the entity that he is part of, his username and his password.

```
source ./userLoginCOne.sh
```

In order to log into the channel 2, where there is a network between the Hospital and the Pharmacy, we have to run the following command.

```
source ./userLoginCTwo.sh
```
