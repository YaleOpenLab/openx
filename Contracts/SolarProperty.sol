pragma solidity ^0.4.0;
contract SolarProperty {

    uint constant PREPA_PRICE = 10; //$/kWh
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD}
    
    struct Holder {
        uint percentageHeld; // must be maintained that the percentageHeld for all holders sums to 100
        HoldingStatus holdingStatus;
        uint lastFullPaymentTimestamp;
        uint unpaidUsage; // measured in kWh
    }

    struct SolarSystem {
        uint totalValue;
        mapping(address => Holder) holders;
    }

    struct ProposedDeployment {
        uint panelSize; // insert other physical properties here
        uint totalValue;
        address contractor;
        address consumer;
        boolean isConfirmedByContractor;
    }
    
    /* public variables */
    address admin;
    mapping(address => SolarSystem) public solarSystems;
    mapping(address => ProposedDeployment) public proposedDeployments;

    /* runs at initialization when contract is executed */
    constructor() public {
        admin = msg.sender;
    }

    // CONFIRMATION
    // agreement between contractor, participant, and investor
    function proposeDeployment(address _ssAddress, uint _panelSize, uint _totalValue, address _contractor, address _consumer) public {
        // TODO: this should be called by a verified investor, so need some way to onboard investors/contractors
        ProposedDeployment memory newDeployment;
        newDeployment.panelSize = _panelSize;
        newDeployment.totalValue = _totalValue;
        newDeployment.contractor = _contractor;
        newDeployment.consumer = _consumer;
        newDeployment.isConfirmedByContractor = false;
    }

    function confirmDeployment(address _ssAddress, address confirmer) public {
        if (proposedDeployments[_ssAddress].contractor == confirmer) {
            isConfirmedByContractor = true;
        } else if (isConfirmedByContractor && proposedDeployments[_ssAddress].consumer == confirmer) {
            ProposedDeployment newDeployment = proposedDeployments[_ssAddress];
            addSolarSystem(newDeployment.totalValue, _ssAddress);
            payContractor();
        }
    }

    function payContractor() {
        //TODO ask neighborly?
    }

    
    // This function will be called once there is confirmation of this system
    // add a new solar panel to our platform and set its properties (_totalValue), 
    // and its uniquely identifying address on the chain
    function addSolarSystem(uint _totalValue, address _ssAddress) public {
        require(msg.sender == admin);

        Holder memory adminHolder = Holder({
            percentageHeld: 100,
            holdingStatus: HoldingStatus.HELD,
            lastFullPaymentTimestamp: now,
            unpaidUsage: 0
        });
        
        SolarSystem memory newSystem;
        newSystem.totalValue = _totalValue;

        solarSystems[_ssAddress] = newSystem;
        
        // need to access holder in this way for access to storage
        solarSystems[_ssAddress].holders[admin] = adminHolder; 
    }

    // Record energy consumed by panel at targetSSAddress by consumer (called by solar panel)
    function energyConsumed(address ssAddress, address consumer, uint energyConsumed) {
        require((msg.sender == admin) || (msg.sender == ssAddress));

        solarSystems[ssAddress].holders[consumer].unpaidUsage += energyConsumed;
    }
    

    // Make a payment from consumer toward any unpaid balance on the panel at ssAddress
    function makePayment(address ssAddress, address consumer, uint amountPaid) public {
        require((msg.sender == admin) || (msg.sender == consumer));

        uint amountPaidInEnergy = amountPaid/PREPA_PRICE; 

        Holder storage consumerHolder = solarSystems[ssAddress].holders[consumer];
        if (amountPaidInEnergy >= consumerHolder.unpaidUsage) {
            consumerHolder.lastFullPaymentTimestamp = now;
        }
        
        addSSHolding(amountPaid/solarSystems[ssAddress].totalValue*100, ssAddress, consumer); // transfer over a portion of ownership
        consumerHolder.unpaidUsage -= amountPaidInEnergy; // update the unpaid balance 
    }

    // Transfer percentTransfer percent of holding of solar system at targetSSAddress to the user with address 'to'
    function addSSHolding(uint percentTransfer, address targetSSAddress, address to) public {
        require(msg.sender == admin);

        mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
        require(targetSSHolders[admin].percentageHeld >= percentTransfer);

        if (targetSSHolders[to].holdingStatus == HoldingStatus.HELD) {
            if (targetSSHolders[to].percentageHeld + percentTransfer >= 100) { // fully paid off!
                grantOwnership(targetSSAddress, to);
            } else {
                targetSSHolders[to].percentageHeld += percentTransfer;
                targetSSHolders[admin].percentageHeld -= percentTransfer;
            }
        } else { // this is their first payment towards this solar system
            targetSSHolders[to] = Holder({
                percentageHeld: percentTransfer,
                holdingStatus: HoldingStatus.HELD,
                lastFullPaymentTimestamp: now,
                unpaidUsage: 0
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


    // functions that require no gas, but will check the state of the system

    // Returns true if consumer has completely paid of any outstanding balance on the panel at ssAddress within their consumerBuffer period
    function isConsumerLiquidForSystem(address ssAddress, address consumer, uint consumerBuffer) view public returns (bool) {
        return (now - solarSystems[ssAddress].holders[consumer].lastFullPaymentTimestamp) <= consumerBuffer;
    }

}