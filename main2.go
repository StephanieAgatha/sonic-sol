package main

//
//import (
//	"bufio"
//	"context"
//	"fmt"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//
//	"github.com/davecgh/go-spew/spew"
//	"github.com/gagliardetto/solana-go"
//	"github.com/gagliardetto/solana-go/programs/system"
//	"github.com/gagliardetto/solana-go/rpc"
//	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
//	"github.com/gagliardetto/solana-go/rpc/ws"
//	"github.com/joho/godotenv"
//)
//
//const solAmount = uint64(1000000) // 0.001 SOL in lamports
//
//func main() {
//	err := godotenv.Load()
//	if err != nil {
//		log.Fatalf("Error loading .env file")
//	}
//
//	privateKeyBase58 := os.Getenv("PRIVATE_KEY")
//	if privateKeyBase58 == "" {
//		log.Fatalf("Private key is not set in .env file")
//	}
//
//	// Generate keypair
//	accountFrom := solana.MustPrivateKeyFromBase58(privateKeyBase58)
//
//	rpcClient := rpc.New("https://devnet.sonic.game")
//	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
//	if err != nil {
//		panic(err)
//	}
//
//	// Check balance
//	balanceResult, err := rpcClient.GetBalance(
//		context.TODO(),
//		accountFrom.PublicKey(),
//		rpc.CommitmentFinalized,
//	)
//	if err != nil {
//		log.Fatalf("Failed to get balance: %v", err)
//	}
//
//	balance := balanceResult.Value
//	if balance == 0 {
//		log.Fatalf("No balance available")
//	}
//
//	// print balance
//	fmt.Printf("Balance: %.9f SOL\n", float64(balance)/1_000_000_000)
//
//	// read
//	addresses, err := readAddresses("address.txt")
//	if err != nil {
//		log.Fatalf("Failed to read address file: %v", err)
//	}
//
//	// check if address have enough balance -_-
//	requiredBalance := solAmount * uint64(len(addresses))
//	if balance < requiredBalance {
//		log.Fatalf("Insufficient balance. Required: %d, Available: %d", requiredBalance, balance)
//	}
//
//	// Get delay
//	fmt.Print("Masukkan delay (dalam detik): ")
//	reader := bufio.NewReader(os.Stdin)
//	delayInput, _ := reader.ReadString('\n')
//	delayInput = strings.TrimSpace(delayInput)
//	delay, err := strconv.Atoi(delayInput)
//	if err != nil {
//		log.Fatalf("Invalid delay input: %v", err)
//	}
//
//	// Send 0.001 SOL to each address
//	for _, addressStr := range addresses {
//		address := solana.MustPublicKeyFromBase58(addressStr)
//
//		recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
//		if err != nil {
//			panic(err)
//		}
//
//		tx, err := solana.NewTransaction(
//			[]solana.Instruction{
//				system.NewTransferInstruction(
//					solAmount, // 0.001 SOL in lamports
//					accountFrom.PublicKey(),
//					address,
//				).Build(),
//			},
//			recent.Value.Blockhash,
//			solana.TransactionPayer(accountFrom.PublicKey()),
//		)
//		if err != nil {
//			panic(err)
//		}
//
//		_, err = tx.Sign(
//			func(key solana.PublicKey) *solana.PrivateKey {
//				if accountFrom.PublicKey().Equals(key) {
//					return &accountFrom
//				}
//				return nil
//			},
//		)
//		if err != nil {
//			panic(fmt.Errorf("unable to sign transaction: %w", err))
//		}
//
//		// print it
//		fmt.Printf("Sending 0.001 SOL to %s, waiting for confirmation...\n", address)
//
//		// Send transaction, and wait for confirmation
//		sig, err := confirm.SendAndConfirmTransaction(
//			context.TODO(),
//			rpcClient,
//			wsClient,
//			tx,
//		)
//		if err != nil {
//			log.Printf("Failed to send transaction to %s: %v", address, err)
//		} else {
//			fmt.Printf("Success sending to %s with signature %s\n", address, sig)
//		}
//		spew.Dump(sig)
//
//		// Delay
//		time.Sleep(time.Duration(delay) * time.Second)
//	}
//}
//
//func readAddresses(filename string) ([]string, error) {
//	file, err := os.Open(filename)
//	if err != nil {
//		return nil, err
//	}
//	defer file.Close()
//
//	var addresses []string
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		address := strings.TrimSpace(scanner.Text())
//		if address != "" {
//			addresses = append(addresses, address)
//		}
//	}
//	if err := scanner.Err(); err != nil {
//		return nil, err
//	}
//	return addresses, nil
//}
