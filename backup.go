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
//	"github.com/gagliardetto/solana-go/text"
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
//	accountFrom := solana.MustPrivateKeyFromBase58(privateKeyBase58)
//
//	rpcClient := rpc.New(rpc.DevNet_RPC)
//	wsClient, err := ws.Connect(context.Background(), "wss://devnet.sonic.game/")
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
//	fmt.Printf("Balance: %.9f SOL\n", float64(balance)/1_000_000_000)
//
//	//read address nya
//	addresses, err := readAddresses("address.txt")
//	if err != nil {
//		log.Fatalf("Failed to read addresses: %v", err)
//	}
//
//	// ensure enough balance
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
//	// Send 0.001 SOL to each address with delay
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
//		spew.Dump(tx)
//		//print it !
//		tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))
//
//		sig, err := confirm.SendAndConfirmTransaction(
//			context.TODO(),
//			rpcClient,
//			wsClient,
//			tx,
//		)
//		if err != nil {
//			log.Printf("Failed to send transaction to %s: %v", address, err)
//		} else {
//			fmt.Printf("Sent %d lamports to %s\n", solAmount, address)
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
