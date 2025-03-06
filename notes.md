## Notes related to this project

## Step 1: intermachine communication
### Goal
The goal is to make simple script to communicate from a machine A to a remote machine B, client-server principle.

I successfuly simulate connexion between clients and a server, with simple handling function server side.

Now I will try to handle server response.
The goal is to simulate little calculus client side and return the result to server.
The server want to make his clients count occurrence of 'e' in a word.

Steps:
  - server started waiting for connections
  - client connect to server
  - client is waiting for server response
  - server send a response (a word)
  - client, by using a little function will return the number of 'e' in the word

### Code Implementation

What has been made:
  - server can wait for a number of connections
  - after some client connections, server send batch of random letters
  - clients count 'e' occurrences and return the result to main server

#### Multifunction execution
I also started a multifunction execution client. For example, to count 'e' occurrences, the server sends `1eMyWordToCountE`.
`1` stand for executing the function to count occurrences and the following letter 'e', the letter to count.


### Results
For initial tests I launched clients in other terminals and one for the central server.
After it works well, I configured 2 old laptop with ubuntu server and lubuntu to test the intermachine communication. All machines were connected in a local subnet.

In first I was sending a lot of letters (100000000), it took longer to send over the network than to process the data.


## Step 2: Multifunction execution
The goal of this one is to be able to execute different functions on clients. In the previous step, client are counting letter occurence letter in a sentence.

### Idea
We're going to start with a simple message parsing to choose the right function
So the idea is to use a structure in this style:
```DIGITargs```
- a digit to execute a certain function (eg: `1` for occurence counting)
- an optional argument (eg: `e` for counting 'e' in the sentence)

### Ping
Implement a function responding "pong" or another message, the goal is just to check if the client is always alive (connected).

### Change the way to wait goroutines
```go
func (ms Marmots) performAction(fctToExecute func(*Marmot)) {
	for _, m := range ms {
		go fctToExecute(m)
	}
	// All clients are connected we can start calculations
	for _, marmot := range ms {
		marmot.start <- true
	}
	for _, marmot := range ms {
		<-marmot.end
	}

}


func (m *Marmot) Ping() {
	<-m.start
	// send 'ping' to client
	m.data = "0Ping\n"
	m.writeData()
	m.readResponse()
	m.end <- true
}

```

Here when we send a ping to the client I'm starting new goroutine, currently the way to wait the end of the goroutine, is by using a channel `end`. 

But if the go routine crash / take too many time, the performAction function will never end.

> Solution: by using `sync.WaitGroup`


## Step 3: Retry connection
The goal is to make a little loop for client, to always try to connect to the server if it's not, every minute for example.
Also catch potential network issue with the retry connection system

### Implementation
The client try to connect to server and handle the connection, as long as the server does not send exit request.

## Step 4: Create calculation menu
Just created basic menu to perform some calculations to clients.

### Counting letter occurences 
User can select a letter and it will send to clients batch of random letters and count occurrence of the letter.

### Prime Number calculation
A given number is a natural number greater than 1 that is not a product of two smallers number.

So with this we can imagine a scenario:
- we have 3 clients connected
- we want to know if the number 90000 is prime
- sqrt(90000) = 300
- we will create 3 range from 0 to 300 [(2,100), (100,200), (200,300)]
- each client will check is the number is not a product of a number in the range
- if client do not found product number, returns -1 else the number
- foreach client, server checks if no product number found


### Message format
For the second client, with the range of [100,200] it will receive: 
`290000@100@200`
`2` is the id for client to know which function to execute
We're using `@` as separator because it should never be used by server (in this context)

You can test it with this numbers: 
```
499139582359084529
943349775459380173
440149999829818483
761795777243436499
403111596387917561
811508878081456537
341705128992983771
178225203054306761
446379149495008381
283762230407893121
880354649778912289
```


### TODO: Stop client calculation if product number found
if a client found a product number, should be good idea to send signals to other clients to stop their calculations.

#### TODO: Potential issue
Currently we have fault tolerance (network distrub) if a client disconnected during the option menu.
But if one is crashing during the calculation, we need to report his calculation to another one client.


## Create a homemade lab
Let's create a homemade lab to simulate a datacenter with multiples servers.

I took 2 old laptop and the one im using to work.
2 old laptops are using lubuntu.

I found an old router and connect 3 laptops to it.

I setup fix ip address:
- Server `$  sudo ip addr add 192.168.1.25/24 dev enp2s0`
- Client 1 `$  sudo ip addr add 192.168.1.26/24 dev enp2s0`
- Client 2 `$  sudo ip addr add 192.168.1.27/24 dev enp2s0`

By using `ping` i can check if all laptops are interconnected.

Now i need to transfer my executable file `./marmotReduce` to launch distant client.
To do this i will use ssh:
`$ scp marmotReduce lubuntu@192.168.1.27:/home/lubuntu`

Unfortunately my 2 old laptops did not `ssh server` installed, so i need to install it:
```
$ sudo apt update
$ sudo apt install openssh-server -y
$ sudo systemctl start ssh
$ sudo systemctl enable ssh
```

I also added my public ssh key to all laptop to gain some time foreach `scp`:
`$ cat ~/.ssh/id_ed25519.pub | ssh lubuntu@192.168.1.27 "mkdir -p ~/.ssh && chmod 700 ~/.ssh && cat >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
`

### The Home made lab:
![Home made lab image](assets/homemadeLab1.jpeg)

## Observations:
To use all processor capacity, the optimal way is to open X clients by X hearts by server.

Error: if I disconnect the server ethernet cable and i try to Ping clients, the function is `Pings` is blocked
	-> maybe add a timeout in each `Ping` method


## TODO: Make .env config for all constants

## TODO: Create a power score for all clients
This score can be used to better distribute tasks, across the system.
Currently task are evenly distributed with all clients, but it should be better if server can distribute tasks in function of computing power.

Idea: create a function asking for some computing power and each client returns the time it took. 


## TODO: Issue if new client connect during calculation
The Wait() methods used to wait all <-m.end waits infinite time if a new client is connecting.
-> maybe add attribut like "isWorking" to check if we need to wait this client or not

## PI estimation
The goal is to have an estimation of PI.
We will use the MonteCarlo method.
To achieve this, we will distribute in equal range the chunk to compute.

### Mesage format
The message sent by server will be: `3N` where `3` is the function id and `N` is the number of random point to calculate.

### Observation
Tests has been run on my homemade lab.
Each "remotes/locals clients" is running on a single laptop.

For the samples amount: `1000000000`
- With total of 3 clients (3 locals): `13.54 seconds`
- With total of 4 clients (3 locals + 4 remotes ): `8.54 seconds`
- With total of 6 clients (3 locals + 4 remotes + 2 remotes): `7.29 seconds`
- With total of 10 clients (3 locals + 4 remotes + 2 remotes + 4 remotes): `5 seconds`


For the samples amount: `10000000000`
- With total of 3 clients (3 locals): `83 seconds`
- With total of 4 clients (3 locals + 4 remotes ): `38.78 seconds`
- With total of 10 clients (3 locals + 4 remotes + 2 remotes + 4 remotes): `49.59 seconds`

## Timeout implementation
### Server side
Implement timeout server side when sending / receiving data, using `context`.
```go

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
		go fctToExecute(ctx, resultChan)
		// timeout implementation:
		select {
		case res := <-resultChan:
			return res
		case <-ctx.Done():
			printError(errorMessage)
			return false
		}
```
### Current issue: send to client signal to stop calculation
Send signal to client to stop calculation, or send exit signal to client 



### Refactor client code to use Marmot type
Refactor client code with `Marmot` type, to use functions like `executeFunctionWithTimeout`, `readData` ect..



## Self-update client
The goal is to have a client with capablity of self-updating.
The server can send new client version, client save the new executable, start it, and kill himself.

### New implementation for messages communication
Now we need to send binary file, not only string message.
So the problem is we can no longer use `ReadString("\n")`.

To address this issue, I will create a `struct` with an `id` and `data` attribut, that will be send between client and server. 

ID: represents the action id
0: Ping
1: Close connection (exit)
2: Counting letter
3: Calculate if a number is prime
4: Calculate pi estimation


#### Buffer size issue
Current issue: i need to know the size of the message received to read it with sized buffer.

So before sending `Message` i will put at the head of `byte` array, the size of `Message`, so i can in first read the size to create right sized buffer.

Solution found: before sending `Message`, i will send the length of it


### Implementation

The server can update all clients from menu.
It sends Struct with current client version and the binary file associated.
The client receive it, check if his current version is outdated, then he will write the binary file and execute it and kill the old one.


## Free fall simulation
Goal: I want to implement "real use cases" to my little distributed system.

So we will calculate free fall simulation.
```
	y = y0 + v0.t - (1/2)g.t²

	y0 = initial height (m)
	v0 = initial velocity (m/s)
	g = 9.81 m / s² (gravity)
	t = time (s)
	
```

### Implementation

Created a file with struct and methods for simulating free fall.
Currently working on a single machine.

Use: 
```go
	f := NewFreeFallEarth(100000, 10000, 0.001)
	f.generateFreeFallData(true, true)
```

Result: 
```
1741250258628| DEBUG: Start generating free fall data
1741250259147| DEBUG: Time to generate free fall points: 519.332332ms
1741250259372| DEBUG: Start generating plot
1741250371886| DEBUG: Graphique generated: 'free_fall_simulation.png' in 1m52.51485s
1741250371886| DEBUG: Start saving results to file
1741250392384| DEBUG: File created: 'freefall.dat' in 20.49809s
```

We can see it take too many time to generate plot, if you need plot, just use the generate results in a file and open a plot using `gnuplot`.

I put 2 images (`free_fall_gnuplot.png` and `free_fall_simulation.png`) inside `assets` folder.

Let's start the multiclient version.

### Multi client
The server will send struct with a range
Client will calculate position inside the range and returns the struct with `FreeFallResults`


### Experimentation results
To be efficient we need to compute a lots of points but it's breaking my RAM limitation (16GO)

The free wall simulation is not a good test pour distributed system.
-> may be add more compute foreach points, to decrease number of point to compute so gain memory.
