# distrubuted_mandatory_2

a) we use struct to imitate packets. They contain the metadata in the form of hashcode, sender, sync and ack count.
the data it self, we use strings to send a message.

b) we use threads. its not realistic cause they are part of the same process, where as with a real network consists of multiple processes.

c) 

d) we hash the complete message string, and send it with the packet. the server then hashses the recieved message and compare the 2 hash results. if they are the same, it send a confirm order to the clinet, otherwise the clinet resend the packet until a confirmation or timeout.

e) to make sure that the connection works both ways and to ensure continues streaming of data. 


