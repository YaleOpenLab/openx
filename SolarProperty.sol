pragma solidity ^0.4.0;
contract SolarProperty {
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD}
    enum PaymentStatus {PAID, OVERDUE}
    
    struct HolderForSS {
        uint percentageHeld; // must be maintained that the percentageHeld for all holders sums to 100
        HoldingStatus holdingStatus;
        uint lastFullPaymentTimestamp;
        uint unpaidBalance;
    }

    struct SolarSystem {
        string name;
        uint pricePerKWH;
        mapping(address => HolderForSS) holders; // address of the holder
    }
    
    /* public variables */
    address approver;
    mapping(address => SolarSystem) public solarSystems;

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

        HolderForSS storage approverHolder = HolderForSS({
            percentageHeld: 100,
            holdingStatus: HoldingStatus.HELD,
            lastFullPaymentTimestamp: now,
            unpaidBalance: 0
        })

        mapping(address => HolderForSS) holderMapping;
        holderMapping[approver] = approverHolder;

        SolarSystem storage newSystem = SolarSystem({
            name: _name,
            pricePerKWH: _pricePerKWH,
            holders: holderMapping
        });

        solarSystems.push(newSystem);

        emit AddSolarSystem(_name);
    }

    /* Transfer _percentTransfer perent of holding of solar system at _targetSSIndex to _to */ 
    function transferPanelHolder(uint _percentTransfer, uint _targetSSIndex, address _to) public {
        require(msg.sender == approver);

        SolarSystem storage targetSS = solarSystems[_targetSSIndex];
        Holders[] storage holders = targetSS.holders;

        targetSS.holdingStatus = HoldingStatus.HELD;
        targetSS.currentHolder = _to;
        
        emit TransferHolder(_to, _targetSSIndex);
    }

    function removePanelHolder(uint _targetSSIndex) public {
        SolarSystem storage targetSS = solarSystems[_targetSSIndex];
        require((msg.sender == approver) || (msg.senderc == targetSS.currentHolder));

        targetSS.currentHolder = 0; // resetting, not used
        targetSS.holdingStatus = HoldingStatus.AVAILABLE;
    }


    function energyProduced(uint _ssIndex, uint _kWhProduced) public {
        SolarSystem storage producingSS = solarSystems[_ssIndex];

        require(producingSS.currentHolder == msg.sender);

        producingSS.unpaidBalance += _kWhProduced*producingSS.pricePerKWH;
        //TODO issue Swytch token here
    }

    /* payment by the currentHolder for the energy consumed */
    function pay(uint _ssIndex) payable public {
        SolarSystem storage targetSS = solarSystems[_ssIndex];
        targetSS.unpaidBalance -= msg.value;
        emit Payment(targetSS.currentHolder, targetSS.unpaidBalance);
    }

    /* transfer of ownership away from currentHolder if fails to pay */
    function repo(uint _ssIndex) public {
        require(msg.sender == approver);

        address overdueHolder = solarSystems[_ssIndex].currentHolder;
        emit Repo(_ssIndex, overdueHolder);
        
        removePanelHolder(_ssIndex);
        
    }
}