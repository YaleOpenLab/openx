contract SolarProperty {
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD, AVAILABLE}
    enum PaymentStatus {Paid, Overdue}
    
    struct SolarSystem {
        string name;
        HoldingStatus holdingStatus;
        address currentHolder;
        PaymentStatus paymentStatus;
        uint unpaidBalance;

        uint pricePerKWH;
    }
    
    /* public variables */
    SolarSystem[] public solarSystems;

    /* public event on the blockchain, clients notified */
    event AddHolder(address holder, uint ssIndex);
    event Payment(address payer, uint unpaidBalance); // if paid in full, unpaidBalance = 0
    event Repo(uint ssIndex, address lateHolder); // reposessing unpaid system, hardware can listen for this

    /* runs at initialization when contract is executed */
    constructor() {
        
    }

    /* if SS at index _targetSSIndex is currently not held by another user, transfer ownership to _holder*/ 
    function addPanelHolder(uint _targetSSIndex, address _holder) {
        SolarSystem memory targetSS = solarSystems[_targetSSIndex];
        require(targetSS.holdingStatus != HoldingStatus.AVAILABLE);

        targetSS.currentHolder = _holder;
        
        AddHolder(_holder, _targetSSIndex);
    }

    function removePanelHolder(uint _targetSSIndex) {
        SolarSystem memory targetSS = solarSystems[_targetSSIndex];
        require(targetSS.holdingStatus != HoldingStatus.HELD);
        require(targetSS.currentHolder != msg.sender);

        targetSS.currentHolder = 0; // resetting, not used
        targetSS.holdingStatus = HoldingStatus.AVAILABLE;
    }

    /* TODO this can only be triggered by the device on the panel */
    /* TODO do this once at the end of each day */
    function energyProduced(uint _ssIndex, uint _kWhProduced) {
        SolarSystem memory producingSS = solarSystems[_ssIndex];
        producingSS.unpaidBalance += _kWhProduced*producingSS.pricePerKWH;
    }

    /* payment by the currentHolder for the energy consumed */
    function pay(uint _ssIndex) payable {
        SolarSystem memory targetSS = solarSystems[_ssIndex];
        targetSS.unpaidBalance -= msg.value;
        Payment(targetSS.currentHolder, targetSS.unpaidBalance);
    }

    /* transfer of ownership away from currentHolder if fails to pay */
    function repo(uint _ssIndex) {
        address overdueHolder = solarSystems[_ssIndex].currentHolder;
        Repo(_ssIndex, overdueHolder);
        
        removePanelHolder(_ssIndex);
        
    }
}