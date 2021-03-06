pragma solidity 0.4.24;

/**
 * @title The LINK exchange contract
 */
contract LinkEx {

  uint256 private historicRate;
  uint256 private rate;
  uint256 private rateHeight;

  function currentRate() public view returns (uint256) {
    if (isFutureBlock()) {
      return rate;
    }
    return historicRate;
  }

  function update(uint256 _rate) public {
    if (isFutureBlock()) {
      historicRate = rate;
      rateHeight = block.number;
    }
    rate = _rate;
  }

  function isFutureBlock() internal view returns (bool) {
    return block.number > rateHeight;
  }
}
