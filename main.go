package main

import (
	"bufio"
	"fmt"
	"keyvaluedb/domain"
	"keyvaluedb/storage"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {

	// Handle interrupt signal
	handleInterruptSignal()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	dbCntStr := os.Getenv("DB_COUNT")

	storage := storage.NewInMemory(dbCntStr)
	kvdb := domain.NewKeyValueDB(storage)

	// Start TCP server
	listener, err := startTcpServer(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v\n", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		// Handle connection in a separate goroutine
		go handleConnection(conn, kvdb)
	}
}

func handleInterruptSignal() {
	// Create an interrupt channel to listen for the interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("Interrupt signal received. Gracefully stopping...")

		os.Exit(0)
	}()
}

func startTcpServer(port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	fmt.Println("TCP server started. Listening on port", port)
	return listener, nil
}

func handleConnection(conn net.Conn, kvdb domain.KeyValueDB) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	dbIndex := 0
	for {
		if dbIndex > 0 {
			fmt.Fprintf(writer, "[%d]$", dbIndex)
		} else {
			fmt.Fprintf(writer, "$")
		}
		writer.Flush()

		// Read client input
		command, err := readCommand(reader)
		if err != nil {
			printResult(writer, err)
			break
		}

		var result interface{}
		dbIndex, result = kvdb.Execute(dbIndex, command)
		printResult(writer, result)
	}
}

func readCommand(reader *bufio.Reader) (domain.Command, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return domain.Command{}, err
	}
	// Trim any leading/trailing whitespace and newline characters
	line = strings.TrimSpace(line)

	// Split the line into words
	words := strings.Split(line, " ")

	// Separate the number of words within the command including double quotes
	args := make([]interface{}, 0, 10)
	count := 0
	inQuotes := false
	isQuoteCompletes := true
	word := ""
	for _, tempWord := range words {
		if strings.HasPrefix(tempWord, `"`) {
			inQuotes = true
			isQuoteCompletes = false
		}

		if inQuotes {
			word = fmt.Sprintf("%s %s", word, tempWord)

			if strings.HasSuffix(tempWord, `"`) {
				args = append(args, strings.ReplaceAll(strings.TrimSpace(word), `"`, ""))
				word = ""
				inQuotes = false
				count++
				isQuoteCompletes = true
			}
		} else {
			args = append(args, strings.ReplaceAll(tempWord, `"`, ""))
			count++
		}
	}

	if count < 1 {
		return domain.Command{}, fmt.Errorf("invalid command")
	}

	if !isQuoteCompletes {
		return domain.Command{}, fmt.Errorf("(error) ERR Protocol error: unbalanced quotes in request")
	}

	cmd := strings.ToUpper(strings.TrimSpace(args[0].(string)))

	command := domain.NewCommand(cmd, args[1:]...)
	return command, nil
}

func printResult(writer *bufio.Writer, result interface{}) {
	switch res := result.(type) {
	case []interface{}:
		for i, item := range res {
			fmt.Fprintf(writer, "%d) %v\n", i+1, item)
		}
	default:
		fmt.Fprintf(writer, "%v\n", result)
	}
	writer.Flush()
}
