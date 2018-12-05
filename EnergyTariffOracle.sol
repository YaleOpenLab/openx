/*
   Energy Price Peg

   This contract keeps in storage a reference
   to the Energy Price in USD
*/


pragma solidity ^0.4.0;
import "github.com/oraclize/ethereum-api/oraclizeAPI.sol";

contract EnergyTariff is usingOraclize {
    
    uint public EnergyTariffUSD;

    event newOraclizeQuery(string description);
    event newEnergyTariff(string price);

    function EnergyTariff() {
        update(); // first check at contract creation
    }

    function __callback(bytes32 myid, string result) {
        if (msg.sender != oraclize_cbAddress()) throw;
        newEnergyTariff(result);
        EnergyTariffUSD = viewInt(result, 2); // let's save it as $ cents
    }
    
    function update() payable {
        newOraclizeQuery("Oraclize query was sent, standing by for the answer..");
        oraclize_query("URL", "xml(https://www.eia.gov/state/print.php?sid=RQ");
    }
    
}
