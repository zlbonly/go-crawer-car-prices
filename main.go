package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go-crawer-car-prices/downloader"
	"go-crawer-car-prices/model"
	"go-crawer-car-prices/scheduler"
	"go-crawer-car-prices/spiders"
	"log"
	"time"
)

var (
	StartUrl     = "/2sc/%s/a0_0msdgscncgpi1ltocsp1exb4/"
	BaseUrl      = "https://car.autohome.com.cn"
	maxPage  int = 99
	cars     []spiders.QcCar
)

func Start(url string, ch chan []spiders.QcCar) {
	body := downloader.Get(BaseUrl + url)
	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		log.Printf("Downloader.Get err: %v", err)
	}

	currentPage := spiders.GetCurrentPage(doc)
	log.Print(currentPage)
	nextPageUrl, _ := spiders.GetNextPageUrl(doc)

	if currentPage > 0 && currentPage < maxPage {
		cars := spiders.GetCars(doc)
		log.Println(cars)
		ch <- cars
		if url := nextPageUrl; url != "" {
			scheduler.AppendUrl(url)
		}
		log.Println(url)

	} else {
		log.Printf("Max Page !!!!")
	}
}

func main() {
	citys := spiders.GetCitys()
	for _, v := range citys {
		scheduler.AppendUrl(fmt.Sprintf(StartUrl, v.Pinyin))
	}
	start := time.Now()
	delayTime := time.Second * 6
	ch := make(chan []spiders.QcCar)

L:
	for {
		if url := scheduler.PopUrl(); url != "" {
			go Start(url, ch)
		}

		select {
		case r := <-ch:
			cars = append(cars, r...)
			go Start(scheduler.PopUrl(), ch)
		case <-time.After(delayTime):
			log.Println("Timeout...")
			break L
		}
	}

	if len(cars) > 0 {
		model.AddCars(cars)
	}

	log.Printf("Time: %s", time.Since(start)-delayTime)

}
