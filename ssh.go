package main

import (
	"log"
	"net"
	"os"

	"code.google.com/p/go.crypto/ssh"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("usage: go-ssh <username> <host>:<port> <cmd>")
	}

	agent, e := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if e != nil {
		panic(e)
	}
	defer agent.Close()

	auths := []ssh.ClientAuth{
		ssh.ClientAuthAgent(ssh.NewAgentClient(agent)),
	}

	clientConfig := &ssh.ClientConfig{
		User: os.Args[1],
		Auth: auths,
	}

	client, err := ssh.Dial("tcp", os.Args[2], clientConfig)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run(os.Args[3]); err != nil {
		panic("Failed to run: " + err.Error())
	}
}
