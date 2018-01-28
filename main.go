package main

import (
    "flag"
    "fmt"
    "strconv"
    "sync"
    "time"
)

type Proposal struct {
    value int
    done chan<- Response
}

type Response struct {
    value int
    response bool
}


func respondSync(channel chan<- Response, message Response) {
    channel <- message
}

func respond(channel chan<- Response, message Response) {
    select {
    case channel <- message:
        break
    default:
        break
    }
}

func acceptor(n int, proposals <-chan Proposal, end <-chan bool, wg sync.WaitGroup) {
    wg.Add(1)
    defer wg.Done()

    for {
        select {
        case proposal := <-proposals:
            if proposal.value % n == 0 {
                respond(proposal.done, Response{proposal.value, true})
                fmt.Printf("Acceptor with key %d: accepted proposal %d\n", n, proposal.value)
            } else {
                respond(proposal.done, Response{proposal.value, false})
                fmt.Printf("Acceptor with key %d: rejected proposal %d\n", n, proposal.value)
            }
        case <-end:
            fmt.Printf("Acceptor with key %d: Killing self\n", n)
            break
        }
    }
}

func main() {
    var wg sync.WaitGroup

    flag.Parse()
    rawParameters := flag.Args()

    numAcceptors := len(rawParameters)

    acceptorParameters := make([]int, numAcceptors, numAcceptors)
    for i, x := range flag.Args() {
        acceptorParameters[i], _ = strconv.Atoi(x)
    }

    proposals := make([]chan Proposal, numAcceptors, numAcceptors)
    killers := make([]chan bool, numAcceptors, numAcceptors)
    for i := 0; i < numAcceptors ; i++ {
        proposals[i] = make(chan Proposal)
        killers[i] = make(chan bool)
        go acceptor(acceptorParameters[i], proposals[i], killers[i], wg)
    }

    responseChannel := make(chan Response)

    suicide := make(chan bool)

    go (func (transmitters []chan Proposal, responseChannel chan<- Response, end <-chan bool, wg sync.WaitGroup) {
        wg.Add(1)
        defer wg.Done()

        var i int = 1
        for {
            select {
            case <-suicide:
                fmt.Printf("Proposer: Killing self\n")
                break
            default:
                proposal := Proposal{i, responseChannel}
                for _, transmitter := range transmitters {
                    transmitter<- proposal
                }
            }
            i++
            time.Sleep(time.Millisecond)
        }
    })(proposals, responseChannel, suicide, wg)

    counts := make(map[int]int)

    for {
        response := <-responseChannel
        value := response.value
        accepted := response.response
        if accepted {
            counts[value]++
        }

        if counts[value] > numAcceptors/2 {
            fmt.Printf("CONSENSUS: %d\n", value)
            suicide <- true
            for _, killer := range killers {
                killer <- true
            }
            break
        }
    }

    wg.Wait()
}
