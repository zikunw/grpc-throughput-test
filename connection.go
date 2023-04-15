package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/tjarratt/babble"
	message "github.com/zikunw/grpc-throughput-test/message"
	"google.golang.org/grpc"
)

// ==========================
// Connection benchmark
// ==========================

type server struct {
	message.UnimplementedMessageServer

	msgCount int
	mu       sync.Mutex
}

func createServer() {
	// Establish server
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Not able to establish tcp server")
	}

	s := grpc.NewServer()
	message.RegisterMessageServer(s, &server{msgCount: 0})
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func (s *server) Send(ctx context.Context, req *message.MessageRequest) (*message.MessageResponse, error) {
	return &message.MessageResponse{Message: req.Message}, nil
}

func (s *server) Stream(stream message.Message_StreamServer) error {
	counter := 0
	for {
		_, err := stream.Recv()
		if err != nil {
			return err
		}
		counter++
		// if counter%1000 == 0 {
		// 	log.Printf("Received %d messages", counter)
		// }
		stream.Send(&message.MessageResponse{Message: "1"})
	}
}

func (s *server) SendRepeated(ctx context.Context, req *message.RepeatedMessageRequest) (*message.MessageResponse, error) {
	// Recieve the messages
	// for _, msg := range req.Messages {
	// 	log.Printf("Received: %s", msg.Message)
	// }
	// record the length
	s.mu.Lock()
	s.msgCount += len(req.Messages)
	//log.Println("Server received", s.msgCount, "messages")
	s.mu.Unlock()

	return &message.MessageResponse{Message: "1"}, nil
}

func createClient() {
	// Establish client
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Not able to establish tcp client")
	}
	defer conn.Close()

	client := message.NewMessageClient(conn)

	// Init babbler
	babbler := babble.NewBabbler()
	babbler.Count = 1

	// record time
	start := time.Now()
	// Send message
	for i := 0; i < WORDS; i++ {
		_, err := client.Send(context.Background(), &message.MessageRequest{Message: "Hello, World!"})
		if err != nil {
			log.Fatalln("Not able to send message")
		}

		if i%10000 == 0 {
			log.Printf("Sent %d messages", i)
		}
	}

	// log time
	log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
}

func createStreamClient() {
	// Establish client
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Not able to establish tcp client")
	}
	defer conn.Close()

	client := message.NewMessageClient(conn)

	// Init babbler
	babbler := babble.NewBabbler()
	babbler.Count = 1

	// record time
	start := time.Now()

	// Send message
	stream, err := client.Stream(context.Background())
	if err != nil {
		log.Fatalln("Not able to create stream")
	}

	for i := 0; i < WORDS; i++ {
		err := stream.Send(&message.MessageRequest{Message: "Hello, World!"})
		if err != nil {
			log.Fatalln("Not able to send message")
		}

		if i%10000 == 0 {
			log.Printf("Sent %d messages", i)
		}

		// Receive response
		_, err = stream.Recv()
		if err != nil {
			log.Fatalln("Not able to receive message")
		}
	}

	// log time
	log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
}

func createRepeatedClient() {
	BATCH_SIZE := 10000
	// Establish client
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Not able to establish tcp client")
	}
	defer conn.Close()

	client := message.NewMessageClient(conn)

	// Init babbler
	babbler := babble.NewBabbler()
	babbler.Count = 1

	// record time
	start := time.Now()

	// Create messages
	sent := 0
	for sent < WORDS {
		// Create batch
		batch := make([]*message.MessageRequest, 0)
		for i := 0; i < BATCH_SIZE; i++ {
			batch = append(batch, &message.MessageRequest{Message: "Hello, World!"})
		}

		// Send batch
		_, err := client.SendRepeated(context.Background(), &message.RepeatedMessageRequest{Messages: batch})
		if err != nil {
			log.Fatalln("Not able to send message")
		}

		sent += BATCH_SIZE
	}

	// log time
	log.Printf("Took %v to send %d messages", time.Since(start), WORDS)

}

// ==========================
// Singular TCP connections
// ==========================

func createOneWriteTCPserver() {
	// Establish server
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Not able to establish tcp server")
	}

	log.Printf("server listening at %v", lis.Addr())
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Fatalln("Not able to accept connection")
			}

			go handleOneWriteTCPRequest(conn)
		}
	}()
}

// Handle
func handleOneWriteTCPRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2)

	for {
		conn.SetDeadline(time.Now().Add(2 * time.Second))
		_, err := conn.Read(buf)
		if err != nil {
			log.Println("Not able to read message")
		}
		// confirm message
		_, err = conn.Write([]byte("1"))
	}
}

// Multiple writes in one TCP connection
func createOneWriteTCPclient() {
	// Establish client
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Not able to establish tcp client")
	}
	defer conn.Close()

	// Init babbler
	babbler := babble.NewBabbler()
	babbler.Count = 1

	// record time
	start := time.Now()
	// Send message
	for i := 0; i < WORDS; i++ {
		//log.Println("Sent", i)
		_, err := conn.Write([]byte(babbler.Babble()))
		if err != nil {
			log.Fatalln("Not able to send message")
		}

		if i%10000 == 0 {
			log.Printf("Sent %d messages", i)
		}

		// Receive response
		buf := make([]byte, 8)
		_, err = conn.Read(buf)
	}
	// log time
	log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
}

// ==========================
// Multiple TCP connections
// ==========================

// Multiple TCP connection writes
func createMultiWriteTCPserver() {
	// Establish server
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Not able to establish tcp server")
	}

	log.Printf("server listening at %v", lis.Addr())
	go func() {
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Fatalln("Not able to accept connection")
			}

			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				log.Fatalln("Not able to read message")
			}

			conn.Close()
		}
	}()
}

// Multiple writes in one TCP connection
func createMultiWriteTCPclient() {

	// Init babbler
	babbler := babble.NewBabbler()
	babbler.Count = 1

	// record time
	start := time.Now()

	for i := 0; i < WORDS; i++ {
		// Establish client
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatalln("Not able to establish tcp client")
		}
		defer conn.Close()

		// Send message
		_, err = conn.Write([]byte(babbler.Babble()))
		if err != nil {
			log.Fatalln("Not able to send message")
		}
		if i%1000 == 0 {
			log.Printf("Sent %d messages", i)
		}
	}

	// log time
	log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
}

// ==========================
// Singular TCP connections
// With one write for all words
// (JSON)
// ==========================

type Message []string

func createAllWordsOneWriteTCPserver() {
	// Establish server
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Not able to establish tcp server")
	}

	log.Printf("server listening at %v", lis.Addr())

	// Accept connection
	conn, err := lis.Accept()
	if err != nil {
		log.Fatalln("Not able to accept connection")
	}

	log.Println("Server accepted connection")

	// Read the size of the message
	buf := make([]byte, 4)
	_, err = conn.Read(buf)
	if err != nil {
		log.Fatalln("Not able to read message")
	}

	log.Println("Server read message size")

	// repeatedly read the message
	// util the size is reached
	msgSize := binary.BigEndian.Uint32(buf)
	// each buffer is 4096 bytes
	buf = make([]byte, 4096)
	// the message
	var msg []byte
	// the size of the message
	var size uint32
	for size < msgSize {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("Not able to read message")
		}
		msg = append(msg, buf[:n]...)
		size += uint32(n)
	}

	log.Println("Server read message")

	// Deserialize message
	var message Message
	err = json.Unmarshal(msg, &message)
	if err != nil {
		panic(err)
	}

	// Print message
	fmt.Println("Server recieved", len(message))

	// respond ok
	_, err = conn.Write([]byte(fmt.Sprint(len(message))))
	if err != nil {
		log.Fatalln("Not able to send message")
	}

	// Close connection
	conn.Close()
}

// Multiple writes in one TCP connection
func createAllWordsOneWriteTCPclient() {
	// Establish client
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln("Not able to establish tcp client")
	}
	defer conn.Close()

	// Init babbler
	babbler := babble.NewBabbler()
	babbler.Count = 1

	// Create a message with WORDS
	log.Println("Creating message")
	words := make([]string, WORDS)
	for i := 0; i < WORDS; i++ {
		words[i] = babbler.Babble()
	}

	log.Println("Created message")

	// Serialize msg to JSON format
	message, err := json.Marshal(words)
	if err != nil {
		panic(err)
	}

	log.Println(len(message))

	// record time
	start := time.Now()

	// First send the size of the message
	size := make([]byte, 4)
	binary.BigEndian.PutUint32(size, uint32(len(message)))
	_, err = conn.Write(size)
	if err != nil {
		log.Fatalln("Not able to send message")
	}

	// Send message in a for loop
	for i := 0; i < len(message); i += 4096 {
		_, err := conn.Write(message[i:min(i+4096, len(message))])
		if err != nil {
			log.Fatalln("Not able to send message")
		}
	}

	// Wait for ok message
	buf := make([]byte, 8)
	l, err := conn.Read(buf)
	if err != nil {
		log.Fatalln("Not able to read message")
	}
	log.Println("Client recieved", string(buf[:l]))

	// log time
	log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
