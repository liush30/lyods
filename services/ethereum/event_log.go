package ethereum

import (
	"lyods-adsTool/domain"
	"lyods-adsTool/pkg/constants"
	"math/big"
	"strings"
)

// IsDepositEvent 判断是否为传统Deposit事件
// Deposit (index_topic_1 address sender, uint256 value)
func IsDepositEvent(interParam *InternalTxnParam, txList *[]domain.InternalTxn) bool {
	var internalTxn domain.InternalTxn
	if !strings.EqualFold(interParam.eventName, constants.EVENT_NAME_DEPOSIT) {
		return false
	}
	_, uintOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT]
	_, addressOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]
	if interParam.length == 2 && addressOk && uintOk {
		for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
			internalTxn.FromAddr = v
		}
		for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT] {
			//将v转成big.int类型
			value, _ := new(big.Int).SetString(v, 10)
			internalTxn.Amount = *value
		}
		internalTxn.ToAddr = interParam.contractAddress
		internalTxn.ContractAddr = interParam.contractAddress
		internalTxn.IsErc20 = interParam.isErc20
		internalTxn.TokenDecimal = interParam.tokenDecimal
		internalTxn.Token = interParam.token
		*txList = append(*txList, internalTxn)
		return true
	}
	return false
}

// IsWithdrawalEvent 判断是否为传统Withdrawal事件
// Withdrawal (index_topic_1 address sender, uint256 value)
// Withdrawal (address to, bytes32 nullifierHash, index_topic_1 address relayer, uint256 fee)
func IsWithdrawalEvent(interParam *InternalTxnParam, txList *[]domain.InternalTxn) bool {

	//if !strings.EqualFold(interParam.eventName, constants.EVENT_NAME_WITHDRAWAL) {
	//	return false
	//}
	//uintVal, uintOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT]
	//addrVal, addressOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]
	////判断interParam.eventNameToValueByAddress 参数值是否符合Withdrawal (index_topic_1 address sender, uint256 value)
	//if interParam.length == 2 && addressOk && uintOk {
	//	var internalTxn domain.InternalTxn
	//	for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
	//		internalTxn.ToAddr = v
	//	}
	//	for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT] {
	//		//将v转成big.int类型
	//		value, _ := new(big.Int).SetString(v, 10)
	//		internalTxn.Amount = *value
	//	}
	//	internalTxn.FromAddr = interParam.contractAddress
	//	internalTxn.ContractAddr = interParam.contractAddress
	//	internalTxn.IsErc20 = interParam.isErc20
	//	internalTxn.TokenDecimal = interParam.tokenDecimal
	//	internalTxn.Token = interParam.token
	//	*txList = append(*txList, internalTxn)
	//	return true
	//}
	//toVal, toOk := addrVal[constants.EVENT_PARAM_TO]
	//relayerVal, relayerOk := addrVal[constants.EVENT_PARAM_RELAYER] //中转者
	//feeVal, feeOk := uintVal[constants.EVENT_PARAM_FEE]             //中转者收取的费用
	////判断interParam.eventNameToValueByAddress 参数值是否符合 Withdrawal (address to, bytes32 nullifierHash, index_topic_1 address relayer, uint256 fee)
	//if interParam.length == 4 && toOk && relayerOk && feeOk {
	//	var internalTxn domain.InternalTxn
	//	//将中转者，和直接to作为两个内部交易存储
	//	return true
	//}
	return false
}

// IsTransferEvent 判断是否为传统Transfer事件
// Transfer(address indexed from, address indexed to, uint256 value/amount/money)
func IsTransferEvent(interParam *InternalTxnParam, txList *[]domain.InternalTxn) bool {
	var internalTxn domain.InternalTxn
	//判断eventName是否为Transfer
	if !strings.EqualFold(interParam.eventName, constants.EVENT_NAME_TRANSFER) {
		return false
	}
	fromValue, fromOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS][constants.EVENT_PARAM_FROM]
	toValue, toOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS][constants.EVENT_PARAM_TO]
	_, uintOk := interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT]
	//判断是否存在(address indexed from, address indexed to, uint256 value)结构
	if interParam.length == 3 && fromOk && toOk && uintOk {
		for _, v := range interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT] {
			//将v转成big.int类型
			value, _ := new(big.Int).SetString(v, 10)
			internalTxn.Amount = *value
		}
		internalTxn.FromAddr = fromValue
		internalTxn.ToAddr = toValue
		internalTxn.ContractAddr = interParam.contractAddress
		internalTxn.IsErc20 = interParam.isErc20
		internalTxn.TokenDecimal = interParam.tokenDecimal
		internalTxn.Token = interParam.token
		*txList = append(*txList, internalTxn)
		return true
	}
	return false
}

//if len(interParam.eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS]) == 2 && len(interParam.eventNameToValueByAddress[constants.EVENT_TYPR_UINT]) == 1 {
//	//from,to都存在,根据字段名称存储
//	if fromOk && toOk {
//
//		//to字段存在，from字段不存在
//	} else if fromOk && !toOk {
//		internalTxn.FromAddr = fromValue
//		delete(eventNameToValueByAddress, constants.EVENT_PARAM_FROM)
//		for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
//			internalTxn.ToAddr = v
//		}
//		//from字段存在，to字段存在
//	} else if !fromOk && toOk {
//		internalTxn.ToAddr = toValue
//		delete(eventNameToValueByAddress, constants.EVENT_PARAM_TO)
//		for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
//			internalTxn.FromAddr = v
//		}
//		//from,to都不存在
//	} else {
//		//判断地址中哪个是from,哪个是to
//		if isErc20 {
//			for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
//				changeMoney, err := e.GetERC20TokenBalanceChange(contractAddress, v, blockNumber)
//				if err != nil {
//					log.Fatal("Fail GetERC20TokenBalanceChange:", err)
//				}
//				//判断changeMoney是否大于0
//				if changeMoney.Cmp(big.NewInt(0)) > 0 {
//					internalTxn.ToAddr = v
//				} else {
//					internalTxn.FromAddr = v
//				}
//			}
//		} else {
//			for _, v := range eventNameToValueByAddress[constants.EVENT_TYPE_ADDRESS] {
//				changeMoney, err := e.GetBalanceChange(blockNumber, v)
//				if err != nil {
//					log.Fatal("Fail GetERC20TokenBalanceChange:", err)
//				}
//				//判断changeMoney是否大于0
//				if changeMoney.Cmp(big.NewInt(0)) > 0 {
//					internalTxn.ToAddr = v
//				} else {
//					internalTxn.FromAddr = v
//				}
//			}
//		}
//	}
//	//判断internalTxn.FromAddr、internalTxn.ToAddr是否为空
//	if len(internalTxn.FromAddr) == 0 || len(internalTxn.ToAddr) == 0 {
//		log.Fatal("Fail internalTxn.FromAddr or internalTxn.ToAddr is empty")
//	}
//	for _, v := range eventNameToValueByAddress[constants.EVENT_TYPR_UINT] {
//		//将v转成big.int类型
//		value, _ := new(big.Int).SetString(v, 10)
//		internalTxn.Amount = value
//	}
//	return
//} else {
//
//}
