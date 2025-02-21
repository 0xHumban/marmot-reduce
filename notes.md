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

## Implementation
The client try to connect to server and handle the connection, as long as the server does not send exit request.
