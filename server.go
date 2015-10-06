package main

import (
//	"fmt"
	"flag"
	"log"
	"strings"
	"net/http"
	//"net/url"
	"github.com/gorilla/websocket"
	"gopkg.in/redis.v3"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	client = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // no password set
    DB:       0,  // use default DB
  })
)


func getDataForOrganizer(organizer string) []byte {
	visitsKey := strings.Join([]string{"Organizer", organizer}, ":")
	visits , err := client.HGetAll(visitsKey).Result()
	if err != nil {
		log.Println("hgetall: ", err)
	}
	return []byte(visits[1])

	/*
  return []byte(`
		{
				"date": [
						"2013-01-01",
						"2013-01-02",
						"2013-01-03",
						"2013-01-04",
						"2013-01-05",
						"2013-01-06",
						"2013-01-07"
				],
				"amount": [
						1500,
						1000,
						3000,
						4000,
						0,
						2500,
						3000
				],
				"children ticket": [
						10,
						50,
						25,
						100,
						0,
						0,
						100
				],
				"adult ticket": [
						20,
						150,
						75,
						300,
						150,
						250,
						250
				]
		}
  `)
  */
}


func sendData(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	values := r.URL.Query()
	log.Println("Values:", values)
	organizer := values.Get("organizer")

	//TODO
	initialData, err := Parse() //getDataForOrganizer(organizer)
	err = c.WriteMessage(websocket.TextMessage, initialData)
	if err != nil {
		log.Println("write:", err)
		return
	}

	pubsub, err := client.Subscribe(organizer)
	log.Println("channel:", organizer)
	if err != nil {
		log.Print("subscribe:", err)
		return
	}

	defer c.Close()
	defer pubsub.Close()
	for {
		message, err := pubsub.ReceiveMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message.Payload)

		data, err := Parse() //getDataForOrganizer(organizer)
		err = c.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", sendData)

  if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
		client.Close()
	}
}
