# distrubuted_mandatory_2

How to use.
When the program starts, it asks for an inputcommand.
Press 1 to enter af message and send it to the server.
Press 2 to emulate the handling of message reordering.

a) We use struct to imitate packets. They contain the metadata in the form of hashcode, sender, sync and ack count.
the data it self, we use strings to send a message.

b) We use threads. its not realistic cause they are part of the same process, where as with a real network consists of multiple processes.

c) We are not actually handling message reordering. Instead we are simulating it through a premade int array that we reverse. The server will then sort this array when it has been received.

d) We hash the complete message string, and send it with the packet. the server then hashses the recieved message and compare the 2 hash results. if they are the same, it send a confirm order to the clinet, otherwise the clinet resend the packet until a confirmation or timeout.

e) To make sure that the connection works both ways and to ensure continues streaming of data.