package mpt

import (
    "encoding/hex"
    "encoding/json"
    "os"
    "strings"
    "testing"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
)

// TestVerifyAccountProof loads fixtures and ensures the MPT proof is valid.
func TestVerifyAccountProof(t *testing.T) {
 
    blockData, err := os.ReadFile("fixtures/block.json")
    if err != nil {
        t.Fatalf("failed to read block.json: %v", err)
    }
    var blk struct { StateRoot string `json:"stateRoot"` }
    if err := json.Unmarshal(blockData, &blk); err != nil {
        t.Fatalf("failed to unmarshal block.json: %v", err)
    }
    rootBytes, err := hex.DecodeString(strings.TrimPrefix(blk.StateRoot, "0x"))
    if err != nil {
        t.Fatalf("failed to decode stateRoot: %v", err)
    }
    var rootHash common.Hash
    copy(rootHash[:], rootBytes)

  
    proofData, err := os.ReadFile("fixtures/proof.json")
    if err != nil {
        t.Fatalf("failed to read proof.json: %v", err)
    }
    var pf struct {
        Address      string   `json:"address"`
        AccountProof []string `json:"accountProof"`
    }
    if err := json.Unmarshal(proofData, &pf); err != nil {
        t.Fatalf("failed to unmarshal proof.json: %v", err)
    }

 
    nodes := make([][]byte, len(pf.AccountProof))
    for i, n := range pf.AccountProof {
        b, err := hex.DecodeString(strings.TrimPrefix(n, "0x"))
        if err != nil {
            t.Fatalf("failed to decode proof node %d: %v", i, err)
        }
        nodes[i] = b
    }

    
    addr := common.HexToAddress(pf.Address)
    key := crypto.Keccak256(addr.Bytes())

   
    val, err := VerifyProof(rootHash, key, nodes)
    if err != nil {
        t.Fatalf("proof verification failed: %v", err)
    }
    if len(val) == 0 {
        t.Fatal("no account value returned by proof")
    }
}
