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

    // add a new solar panel to our platform and set its name _name, physical properties (_pricePerKWH), 
    // and its uniquely identifying address on the chain
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

    // Transfer percentTransfer percent of holding of solar system at targetSSAddress to the user with address 'to'
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

    // Reclaim percentTransfer percent of holding of solar system at targetSSAddress currently held by user with address 'from' 
    function removeSSHolding(uint percentTransfer, address targetSSAddress, address from) public {
        require((msg.sender == admin) || (msg.sender == from));

        mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
        require(targetSSHolders[from].percentageHeld >= percentTransfer);

        targetSSHolders[from].percentageHeld -= percentTransfer;
        targetSSHolders[admin].percentageHeld += percentTransfer;
    }

    // Grant ownership to newOwner of whatever portion of the solar panel at targetSSAddress they currently hold
    function grantOwnership(address targetSSAddress, address newOwner) public {
        require(msg.sender == admin);

        mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
        targetSSHolders[newOwner].holdingStatus = HoldingStatus.OWNED;
    }

    // Given a solar panel at ssAddress and the kWh produced, and given a consumer and their kWhConsumed, 
    // update the balance they owe on this solar panel
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

    // Make a payment from consumer toward any unpaid balance on the panel at ssAddress
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

    // Returns true if consumer has completely paid of any outstanding balance on the panel at ssAddress within their consumerBuffer period
    function isConsumerLiquidForSystem(address ssAddress, address consumer, uint consumerBuffer) view public returns (bool) {
        return (now - solarSystems[ssAddress].holders[consumer].lastFullPaymentTimestamp) <= consumerBuffer;
    }

}