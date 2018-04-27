package main

import (
	"io/ioutil"
	"log"
	"gossh/ssh"
	"encoding/json"
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
)

func main() {
	b,err := ioutil.ReadFile("config/server.json")
	if err != nil {
		log.Fatalf("[Faild]: read server.json got err: %v",err)
	}

	var data []ssh.Server
	if err := json.Unmarshal(b,&data);err != nil {
		log.Fatal(err)
	}

	fmt.Println("show server list table:")
	fmt.Println("========================")
	for i,s := range data {
		fmt.Printf("%d: %v\n",i,s.Name)
	}
	fmt.Println("========================")

	inputReader := bufio.NewReader(os.Stdin)

	fmt.Println("Please input your choice:")
	input,err := inputReader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}
	choice, err := strconv.Atoi(strings.Replace(input,"\n","",-1))
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Your input is: %v", choice)
	data[choice].ClientConnection()
}
