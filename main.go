package main

import (
	"fmt"
	"time"
    "net/http"
	"github.com/barcatfigaro/moodmessage/search"
	"github.com/barcatfigaro/moodmessage/stringsem"
	"github.com/barcatfigaro/moodmessage/twPull"
	twilio "github.com/saintpete/twilio-go"
)
import "encoding/json"

//var goodmessages = []twilio.Message{}
var goodmessages = []string{}

type response struct {
    Messages []string `json:"messages"`
}


func handler(w http.ResponseWriter, r *http.Request) {
//    messageStrings := []string{}
//    for _,msg := range goodmessages {
//        messageStrings = append(messageStrings, msg.Body)
//    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
    messages := &response{
//        Messages: messageStrings,
        Messages: goodmessages,
    }
    a, err := json.Marshal(messages)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    w.Write(a)
}

func main() {
	now := time.Now()
	spider := search.NewSpider()
	search.RunSpider(spider, "https://www.reddit.com/r/UBC/")

    http.HandleFunc("/messages", handler)
    go http.ListenAndServe(":8080",nil)

    messages := []twilio.Message{}
    newmessages := []string{}
    fmt.Println("test")
    for {
        messages, newmessages = twPull.GetMessages(now,messages)
        for _,msg := range newmessages {
            fmt.Println(msg)
            if stringsem.IsGood(msg){
                fmt.Println("that was a good message")
                goodmessages = append(goodmessages,msg)
            } else {
                fmt.Println("that was bad message")
            }
        }
    }
}
