package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

var redisPool *redis.Pool

const ORGANIZER = "1"
const EVENT = "1"
const SESSION = "1"
const CHANNEL = "1"

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Parse() []byte {
	const dateValuesScript = `
        local values = redis.call("HVALS", KEYS[1])
        local sum = 0

        for _,val in ipairs(values) do
            sum = sum + tonumber(val)
        end

        return sum
    `

	ticketTypes := make(map[int]string)
	ticketTypes[1] = "General"
	channelTypes := make(map[int]string)
	channelTypes[1] = "Online"
	channelTypes[2] = "BoxOffice"
	channelTypes[3] = "IFrame"

	result := make(map[string]map[string][]string)
	for i := 1; i < 7; i++ {
		result[strconv.Itoa(i)] = make(map[string][]string)
	}
	var dates []string
	var amounts []string
	weekAmount := 0
	weekQuantity := 0
	channelWeekQuantity := make(map[int]int)

	redisScript := redis.NewScript(1, dateValuesScript)
	redisPool := newPool("localhost:6379")
	redisConn := redisPool.Get()
	defer redisConn.Close()

	counter := 0
	for counter < 7 {
		dateTime := time.Now().AddDate(0, 0, -1*counter)
		date := strings.Split(dateTime.String(), " ")
		dates = append(dates, date[0])
		counter += 1
	}

	result["5"]["date"] = dates

	for _, date := range dates {
		// Get total amount sales group by date
		key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION + ":Date:" + date + ":Amount"
		values, values_err := redis.Int(redisScript.Do(redisConn, key))
		if values_err != nil {
			fmt.Println(values_err)
			values = 0
		}
		amounts = append(amounts, strconv.Itoa(values))
		result["5"]["amount"] = amounts

		// Increment week amount
		weekAmount += values

		// // Increment week quantity
		key = "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION + ":Date:" + date + ":Quantity"
		dayQuantity, dayQuantity_err := redis.Int(redisScript.Do(redisConn, key))
		if dayQuantity_err != nil {
			fmt.Println(values_err)
			dayQuantity = 0
		}
		weekQuantity += dayQuantity

		for id, name := range ticketTypes {
			// Get total quantity group by ticket type
			ticketTypeKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION + ":TicketType:" + strconv.Itoa(id) + ":Date:" + date + ":Amount"
			values, values_err := redis.Int(redisScript.Do(redisConn, ticketTypeKey))
			if values_err != nil {
				fmt.Println(values_err)
				values = 0
			}
			result["5"][name] = append(result["5"][name], strconv.Itoa(values))
		}

		for channel, _ := range channelTypes {
			channelTypeKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + strconv.Itoa(channel) + ":Session:" + SESSION + ":Date:" + date + ":Quantity"
			channelQuantity, channelQuantity_err := redis.Int(redisConn.Do("GET", channelTypeKey))
			if channelQuantity_err != nil {
				channelQuantity = 0
			}
			channelWeekQuantity[channel] += channelQuantity
		}
	}

	// Week amount and quantity
	result["1"]["Value"] = append(result["1"]["Value"], strconv.Itoa(weekQuantity))
	result["2"]["Value"] = append(result["2"]["Value"], strconv.Itoa(weekAmount))
	// Event total quantity
	eventTotalQuantityKey := "Organizer" + ORGANIZER + ":Event:" + EVENT + "TotalQuantity"
	totalQuantity, _ := redis.Int(redisConn.Do("GET", eventTotalQuantityKey))
	result["3"]["Value"] = append(result["3"]["Value"], strconv.Itoa(totalQuantity))
	// Event total amount
	eventTotalAmountKey := "Organizer" + ORGANIZER + ":Event:" + EVENT + "TotalAmount"
	totalAmount, _ := redis.Int(redisConn.Do("GET", eventTotalAmountKey))
	result["4"]["Value"] = append(result["4"]["Value"], strconv.Itoa(totalAmount))
	// Channel quantity distribution
	for channel, channelName := range channelTypes {
		result["6"][channelName] = append(result["6"][channelName], strconv.Itoa(channelWeekQuantity[channel]))
	}

	output, o_err := json.Marshal(result)
	if o_err != nil {
		panic("Error generating JSON")
	}
	redisPool.Close()
	fmt.Println(result)
	return output
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     25,
		MaxActive:   12500,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
