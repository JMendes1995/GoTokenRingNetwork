package calculator

import (
	"context"
	"encoding/json"
	"fmt"
	cpb "goTokenRingNetwork/calculator/proto"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	//"time"
)



func (queue *CalcQueue) Client(mutex *sync.Mutex) {
	conn, err := grpc.NewClient(ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error connecting to server %s", err)
	}

	c := cpb.NewCalculatorClient(conn)
	defer conn.Close()
	// setupt connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mutex.Lock()

	log.Print("Local queue: ", queue.CalcOperations)
	if len(queue.CalcOperations) > 0 {
		// reverse the list of operations to remove from the highest index to the lowest to ensure that the correct calc is being remove from the queue
		for i := len(queue.CalcOperations) - 1; i >= 0; i-- {
			calcJson, _ := json.Marshal(queue.CalcOperations[i])
			response, err := c.Calculate(ctx, &cpb.CalculateRequest{Calc: string(calcJson)})
			if err != nil {
				log.Fatalf("Error requesting message from server %v", err)
			}
			log.Printf("%s Response: Calculation %s%s%s=%s\n", ServerAddress, fmt.Sprintf("%d", queue.CalcOperations[i].Operand1),
				queue.CalcOperations[i].Operator, fmt.Sprintf("%d", queue.CalcOperations[i].Operand2), response.Result)
			queue.Delete(queue.CalcOperations[i])
		}
	}
	log.Print("The queue is empty: ", queue.CalcOperations)
	mutex.Unlock()
}
