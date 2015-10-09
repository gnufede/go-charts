package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	updated = false
)

func getDataForOrganizer(organizer string) []byte {
	visitsKey := strings.Join([]string{"Organizer", organizer}, ":")
	visits, err := client.HGetAll(visitsKey).Result()
	if err != nil {
		log.Println("hgetall: ", err)
	}
	return []byte(visits[1])
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

	initialData := Parse()
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

		if message.Payload == "send" {
			data := Parse() //getDataForOrganizer(organizer)
			err = c.WriteMessage(websocket.TextMessage, data)
			//err = c.WriteJSON(data)
			if err != nil {
				log.Println("write:", err)
				break
			}

		}
	}

}

func date_str() string {
	return time.Now().Format("2006-01-02")
}

func time_str() string {
	return time.Now().Format("15:04")
}

func update_ticket(w http.ResponseWriter, r *http.Request) {
	ticket_id := r.FormValue("ticket_id")
	channelType := r.FormValue("channel")

	price, _ := strconv.Atoi(r.FormValue("price"))
	price_it64 := int64(price)

	fake_session_key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + channelType + ":Session:" + SESSION
	fake_ticket_key := fake_session_key + ":TicketType:" + ticket_id + ":Date:" + date_str()

	pipe := client.Pipeline()

	// Add ticket to total quantity
	pipe.HIncrBy(fake_session_key+":Date:"+date_str(), time_str(), 1)

	// Add ticket to ticket type quantity
	pipe.HIncrBy(fake_ticket_key, time_str(), 1)

	// Increment ticket type
	pipe.HIncrBy(fake_ticket_key+":Quantity", time_str(), 1)
	pipe.HIncrBy(fake_ticket_key+":Amount", time_str(), price_it64)

	// Increment session quantity and amount
	pipe.HIncrBy(fake_session_key+":Date:"+date_str()+":Quantity", time_str(), 1)
	pipe.HIncrBy(fake_session_key+":Date:"+date_str()+":Amount", time_str(), price_it64)

	// Increment event totals
	const fake_event_key = "Organizer:" + ORGANIZER + ":Event:" + EVENT
	pipe.IncrBy(fake_event_key+":TotalQuantity", 1)
	pipe.IncrBy(fake_event_key+":TotalAmount", price_it64)

	// Increment event totals per channel
	fake_channel_with_date_key := fake_event_key + ":Channel:" + channelType + ":Date:" + date_str()
	pipe.IncrBy(fake_channel_with_date_key+":Quantity", 1)
	pipe.IncrBy(fake_channel_with_date_key+":Amount", price_it64)

	pipe.Exec()

	client.Publish("1", "lets_go")
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	configure_routes()

	go read_updates()
	go update_data_500()
	go update_data_1s()

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
		client.Close()
	}
}

func read_updates() {
	pubsub, err := client.Subscribe("1")
	if err != nil {
		log.Print("subscribe:", err)
		return
	}
	for {
		message, err := pubsub.ReceiveMessage()
		if err != nil {
			log.Println("read:", err)
		}
		if message.Payload == "lets_go" {
			log.Printf("reader recv: %s", message.Payload)
			updated = true
		}
	}
}

func update_data_500() {
	c := time.Tick(500 * time.Millisecond)
	for range c {
		if updated {
			client.Publish("1", "send")
			log.Println("sent")
			updated = false
		}
	}
}

func update_data_1s() {
	c := time.Tick(1 * time.Minute)
	for range c {
		client.Publish("1", "send")
	}
}

func configure_routes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/simulator", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "simulator.html")
	})

	http.HandleFunc("/update_ticket", update_ticket)
	http.HandleFunc("/ws", sendData)

	http.Handle("/sounds/", http.StripPrefix("/sounds/", http.FileServer(http.Dir("sounds"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
}
