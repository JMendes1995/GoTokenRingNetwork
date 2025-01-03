package calculator

import (
	"errors"
	"fmt"
	"goTokenRingNetwork/poisson"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Calc struct {
	Operand1 int    `json:"Operand1"`
	Operand2 int    `json:"Operand2"`
	Operator string `json:"Operator"`
}

type CalcQueue struct {
	CalcOperations []Calc `json:"CalcOperations"`
}

func (q *CalcQueue) Add(op Calc) {
	q.CalcOperations = append(q.CalcOperations, op)
}

func (q *CalcQueue) Delete(op Calc) {
	for i, c := range q.CalcOperations {
		if c == op {
			q.CalcOperations = append(q.CalcOperations[:i], q.CalcOperations[i+1:]...)
		}
	}
}

const Port = ":50001"

var (
	Queue         CalcQueue
	ServerAddress string
)

func RandomCalcOperaton(rng rand.Rand) Calc {
	operators := []string{"+", "-", "*", "/"}
	operator := operators[rng.Intn(len(operators))]
	return Calc{Operand1: rng.Intn(1000), Operand2: rng.Intn(1000), Operator: operator}
}
func CalcOperation(operation Calc) (float64, error) {
	var result float64
	if operation.Operator == "+" {
		result = float64(operation.Operand1 + operation.Operand2)
		return result, nil
	} else if operation.Operator == "-" {
		result = float64(operation.Operand1 - operation.Operand2)
		return result, nil
	} else if operation.Operator == "/" {
		result = float64(operation.Operand1) / float64(operation.Operand2)
		return result, nil
	} else if operation.Operator == "*" {
		result = float64(operation.Operand1 * operation.Operand2)
		return result, nil
	} else {
		return 0, errors.New("operator not valid, please choose one operator between (+,-,*,/)")
	}
}

func (queue *CalcQueue) EventGenerator(mutex *sync.Mutex) {
	time.Sleep(time.Second * 10)
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	lambda := 4.0 // rate of 4 events per minute
	poissonProcess := poisson.PoissonProcess{Lambda: lambda, Rng: rng}
	totalRequests := 0
	t := 1 // current minute

	for {
		nRequests := poissonProcess.PoissonRandom()
		log.Printf("Minute:%d Nrequests:%d", t, nRequests)
		currentTime := 0.0
		previousTime := 0.0
		for i := 1; i <= nRequests; i++ {

			totalRequests++
			newCalc := RandomCalcOperaton(*rng)

			mutex.Lock()
			queue.Add(newCalc)
			mutex.Unlock()

			// get the time for the next request to be executed
			interArrivalTime := poissonProcess.ExponentialRandom()
			previousTime = currentTime
			currentTime += interArrivalTime * 60

			log.Printf("Request %d at %f seconds :::: calculation: %d%s%d\n", i, currentTime, newCalc.Operand1, newCalc.Operator, newCalc.Operand2)
			log.Printf("Sleep %.5f seconds...\n", float64(currentTime-previousTime))
			delta := time.Duration(currentTime-previousTime) * time.Second
			time.Sleep(delta)
			if i == nRequests && currentTime < 60 {
				log.Printf("Requests for the minute %d endend before finishing the 60s.\nWaiting %f seconds to complete the cycle of 60s....", t, float64(60-currentTime))
				time.Sleep((time.Duration(60-currentTime) * time.Second))
			}
		}
		fmt.Println()
		log.Printf("Statistics: Total requests: %d Minutes spent: %d rate:%f\n", totalRequests, t, float64(totalRequests)/float64(t))
		t++
	}
}
