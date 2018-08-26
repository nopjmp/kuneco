package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"

	"cwtch.im/cwtch/connectivity/tor"
)

type Application struct {
	TorManager *tor.Manager
}

var app Application

//Page is a struct
type Page struct {
	Title   string
	Welcome string
	User    User
	Feed    []Post
}

var title = "Kuneco"
var homeENV string
var savePath, torPath string
var err error

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Println("> ", sig)
			if app.TorManager != nil {
				fmt.Println("Shutting down Tor process...")
				app.TorManager.Shutdown()
			}
			os.Exit(0)

		}
	}()
	if runtime.GOOS == "windows" {
		homeENV = os.Getenv("APPDATA")
		os.Mkdir(homeENV+"/kuneco/", os.ModePerm)
		os.Mkdir(homeENV+"/kuneco/tor", os.ModePerm)
	} else {
		homeENV = os.Getenv("HOME")
	}

	savePath = homeENV + "/kuneco/"
	if torPath == "" {
		torPath, err = exec.LookPath("tor")
	}

	if err != nil {
		log.Fatal("ERROR: tor could not be found on this system. Please install it in the system $PATH")
	}
	fmt.Println("> Tor found: ", torPath)

	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/follow", addFriend)
	http.HandleFunc("/friendRequest", friendRequest)
	http.HandleFunc("/accept", acceptFriendRequest)
	http.HandleFunc("/feed", feed)
	http.HandleFunc("/post", newPost)
	err = http.ListenAndServe("127.0.0.1:48486", nil)

	if err != nil {
		log.Fatal("ERROR: Listen failed")
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: title}
	t, _ := template.ParseFiles("index.tmpl")
	t.Execute(w, p)
}
