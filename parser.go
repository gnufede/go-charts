package main

import (
    "fmt"
    "strconv"
    // log "github.com/Sirupsen/logrus"
    "github.com/garyburd/redigo/redis"
    "strings"
    "time"
		"encoding/gob"
		"bytes"
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
        fmt.Println(date[0])
        counter += 1
    }

    result["date"] = dates

    for _, date := range dates {
        key := "Organizer:" + ORGANIZER + ":Event:" + EVENT + ":Channel:" + CHANNEL + ":Session:" + SESSION + ":Date:" + date

        values, values_err := redis.Int(redisScript.Do(redisConn, key))
        if values_err != nil {
            fmt.Println(values_err)
        }
        fmt.Println(values)
        amounts = append(amounts, strconv.Itoa(values))
    }
    result["amount"] = amounts

    fmt.Println(result)
    redisPool.Close()
	return GetBytes(result)
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
