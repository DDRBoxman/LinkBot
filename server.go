package main

type Server struct {
	
	Server string
	Nick string
	Port int
	Password string
	SSL bool
	Channels []string
	Blacklist []string
	MinLength int

}

