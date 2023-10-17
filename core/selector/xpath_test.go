package selector

import (
	"fmt"
	"github.com/antchfx/xpath"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"testing"
)

func getBody() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start=200filter=", nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

func TestXpath(t *testing.T) {
	selector, err := NewSelector(getBody())
	if err != nil {
		fmt.Println(err)
		return
	}

	allMovies := selector.Xpath(`//ol[@class="grid_view"]li`)
	if len(allMovies) == 0 {
		logx.Info("没有获取到任何子节点")
		return
	}

	for _, v := range allMovies {
		title := v.FirstNode(`.//span[@class="title"]`)
		quote := v.FirstNode(`.//p[@class="quote"]/span`)
		link := v.FirstNode(`.//div[@class="hd"]/as`).GetAttribute("href")
		fmt.Println(fmt.Sprintf("链接:%v,title:%s,引言:%s", link, title.Text(), quote.Text()))
	}

}

func TestCheckXpathIllegible(t *testing.T) {
	compile, err := xpath.Compile("aaaa")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(compile)
}
