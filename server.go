package main

import (
//	"fmt"
	"flag"
	"log"
	"strings"
	"net/http"

	"time"
	"strconv"
	"github.com/rs/cors" // used for fuck the cors rules :)

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
	initialData := Parse() //getDataForOrganizer(organizer)
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

		data := Parse() //getDataForOrganizer(organizer)
		err = c.WriteMessage(websocket.TextMessage, data)
		//err = c.WriteJSON(data)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}

}

func leadingZeros(origin int) string {
	origin_str := strconv.FormatInt(int64(origin), 10)

	if len(origin_str) == 1 {
		return "0" + origin_str
	}	
	return origin_str
}



func date_str() string {
	t := time.Now()

	result := strconv.FormatInt(int64(t.Year()), 10)
	result += "-"
	result += leadingZeros(int(t.Month()))
	result += "-"
	result += leadingZeros(t.Day())

	return result
}

func time_str() string {
	t := time.Now()	

	result := strconv.FormatInt(int64(t.Hour()), 10)
	result += ":"
	result += leadingZeros(t.Minute())

	return result
}



func update_ticket(w http.ResponseWriter, r *http.Request) {
  ticket_id := r.FormValue("ticket_id")
  price, _ := strconv.Atoi(r.FormValue("price"))
  price_it64 := int64(price)
	
	const fake_session_key = "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION
  fake_ticket_key := fake_session_key + ":TicketType:" + ticket_id + ":Date:" + date_str()


	pipe := client.Pipeline()

	// Add ticket to total quantity
	pipe.HIncrBy(fake_session_key +	":Date:" + date_str(), time_str(),			 1)

	// Add ticket to ticket type quantity
	pipe.HIncrBy(fake_ticket_key, time_str(), 1)

	// Increment ticket type
	pipe.HIncrBy(fake_ticket_key + ":Quantity", time_str(), 1)
	pipe.HIncrBy(fake_ticket_key + ":Amount", time_str(), price_it64)

  // Increment event totals
  const fake_event_key = "Organizer:" + ORGANIZER + ":Event:" + EVENT
  pipe.IncrBy(fake_event_key + ":TotalQuantity", 1)
  pipe.IncrBy(fake_event_key + ":TotalAmount", price_it64)

  // Increment event totals per channel
  fake_channel_with_date_key := fake_event_key + ":Channel:" + CHANNEL + ":Date:" + date_str()
  pipe.IncrBy(fake_channel_with_date_key + ":Quantity", 1)
  pipe.IncrBy(fake_channel_with_date_key + ":Amount", price_it64)

	pipe.Exec()


	client.Publish("1", "lets_go")
}



func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", sendData)

	mux := http.NewServeMux()
	mux.HandleFunc("/update_ticket", update_ticket)
	handlerCors := cors.Default().Handler(mux)

  if err := http.ListenAndServe(*addr, handlerCors); err != nil {
		log.Fatal("ListenAndServe:", err)
		client.Close()
	}
}
