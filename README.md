# SmartProperty MVP

#### Main Idea 
Each Microgrid energy producer/component will serve as a “lite client” that has access to/exists as an address on the Ethereum blockchain. This address/client will hold funds and continuously execute contract calls on the solar property contract and observe the returned data. With this data, the client will determine if funds are sufficient to continue providing energy from the connected renewable energy source.

#### Key Terms

Infura - An API built by a spoke of Consensys in order to open access to the blockchain through familiar REST patterns.

CloudMQTT - Message queueing service. This system acts as a buffer for data and actions that need to be forwarded. Received messages are organized under “topics” that other devices/servers can listen to. These listeners will be pinged when the topic receives new messages. It’s resources are hosted on AWS

Blockchain - A smart contract will be published on the Ethereum blockchain. This contract will contain the logic for payments by microgrid participants. It will keep track of ownership and transfer ownership as needed (when payments aren’t fulfilled). 

#### Current Implementation
The script found in ____ can be run on the photon which will publish a message to a topic “” on the MQ. The server code written in go (based on _____) can be run on a local machine to simulate how a cloud hosted server would respond. The server will listen to the topic ”” and when pinged will make an infura call to a dummy contract that was published _____. The current transaction doesn’t need gas/to be signed because it’s an eth_call which doesn’t alter or publish to the contract but simply requests information. 

#### Next Steps
The software components are connected and functional. Now we need to begin writing the necessary code that will work with the energy data. We need to parse and hash energy data that is received by the hardware system. Then, write photon script that will publish this data on the MQ service. At the same time, we are working to publish a prototype smart contract (`SolarProperty.sol`) on the testnet that will have the property transaction logic. 

#### Possible Alternative system: 