# Teller

Teller runs on the IoT hub in a particular project in case of the opensolar platform. The teller is responsible for assimilating data from the Zigbee devices installed on the solar panels, filtering out noise and transmitting resulting data to external third party providers who may be interested to partner on providing various services. This would also be used to get data about the amount that the recipient owes per month and is used to trigger payback.

We need the teller to be secure and tramper resistant, because it is going to face the recipient of a particular project. Broadly, what the teller needs is:

1. The code that is running on the device is same as the one we want to run on the device.
2. The teller doesn't randomly stop and hence reduce the amount owed by the recipient.

These can be combined if we manage to track the uptime of a device and generate an unique id each time it (re)starts. In this prototype, this is the publickey of the device is used as this unique id and we track uptime by committing the ipfs hash of the start and end times to the blockchain in two unique stellar transactions. Any public entity can try to retrieve these details and make sure that the device's uptime is as expected, hence providing both guarantees. This can be improved with reputation based systems like the one Atonomi provides and we will be using their service to improve this guarantee.
