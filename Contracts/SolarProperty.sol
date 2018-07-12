pragma solidity ^0.4.0;
contract SolarProperty {
    
    /* declaration of specialized data types */
    enum HoldingStatus {OWNED, HELD}
    enum PaymentStatus {PAID, OVERDUE}
    
    struct Holder {
        uint percentageHeld; // must be maintained that the percentageHeld for all holders sums to 100
        HoldingStatus holdingStatus;
        uint lastFullPaymentTimestamp;
        uint unpaidBalance;
    }

    struct SolarSystem {
        string name;
        uint pricePerKWH;
        mapping(address => Holder) holders;
    }
    
    /* public variables */
    address admin;
    mapping(address => SolarSystem) public solarSystems;

    /* public event on the blockchain, clients notified */
    // event AddSolarSystem(string name);
    // event AddSSHolding(address holder, uint ssIndex);
    // event Payment(address payer, uint unpaidBalance); // if paid in full, unpaidBalance = 0
    // event Repo(uint ssIndex, address lateHolder); // reposessing unpaid system, hardware can listen for this

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

    /* Transfer _percentTransfer perent of holding of solar system at _targetSSAddress to _to */ 
    function addSSHolding(uint _percentTransfer, address _targetSSAddress, address _to) public {
        require((msg.sender == admin));

        mapping(address => Holder) targetSSHolders = solarSystems[_targetSSAddress].holders;
        require(targetSSHolders[admin].percentageHeld >= _percentTransfer);

        targetSSHolders[admin].percentageHeld -= _percentTransfer; //TODO Just changedthis back to usr targetSSHoldrs variable, see if still works
        if (targetSSHolders[_to].holdingStatus == HoldingStatus.HELD) {
            targetSSHolders[_to].percentageHeld += _percentTransfer;
        } else {
            targetSSHolders[_to] = Holder({
                percentageHeld: _percentTransfer,
                holdingStatus: HoldingStatus.HELD,
                lastFullPaymentTimestamp: now,
                unpaidBalance: 0
            });
        }
    }

    function removeSSHolding(uint _percentTransfer, address _targetSSAddress, address _from) public {
        require((msg.sender == admin) || (msg.sender == _from));

        mapping(address => Holder) targetSSHolders = solarSystems[_targetSSAddress].holders;
        require(targetSSHolders[_from].percentageHeld >= _percentTransfer);

        targetSSHolders[_from].percentageHeld -= _percentTransfer;
        targetSSHolders[admin].percentageHeld += _percentTransfer;
    }


    // function energyProduced(uint _ssAddress, uint _kWhProduced) public {
    //     SolarSystem storage producingSS = solarSystems[_ssAddress];

    //     require(producingSS.currentHolder == msg.sender);

    //     producingSS.unpaidBalance += _kWhProduced*producingSS.pricePerKWH;
    //     //TODO issue Swytch token here
    // }

    // /* payment by a currnet holder for energy consumed */
    // function pay(uint _ssAddress) payable public {
    //     SolarSystem storage targetSS = solarSystems[_ssAddress];
    //     targetSS.unpaidBalance -= msg.value;
    // }

}