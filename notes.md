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


### Stop client calculation if product number found
if a client found a product number, should be good idea to send signals to other clients to stop their calculations.

#### Potential issue
Currently we have fault tolerance (network distrub) if a client disconnected during the option menu.
But if one is crashing during the calculation, we need to report his calculation to another one client.


