package header

import (
    "encoding/hex"
    "os"
    "strings"

    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/rlp"
)

// Header represents an Ethereum block header (15 fields per Yellow Paper).
type Header struct {
    ParentHash  common.Hash  
    UncleHash   common.Hash  
    Coinbase    common.Address
    Root        common.Hash   
    TxHash      common.Hash   
    ReceiptHash common.Hash   
    Bloom       [256]byte      
    Difficulty  *big.Int       
    Number      *big.Int       
    GasLimit    uint64         
    GasUsed     uint64         
    Time        uint64       
    Extra       []byte       
    MixDigest   common.Hash    
    Nonce       uint64         
}

// DecodeHeaderFromHexFile reads a hex-encoded RLP header file and decodes it into Header.
// It returns the decoded Header, the raw RLP bytes, or an error.
func DecodeHeaderFromHexFile(path string) (*Header, []byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, nil, err
    }
    
    hexData := strings.TrimSpace(string(data))
    hexData = strings.TrimPrefix(hexData, "0x")

    rlpBytes, err := hex.DecodeString(hexData)
    if err != nil {
        return nil, nil, err
    }

    var h Header
    if err := rlp.DecodeBytes(rlpBytes, &h); err != nil {
        return nil, nil, err
    }

    return &h, rlpBytes, nil
}