package main

import (
    "encoding/hex"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "strings"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/rose2221/stateless-spv/pkg/header"
    "github.com/rose2221/stateless-spv/pkg/mpt"
)

// proofFile represents the minimal subset of eth_getProof output needed.
type proofFile struct {
    Address      string   `json:"address"`
    AccountProof []string `json:"accountProof"`
}

func main() {
    // Flags
    hdrPath := flag.String("header", "", "path to hex-encoded RLP header file")
    proofPath := flag.String("proof", "", "path to proof JSON file (eth_getProof result)")
    flag.Parse()

    if *hdrPath == "" || *proofPath == "" {
        fmt.Fprintln(os.Stderr, "Usage: spv-verify --header <file> --proof <file>")
        os.Exit(1)
    }

  
    h, _, err := header.DecodeHeaderFromHexFile(*hdrPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error decoding header: %v\n", err)
        os.Exit(1)
    }

   
    data, err := ioutil.ReadFile(*proofPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading proof file: %v\n", err)
        os.Exit(1)
    }
    var pf proofFile
    if err := json.Unmarshal(data, &pf); err != nil {
        fmt.Fprintf(os.Stderr, "Error unmarshaling proof JSON: %v\n", err)
        os.Exit(1)
    }

    nodes := make([][]byte, len(pf.AccountProof))
    for i, hexNode := range pf.AccountProof {
        hexNode = strings.TrimPrefix(hexNode, "0x")
        b, err := hex.DecodeString(hexNode)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Invalid proof node at index %d: %v\n", i, err)
            os.Exit(1)
        }
        nodes[i] = b
    }

    addr := common.HexToAddress(pf.Address)
    key := crypto.Keccak256(addr.Bytes())

    root := h.Root
    if _, err := mpt.VerifyProof(root, key, nodes); err != nil {
        fmt.Fprintf(os.Stderr, "Proof invalid: %v\n", err)
        os.Exit(1)
    }

    fmt.Println(" Proof valid!")
    os.Exit(0)
}
