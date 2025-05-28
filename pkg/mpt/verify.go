package mpt

import (
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/trie"
	"log"
)

// VerifyProof checks a Merkle-Patricia proof under the given root.
// root should be the stateRoot from the block header.
// key is the path (e.g., keccak256(accountAddress.Bytes())).
// proofNodes is the slice of RLP-encoded trie node bytes forming the proof.
// Returns the RLP-encoded account value (nonce, balance, storageHash, codeHash) or an error.
func VerifyProof(root common.Hash, key []byte, proofNodes [][]byte) ([]byte, error) {

    log.Printf("Verifying proof for key: %x", key)
    return trie.VerifyProof(root, key, proofNodes)
}