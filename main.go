package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/joho/godotenv"
)

const solAmount = uint64(1000000) // 0.001 SOL in lamports

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	privateKeyBase58 := os.Getenv("PRIVATE_KEY")
	if privateKeyBase58 == "" {
		log.Fatalf("Private key is not set in .env file")
	}

	// Generate keypair
	accountFrom := solana.MustPrivateKeyFromBase58(privateKeyBase58)

	rpcClient := rpc.New("https://devnet.sonic.game")
	wsClient, err := ws.Connect(context.Background(), rpc.DevNet_WS)
	if err != nil {
		panic(err)
	}

	// Check balance
	balanceResult, err := rpcClient.GetBalance(
		context.TODO(),
		accountFrom.PublicKey(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	balance := balanceResult.Value
	if balance == 0 {
		log.Fatalf("No balance available")
	}

	// Print balance
	fmt.Printf("Balance: %.9f SOL\n", float64(balance)/1_000_000_000)

	// Get number of addresses to generate
	fmt.Print("How many addresses do you want to generate: ")
	reader := bufio.NewReader(os.Stdin)
	addressCountInput, _ := reader.ReadString('\n')
	addressCountInput = strings.TrimSpace(addressCountInput)
	addressCount, err := strconv.Atoi(addressCountInput)
	if err != nil {
		log.Fatalf("Invalid number of addresses: %v", err)
	}

	// Ensure enough balance to send to all addresses
	requiredBalance := solAmount * uint64(addressCount)
	if balance < requiredBalance {
		log.Fatalf("Insufficient balance. Required: %d, Available: %d", requiredBalance, balance)
	}

	// Get delay
	fmt.Print("Input delay (dalam detik): ")
	delayInput, _ := reader.ReadString('\n')
	delayInput = strings.TrimSpace(delayInput)
	delay, err := strconv.Atoi(delayInput)
	if err != nil {
		log.Fatalf("Invalid delay input: %v", err)
	}

	// Get Authorization value from user
	fmt.Print("Enter Authorization key (or press enter to skip): ")
	authKey, _ := reader.ReadString('\n')
	authKey = strings.TrimSpace(authKey)

	// Generate random addresses
	var addresses []solana.PublicKey
	for i := 0; i < addressCount; i++ {
		newKeypair := generateRandomKeypair()
		addresses = append(addresses, newKeypair.PublicKey())
		fmt.Printf("Generated address %d: %s\n", i+1, newKeypair.PublicKey())
	}

	// Send 0.001 SOL to each address
	for _, address := range addresses {
		recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
		if err != nil {
			panic(err)
		}

		tx, err := solana.NewTransaction(
			[]solana.Instruction{
				system.NewTransferInstruction(
					solAmount, // 0.001 SOL in lamports
					accountFrom.PublicKey(),
					address,
				).Build(),
			},
			recent.Value.Blockhash,
			solana.TransactionPayer(accountFrom.PublicKey()),
		)
		if err != nil {
			panic(err)
		}

		_, err = tx.Sign(
			func(key solana.PublicKey) *solana.PrivateKey {
				if accountFrom.PublicKey().Equals(key) {
					return &accountFrom
				}
				return nil
			},
		)
		if err != nil {
			panic(fmt.Errorf("unable to sign transaction: %w", err))
		}

		// Print sending status
		fmt.Printf("Sending 0.001 SOL to %s, waiting for confirmation...\n", address)

		// Send transaction, and wait for confirmation
		sig, err := confirm.SendAndConfirmTransaction(
			context.TODO(),
			rpcClient,
			wsClient,
			tx,
		)
		if err != nil {
			log.Printf("Failed to send transaction to %s: %v", address, err)
		} else {
			fmt.Printf("Success sending to %s with signature %s\n", address, sig)
		}
		spew.Dump(sig)

		//if user input a jwt token then fetch tx total
		if authKey != "" {
			fmt.Println("==================")
			fmt.Println("Fetch transaction from sonic server ...")
			getTxMilestone(authKey)
			fmt.Println("==================")
		}

		// Delay
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func generateRandomKeypair() solana.PrivateKey {
	keypair, err := solana.NewRandomPrivateKey()
	if err != nil {
		log.Fatalf("Failed to generate random keypair: %v", err)
	}
	return keypair
}

func readAddresses(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var addresses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		address := strings.TrimSpace(scanner.Text())
		if address != "" {
			addresses = append(addresses, address)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return addresses, nil
}

func getTxMilestone(authKey string) {
	url := "https://odyssey-api.sonic.game/user/transactions/state/daily"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", authKey)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Failed to unmarshal response:", err)
		return
	}

	if data, ok := result["data"].(map[string]interface{}); ok {
		if totalTransactions, ok := data["total_transactions"].(float64); ok {
			fmt.Printf("Total transactions: %.0f\n", totalTransactions)
		} else {
			fmt.Println("total_transactions not found ", err.Error())
		}
	} else {
		fmt.Println("failed to fetch data", err.Error())
	}
}
