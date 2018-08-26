package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//Post contains information about a post
type Post struct {
	PostedBy string
	Content  string
	Time     time.Time
}

func newPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.FormValue("content")
	post := new(Post)
	post.Content = content
	post.PostedBy = loadedUserProfile.Hostname
	post.Time = time.Now()
	loadedUserProfile.Posts = append(loadedUserProfile.Posts, *post)
	saveUserProfile()
	if len(latestPosts) < 10 {
		latestPosts = append(latestPosts, *post)
	} else {
		latestPosts = append(latestPosts[:0], latestPosts[1:]...)
		latestPosts = append(latestPosts, *post)
	}
}

var latestPosts []Post

func feed(w http.ResponseWriter, r *http.Request) {
	feed := loadedUserProfile.Posts
	JSON, _ := json.Marshal(feed)
	w.Write(JSON)
}

func getPosts() []Post {
	var friendPosts []Post
	for _, friend := range loadedUserProfile.Friends {
		var posts []Post
		resp, err := sendTorRequest("GET", "http://"+friend.Hostname+"/feed")
		if err == nil {
			json.Unmarshal(resp, &posts)
			friendPosts = append(friendPosts, posts...)
		} else {
			fmt.Println("> ERROR: ", err)
		}

	}

	return friendPosts
}
