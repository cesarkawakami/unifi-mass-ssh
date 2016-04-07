package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage: unifi-mass-ssh USER PASS CMD IP...")
	}

	user := os.Args[1]
	pass := os.Args[2]
	cmd := os.Args[3]
	ips := os.Args[4:]

	log.Printf("Starting sending '%s' to IPs %s", cmd, strings.Join(ips, ", "))

	wg := &sync.WaitGroup{}
	for _, ip := range ips {
		wg.Add(1)
		go do(ip, user, pass, cmd, wg)
	}
	wg.Wait()

	log.Printf("Done.")
}

func do(ip, user, pass, cmd string, wg *sync.WaitGroup) {
	defer wg.Done()

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout: 30 * time.Second,
	}

	conn, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		log.Printf("Error dialing '%s': %s", ip, err)
		return
	}

	sess, err := conn.NewSession()
	if err != nil {
		log.Printf("Error creating a new session to '%s': %s", ip, err)
		return
	}

	out, err := sess.CombinedOutput(cmd)
	if err != nil {
		log.Printf("Failed running '%s' on '%s': %s", cmd, ip, err)
		return
	}

	log.Printf("Command '%s' on '%s' ran successfully. Output:", cmd, ip)
	log.Printf("%s", out)
}
