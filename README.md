# Tracebook Messenger
Not just for the pun of it! This repository houses my blockchain inspired peer-to-peer command line interface chat application. I have a directory for my client/webserver implementation as well.

* Assign yourself a screen name
* Have access to the entire chat history that occured before your connection
* Chat with your peers!

### How does it work?
When a node is run, you can choose to connect to a peer or not. If you connect to a peer, your node is updated with the chat history and list of peers. At this point all connected peers are simultaneously updated of the existence of the new node. Peers can come and go, but as long as there is one live node in the chat, any node that connects to it will be able to view the entire history of the chat.

### Blockchain inspired?
Each node in a blockchain contains every blockchain interaction that ever occured. The chain itself. For a node to add a block to the chain, the other nodes need to agree upon the validity of the transaction. The transaction then gets added to the chain data of each participating node.

This chat application behaves in a similar manner with its chat history. When a node wants to send a message to its peers, it adds a message to its version of the chat history (or chain data), then sends an update to its peers. A node accepts an update from a peer if the history is new.