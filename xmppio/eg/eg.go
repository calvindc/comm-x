package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/calvindc/comm-x/xmppio"
)

var (
	server        = flag.String("server", "3.0.24.15:5222", "server")
	username      = flag.String("username", "user888@3.0.24.15", "username")
	password      = flag.String("password", "123456", "password")
	status        = flag.String("status", "xa", "status")
	statusMessage = flag.String("status-msg", "I for one welcome our new codebot overlords.", "status message")
	notls         = flag.Bool("notls", true, "No TLS")
	debug         = flag.Bool("debug", false, "debug output")
	session       = flag.Bool("session", false, "use server session")
)

func serverName(host string) string {
	return strings.Split(host, ":")[0]
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: example [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *username == "" || *password == "" {
		if *debug && *username == "" && *password == "" {
			fmt.Fprintf(os.Stderr, "no username or password were given; attempting ANONYMOUS auth\n")
		} else if *username != "" || *password != "" {
			flag.Usage()
		}
	}

	if !*notls {
		xmppio.DefaultConfig = &tls.Config{
			ServerName:         serverName(*server),
			InsecureSkipVerify: false,
		}
	}

	var xmppclient *xmppio.Client
	var err error
	options := xmppio.Options{
		Host:                         *server,
		User:                         *username,
		Password:                     *password,
		NoTLS:                        *notls,
		Debug:                        *debug,
		Session:                      *session,
		Status:                       *status,
		StatusMessage:                *statusMessage,
		DialTimeout:                  5 * time.Second,
		InsecureAllowUnencryptedAuth: true,
	}

	xmppclient, err = options.NewClient()
	if err != nil {
		log.Fatal(err)

	}
	xmppclient.PubsubSubscribeNode("", "chenyun2@16.163.154.147")

	go func() {
		for {
			chat, err := xmppclient.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmppio.Chat:
				fmt.Println("xmpp.Chat")
				fmt.Println(v.Remote, v.Text)
			case xmppio.Presence:
				fmt.Println("xmpp.Presence")
				fmt.Println(v.From, v.Show)
			}
		}
	}()

	for {
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimRight(line, "\n")

		tokens := strings.SplitN(line, " ", 2)
		fmt.Println(len(tokens))
		fmt.Println(tokens[0]) //eg:user888@3.0.24.15
		fmt.Println(tokens[1])
		if len(tokens) == 2 {
			xmppclient.Send(xmppio.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
			fmt.Println(xmppclient.JID())
		}

		//define presence parms for subSubscribe node
	}
}
