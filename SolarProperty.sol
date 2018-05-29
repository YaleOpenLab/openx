pragma solidity ^0.4.0;
contract SolarProperty {
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD, AVAILABLE}
    enum PaymentStatus {PAID, OVERDUE}
    
    struct SolarSystem {
        string name;
        HoldingStatus holdingStatus;
        address currentHolder;
        PaymentStatus paymentStatus;
        uint unpaidBalance;

        uint pricePerKWH;
    }
    
    /* public variables */
    address approver;
    SolarSystem[] public solarSystems;

    /* public event on the blockchain, clients notified */
    event AddSolarSystem(string name);
    event TransferHolder(address holder, uint ssIndex);
    event Payment(address payer, uint unpaidBalance); // if paid in full, unpaidBalance = 0
    event Repo(uint ssIndex, address lateHolder); // reposessing unpaid system, hardware can listen for this

    /* runs at initialization when contract is executed */
    constructor() public {
        approver = msg.sender;
    }

    
    function addSolarSystem(string _name, uint _pricePerKWH) public {
        require(msg.sender == approver);

        SolarSystem memory newSystem = SolarSystem({
            name: _name,
            holdingStatus: HoldingStatus.AVAILABLE,
            currentHolder: approver,
            paymentStatus: PaymentStatus.PAID,
            unpaidBalance: 0,
            pricePerKWH: _pricePerKWH
        });

        solarSystems.push(newSystem);

        emit AddSolarSystem(_name);
    }

    /* Transfer ownership of solar system at _targetSSIndex to _to */ 
    function transferPanelHolder(uint _targetSSIndex, address _to) public {
        SolarSystem storage targetSS = solarSystems[_targetSSIndex];
        if (msg.sender != approver) {
            require((targetSS.holdingStatus == HoldingStatus.HELD) && (targetSS.currentHolder == msg.sender));
        }

        targetSS.holdingStatus = HoldingStatus.HELD;
        targetSS.currentHolder = _to;
        
        emit TransferHolder(_to, _targetSSIndex);
    }

    function removePanelHolder(uint _targetSSIndex) public {
        SolarSystem storage targetSS = solarSystems[_targetSSIndex];
        require(targetSS.holdingStatus != HoldingStatus.HELD);
        require(targetSS.currentHolder != msg.sender);

        targetSS.currentHolder = 0; // resetting, not used
        targetSS.holdingStatus = HoldingStatus.AVAILABLE;
    }

    /* TODO this can only be triggered by the device on the panel */
    /* TODO do this once at the end of each day */
    function energyProduced(uint _ssIndex, uint _kWhProduced) public {
        SolarSystem storage producingSS = solarSystems[_ssIndex];
        producingSS.unpaidBalance += _kWhProduced*producingSS.pricePerKWH;
    }

    /* payment by the currentHolder for the energy consumed */
    function pay(uint _ssIndex) payable public {
        SolarSystem storage targetSS = solarSystems[_ssIndex];
        targetSS.unpaidBalance -= msg.value;
        emit Payment(targetSS.currentHolder, targetSS.unpaidBalance);
    }

    /* transfer of ownership away from currentHolder if fails to pay */
    function repo(uint _ssIndex) public {
        address overdueHolder = solarSystems[_ssIndex].currentHolder;
        emit Repo(_ssIndex, overdueHolder);
        
        removePanelHolder(_ssIndex);
        
    }
}