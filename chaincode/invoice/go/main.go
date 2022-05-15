package main

import (
	"encoding/json"
	"fmt"

	"github.com/zghh/redactable-fabric/core/chaincode/shim"
	"github.com/zghh/redactable-fabric/protos/peer"
)

const (
	submitter                int = 1                                                   // 上链方
	superintendent           int = 2                                                   // 监管方
	belonger                 int = 3                                                   // 归属方
	inquirer                 int = 4                                                   // 查询方
	recoder                  int = 5                                                   // 入账方
	locker                   int = 6                                                   // 锁定方
	producer                 int = 7                                                   // 监制方
	submitterState           int = 1                                                   // 上链方状态
	belongerState            int = 2                                                   // 归属方状态
	recoderState             int = 4                                                   // 入账方状态
	inquirerState            int = 8                                                   // 查询方状态
	superintendentState      int = 16                                                  // 监管方状态
	lockerState              int = 32                                                  // 锁定方状态
	producerState            int = 64                                                  // 监制方状态
	printAuthority           int = submitterState                                      // 打印
	accreditAuthority        int = belongerState                                       // 授权
	assignAuthority          int = belongerState                                       // 转让
	lockAuthority            int = belongerState | inquirerState                       // 锁定
	unlockAuthority          int = lockerState                                         // 解锁
	recordAuthority          int = belongerState | inquirerState | lockerState         // 入账
	unrecordAuthority        int = recoderState                                        // 入账撤销
	queryCiphertextAuthority int = belongerState | superintendentState | producerState // 查询隐私密文

	publicAccountType         string = "1" // 公共代管账户
	financeAccountType        string = "2" // 财政初始账户
	normalAccountType         string = "3" // 一般财政账户
	socialAccountType         string = "4" // 社会化服务平台账户
	insuranceAccountType      string = "5" // 医保账户
	reimbursementAccountType  string = "6" // 一般报销单位账户
	invoiceAccountType        string = "7" // 开票单位账户
	publicAccountState        int    = 1   // 公共代管账户状态
	financeAccountState       int    = 2   // 财政初始账户状态
	normalAccountState        int    = 4   // 一般财政账户状态
	socialAccountState        int    = 8   // 社会化服务平台账户状态
	insuranceAccountState     int    = 16  // 医保账户状态
	reimbursementAccountState int    = 32  // 一般报销单位账户状态
	invoiceAccountState       int    = 64  // 开票单位账户状态
)

type Invoice struct {
}

type EInvoiceInfo struct {
	EInvoiceCode          string             `json:"eInvoiceCode"`                    // 电子票据代码
	EInvoiceNumber        string             `json:"eInvoiceNumber"`                  // 电子票据号码
	EInvoiceAddress       string             `json:"eInvoiceAddress,omitempty"`       // 电子票据地址
	IsRed                 bool               `json:"isRed"`                           // 是否为红票
	RelatedEInvoiceCode   string             `json:"relatedEInvoiceCode,omitempty"`   // 关联票据代码
	RelatedEInvoiceNumber string             `json:"relatedEInvoiceNumber,omitempty"` // 关联票据号码
	Lock                  bool               `json:"lock"`                            // 是否锁定
	PaperInfos            []PaperInfo        `json:"paperInfos,omitempty"`            // 打印信息
	RecordInfos           []RecordInfo       `json:"recordInfos,omitempty"`           // 入账信息
	UnrecordInfos         []RecordInfo       `json:"unrecordInfos,omitempty"`         // 入账撤销信息
	AuthorizedInfos       []AuthorizedInfo   `json:"authorizedInfos,omitempty"`       // 授权信息
	AssignInfos           []AssignInfo       `json:"assignInfos,omitempty"`           // 转让信息
	LockInfos             []LockOrUnlockInfo `json:"lockInfos,omitempty"`             // 锁定信息
	UnlockInfos           []LockOrUnlockInfo `json:"unlockInfos,omitempty"`           // 解锁信息
}

type PaperInfo struct {
	PaperCode   string `json:"paperCode"`
	PaperNumber string `json:"paperNumber"`
}

type RecordInfo struct {
	AgencyCode  string `json:"agencyCode"`  // 单位代码
	AgencyName  string `json:"agencyName"`  // 单位名称
	AccNumber   string `json:"accNumber"`   // 入账凭证号
	AccAmount   int64  `json:"accAmount"`   // 入账金额
	AccDateTime string `json:"accDateTime"` //入账时间
}

type AuthorizedInfo struct {
	AccountId           string `json:"accountId"`           // 账户id
	AuthorizerName      string `json:"authorizerName"`      // 授权人名称
	AuthorizerCode      string `json:"authorizerCode"`      // 授权人代码
	AuthorizerCodeType  string `json:"authorizerCodeType"`  // 授权人代码类型
	TargetAccountId     string `json:"targetAccountId"`     // 目标账户id
	AuthorizedPartyName string `json:"authorizedPartyName"` // 被授权单位名称
	AuthorizedPartyCode string `json:"authorizedPartyCode"` // 被授权单位代码
}

type AssignInfo struct {
	AccountId          string `json:"accountId"`          // 账户id
	AssignorType       string `json:"assignorType"`       // 转让方身份
	AssignorName       string `json:"assignorName"`       // 转让人名称
	AssignorCode       string `json:"assignorCode"`       // 转让人代码
	AssignorCodeType   string `json:"assignorCodeType"`   // 转让方代码类型
	TargetAccountId    string `json:"targetAccountId"`    // 目标账户id
	TransfereeName     string `json:"transfereeName"`     // 被转让人名称
	TransfereeCode     string `json:"transfereeCode"`     // 被转让人代码
	TransfereeCodeType string `json:"transfereeCodeType"` // 被转让方代码类型
}

type LockOrUnlockInfo struct {
	AgencyCode string `json:"agencyCode"` // 单位代码
	AgencyName string `json:"agencyName"` // 单位名称
}

type BaseData struct {
	EInvoiceCode   string `json:"eInvoiceCode"`
	EInvoiceNumber string `json:"eInvoiceNumber"`
}

type SaveEInvoiceData struct {
	BaseData
	EInvoice              string         `json:"eInvoice"`
	RelatedEInvoiceCode   string         `json:"relatedEInvoiceCode,omitempty"`
	RelatedEInvoiceNumber string         `json:"relatedEInvoiceNumber,omitempty"`
	RelevantList          []RelevantData `json:"relevantList,omitempty"`
	EInvoiceAddress       string         `json:"eInvoiceAddress,omitempty"`
}

type PrintEInvoiceData struct {
	BaseData
	PaperCode   string `json:"paperCode"`
	PaperNumber string `json:"paperNumber"`
}

type AccreditEInvoiceData struct {
	BaseData
	AuthorizerName      string `json:"authorizerName"`
	AuthorizerCode      string `json:"authorizerCode"`
	AuthorizerCodeType  string `json:"authorizerCodeType"`
	TargetAccountId     string `json:"targetAccountId"`
	AuthorizedPartyName string `json:"authorizedPartyName"`
	AuthorizedPartyCode string `json:"authorizedPartyCode"`
	Ciphertext          string `json:"ciphertext"` // 增加的查询方密文
}

type AssignEInvoiceData struct {
	BaseData
	AssignorType       string `json:"assignorType"`
	AssignorName       string `json:"assignorName"`
	AssignorCode       string `json:"assignorCode"`
	AssignorCodeType   string `json:"assignorCodeType"`
	TargetAccountId    string `json:"targetAccountId"`
	TransfereeName     string `json:"transfereeName"`
	TransfereeCode     string `json:"transfereeCode"`
	TransfereeCodeType string `json:"transfereeCodeType"`
	Ciphertext         string `json:"ciphertext"` // 目标账户密文
}

type LockOrUnlockEInvoiceData struct {
	BaseData
	AgencyCode string `json:"agencyCode"`
	AgencyName string `json:"agencyName"`
}

type BookedEInvoiceData struct {
	BaseData
	AgencyCode  string `json:"agencyCode"`
	AgencyName  string `json:"agencyName"`
	AccNumber   string `json:"accNumber"`
	AccAmount   int64  `json:"accAmount"`
	AccDateTime string `json:"accDateTime"`
}

type SaveTemplateFileData struct {
	BaseData
	TemplateId   string `json:"templateId"`
	Template     []byte `json:"template"`
	TemplateType string `json:"templateType"`
}

type TemplateData struct {
	TemplateId   string     `json:"templateId"`
	TemplateList []Template `json:"templateList"`
}

type Template struct {
	Template     []byte `json:"template"`
	TemplateType string `json:"templateType"`
}

type EInvoiceData struct {
	EInvoiceInfo EInvoiceInfo `json:"eInvoiceInfo"`
	EInvoiceXML  string       `json:"eInvoiceXML"`
}

type TrajectoryData struct {
	Trajectories []Trajectory `json:"trajectories"`
}

type Trajectory struct {
	Hash             string `json:"hash"`
	TransactionBytes []byte `json:"transactionBytes"`
}

type BlockData struct {
	BlockList []Block `json:"blockList"`
}

type Block struct {
	BlockId    string `json:"blockId"`
	BlockBytes []byte `json:"blockBytes"`
}

type RelevantData struct {
	Identity string `json:"identity"`
	Type     int    `json:"type"`
	Value    string `json:"value,omitempty"`
}

func (t *Invoice) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *Invoice) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "queryEInvoice" {
		result, err = queryEInvoice(stub, args)
	} else if fn == "exist" {
		result, err = exist(stub, args)
	} else {
		if len(args) != 1 {
			return shim.Error(fmt.Errorf("Incorrect arguments").Error())
		}

		if fn == "saveNormalEInvoice" {
			result, err = saveNormalEInvoice(stub, args)
		} else if fn == "saveRedEInvoice" {
			result, err = saveRedEInvoice(stub, args)
		} else if fn == "printEInvoice" {
			result, err = printEInvoice(stub, args)
		} else if fn == "accreditEInvoice" {
			result, err = accreditEInvoice(stub, args)
		} else if fn == "assignEInvoice" {
			result, err = assignEInvoice(stub, args)
		} else if fn == "lockEInvoice" {
			result, err = lockEInvoice(stub, args)
		} else if fn == "unlockEInvoice" {
			result, err = unlockEInvoice(stub, args)
		} else if fn == "bookedEInvoice" {
			result, err = bookedEInvoice(stub, args)
		} else if fn == "unbookedEInvoice" {
			result, err = unbookedEInvoice(stub, args)
		} else {
			err = fmt.Errorf("Function is not exists: %s", fn)
		}
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

/**
 * 蓝票存证
 */
func saveNormalEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	return saveEInvoice(stub, args, false)
}

/**
 * 红票存证
 */
func saveRedEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	return saveEInvoice(stub, args, true)
}

/**
 * 打印
 */
func printEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	printEInvoiceData := PrintEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &printEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("PrintEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, printEInvoiceData.EInvoiceCode, printEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	flags := make([]bool, len(eInvoiceInfo.RecordInfos))
	flag := false
	for i, ri := range eInvoiceInfo.RecordInfos {
		if flags[i] {
			continue
		}
		flags[i] = true
		count := 1
		for j := i + 1; j < len(eInvoiceInfo.RecordInfos); j++ {
			r := eInvoiceInfo.RecordInfos[j]
			if r.AccNumber == ri.AccNumber && r.AccDateTime == ri.AccDateTime && r.AccAmount == ri.AccAmount {
				count++
				flags[j] = true
			}
		}
		for _, ui := range eInvoiceInfo.UnrecordInfos {
			if ui.AccNumber == ri.AccNumber && ui.AccDateTime == ri.AccDateTime && ui.AccAmount == ri.AccAmount {
				count--
			}
		}
		if count > 0 {
			flag = true
			break
		}
	}
	if flag {
		return "", fmt.Errorf("Has record: %s, %s", printEInvoiceData.PaperCode, printEInvoiceData.PaperNumber)
	}

	if eInvoiceInfo.Lock {
		return "", fmt.Errorf("Lock: %s, %s", printEInvoiceData.PaperCode, printEInvoiceData.PaperNumber)
	}

	paperInfo := PaperInfo{PaperCode: printEInvoiceData.PaperCode, PaperNumber: printEInvoiceData.PaperNumber}
	eInvoiceInfo.PaperInfos = append(eInvoiceInfo.PaperInfos, paperInfo)
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", printEInvoiceData.EInvoiceCode, printEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", printEInvoiceData.EInvoiceCode, printEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 授权
 */
func accreditEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	accreditEInvoiceData := AccreditEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &accreditEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("AccreditEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, accreditEInvoiceData.EInvoiceCode, accreditEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	if eInvoiceInfo.Lock {
		return "", fmt.Errorf("Lock: %s, %s", accreditEInvoiceData.EInvoiceCode, accreditEInvoiceData.EInvoiceNumber)
	}

	authorizedInfo := AuthorizedInfo{AuthorizerName: accreditEInvoiceData.AuthorizerName, AuthorizerCode: accreditEInvoiceData.AuthorizerCode, AuthorizerCodeType: accreditEInvoiceData.AuthorizerCodeType, TargetAccountId: accreditEInvoiceData.TargetAccountId, AuthorizedPartyCode: accreditEInvoiceData.AuthorizedPartyCode, AuthorizedPartyName: accreditEInvoiceData.AuthorizedPartyName}
	eInvoiceInfo.AuthorizedInfos = append(eInvoiceInfo.AuthorizedInfos, authorizedInfo)
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", accreditEInvoiceData.EInvoiceCode, accreditEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", accreditEInvoiceData.EInvoiceCode, accreditEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 转让
 */
func assignEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	assignEInvoiceData := AssignEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &assignEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("AssignEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, assignEInvoiceData.EInvoiceCode, assignEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	if eInvoiceInfo.Lock {
		return "", fmt.Errorf("Lock: %s, %s", assignEInvoiceData.EInvoiceCode, assignEInvoiceData.EInvoiceNumber)
	}

	assignInfo := AssignInfo{AssignorType: assignEInvoiceData.AssignorType, AssignorName: assignEInvoiceData.AssignorName, AssignorCode: assignEInvoiceData.AssignorCode, AssignorCodeType: assignEInvoiceData.AssignorCodeType, TargetAccountId: assignEInvoiceData.TargetAccountId, TransfereeCode: assignEInvoiceData.TransfereeCode, TransfereeName: assignEInvoiceData.TransfereeName, TransfereeCodeType: assignEInvoiceData.TransfereeCodeType}
	eInvoiceInfo.AssignInfos = append(eInvoiceInfo.AssignInfos, assignInfo)
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", assignEInvoiceData.EInvoiceCode, assignEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", assignEInvoiceData.EInvoiceCode, assignEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 锁定
 */
func lockEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	lockEInvoiceData := LockOrUnlockEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &lockEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("LockOrUnlockEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, lockEInvoiceData.EInvoiceCode, lockEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	if eInvoiceInfo.Lock {
		return "", fmt.Errorf("Has Locked: %s, %s", lockEInvoiceData.EInvoiceCode, lockEInvoiceData.EInvoiceNumber)
	}

	lockInfo := LockOrUnlockInfo{AgencyCode: lockEInvoiceData.AgencyCode, AgencyName: lockEInvoiceData.AgencyName}
	eInvoiceInfo.LockInfos = append(eInvoiceInfo.LockInfos, lockInfo)
	eInvoiceInfo.Lock = true
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", lockEInvoiceData.EInvoiceCode, lockEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", lockEInvoiceData.EInvoiceCode, lockEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 解锁
 */
func unlockEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	unlockEInvoiceData := LockOrUnlockEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &unlockEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("LockOrUnlockEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, unlockEInvoiceData.EInvoiceCode, unlockEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	if !eInvoiceInfo.Lock {
		return "", fmt.Errorf("Unloced: %s, %s", unlockEInvoiceData.EInvoiceCode, unlockEInvoiceData.EInvoiceNumber)
	}

	unlockInfo := LockOrUnlockInfo{AgencyCode: unlockEInvoiceData.AgencyCode, AgencyName: unlockEInvoiceData.AgencyName}
	eInvoiceInfo.UnlockInfos = append(eInvoiceInfo.LockInfos, unlockInfo)
	eInvoiceInfo.Lock = false
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", unlockEInvoiceData.EInvoiceCode, unlockEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", unlockEInvoiceData.EInvoiceCode, unlockEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 入账
 */
func bookedEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	bookedEInvoiceData := BookedEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &bookedEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("BookedEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	if eInvoiceInfo.PaperInfos != nil {
		return "", fmt.Errorf("Has printed: %s, %s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	}

	recordInfo := RecordInfo{AgencyCode: bookedEInvoiceData.AgencyCode, AgencyName: bookedEInvoiceData.AgencyName, AccNumber: bookedEInvoiceData.AccNumber, AccAmount: bookedEInvoiceData.AccAmount, AccDateTime: bookedEInvoiceData.AccDateTime}
	count := 0
	for _, ri := range eInvoiceInfo.RecordInfos {
		if recordInfo.AccNumber == ri.AccNumber && recordInfo.AccDateTime == ri.AccDateTime && recordInfo.AccAmount == ri.AccAmount {
			count++
		}
	}
	for _, ui := range eInvoiceInfo.UnrecordInfos {
		if ui.AccNumber == recordInfo.AccNumber && ui.AccDateTime == recordInfo.AccDateTime && ui.AccAmount == recordInfo.AccAmount {
			count--
		}
	}
	if count != 0 {
		return "", fmt.Errorf("Has recorded: %s, %s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	}

	eInvoiceInfo.RecordInfos = append(eInvoiceInfo.RecordInfos, recordInfo)
	if eInvoiceInfo.Lock {
		eInvoiceInfo.Lock = false
	}
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 入账撤销
 */
func unbookedEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	bookedEInvoiceData := BookedEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &bookedEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("BookedEInvoiceData error: %s", args[0])
	}
	eInvoiceData, err := getEInvoice(stub, bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	eInvoiceInfo := eInvoiceData.EInvoiceInfo

	unrecordInfo := RecordInfo{AgencyCode: bookedEInvoiceData.AgencyCode, AgencyName: bookedEInvoiceData.AgencyName, AccNumber: bookedEInvoiceData.AccNumber, AccAmount: bookedEInvoiceData.AccAmount, AccDateTime: bookedEInvoiceData.AccDateTime}
	count := 0
	for _, recordInfo := range eInvoiceInfo.RecordInfos {
		if recordInfo.AccNumber == unrecordInfo.AccNumber && recordInfo.AccDateTime == unrecordInfo.AccDateTime && recordInfo.AccAmount == unrecordInfo.AccAmount {
			count++
		}
	}
	for _, ui := range eInvoiceInfo.UnrecordInfos {
		if ui.AccNumber == unrecordInfo.AccNumber && ui.AccDateTime == unrecordInfo.AccDateTime && ui.AccAmount == unrecordInfo.AccAmount {
			count--
		}
	}
	if count != 1 {
		return "", fmt.Errorf("Has not recorded: %s, %s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	}

	eInvoiceInfo.UnrecordInfos = append(eInvoiceInfo.UnrecordInfos, unrecordInfo)
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", bookedEInvoiceData.EInvoiceCode, bookedEInvoiceData.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 票据查询
 * args[0] eInvoiceCode
 * args[1] eInvoiceNumber
 */
func queryEInvoice(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments")
	}

	return getEInvoiceJson(stub, args[0], args[1])
}

/**
 * 查询票据是否存在
 * args[0] eInvoiceCode
 * args[1] eInvoiceNumber
 */
func exist(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments")
	}
	value, err := stub.GetState(fmt.Sprintf("eInvoiceInfo_%s_%s", args[0], args[1]))
	if err != nil {
		return "", fmt.Errorf("GetState error: eInvoiceInfo_%s_%s", args[0], args[1])
	}
	if value == nil {
		return "false", nil
	}
	return "true", nil
}

/**
 * 票据存证
 */
func saveEInvoice(stub shim.ChaincodeStubInterface, args []string, isRed bool) (string, error) {
	saveEInvoiceData := SaveEInvoiceData{}
	err := json.Unmarshal([]byte(args[0]), &saveEInvoiceData)
	if err != nil {
		return "", fmt.Errorf("SaveEInvoiceData error: %s", args[0])
	}
	value, err := stub.GetState(fmt.Sprintf("eInvoiceInfo_%s_%s", saveEInvoiceData.EInvoiceCode, saveEInvoiceData.EInvoiceNumber))
	if err != nil {
		return "", fmt.Errorf("Get state error: eInvoiceInfo_%s_%s", saveEInvoiceData.EInvoiceCode, saveEInvoiceData.EInvoiceNumber)
	}
	if value != nil {
		return "", fmt.Errorf("EInvoice is exists: %s, %s", saveEInvoiceData.EInvoiceCode, saveEInvoiceData.EInvoiceNumber)
	}

	eInvoiceInfo := EInvoiceInfo{EInvoiceCode: saveEInvoiceData.EInvoiceCode, EInvoiceNumber: saveEInvoiceData.EInvoiceNumber, EInvoiceAddress: saveEInvoiceData.EInvoiceAddress, IsRed: false, Lock: false}
	if isRed {
		relatedEInvoiceData, err := getEInvoice(stub, saveEInvoiceData.RelatedEInvoiceCode, saveEInvoiceData.RelatedEInvoiceNumber)
		if err != nil {
			return "", err
		}
		relatedEInvoiceInfo := relatedEInvoiceData.EInvoiceInfo
		relatedEInvoiceInfo.RelatedEInvoiceCode = saveEInvoiceData.EInvoiceCode
		relatedEInvoiceInfo.RelatedEInvoiceNumber = saveEInvoiceData.EInvoiceNumber

		jsonBytes, err := json.Marshal(relatedEInvoiceInfo)
		if err != nil {
			return "", fmt.Errorf("Marshal error: %s", relatedEInvoiceInfo)
		}
		err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", relatedEInvoiceInfo.EInvoiceCode, relatedEInvoiceInfo.EInvoiceNumber), jsonBytes)
		if err != nil {
			return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", relatedEInvoiceInfo.EInvoiceCode, relatedEInvoiceInfo.EInvoiceNumber)
		}

		eInvoiceInfo.RelatedEInvoiceCode = saveEInvoiceData.RelatedEInvoiceCode
		eInvoiceInfo.RelatedEInvoiceNumber = saveEInvoiceData.RelatedEInvoiceNumber
		eInvoiceInfo.IsRed = true
	}
	jsonBytes, err := json.Marshal(eInvoiceInfo)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceInfo)
	}
	err = stub.PutState(fmt.Sprintf("eInvoiceInfo_%s_%s", eInvoiceInfo.EInvoiceCode, eInvoiceInfo.EInvoiceNumber), jsonBytes)
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceInfo_%s_%s", eInvoiceInfo.EInvoiceCode, eInvoiceInfo.EInvoiceNumber)
	}

	err = stub.PutState(fmt.Sprintf("eInvoiceXML_%s_%s", eInvoiceInfo.EInvoiceCode, eInvoiceInfo.EInvoiceNumber), []byte(saveEInvoiceData.EInvoice))
	if err != nil {
		return "", fmt.Errorf("PutState error: eInvoiceXML_%s_%s", eInvoiceInfo.EInvoiceCode, eInvoiceInfo.EInvoiceNumber)
	}

	return "", nil
}

/**
 * 查询票据信息
 * 返回类型为Json
 */
func getEInvoiceJson(stub shim.ChaincodeStubInterface, eInvoiceCode string, eInvoiceNumber string) (string, error) {
	eInvoiceData, err := getEInvoice(stub, eInvoiceCode, eInvoiceNumber)
	if err != nil {
		return "", err
	}
	jsonBytes, err := json.Marshal(eInvoiceData)
	if err != nil {
		return "", fmt.Errorf("Marshal error: %s", eInvoiceData)
	}
	return string(jsonBytes), nil
}

/**
 * 查询票据信息
 */
func getEInvoice(stub shim.ChaincodeStubInterface, eInvoiceCode string, eInvoiceNumber string) (EInvoiceData, error) {
	eInvoiceData := EInvoiceData{}

	value, err := stub.GetState(fmt.Sprintf("eInvoiceInfo_%s_%s", eInvoiceCode, eInvoiceNumber))
	if err != nil {
		return eInvoiceData, fmt.Errorf("GetState error: eInvoiceInfo_%s_%s", eInvoiceCode, eInvoiceNumber)
	}
	if value == nil {
		return eInvoiceData, fmt.Errorf("Value is nil: eInvoiceInfo_%s_%s", eInvoiceCode, eInvoiceNumber)
	}
	eInvoiceInfo := EInvoiceInfo{}
	err = json.Unmarshal(value, &eInvoiceInfo)
	if err != nil {
		return eInvoiceData, fmt.Errorf("EInvoiceInfo error: %s", string(value))
	}
	eInvoiceData.EInvoiceInfo = eInvoiceInfo

	value, err = stub.GetState(fmt.Sprintf("eInvoiceXML_%s_%s", eInvoiceCode, eInvoiceNumber))
	if err != nil {
		return eInvoiceData, fmt.Errorf("GetState error: eInvoiceXML_%s_%s", eInvoiceCode, eInvoiceNumber)
	}
	if value == nil {
		return eInvoiceData, fmt.Errorf("Value is nil: eInvoiceXML_%s_%s", eInvoiceCode, eInvoiceNumber)
	}
	eInvoiceData.EInvoiceXML = string(value)

	return eInvoiceData, nil
}

func main() {
	if err := shim.Start(new(Invoice)); err != nil {
		fmt.Printf("Error starting Invoice chaincode: %s", err)
	}
}
