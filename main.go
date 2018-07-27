package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"

	"github.com/arkezar/evehub/db"
	"github.com/arkezar/evehub/esi"
	"github.com/arkezar/evehub/model"
	"github.com/globalsign/mgo"
)

func main() {
	fmt.Println(Configuration.Get("database.host"))
	fmt.Println(Configuration.Get("database.name"))
	fmt.Println(Configuration.Get("database.user"))
	fmt.Println(Configuration.Get("database.pass"))
}

func processEmptyKills() {
	var processed int32

	db := db.New(Configuration)
	done := make(chan bool)
	errors := make(chan int32)
	kills := make(chan model.KillId)
	esi := esi.New(Configuration)

	go func() {
		for k := range kills {
			km := esi.GetKillmail(errors, k)
			k.Killmail = km[0]
			db.UpdateKillmail(k)
			atomic.AddInt32(&processed, 1)
		}
		done <- true
	}()

	go func() {
		timer := time.NewTicker(time.Minute)
		var total int32

		go func() {
			for t := range timer.C {
				fmt.Printf(" (%s Errors last minute: %d)", t.String(), atomic.LoadInt32(&total))
				atomic.StoreInt32(&total, 0)
			}
		}()

		for i := range errors {
			atomic.AddInt32(&total, i)
			if atomic.LoadInt32(&total) > 90 {
				panic("Too many errors")
			}
		}
	}()

	go func() {
		var avg float32
		for {
			i := atomic.SwapInt32(&processed, 0)
			avg = (avg + float32(i)) / 2
			fmt.Printf("\rProcessing: %f/s", avg)
			time.Sleep(time.Second)
		}
	}()

	test := db.GetUnfetchedKillmails()
	var item model.KillId
	for test.Next(&item) {
		kills <- item
	}

	for {
		<-done
		close(kills)
	}
}

func zkbWS() {
	url := url.URL{Scheme: "wss", Host: "zkillboard.com:2096"}
	ws, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	ws.WriteMessage(websocket.TextMessage, []byte("{\"action\":\"sub\",\"channel\":\"killstream\"}"))

	done := make(chan struct{})
	go func() {
		defer close(done)
		var km model.Killmail
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			json.Unmarshal([]byte(message), &km)
			fmt.Println(km)
		}
	}()

	for {
		select {
		case <-done:
			return
		}
	}
}

func parseCsvToMongo() {
	file, err := os.Open("C:\\Users\\arkez\\kills_esi.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	conn, err := mgo.Dial(Configuration.GetString("database.host"))

	if err != nil {
		panic(err)
	}

	defer conn.Close()
	db := conn.DB(Configuration.GetString("database.name"))

	err = db.Login(Configuration.GetString("database.user"), Configuration.GetString("database.pass"))

	if err != nil {
		panic(err)
	}

	c := db.C("kills")

	scanner := bufio.NewScanner(file)

	bulk := c.Bulk()
	i := 0
	for scanner.Scan() {
		if i == 10000 {
			fmt.Println("Insert")
			_, err := bulk.Run()
			if err != nil {
				panic(err)
			}
			i = 0
		}
		bulk.Insert(getKillModel(scanner.Text()))
		i++
	}
}

func getKillModel(s string) model.KillId {
	tokens := strings.Split(s, ",")
	id, _ := strconv.ParseInt(tokens[0], 10, 32)
	return model.KillId{EsiID: int(id), EsiHash: strings.Trim(tokens[1], "\"")}
}
