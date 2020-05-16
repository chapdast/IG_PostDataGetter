package main

import (
	"fmt"
	"github.com/ahmdrz/goinsta/v2"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var i *goinsta.Instagram
var err error

type PostData struct {
	ID           int
	Username     string
	Link         string
	LikeCount    int
	CommentCount int
}


func LoginToIG(username, password string) {

	//Login to Instagram
	i, err = goinsta.Import(fmt.Sprintf("./_ipguser_%s.json", username))
	if err != nil {
		i = goinsta.New(username, password)
		err = i.Login()
		if err != nil {
			log.Fatalf("cant login to insta > %v", err)
		}
		i.Export(fmt.Sprintf("./_ipguser_%s.json", username));
	}
}

var ID = 0
var out *os.File
var tag string

func main() {
	fmt.Print(
		`*************************************************
Instagram Post Analytics Data Extractor
*************************************************
`)
	if len(os.Args) < 2 {
		log.Fatal("IG TAG is not Defined. Exiting.")
	}

	tag = os.Args[1]
	out, err = os.Create("./_" + tag + ".csv")
	if err != nil {
		log.Fatal("Can't Create Output file. Exiting.")
	}
	out.WriteString("id, username,link,likes,comments\n")
	pc := make(chan PostData)

	username:= os.Getenv("IPG_USERNAME")
	password:= os.Getenv("IPG_PASSWORD")
	if (username == "" || password == ""){
		log.Fatal("IPG_USERNAME/IPG_PASSWORD is not defined in ENV. ")
	}

	LoginToIG(username, password)
	for i == nil {
		fmt.Println("Not Initialized")
		time.Sleep(4 * time.Second)
	}

	go func() {
		GetPosts(tag, pc)
		close(pc)
		pc = nil
	}()
	for pc != nil {
		select {
		case pd := <-pc:
			InsertPost(pd)
		}
	}
	defer out.Close()
}
func InsertPost(pd PostData) {
	data := fmt.Sprintf("%d, %s, %s, %d, %d\n", pd.ID, pd.Username, pd.Link, pd.LikeCount, pd.CommentCount)
	if err != nil {
		fmt.Println("cant read")
	}
	n, err := out.WriteString(data)
	if err != nil {
		fmt.Println("cant write ", n, err.Error())
	}
}
func GetPosts(username string, pc chan PostData) {
	fmt.Println("Get Feed Posts of: ", strings.ToUpper(username))
	user, err := i.Profiles.ByName(username)
	if err != nil {
		log.Fatalf("Can't Get User > %v")
	}
	feed := user.Feed()

	for feed.Next() {
		for _, item := range feed.Items {
			pc <- PostData{
				ID:           ID,
				Username:     username,
				Link:         "https://www.instagram.com/p/" + item.Code,
				LikeCount:    item.Likes,
				CommentCount: item.CommentCount,
			}
			ID++
		}

		randomTime := rand.Intn(10)
		time.Sleep(time.Duration(randomTime) * time.Second)
		fmt.Printf("Hold for %d second(s). ", randomTime)
		fmt.Printf("Written %d of %d Media | %d%%\n", ID, user.MediaCount,ID*100/user.MediaCount)

	}
}
