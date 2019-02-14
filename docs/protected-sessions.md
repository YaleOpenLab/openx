# Protected Sessions

Protected sessions allows a server and client to communicate securely using MQTT. This is done by generating the relevant public and private keys for encrypting traffic on both ends. Right now, the openx platform uses Swytch.io's MQTT broker but in the future may use a more generalized construction. THe process for setting up this is relatively straightforward:

1. Download the tar blob that Atonomi provides one with (not linked here because that part is not open source yet)
2. Extract the blob
3. For your architecture, create files that end in `1` rather than the `1.0` that the tar blob ships with.
3. Point the `LD_LIBRARY_PATH` variable to the centri/lib subfolder in your extracted directory
4. `make` from the project's root directory and run the `demoem` executable.
