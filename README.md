# tcppoc

for peerstream interview assignment

go tcp chat server PoC

TCP server application performing as a chat server. Clients should be able to connect to the listening port using plaintext protocol and be able to communicate with each other. Messages are separated by <LF>, when connected user should be presented with a list of users currently online, everyone see messages from everyone, server should support /nick <nickname> command for users to be able to redefine the default auto-assigned nickname "GuestNNN" and /register command for users to be able to protect their nickname from being taken by other user with a password.


Instructions to build and run:

- go build main.go
- ./main


on machine to connect:

telnet [ip or localhost] 9009


commands:

/nick to change nick
/login password to login with your password
/register password to register username with password
