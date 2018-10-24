pragma solidity ^0.4.11;
contract SolarProperty {

    uint constant PREPA_PRICE = 10; //$/kWh
    
    struct SolarSystem {
        uint panelSize;
        uint totalValue;
        address consumer;
        uint percentageHeld;
        uint lastPaymentTimestamp;
        uint unpaidUsage; // measured in kWh
    }

    struct ProposedDeployment {
        uint panelSize; // insert other physical properties here
        uint totalValue;
        uint256 contractorPayout;
        
        address contractor;
        address consumer;
        address investor;
        
        //must have both confirmations to be valid
        bool isConfirmedByContractor;
        bool isConfirmedByConsumer;
    }
    
    /* public variables */
    address admin;
    mapping(uint => SolarSystem) public solarSystems;
    mapping(uint => ProposedDeployment) public proposedDeployments;

    /* runs at initialization when contract is executed */
    constructor() public {
        admin = msg.sender;
    }
    
    //getter for proposed deployment
    function getProposedDeploymentDetails(uint _ssAddress) view public returns(uint, uint, uint256, address, address, address, bool){
        return(
            proposedDeployments[_ssAddress].panelSize, 
            proposedDeployments[_ssAddress].totalValue, 
            proposedDeployments[_ssAddress].contractorPayout, 
            proposedDeployments[_ssAddress].consumer,
            proposedDeployments[_ssAddress].contractor,
            proposedDeployments[_ssAddress].investor,
            proposedDeployments[_ssAddress].isConfirmedByContractor
        );
    }
    
    function getSolarSystemDetails(uint _ssAddress) view public returns(uint, uint, address, uint, uint){
        return(
            solarSystems[_ssAddress].panelSize, 
            solarSystems[_ssAddress].totalValue,  
            solarSystems[_ssAddress].consumer,
            solarSystems[_ssAddress].percentageHeld,
            solarSystems[_ssAddress].unpaidUsage
        );
    }

    // CONFIRMATION
    // agreement between contractor, participant, and investor
    function proposeDeployment(uint _ssAddress, uint _payout,  uint _panelSize, uint _totalValue, address _contractor, address _consumer) payable public {
        // TODO: this should be called by a verified investor, so need some way to onboard investors/contractors
        ProposedDeployment memory newDeployment;
        newDeployment.panelSize = _panelSize;
        newDeployment.totalValue = _totalValue;
        newDeployment.contractorPayout = _payout;
        newDeployment.contractor = _contractor;
        newDeployment.consumer = _consumer;
        newDeployment.investor = msg.sender;
        newDeployment.isConfirmedByContractor = false;
        newDeployment.isConfirmedByConsumer = false;
        proposedDeployments[_ssAddress] = newDeployment;
    }

    function confirmDeployment(uint _ssAddress) public {
        if (proposedDeployments[_ssAddress].contractor == msg.sender) {
            proposedDeployments[_ssAddress].isConfirmedByContractor = true;
        } else if (proposedDeployments[_ssAddress].isConfirmedByContractor && proposedDeployments[_ssAddress].consumer == msg.sender) {
            
            //for record keeping
            proposedDeployments[_ssAddress].isConfirmedByConsumer = true;
            
            //create a solarsystem
            addSolarSystem(_ssAddress, msg.sender);
        }
    }
    
    //allow contractor to collect payout on system
    function collectPayout(uint _ssAddress) public {
        require(msg.sender == proposedDeployments[_ssAddress].contractor);
        require(proposedDeployments[_ssAddress].isConfirmedByConsumer && proposedDeployments[_ssAddress].isConfirmedByContractor);
        msg.sender.transfer(proposedDeployments[_ssAddress].contractorPayout);
    }

    // This function will be called once there is confirmation of this system
    // add a new solar panel to our platform and set its properties (_totalValue), 
    // and its uniquely identifying address on the chain
    function addSolarSystem(uint _ssAddress, address _consumer) private {
        
        SolarSystem memory newSystem;
        newSystem.totalValue = proposedDeployments[_ssAddress].totalValue;
        
        //TODO: must change this because now refers to block timestamp not current time
        newSystem.lastPaymentTimestamp = now;
        newSystem.unpaidUsage = 0;
        newSystem.percentageHeld = 0;
        newSystem.panelSize = proposedDeployments[_ssAddress].panelSize;
        newSystem.consumer = _consumer;

        solarSystems[_ssAddress] = newSystem;
    }

    // Record energy consumed by panel at targetSSAddress by consumer (called by solar panel)
    function energyConsumed(uint _ssAddress, uint _energyConsumed) public {
        //TODO: require sent from specific address
        require((proposedDeployments[_ssAddress].isConfirmedByContractor) && (proposedDeployments[_ssAddress].isConfirmedByConsumer));
        solarSystems[_ssAddress].unpaidUsage += _energyConsumed;
        
    }
    

    // Make a payment from consumer toward any unpaid balance on the panel at ssAddress
    function makePayment(uint _ssAddress, uint _amountPaid) payable public {
        require(_amountPaid == msg.value);

        uint amountPaidInEnergy = _amountPaid/PREPA_PRICE; 

        solarSystems[_ssAddress].lastPaymentTimestamp = now;
        uint percentageChange = getPercent(_amountPaid, solarSystems[_ssAddress].totalValue);
        solarSystems[_ssAddress].percentageHeld += percentageChange; // transfer over a portion of ownership
        solarSystems[_ssAddress].unpaidUsage -= amountPaidInEnergy; // update the unpaid balance 
    }
    
    //logic here should should
    function getPercent(uint part, uint whole) public pure returns(uint percent) {
        uint numerator = part * 1000;
        require(numerator > part); // overflow. Should use SafeMath throughout if this was a real implementation. 
        uint temp = numerator / whole + 5; // proper rounding up
        return temp / 10;
    }

    // // Transfer percentTransfer percent of holding of solar system at targetSSAddress to the user with address 'to'
    // function addSSHolding(uint _percentTransfer, uint _targetSSAddress) private {
    //     solarSystems[_targetSSAddress].percentageHeld += _percentTransfer;
    // }

    // // Reclaim percentTransfer percent of holding of solar system at targetSSAddress currently held by user with address 'from' 
    // function removeSSHolding(uint percentTransfer, address targetSSAddress, address from) public {
    //     require((msg.sender == admin) || (msg.sender == from));

    //     mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
    //     require(targetSSHolders[from].percentageHeld >= percentTransfer);

    //     targetSSHolders[from].percentageHeld -= percentTransfer;
    //     targetSSHolders[admin].percentageHeld += percentTransfer;
    // }

    // // Grant ownership to newOwner of whatever portion of the solar panel at targetSSAddress they currently hold
    // function grantOwnership(address targetSSAddress, address newOwner) public {
    //     require(msg.sender == admin);

    //     mapping(address => Holder) targetSSHolders = solarSystems[targetSSAddress].holders;
    //     targetSSHolders[newOwner].holdingStatus = HoldingStatus.OWNED;
    // }


    // // functions that require no gas, but will check the state of the system

    // // Returns true if consumer has completely paid of any outstanding balance on the panel at ssAddress within their consumerBuffer period
    // function isConsumerLiquidForSystem(address ssAddress, address consumer, uint consumerBuffer) view public returns (bool) {
    //     return (now - solarSystems[ssAddress].holders[consumer].lastFullPaymentTimestamp) <= consumerBuffer;
    // }

}