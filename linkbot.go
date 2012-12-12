package main

import (
	"flag"
	"encoding/json"
	"os"
	"io/ioutil"
	"fmt"
	"crypto/tls"
	"regexp"
	"strings"
	irc "github.com/fluffle/goirc/client"
	"github.com/NickPresta/GoURLShortener"
) 

var urlRegex *regexp.Regexp

func initServerConnection(server Server, quit chan bool) {
	c := irc.SimpleClient(server.Nick)
    	
    	c.SSL = server.SSL
	c.SSLConfig = &tls.Config{InsecureSkipVerify: true} 

	c.AddHandler("connected",
	func(conn *irc.Conn, line *irc.Line) { 
		for _,channel := range server.Channels {
			conn.Join(channel) 
		}
	})
	c.AddHandler("disconnected", func(conn *irc.Conn, line *irc.Line) { quit <- true })
	c.AddHandler("privmsg", func(conn *irc.Conn, line *irc.Line) {
		matches := urlRegex.FindAllString(line.Args[1] ,-1)
		for _,match := range matches {
			if (len(match) >= server.MinLength) {
				blacklist := false 

				for _,item := range server.Blacklist {

					if strings.Contains(match, item) {
						blacklist = true
						continue
					}
				}

				if blacklist {
					continue
				}	
				
				uri, err := goisgd.Shorten(match)
				if err != nil {
					continue
				}
				conn.Privmsg(line.Args[0], uri)
			}
		}
	})

	// Tell client to connect
	if err := c.Connect(fmt.Sprintf("%s:%d", server.Server, server.Port), server.Password)
	err != nil {
		fmt.Printf("Connection error: %v\n", err)
	}
}

func main() {
	var configFile = flag.String("config", "./channels.config", "The file containing the server info")
	flag.Parse()

	var err error
	urlRegex,err = regexp.Compile(`(http|https|ftp|ftps)\://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,4}(/\S*)?`)
	if err != nil {
		fmt.Printf("Regex failed to compile: %v", err)
		os.Exit(1)
	}

	quit := make(chan bool)

	file, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)	
		os.Exit(1)
	}

	var config []Server
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Printf("Failed to parse config. %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Results: %v\n", config)

	for _,server := range config {
		initServerConnection(server, quit)
	}

	// wait on quit channel
	<-quit
}
