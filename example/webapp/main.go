package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
)

// TODO: server with user create/edit/delete
const addr = "localhost:1985"

var (
	isServer = flag.Bool("s", false, "is server")
	isClient = flag.Bool("c", true, "is client")
	debug    = flag.Bool("d", false, "debug / log")

	doSignup  = flag.Bool("signup", false, "signup")
	doProfile = flag.Bool("profile", false, "profile")
	signup    = flag.String("signupInfo", "one:1234:name:phone:agency", "username:password:name:phone:user-type")
	login     = flag.String("login", "one:1234", "username:password")
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()

	if *isServer {
		srv := newServer()
		log.Fatal(srv.s.Run(addr))
	}

	if *isClient {
		c := newClient()
		if *doSignup {
			j, _ := json.Marshal(c.Signup())
			fmt.Printf("%s\n", j)
		}
		if *doProfile {
			j, _ := json.Marshal(c.Login())
			fmt.Printf("%s\n", j)

			j, _ = json.Marshal(c.Profile())
			fmt.Printf("%s\n", j)

		}
	}

}
