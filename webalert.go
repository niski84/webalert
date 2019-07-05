package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

var to string
var from string
var pwd string
var subject string
var msg string
var searchTextInclude string
var searchTextExclde string
var website string

func getFlags() {
	flag.StringVar(&to, "t", "", "destination Internet mail address")
	flag.StringVar(&from, "f", "", "the sender's GMail address")
	flag.StringVar(&pwd, "p", "", "the sender's password")
	flag.StringVar(&subject, "s", "", "subject line of email")
	flag.StringVar(&msg, "m", "", "a one-line email message")
	flag.StringVar(&searchTextExclde, "x", "", "alert when this text is excluded on webpage")
	flag.StringVar(&searchTextInclude, "i", "", "alert when this text is included on webpage")
	flag.StringVar(&website, "w", "", "website example: http://example.com")
	flag.Usage = func() {
		fmt.Printf("Syntax:\n\tWebalert [flags]\nwhere flags are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if searchTextExclde == "" && searchTextInclude == "" {
		fmt.Println("At least -i or -x required, but not both.")
		flag.Usage()
		os.Exit(1)
	}

	if website == "" {
		fmt.Println("website to parse must be included (-w)")
		flag.Usage()
		os.Exit(1)
	}



}

func sendEmail() {

	emailbody := "To: " + to + "\r\nSubject: " +
		subject + "\r\n\r\n" + msg
	auth := smtp.PlainAuth("", from, pwd, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from,
		[]string{to}, []byte(emailbody))
	if err != nil {
		log.Fatal(err)
	}

}
func main() {

	getFlags()

	fmt.Printf("webalert v 1.0\n")

	for {
		fmt.Printf("fetching %v..\n", website)
		textFound := CheckWebsite(website, searchTextInclude, searchTextExclde)
		if textFound {
			fmt.Printf("match!\n")
			if to != "" {
				sendEmail()
			} else {
				fmt.Printf("skip sending email. missing arguments")
			}
			os.Exit(0)
		}
		time.Sleep(10 * time.Second)
	}

}

func CheckWebsite(website string, searchTextInclude string, searchTextExclde string) bool {

	resp, err := http.Get(website)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if searchTextExclde != "" {
		if !bytes.Contains(body, []byte(searchTextExclde)) {
			return true
		} else {
			return false
		}
	}

	if searchTextInclude != "" {
		if bytes.Contains(body, []byte(searchTextInclude)) {
			return true
		} else {
			return false
		}
	}
	return false
}
