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
const DATE_KEY = "2015-10-09"

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

        redis.log(redis.LOG_WARNING, "SUM TOTAL ".. KEYS[1] .. ":" .. sum)
        return sum
    `

	ticketTypes := make(map[int]string)
	ticketTypes[1] = "General"
	ticketTypes[2] = "Infantil"
	ticketTypes[3] = "Jubilados"
	ticketTypes[4] = "Gratuita"
	channelTypes := make(map[int]string)
	channelTypes[1] = "Online"
	channelTypes[2] = "BoxOffice"
	channelTypes[3] = "IFrame"

	result := make(map[string]map[string][]string)
	for i := 1; i < 8; i++ {
		result[strconv.Itoa(i)] = make(map[string][]string)
	}

	var dates []string
	var minutes []string
	var minutesShow []string

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

	start := time.Now()
	minutes = append(minutes, start.Format("15:04"))
	minutesShow = append(minutesShow, start.Format("15:04"))
	for i := 0; i < 4; i++ {
		start = start.Add(-1 * time.Minute)
		minutes = append(minutes, start.Format("15:04"))
		minutesShow = append(minutesShow, start.Format("15:04"))
	}

	result["5"]["date"] = dates
	result["7"]["minutes"] = minutesShow

	for _, date := range dates {
		valueA := getChannelDataTotals(date, "Amount", channelTypes, redisScript, redisConn)
		sessionAmount := valueA / 100
		result["5"]["Amount"] = append(result["5"]["Amount"], strconv.Itoa(sessionAmount))

		// Increment week amount
		weekAmount += sessionAmount

		valueQ := getChannelDataTotals(date, "Quantity", channelTypes, redisScript, redisConn)
		sessionQuantity := valueQ
		result["5"]["Total"] = append(result["5"]["Total"], strconv.Itoa(sessionQuantity))
		weekQuantity += sessionQuantity

		for id, name := range ticketTypes {
			ticketQuantity := getTicketTypeTotals(date, "Quantity", id, channelTypes, redisScript, redisConn)
			result["5"][name] = append(result["5"][name], strconv.Itoa(ticketQuantity))
		}

		for channel, _ := range channelTypes {
			channelQuantity := getTotalPerChannel(date, "Quantity", channel, redisScript, redisConn)
			channelWeekQuantity[channel] += channelQuantity
		}

	}

	for _, minute := range minutes {
		minuteAmount := getChannelDataTotalsPerMinute(minute, "Amount", channelTypes, redisConn)
		minuteAmount = minuteAmount / 100
		result["7"]["Amount"] = append(result["7"]["Amount"], strconv.Itoa(minuteAmount))

		minuteQuantity := getChannelDataTotalsPerMinute(minute, "Quantity", channelTypes, redisConn)
		result["7"]["Total"] = append(result["7"]["Total"], strconv.Itoa(minuteQuantity))

		for id, name := range ticketTypes {
			ticketQuantity := getTicketTypeTotalsPerMinute(minute, "Quantity", id, channelTypes, redisConn)
			result["7"][name] = append(result["7"][name], strconv.Itoa(ticketQuantity))
		}

	}

	// Week amount and quantity
	result["1"]["Value"] = append(result["1"]["Value"], strconv.Itoa(weekQuantity))
	result["2"]["Value"] = append(result["2"]["Value"], strconv.Itoa(weekAmount))
	// Event total quantity
	eventTotalQuantityKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":TotalQuantity"
	totalQuantity, _ := redis.Int(redisConn.Do("GET", eventTotalQuantityKey))
	result["3"]["Value"] = append(result["3"]["Value"], strconv.Itoa(totalQuantity))
	// Event total amount
	eventTotalAmountKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":TotalAmount"
	totalAmount, _ := redis.Int(redisConn.Do("GET", eventTotalAmountKey))
	result["4"]["Value"] = append(result["4"]["Value"], strconv.Itoa(totalAmount/100))
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

func getChannelDataTotals(date string, keyFragments string, channelTypes map[int]string, redisScript *redis.Script, redisConn redis.Conn) int {
	result := 0
	for channel, _ := range channelTypes {
		key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + strconv.Itoa(channel) + ":Session:" + SESSION + ":Date:" + date + ":" + keyFragments
		values, values_err := redis.Int(redisScript.Do(redisConn, key))
		if values_err != nil {
			values = 0
		}
		result += values
	}
	return result
}

func getChannelDataTotalsPerMinute(minute string, keyFragments string, channelTypes map[int]string, redisConn redis.Conn) int {
	result := 0
	for channel, _ := range channelTypes {
		key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + strconv.Itoa(channel) + ":Session:" + SESSION + ":Date:" + DATE_KEY + ":" + keyFragments
		values, values_err := redis.Int(redisConn.Do("HGET", key, minute))
		if values_err != nil {
			values = 0
		}
		result += values
	}
	return result
}

func getTotalPerChannel(date string, keyFragments string, channel int, redisScript *redis.Script, redisConn redis.Conn) int {
	key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + strconv.Itoa(channel) + ":Session:" + SESSION + ":Date:" + date + ":" + keyFragments
	values, values_err := redis.Int(redisScript.Do(redisConn, key))
	if values_err != nil {
		values = 0
	}

	return values
}

func getTicketTypeTotals(date string, keyType string, ticketId int, channelTypes map[int]string, redisScript *redis.Script, redisConn redis.Conn) int {
	result := 0
	for channel, _ := range channelTypes {
		ticketTypeKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + strconv.Itoa(channel) + ":Session:" + SESSION + ":TicketType:" + strconv.Itoa(ticketId) + ":Date:" + date + ":" + keyType
		values, values_err := redis.Int(redisScript.Do(redisConn, ticketTypeKey))
		if values_err != nil {
			values = 0
		}
		result += values
	}
	return result
}

func getTicketTypeTotalsPerMinute(minute string, keyType string, ticketId int, channelTypes map[int]string, redisConn redis.Conn) int {
	result := 0
	for channel, _ := range channelTypes {
		ticketTypeKey := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + strconv.Itoa(channel) + ":Session:" + SESSION + ":TicketType:" + strconv.Itoa(ticketId) + ":Date:" + DATE_KEY + ":" + keyType
		values, values_err := redis.Int(redisConn.Do("HGET", ticketTypeKey, minute))
		if values_err != nil {
			values = 0
		}
		result += values
	}
	return result
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
