pragma solidity ^0.4.0;
contract SolarProperty {
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD}
    
    struct Holder {
        uint percentageHeld; // must be maintained that the percentageHeld for all holders sums to 100
        HoldingStatus holdingStatus;
        uint lastFullPaymentTimestamp;
        uint unpaidBalance;
    }

    struct SolarSystem {
        string name;
        uint pricePerKWH; // add any other properties of the system here
        mapping(address => Holder) holders;
    }
    
    /* public variables */
    address admin;
    mapping(address => SolarSystem) public solarSystems;

    /* runs at initialization when contract is executed */
    constructor() public {
        admin = msg.sender;
    }

    
    function addSolarSystem(string _name, uint _pricePerKWH, address _ssAddress) public {
        require(msg.sender == admin);

        Holder memory adminHolder = Holder({
            percentageHeld: 100,
            holdingStatus: HoldingStatus.HELD,
            lastFullPaymentTimestamp: now,
            unpaidBalance: 0
        });
        
        SolarSystem memory newSystem;
        newSystem.name = _name;
        newSystem.pricePerKWH = _pricePerKWH;

        solarSystems[_ssAddress] = newSystem;
        
        // need to access holder in this way for access to storage
        solarSystems[_ssAddress].holders[admin] = adminHolder; 
    }

    /* Transfer percentTransfer perent of holding of solar system at targetSSAddress to to */ 
    function addSSHolding(uint percentTransfer, address targetSSAddress, address to) public {
        require(msg.sender == admin);

        mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
        require(targetSSHolders[admin].percentageHeld >= percentTransfer);

        targetSSHolders[admin].percentageHeld -= percentTransfer;
        if (targetSSHolders[to].holdingStatus == HoldingStatus.HELD) {
            targetSSHolders[to].percentageHeld += percentTransfer;
        } else {
            targetSSHolders[to] = Holder({
                percentageHeld: percentTransfer,
                holdingStatus: HoldingStatus.HELD,
                lastFullPaymentTimestamp: now,
                unpaidBalance: 0
            });
        }
    }

    function removeSSHolding(uint percentTransfer, address targetSSAddress, address from) public {
        require((msg.sender == admin) || (msg.sender == from));

        mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
        require(targetSSHolders[from].percentageHeld >= percentTransfer);

        targetSSHolders[from].percentageHeld -= percentTransfer;
        targetSSHolders[admin].percentageHeld += percentTransfer;
    }

    function grantOwnership(address targetSSAddress, address newOwner) public {
        require(msg.sender == admin);

        mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
        targetSSHolders[newOwner].holdingStatus = HoldingStatus.OWNED;
    }

    function recordEnergyMatchupForConsumer(address ssAddress, uint kWhProduced, address consumer, uint kWhConsumed) public {
        require(msg.sender == admin);

        Holder storage consumerHolder = solarSystems[ssAddress].holders[consumer];
        uint energyForConsumer = (kWhProduced * consumerHolder.percentageHeld)/100;
        if (kWhConsumed >= energyForConsumer) {
            consumerHolder.unpaidBalance += solarSystems[ssAddress].pricePerKWH * energyForConsumer;
        } else { // didn't use all of the produced energy
            consumerHolder.unpaidBalance += solarSystems[ssAddress].pricePerKWH * kWhConsumed;
        }
    }

    function makePayment(address ssAddress, address consumer) payable public {
        require((msg.sender == admin) || (msg.sender == consumer));

        Holder storage consumerHolder = solarSystems[ssAddress].holders[consumer];
        if (msg.value >= consumerHolder.unpaidBalance) {
            consumerHolder.lastFullPaymentTimestamp = now;
            consumerHolder.unpaidBalance = 0;
        } else {
            consumerHolder.unpaidBalance -= msg.value; 
        }
    }


    // functions that require no gas, but will check the state of the system

    function isConsumerLiquidForSystem(address ssAddress, address consumer, uint consumerBuffer) view public returns (bool) {
        return (now - solarSystems[ssAddress].holders[consumer].lastFullPaymentTimestamp) <= consumerBuffer;
    }

}