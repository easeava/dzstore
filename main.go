package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"strconv"
	"sync"
)

type Store struct {
	title string
}

var spiderWait sync.WaitGroup

// Spider func
func Spider() {
	client := &http.Client{}
	url := "http://www.dianping.com/jinan/ch10"
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(reqest)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(fmt.Sprintf("%T", doc))

	// 分页总数
	pageNumStr := doc.Find(".PageLink").Last().Text()

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		log.Fatal(err)
	}

	//c := make(chan string)
	//fmt.Println(fmt.Sprintf("%T", stores), pageNum)
	for i := 1; i <= pageNum; i++ {
		spiderWait.Add(1)
		go GetPageDatas(i, &spiderWait)
	}

	spiderWait.Wait()
}

func GetPageDatas(num int, wait *sync.WaitGroup) {
	url := fmt.Sprintf("http://www.dianping.com/jinan/ch10/p%d", num)

	doc, err := GetResponse(url)
	if err != nil {
		log.Fatal(err)
		wait.Done()
	}

	doc.Find("#shop-all-list > ul > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".tit h4").Text()
		tag := s.Find(".tag-addr .tag").First().Text()
		point := s.Find(".tag-addr .tag").Last().Text()
		addr := s.Find(".addr").Text()
		fmt.Printf(fmt.Sprintf("第一%d页: title: %s 菜系: %s 地址: %s - %s\n", num, title, tag, point, addr))
		//c <- fmt.Sprintf("第一%d页: title: %s", num, title)
	})

	wait.Done()
}

func sum(num int, c chan int) {
	c <- num
}

// GetResponse func
func GetResponse(url string) (*goquery.Document, error) {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	if err != nil {
		log.Fatal(err)
	}
	res, err := client.Do(reqest)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	return goquery.NewDocumentFromReader(res.Body)
}

func main() {
	Spider()
}
