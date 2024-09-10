package main

import (
	"crypto/md5"    // Importing the MD5 hash function from the crypto package.
	"crypto/sha256" // Importing the SHA256 hash function from the crypto package.
	"encoding/hex"  // Package hex provides functions for encoding and decoding binary data in hexadecimal.
	"encoding/json" // Package json implements encoding and decoding of JSON.
	"fmt"           // Package fmt implements formatted I/O with functions analogous to C's printf and scanf.
	"io"            // Package io provides basic interfaces to I/O primitives.
	"log"           // Package log provides a simple logging package.
	"net/http"      // Package http provides HTTP client and server implementations.

	"github.com/gorilla/mux" // Mux is a powerful URL router and dispatcher.
)

// Block represents each block in the blockchain.
type Block struct {
	Index     int         `json:"index"` // Position of the data record in the blockchain.
	Timestamp string      `json:"timestamp"` // Auto-generated timestamp of when the block was created.
	Data      BookCheckout `json:"data"` // Embedded struct that contains the block data.
	Hash      string      `json:"hash"` // SHA256 hash of the block's contents.
	PrevHash  string      `json:"prev_hash"` // Hash of the previous block in the chain to ensure integrity.
}

// BookCheckout defines the data for a library book checkout system.
type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckOutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"` // Flag to identify the genesis (first) block.
}

// Book represents a book in the library system.
type Book struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Author       string `json:"author"`
	PublishDate  string `json:"publish_date"`
	ISBN         string `json:"isbn"`
}

// Blockchain is a series of validated blocks.
type Blockchain struct {
	blocks []*Block // Slice of pointers to Block, representing the chain.
}

var BlockChain *Blockchain // Global variable that will contain the pointer to the blockchain.

// newBook handles creating new book records via the HTTP API.
func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not create book", err)
		w.Write([]byte("Could not create book"))
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate) // Concatenate ISBN and PublishDate to form a unique ID.
	book.ID = fmt.Sprintf("%x", h.Sum(nil)) // Convert binary data to a hex string.

	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not marshal book", err)
		w.Write([]byte("Could not marshal book"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// generateHash creates a hash for a block.
func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)
	data := fmt.Sprintf("%d%s%s%s", b.Index, b.Timestamp, string(bytes), b.PrevHash)
	hash := sha256.New()
	hash.Write([]byte(data))
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

// CreateBlock creates a new block using previous block's hash.
func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	block := &Block{
		Index:     prevBlock.Index + 1,
		Timestamp: "now", // Placeholder timestamp. Consider using a real timestamp here.
		Data:      data,
		PrevHash:  prevBlock.Hash,
	}
	block.generateHash() // Automatically generate hash during block creation.
	return block
}

// validBlock checks if a block is valid by comparing the index and hashes.
func validBlock(block, prevBlock *Block) bool {
	return prevBlock.Index+1 == block.Index &&
		prevBlock.Hash == block.PrevHash &&
		block.validateHash(block.Hash)
}

// validateHash checks if a block's hash is valid.
func (b *Block) validateHash(hash string) bool {
	b.generateHash()
	return b.Hash == hash
}

// AddBlock adds a block to the blockchain after validation.
func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevBlock, data)
	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

// writeBlock writes a new block to the blockchain via HTTP POST.
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not decode request", err)
		w.Write([]byte("Could not decode request"))
		return
	}
	BlockChain.AddBlock(checkoutItem)
}

// GenesisBlock creates the first block in the blockchain.
func GenesisBlock() *Block {
	return CreateBlock(&Block{}, BookCheckout{IsGenesis: true})
}

// NewBlockChain initializes a new blockchain with a genesis block.
func NewBlockChain() *Blockchain {
	return &Blockchain{blocks: []*Block{GenesisBlock()}}
}

// getBlockchain serves the entire blockchain as JSON via HTTP GET.
func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbytes, err := json.MarshalIndent(BlockChain.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}
	io.WriteString(w, string(jbytes))
}

// main sets up the HTTP server and routes.
func main() {
	BlockChain = NewBlockChain()
	router := mux.NewRouter()
	router.HandleFunc("/", getBlockchain).Methods("GET")
	router.HandleFunc("/", writeBlock).Methods("POST")
	router.HandleFunc("/new", newBook).Methods("POST")

	go func() {
		for _, block := range BlockChain.blocks {
			bytes, _ := json.MarshalIndent(block.Data, "", " ")
			fmt.Printf("Index: %d\n", block.Index)
			fmt.Printf("Timestamp: %s\n", block.Timestamp)
			fmt.Printf("Data: %s\n", string(bytes))
			fmt.Printf("Hash: %s\n", block.Hash)
			fmt.Printf("PrevHash: %s\n", block.PrevHash)
			fmt.Println()
		}
	}()

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
