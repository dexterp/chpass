package main

import (
	"bytes"
	"chpassh/config"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func chpaSsh(user, host, port, oldpass, newpass string) (string, error) {
	// user@host:port: STATUS - Message
	output := "%s@%s:%s: %s - %s"
	var hostBuf bytes.Buffer
	hostBuf.WriteString(host)
	hostBuf.WriteString(":")
	hostBuf.WriteString(port)

	// Check password
	if strings.Compare(oldpass, newpass) == 0 {
		return fmt.Sprintf(output, user, host, port, "SKIPPED", "Old and new password identical"), nil
	}

	// Prepare ssh config
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(oldpass)},
	}

	// Dial a host
	connection, err := ssh.Dial("tcp", hostBuf.String(), sshConfig)
	if err != nil {
		errmsg := fmt.Sprintf("Failed to dial host: %s", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}

	// Create a session
	session, err := connection.NewSession()
	if err != nil {
		errmsg := fmt.Sprintf("Failed to create session: %s", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}

	// Set Pty
	modes := ssh.TerminalModes{
		ssh.ECHO:          0, // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		errmsg := fmt.Sprintf("request for pseudo terminal failed: %s", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}

	// Link STDERR and STDOUT
	stdin, err := session.StdinPipe()
	if err != nil {
		errmsg := fmt.Sprintf("Unable to setup stdin for session: %v", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		errmsg := fmt.Sprintf("Unable to setup stdout for session: %v", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		errmsg := fmt.Sprintf("Unable to setup stderr for session: %v", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}
	go io.Copy(os.Stderr, stderr)

	cmd := "echo \"%s\n%s\n%s\" | /usr/bin/passwd"
	err = session.Run(fmt.Sprintf(cmd, oldpass, newpass, newpass))
	if err != nil {
		errmsg := fmt.Sprintf("Unable to run command: %v", err)
		return fmt.Sprintf(output, user, host, port, "ERROR", errmsg), fmt.Errorf(errmsg)
	}
	return fmt.Sprintf(output, user, host, port, "OK", "Password changed"), nil
}

var confPath string

func init() {
	flag.StringVar(&confPath, "config", "chpassh.json", "Path to JSON configuration file")
}

func main() {
	flag.Parse()
	conf := config.Config{}

	// Load hosts
	err := conf.LoadConfig(confPath)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	hosts := conf.Hosts
	for _, host := range hosts {
		fmt.Printf("%s@%s:%d: PROCESSING\n", conf.User, host.Name, host.Port)
		output, err := chpaSsh(conf.User, host.Name, strconv.Itoa(host.Port), conf.Oldpassword, conf.Newpassword)
		fmt.Printf("%s\n", output)
		if err != nil {
			continue
		}
	}
}
