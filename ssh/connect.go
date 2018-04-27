package ssh

import (
	"golang.org/x/crypto/ssh"
	"log"
	"io/ioutil"
	"os"
	"golang.org/x/crypto/ssh/terminal"
	"strconv"
	"errors"
)

type Server struct {
	Name 		string	`json: name`
	Ip			string	`json: ip`
	Port		int		`json: port`
	User 		string	`json: user`
	Password 	string	`json: password`
	Method 		string	`json: method`
	Key 		string 	`json: key`
}

func parseAuthMethods(s *Server) ([]ssh.AuthMethod,error) {
	sshAuth := []ssh.AuthMethod{}

	switch s.Method {
	case "password":
		sshAuth = append(sshAuth,ssh.Password(s.Password))
		break
	case "pubkey":
		method,err := pemparse(s)
		if err != nil {
			log.Fatal("[Faild]: pemparse err: ",err)
		}
		sshAuth = append(sshAuth,method)
		break
	default:
		err := errors.New("invalid server method or it not empty!")
		return nil,err
	}
	return sshAuth,nil
}

func pemparse(s *Server) (ssh.AuthMethod,error) {
	key,err := ioutil.ReadFile(s.Key)
	if err != nil {
		return nil,err
	}
	var signer ssh.Signer
	if s.Password == "" {
		signer,err = ssh.ParsePrivateKey(key)
	} else {
		signer,err = ssh.ParsePrivateKeyWithPassphrase(key,[]byte(s.Password))
	}
	if err != nil {
		return nil,err
	}
	return ssh.PublicKeys(signer),nil
}

func (s *Server) ClientConnection() {
	auth,err := parseAuthMethods(s)
	if err != nil {
		log.Fatalf("[Faild]: %v",err)
		os.Exit(1)
	}

	config := &ssh.ClientConfig{
		User: s.User,
		Auth: auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := s.Ip + ":" + strconv.Itoa(s.Port)
	client,err := ssh.Dial("tcp",addr,config)
	if err != nil {
		log.Fatal("Failed to dial:",err)
		os.Exit(1)
	}

	session,err := client.NewSession()
	if err != nil {
		log.Fatal("session create err:",err)
		os.Exit(1)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState,err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal("create file descriptor faild: ",err)
	}

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	termWidth,termHeight,err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	defer terminal.Restore(fd,oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:			1,
		ssh.TTY_OP_ISPEED:	14400,
		ssh.TTY_OP_OSPEED:	14400,
	}

	if err := session.RequestPty("xterm-256color",termHeight,termWidth,modes);err != nil {
		log.Fatalf("request for pseudo terminal faild: %v",err)
	}

	if err := session.Shell();err != nil {
		log.Fatalf("faild to start shell: %v",err)
	}

	if err := session.Wait();err != nil {
		return
	}

}

