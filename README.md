# TFM-UPC Pharma

This is a project simulating a pharma industry use case where blockchain can be applied to fasten the medical supply process and prevent medicament falsification.

![Network](https://i.ibb.co/3v4KnGb/Screenshot-2022-11-17-at-21-40-03.png)

## Set the network

In order to set up the network, run the following commands in order to set and run the network.

```
$ cd TFM-UPC/pharma
$ ./startPharma.sh
```

## Init the ledgers and prepare the CLI to run the commands

```
$ source ./setEnvVars.sh
$ ./initLedgers.sh
```

## Login to the channels

In order to log into, we have set some predefined users.

To log in into the Pharmacy - Laboratory network set on channel 1, we have to run the following command. This command will create 3 inputs on the shell so the user can input the entity that he is part of, his username and his password.

```
$ source ./userLoginCOne.sh
```

In order to log into the channel 2, where there is a network between the Hospital and the Pharmacy, we have to run the following command.

```
$ source ./userLoginCTwo.sh
```

## Demo with the network

In order to give a small test on our network, we are going to simulate the full process of a medicine from its creation un till it given to the user, showing its status on each point.

![Demo](https://i.ibb.co/tJV8gdM/Screenshot-2022-11-17-at-22-41-13.png)

### Medicine registration

First of all, we need to log in with the lab user or lab admin so we can register the medicine in channel 1.

In order to register the medicine, we execute the following command and select the option 1.

```
$ ./userLabFunctions.sh
1) 1- Register a medicament
2) 2- Send a medicament to a pharmacy
3) 3- Quit
Which of the following actions you want to do: 1
```

Once we selected the option, we have to introduce the medicine data: Medicine Name, the medicine product code, its serial number, the expiration year and the expiration month.

### Read created medicine / Read all users / Read all medicines

Once we have the medicine registered, and logged in with the admin user, we can check its information and also the information about all the medicines and users related to the pharmacy.

```
$ ./adminLabFunctions.sh
1) 1- Read all medicaments    3) 3- Read a medicament info
2) 2- Read all users          4) 4- Quit
```

### Send medicine to pharmacy

As we did before, we have to be logged in with one of the lab users. Then we select the option 2 and input the medicine product code.

```
$ ./userLabFunctions.sh
1) 1- Register a medicament
2) 2- Send a medicament to a pharmacy
3) 3- Quit
Which of the following actions you want to do: 2
```

### Notify received medicine

When the medicine is received by the pharmacy. The pharmacy, throught a pharmacy user or admin, will have to notify that it was received to both channels.

```
$ ./userPharmacyFunctions.sh
1) 1- Receive a medicament   3) 3- Use prescription
2) 2- Dispense a medicament  4) 4- Quit
Which of the following actions you want to do: 1
```

### Create medicine prescription

Once the medicine is recieved by the pharmacy and notified to the pharmacy-hospital channel2, at the hospital side you will have to create a prescription in order to allow its dispatchment.

To do so, first we will have to log in with a member of the hospital (user or admin), then the hospital member will be able to generate a prescription for the patient.

```
./generatePrescription.sh
```

### User receives the medicine (Consume prescription and medicine dispatchment)

Once the prescription is generated and received at the channel, the patient will be able to go to the pharmacy and ask for the medicine that was asigned to him. In order to get the medicine, the prescription will be consumed and the medicine will be dispatched from the pharmacy.

```
$ ./userPharmacyFunctions.sh
1) 1- Receive a medicament   3) 3- Use prescription
2) 2- Dispense a medicament  4) 4- Quit
Which of the following actions you want to do:
```

### More information

For deeper knowledge of the network we have built, its configuration and its chaincode, we also complemented the project with written memory that was delivered to our teacher.
