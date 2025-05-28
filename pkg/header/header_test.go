package header

import (
    "bytes"
    "crypto/sha3"
    "encoding/hex"
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "testing"

    "github.com/ethereum/go-ethereum/rlp"
)

// TestDecodeAndHash ensures that decoding and re-hashing reproduces the block's hash.
func TestDecodeAndHash(t *testing.T) {
  
    headerPath := filepath.Join("..", "fixtures", "header.txt")
    h, rlpBytes, err := DecodeHeaderFromHexFile(headerPath)
    if err != nil {
        t.Fatalf("failed to decode header: %v", err)
    }

   
    encoded, err := rlp.EncodeToBytes(h)
    if err != nil {
        t.Fatalf("failed to re-encode header: %v", err)
    }
    if !bytes.Equal(encoded, rlpBytes) {
        t.Error("re-encoded RLP differs from original header bytes")
    }

  
    hasher := sha3.NewLegacyKeccak256()
    hasher.Write(rlpBytes)
    gotHash := hasher.Sum(nil)

    blockJSONPath := filepath.Join("..", "fixtures", "block.json")
    blockData, err := os.ReadFile(blockJSONPath)
    if err != nil {
        t.Fatalf("failed to read block.json: %v", err)
    }
    var blk struct { Hash string `json:"hash"` }
    if err := json.Unmarshal(blockData, &blk); err != nil {
        t.Fatalf("failed to unmarshal block.json: %v", err)
    }
    expectedHashBytes, err := hex.DecodeString(strings.TrimPrefix(blk.Hash, "0x"))
    if err != nil {
        t.Fatalf("invalid expected hash hex: %v", err)
    }

   
    if !bytes.Equal(gotHash, expectedHashBytes) {
        t.Errorf("hash mismatch:\n got:      0x%x\n expected: %s", gotHash, blk.Hash)
    }
}
