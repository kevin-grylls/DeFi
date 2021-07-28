// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./ERC20.sol";
import "./SafeMath.sol";
import "./Owned.sol";

abstract contract ApproveAndCallFallBack {
    function receiveApproval(address from, uint256 tokens, address token, bytes memory data) public virtual;
}

contract InsuranceToken is ERC20, Owned {
    uint256 constant WEI = 1000000000000000000;
    using SafeMath for uint256;
    uint256 _totalSupply;
    uint _exchangeRate = 2600000;
    uint _revenueRate = 2613000; // 0.005% rate
    uint _averageRate = 2400000;

    mapping(address => uint256) balances;
    mapping(address => mapping(address => uint)) allowed;
    
    event BuyInsurance();
    event RefundSavings(address to, uint256 tokens);
    

    constructor() ERC20("Bonus","BN")  {
        _totalSupply = 0;
        init();
        balances[owner] = _totalSupply;
        emit Transfer(address(0), owner, _totalSupply);
    }
    
    function setExchangeRate(uint rate) onlyOwner public {
        require(rate > 0);
        _exchangeRate = rate;
    }
    
    function setRevenueRate(uint rate) onlyOwner public {
        require(rate > 0);
        _revenueRate = rate;
    }
    
    function setAvgRate(uint rate) onlyOwner public {
        require(rate > 0);
        _averageRate = rate;
    }
    
    function getAllRates() onlyOwner public view returns (uint[] memory) {
        uint256[] memory rates = new uint256[](3);
        
        rates[0] = _exchangeRate;
        rates[1] = _revenueRate;
        rates[2] = _averageRate;
        
        return rates;
    }

    function totalSupply() public override view returns (uint256) {
        return _totalSupply.sub(balances[address(0)]);
    }


    function balanceOf(address tokenOwner) public override view returns (uint256 balance) {
        return balances[tokenOwner];
    }


    function transfer(address to, uint256 tokens) public override returns (bool success) {
        balances[msg.sender] = SafeMath.sub(balances[msg.sender], tokens);
        balances[to] = SafeMath.add(balances[to], tokens);
        emit Transfer(msg.sender, to, tokens);
        return true;
    }


    function approve(address spender, uint256 tokens) public override returns (bool success) {
        allowed[msg.sender][spender] = tokens;
        emit Approval(msg.sender, spender, tokens);
        return true;
    }


    function transferFrom(address from, address to, uint256 tokens) public override returns (bool success) {
        approve(from, tokens);
        
        uint256 balance = balances[from];
        balances[from] = SafeMath.sub(balance, tokens);
        
        uint256 receiverBalance = balances[to];
        balances[to] = SafeMath.add(receiverBalance, tokens);
        
        emit Transfer(from, to, tokens);
        return true;
    }


    function allowance(address tokenOwner, address spender) public override view returns (uint256 remaining) {
        return allowed[tokenOwner][spender];
    }


    function approveAndCall(address spender, uint tokens, bytes memory data) public returns (bool success) {
        allowed[msg.sender][spender] = tokens;
        emit Approval(msg.sender, spender, tokens);
        ApproveAndCallFallBack(spender).receiveApproval(msg.sender, tokens, address(this), data);
        return true;
    }


     function buyInsurance() external payable {
        require(msg.value > 0.01 ether, "0.01 Ethereum is minimum value to buy bonus.");
        uint256 converedToEther = SafeMath.div(msg.value, WEI);
        uint256 point = SafeMath.mul(converedToEther,_exchangeRate);
        
        balances[msg.sender] = SafeMath.add(balances[msg.sender],point);
        _totalSupply = SafeMath.add(_totalSupply, point);
        emit BuyInsurance();
    }
    
    function refundSavings(uint256 tokens) public {
        require(tokens > 0, "You need to sell at least 1 tokens");
        
        approve(msg.sender, tokens);
        uint256 balance = balanceOf(msg.sender);
        require(balance >= tokens, "Cannot refund tokens than you currently have.");
        
        uint256 remainBalance = SafeMath.sub(balance, tokens);
        balances[msg.sender] = remainBalance;
        
        uint256 remainTotalBalance = SafeMath.sub(_totalSupply, tokens);
        _totalSupply = remainTotalBalance;
 
        uint256 revenue = SafeMath.div(tokens, _revenueRate);
        uint256 finalRevenue = SafeMath.mul(revenue, WEI);
        
        address payable wallet = payable(msg.sender);
        wallet.transfer(finalRevenue);
        emit RefundSavings(msg.sender, tokens);
    }
    
    
    function refundFinalSavings() public {
        uint256 balance = balanceOf(msg.sender);
        require(balance > 0, "Cannot refund tokens than you currently have.");
        
        balances[msg.sender] = 0;
        _totalSupply =  SafeMath.sub(_totalSupply, balance);
        
        uint256 revenue = SafeMath.div(balance, _averageRate);
        uint256 finalRevenue = SafeMath.mul(revenue, WEI);
        
        address payable wallet = payable(msg.sender);
        wallet.transfer(finalRevenue);
        emit RefundSavings(msg.sender, balance);
    }

}