package models

import (
	"douyu/ioclient"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Danmu struct {
	From string `json:"from"`
	User string `json:"user"`
	Text string `json:"text"`
	Time int64  `json:"time"`
}

func (d *Danmu) Key() string {
	return d.From + d.User + d.Text + fmt.Sprintf("%d", d.Time/2)
}

type DanmuLive struct {
	All   []Danmu
	Index int
}

func NewDanmuLive() *DanmuLive {
	return &DanmuLive{
		All: make([]Danmu, 0, 1000),
	}
}

func (dl *DanmuLive) GetCur() []Danmu {
	res := dl.All[dl.Index:]
	dl.Index = len(dl.All)
	return res
}

const (
	MaxChan = 5
)

var (
	Dic     map[string]int
	Total   map[string]*DanmuLive
	command string
	dmutex  sync.RWMutex
	danmuch [MaxChan]chan Danmu
	Count   int
)

func init() {
	Dic = make(map[string]int, 1000)
	Total = make(map[string]*DanmuLive, 5)
	for i := 0; i < MaxChan; i++ {
		danmuch[i] = make(chan Danmu, 100)
	}
	go Broadcast()
}

func Broadcast() {
	for {
		res := make([]Danmu, 0, 100)
		c := time.After(200 * time.Millisecond)
	L:
		for {
			select {
			case v := <-danmuch[rand.Int()%MaxChan]:
				res = append(res, v)
			case <-c:
				break L
			}
		}
		if len(res) != 0 {
			ioclient.BroadCast("danmu", res)
		}
	}
}

func SetCommand(cmd string) {
	command = cmd
	dmutex.Lock()
	defer dmutex.Unlock()
	Dic = make(map[string]int, 1000)
	Total = make(map[string]*DanmuLive, 5)
}

func RecDanmu(ds []Danmu) {
	dmutex.Lock()
	defer dmutex.Unlock()
	from := ""
	for _, v := range ds {
		if command != "" && !strings.Contains(v.Text, command) {
			continue
		}
		from = v.From
		key := v.Key()
		if _, ok := Dic[key]; ok {
			continue
		}
		Dic[key] += 1

		if _, ok := Total[from]; !ok {
			Total[from] = NewDanmuLive()
		}
		Count++
		danmuch[Count%MaxChan] <- v
		Total[from].All = append(Total[from].All, v)
	}

	var l int
	if _, ok := Total[from]; ok {
		l = len(Total[from].All)
	}
	go ioclient.BroadCast("count", map[string]interface{}{
		"from":  from,
		"count": l,
	})

}

func GetDanmu(from string) []Danmu {
	return Total[from].All
}

func StatisticsDanmu(keys []string) string {
	dmutex.RLock()
	defer dmutex.RUnlock()
	keysCnt := make(map[string]int, len(keys))
	total := 0
	for _, v := range keys {
		for _, v2 := range Total {
			for _, v3 := range v2.All {
				if strings.Contains(v3.Text, v) {
					keysCnt[v] += 1
					total += 1
				}
			}
		}
	}

	fileName := path.Join("static", "file", "statistics.txt")
	f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	sum := 100
	if total == 0 {
		total = 1
	}
	for k, v := range keys {
		cnt := keysCnt[v]
		value := cnt * 100 / total
		if k == len(keys)-1 {
			value = sum
		}
		sum -= value
		f.WriteString(fmt.Sprintf("%s=%d\n", v, value))
	}

	return fileName
}
