## Notes related to this project

## Step 1: intermachine communication
The goal is to make simple script to communicate from a machine A to a remote machine B, client-server principle.

I successfuly simulate connexion between clients and a server, with simple handling function server side.

Now I will try to handle server response.
The goal is to simulate little calculus client side and return the result to server.
The server want to make his clients count occurrence of 'a' in a word.

Steps:
  - server started waiting for connections
  - client connect to server
  - client is waiting for server response
  - server send a response (a word)
  - client, by using a little function will return the number of 'a' in the word

