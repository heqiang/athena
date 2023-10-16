package selector

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

const htmlSample = `<!DOCTYPE html><html lang="en-US">
<head>
<title>Hello,World!</title>
</head>
<body>
<div class="container">
<header>
	<!-- Logo -->
   <h1>City Gallery</h1>
</header>  
<nav>
  <ul>
    <li><a href="/London">London</a></li>
    <li><a href="/Paris">Paris</a></li>
    <li><a href="/Tokyo">Tokyo</a></li>
  </ul>
</nav>
<article>
  <h1>London</h1>
  <img src="pic_mountain.jpg" alt="Mountain View" style="width:304px;height:228px;">
  <p>London is the capital city of England. It is the most populous city in the  United Kingdom, with a metropolitan area of over 13 million inhabitants.</p>
  <p>Standing on the River Thames, London has been a major settlement for two millennia, its history going back to its founding by the Romans, who named it Londinium.</p>
</article>
<footer>Copyright &copy; W3Schools.com</footer>
</div>
</body>
</html>
`

func TestXpath(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start=200filter=", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	selector, err := NewSelector(string(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	all_movieds, err := selector.Xpath(`//ol[@class="grid_view"]/li`)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range all_movieds {
		//title, _ := v.FirstNode(`.//span[@class="title"]`)
		//quote, _ := v.FirstNode(`.//p[@class="quote"]/span`)
		link, _ := v.FirstNode(`.//div[@class="hd"]/a`)
		fmt.Println(fmt.Sprintf("链接:%v", link.GetAttribute("href")))
	}

}
