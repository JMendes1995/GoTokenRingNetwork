# GoTokenRingNetwork

Project that implements a Mutual Exclusion with the Token Ring Algorithm

### Implementation Requirements
* Create a ring network with 5 peers (call them p1 to p5) such as the one represented in
the figure above. Each peer must be in a different machine (m1 to m5). Besides these 5
peers, create also a calculatormulti server called server that runs on machine m6.
* Each peer only knows the IP address of the machine where the next peer is located: p1
knows the IP of m2, p2 knows the IP of m3, . . . , p5 has the IP of m1, thus closing the
ring. All peers know the IP address of the machine where the server is located (m6).
These can be passed to the peer in the command line, e.g., for peer p2 you would run
the following command (in machine m2): $ peer m3 m6.
* One of the threads in each peer generates requests for the server following a Poisson
distribution with an average frequency of 4 events per minute. Each request is also
random (both the operation and the arguments). These requests are placed in a local
queue.
* Another thread in that peer runs in a loop waiting for a message (that can only come
from the previous peer). This message is called a token. The token can be a void
message. The peer that holds it at any given moment has exclusive access to server,
effectively implementing mutual exclusion.
* When the peer receives the token, it checks if it has requests for the server in its local
queue. If not, it forwards the token to the next peer in the ring. If so, it holds the
token until the requests are processed. When the results are received and printed on
the terminal, the peer restarts the forwarding of the token.


### How to build the project
In the root directory of the project are present 2 binaries from **Linux (goTokenRingNetwork-linux)** and **MacOS (goTokenRingNetwork-macos)** arm architecture. However is possible to build the project manually by executing the following commands.

```bash
$ cd GoTokenRingNetwork
$ go build
```

### How to execute the project
To execute the project there is 2 possible ways.
1. Execute the code in phsical nodes
2. Build docker environment using docker compose.

#### Execution in physical nodes
Before executing the project, is required to pull the source code into the nodes that will be part of the network.
Furthermore, is also require to execute the command passing input parametes depending if the node is a peer node or is the calculator server

##### Peer input paramenters
* `--mode <peer or server>`
* `--priority <int value (1..n)>`
* `--neighbourAddress <Neighbour peer IP address>`
* `--serverAddress <Calculator server  IP address>`

Example:
Note: this commands should be executed in the respective nodes and not in the same node.

```bash
$ cd GoTokenRingNetwork
$ ./goTokenRingNetwork-linux --mode peer --priority 1 --neighbourAddress 10.10.1.2 --serverAddress 10.10.1.6
$ ./goTokenRingNetwork-linux --mode peer --priority 2 --neighbourAddress 10.10.1.3 --serverAddress 10.10.1.6
$ ./goTokenRingNetwork-linux --mode peer --priority 3 --neighbourAddress 10.10.1.4 --serverAddress 10.10.1.6
$ ./goTokenRingNetwork-linux --mode peer --priority 4 --neighbourAddress 10.10.1.5 --serverAddress 10.10.1.6
$ ./goTokenRingNetwork-linux --mode peer --priority 5 --neighbourAddress 10.10.1.1 --serverAddress 10.10.1.6
```
##### Calculator server input parametes
* `--mode <peer or server>`

Example:
```bash
$ cd GoP2PNetwork
$ ./goTokenRingNetwork-linux --mode server
```


#### Docker

```bash
$ cd GoP2PNetwork
$ docker compose up -d
```

##### Clean up docker environment
```bash
$ docker compose down && docker compose rm && docker rmi -f $(docker images -aq)
```

##### Clean up docker environment
```bash
$ docker compose down && docker compose rm && docker rmi -f $(docker images -aq)
```

### Results intrepertation
##### Tokenring Server startup
```bash
#Node 1
2025/01/03 12:02:39 Token Ring overlay network server listening at 10.10.1.1:50002
```
##### Leadership challenge
All nodes send a leadership challeng to its neighbour when the ring is completed, its determined the leader.
```bash

#Node 1

2025/01/03 12:02:39 Challeging leadership with peer 10.10.1.2:50002
2025/01/03 12:02:41 Leader election Request comming from peer 10.10.1.5 with priority: 5
2025/01/03 12:02:41 Leader election Result Leader: 10.10.1.5 with priority: 5
```
##### Leadership Result
After the completing the challenge the leadership of the network generates a toke which is a hash sha3-256 of the jason node map
```bash
type Peer struct {
	Uid         string `json:"Uid"`
	Address     string `json:"Address"`
	Priority    int    `json:"Priority"`
	Leader      bool   `json:"Leader"`
	TokenHolder bool   `json:"TokenHolder"`
}
```

```bash
#Node 5

2025/01/03 12:02:41 Challeging leadership with peer 10.10.1.1:50002
2025/01/03 12:02:41 Peer 10.10.1.5:50002 is the Leader with priority of 5
2025/01/03 12:02:41 SHA3-256 hash: 3e28f3fe196b545c8e58698db309d24b01b73123307a68a949e9d5048398c0e3
2025/01/03 12:02:41 Generation New TOKEN: 3e28f3fe196b545c8e58698db309d24b01b73123307a68a949e9d5048398c0e3
2025/01/03 12:02:41 Leader election Request comming from peer 10.10.1.4 with priority: 4
2025/01/03 12:02:41 Leader election Result Leader: 10.10.1.5 with priority: 5
```

##### Event generation
After defining the leadership and generationg the token which will give access to send events to the calculator server,
the node parallelly generates events following a Poisson distribution with an average frequency of 4 events per minute.
These events are stored in a local queue to be sent to the server when the node has the token
```bash
#Node 5

2025/01/03 12:02:41 Minute:1 Nrequests:5
2025/01/03 12:02:41 Request 1 at 1.511753 seconds :::: calculation: 658-951
2025/01/03 12:02:41 Sleep 1.51175 seconds...
2025/01/03 12:02:42 Request 2 at 19.707801 seconds :::: calculation: 284/192
2025/01/03 12:02:42 Sleep 18.19605 seconds...
2025/01/03 12:03:00 Request 3 at 29.862512 seconds :::: calculation: 838/557
2025/01/03 12:03:00 Sleep 10.15471 seconds...
```


##### Token transmission
```bash
#Node 5

2025/01/03 12:03:01 Sending token to peer: 10.10.1.1:50002
2025/01/03 12:03:01 Local Peer holding the tokan: false
```
##### Recieving the token
When a node recives the token from its peer, it looks at the local queue, and transmits the stored calculus. If the local queue is empty transmists the token to the following  peer.
```bash
#Node 1

2025/01/03 12:03:01 Recieved token from peer 10.10.1.5 with value: 3e28f3fe196b545c8e58698db309d24b01b73123307a68a949e9d5048398c0e3
2025/01/03 12:03:01 Token Previous holder: 10.10.1.5:50002
2025/01/03 12:03:11 Local Peer holding the token: true
2025/01/03 12:03:11 Local queue: [{771 626 *} {279 177 *} {126 486 *} {607 148 -} {797 676 /} {25 978 -}]
```

##### Results from the calculator server
```bash
#Node 5
2025/01/03 12:03:11 10.10.1.6:50001 Response: Calculation 25-978=-953.000000
2025/01/03 12:03:11 10.10.1.6:50001 Response: Calculation 797/676=1.178994
2025/01/03 12:03:11 10.10.1.6:50001 Response: Calculation 607-148=459.000000
2025/01/03 12:03:11 10.10.1.6:50001 Response: Calculation 126*486=61236.000000
2025/01/03 12:03:11 10.10.1.6:50001 Response: Calculation 279*177=49383.000000
2025/01/03 12:03:11 10.10.1.6:50001 Response: Calculation 771*626=482646.000000
2025/01/03 12:03:01 The queue is empty: []
```

Calculator server
```bash
2025/01/03 12:02:41 Calculator server listening at [::]:50001
2025/01/03 12:03:01 Request from 10.10.1.5: Calculation 838/557=1.504488
2025/01/03 12:03:01 Request from 10.10.1.5: Calculation 284/192=1.479167
2025/01/03 12:03:01 Request from 10.10.1.5: Calculation 658-951=-293.000000
2025/01/03 12:03:11 Request from 10.10.1.1: Calculation 25-978=-953.000000
2025/01/03 12:03:11 Request from 10.10.1.1: Calculation 797/676=1.178994
2025/01/03 12:03:11 Request from 10.10.1.1: Calculation 607-148=459.000000
2025/01/03 12:03:11 Request from 10.10.1.1: Calculation 126*486=61236.000000
2025/01/03 12:03:11 Request from 10.10.1.1: Calculation 279*177=49383.000000
2025/01/03 12:03:11 Request from 10.10.1.1: Calculation 771*626=482646.000000
2025/01/03 12:03:21 Request from 10.10.1.2: Calculation 725-62=663.000000
2025/01/03 12:03:21 Request from 10.10.1.2: Calculation 435*227=98745.000000
2025/01/03 12:03:21 Request from 10.10.1.2: Calculation 435/473=0.919662
```
