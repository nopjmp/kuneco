package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/openpgp"

	"golang.org/x/crypto/bcrypt"
)

//User profile information
type User struct {
	Username       string
	Password       string
	Hostname       string
	PrivateKey     string
	PGPEntity      *openpgp.Entity
	Friends        []Friend
	FriendRequests []Friend
	Posts          []Post
}

var loadedUserProfile User

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.FormValue("username")
	pass := []byte(r.FormValue("password"))
	pass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("> ERROR: bcrypt hashing failed: ", err)
	}
	if findUserProfile(user) {
		loadUserProfile(user, string(pass))
	} else {
		createUserProfile(user, string(pass))
	}
	err = startTor(torPath)
	if err != nil {
		log.Fatal("> ERROR: tor failed to start: ", err)
	} else {
		log.Println("> SUCCESS: Tor proxy successfully started")
	}
	feed := getPosts()
	p := &Page{Title: title, User: loadedUserProfile, Feed: feed}
	t, _ := template.ParseFiles("login.tmpl")
	t.Execute(w, p)

}
func findUserProfile(user string) bool {
	if _, err := os.Stat(savePath + user + "/user.json"); os.IsNotExist(err) {
		return false
	}
	return true
}

func loadUserProfile(user string, password string) {
	userLoad, _ := ioutil.ReadFile(savePath + user + "/user.json")
	json.Unmarshal(userLoad, &loadedUserProfile)
	log.Println("> SUCCESS:", loadedUserProfile.Username, "LOGGED IN!")
	ioutil.WriteFile(savePath+loadedUserProfile.Username+"/private_key", []byte(loadedUserProfile.PrivateKey), 0644)
	ioutil.WriteFile(savePath+loadedUserProfile.Username+"/hostname", []byte(loadedUserProfile.Hostname), 0644)

}

func saveUserProfile() {
	saveJSON, _ := json.Marshal(loadedUserProfile)
	ioutil.WriteFile(savePath+loadedUserProfile.Username+"/user.json", saveJSON, 0644)
}

func createUserProfile(user string, password string) {
	if user != "" {
		userProfile := new(User)
		userProfile.Username = user
		userProfile.Password = password
		hostname, privateKey := generateOnionAddress()
		userProfile.Hostname = hostname
		userProfile.PrivateKey = privateKey
		userProfile.PGPEntity, _ = openpgp.NewEntity(user, "", "", nil)
		loadedUserProfile = *userProfile
		os.Mkdir(savePath+user, os.ModePerm)
		saveUserProfile()
		loadUserProfile(user, password)
	}
}
