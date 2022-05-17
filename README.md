# golimit
Golang 并发限制

### 使用 

`go get -u github.com/arsenalpoll/golimit`

```go
package main

import (
	"github.com/arsenalpoll/golimit"
	"log"
	"time"
)

func main() {

	start := time.Now()
	g := golimit.GoLimit(5)

	SaoMiaoPath := "D:\\ttt\\"
	files, _ := ioutil.ReadDir(SaoMiaoPath)

	for _, file := range files {
		value := file
		g.Add()
		go func(g *golimit.Limit) {
			defer g.Done()
			log.Println(value.Name())
			time.Sleep(5 * time.Second)
		}(g)
	}
	log.Println("the end")
	g.Wait()

	cost := time.Since(start)
	fmt.Printf("\ncost=[%s]\n", cost)

}

```


```go
package main

import (
	"github.com/arsenalpoll/golimit"
	"log"
	"time"
)

func main() {
	start := time.Now()
	g := golimit.GoLimit(5)

	for i := 0; i < 10; i++ {		
		g.Add()
		go func(g *golimit.Limit, i int) {
			defer g.Done() 
			time.Sleep(5 * time.Second)
			log.Println(i, "done")
		}(g, i)
	}
	log.Println("the end")
	g.Wait()

	cost := time.Since(start)
	fmt.Printf("\ncost=[%s]\n", cost)    
}

```