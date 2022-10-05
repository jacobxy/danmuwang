package models

import (
	"douyu/ioclient"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type LuckInfo struct {
	Time string `json:"time"`
	From string `json:"from"`
	User string `json:"user"`
}

func WriteLuck(from string, user string) {
	l := LuckInfo{
		From: from,
		User: user,
	}
	l.Time = time.Now().Format("2006-01-02 15:04:05")
	fileName := path.Join("static", "file", "luck.json")
	b, _ := ioutil.ReadFile(fileName)
	total := make([]LuckInfo, 0, 10)
	json.Unmarshal(b, &total)
	total = append(total, l)

	f, err := os.OpenFile(path.Join("static", "file", "luck.json"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	res, _ := json.Marshal(total)
	f.WriteString(string(res) + "\n")

	if from == luckFrom && user == luckName {
		luckFrom = ""
		luckName = ""
		go ioclient.BroadCast("luck", map[string]interface{}{
			"from": luckFrom,
			"name": luckName,
		})
	}
}

var luckFrom string
var luckName string
var luckText string

func SetLuck(from, name, text string) {
	luckFrom = from
	luckName = name
	luckText = text
	go ioclient.BroadCast("luck", map[string]interface{}{
		"from": luckFrom,
		"name": luckName,
		"text": text,
	})
}
