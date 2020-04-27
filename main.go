package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gen2brain/beeep"
)

var timeNotFound = `Доступное время отсутствует. `

func main() {
	point := flag.Int("point", 0, "point number 0-based")
	interval := flag.Int("interval", 60, "interval in seconds")
	isNotify := flag.Bool("notify", true, "notify")
	flag.Parse()

	if *interval < 5 {
		fmt.Println("interval must be greater than 4")
		os.Exit(0)
	}

	for {
		loop(*point, *isNotify)
		time.Sleep(time.Second * time.Duration(*interval))
	}
}

func loop(point int, isNotify bool) {
	resp, err := http.Get("https://e-dostavka.by/")
	if err != nil {
		fmt.Println("failed to http get", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("failed to get new document", err)
		return
	}

	doc.Find(".item").Each(func(i int, item *goquery.Selection) {
		if i != point {
			return
		}

		fmt.Println(time.Now().Format(time.RFC3339))

		pointName := item.Find("a").First().Text()
		fmt.Println(pointName)

		deliveryTime := item.Contents().Not("a").Text()
		fmt.Println(deliveryTime)

		if deliveryTime != timeNotFound  && isNotify{
			err = beeep.Notify(pointName, deliveryTime, "")
			if err != nil {
				fmt.Println("failed to notify", err)
			}
		}

		fmt.Println()
	})
}