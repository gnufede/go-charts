package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

func Parse() ([]byte, error) {
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
	result := make(map[string][]string)
	var dates []string
	var amounts []string

	redisScript := redis.NewScript(1, dateValuesScript)
	redisPool = newPool("localhost:6379")
	redisConn := redisPool.Get()
	defer redisConn.Close()

	counter := 0
	for counter < 7 {
		dateTime := time.Now().AddDate(0, 0, -1*counter)
		date := strings.Split(dateTime.String(), " ")
		dates = append(dates, date[0])
		counter += 1
	}
	result["date"] = dates

	for _, date := range dates {
		// Get total amount sales group by date
		key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION + ":Date:" + date
		values, values_err := redis.Int(redisScript.Do(redisConn, key))
		if values_err != nil {
			panic(values_err)
		}
		amounts = append(amounts, strconv.Itoa(values))
		result["amount"] = amounts

		for id, name := range ticketTypes {
			// Get total quantity group by ticket type
			ticketTypeKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION + ":TicketType:" + strconv.Itoa(id) + ":Date:" + date
			values, values_err := redis.Int(redisScript.Do(redisConn, ticketTypeKey))
			if values_err != nil {
				panic(values_err)
			}
			result[name] = append(result[name], strconv.Itoa(values))
		}

	}

	output, o_err := json.Marshal(result)
	if o_err != nil {
		panic("Error generating JSON")
	}
	redisPool.Close()
	return GetBytes(output)
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
