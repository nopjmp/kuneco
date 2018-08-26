package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
)

type Friend struct {
	Name      string
	Hostname  string
	PublicKey string
}

func sendFriendRequest(hostname string) {
	pubkey := loadedUserProfile.PGPEntity.PrimaryKey.PublicKey

	pubkeyStr := hex.EncodeToString([]byte(fmt.Sprint(pubkey)))
	log.Println("> SUCCESS encoding public key", pubkeyStr)
	sendTorRequest("GET", "http://"+hostname+"/friendRequest?name="+loadedUserProfile.Username+"&hostname="+loadedUserProfile.Hostname+"&pubkey="+pubkeyStr)

}

func alreadyFriends(hostname string) bool {
	for i := range loadedUserProfile.Friends {
		if loadedUserProfile.Friends[i].Hostname == hostname {
			return true
		}
	}
	return false
}

func alreadyHaveFriendRequest(hostname string) bool {
	for i := range loadedUserProfile.Friends {
		if loadedUserProfile.Friends[i].Hostname == hostname {
			return true
		}
	}
	return false
}
func acceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	for i, friendRequest := range loadedUserProfile.FriendRequests {
		if friendRequest.Hostname == r.FormValue("hostname") {
			loadedUserProfile.Friends = append(loadedUserProfile.Friends, friendRequest)
			loadedUserProfile.FriendRequests = append(loadedUserProfile.FriendRequests[:i], loadedUserProfile.FriendRequests[i+1:]...)
			saveUserProfile()
		}
	}
}

func addFriend(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	sendFriendRequest(r.FormValue("hostname"))
	fmt.Fprintln(w, "done")
}

func friendRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if !alreadyFriends(r.FormValue("hostname")) && !alreadyHaveFriendRequest(r.FormValue("hostname")) {
		fReq := new(Friend)
		fReq.Name = r.FormValue("name")
		fReq.Hostname = r.FormValue("hostname")
		fReq.PublicKey = r.FormValue("pubkey")

		loadedUserProfile.FriendRequests = append(loadedUserProfile.FriendRequests, *fReq)
		saveUserProfile()
		fmt.Println("> Friend request: ", r.FormValue("name"), r.FormValue("hostname"), r.FormValue("pubkey"))
		fmt.Fprintln(w, "done")
	}
}
