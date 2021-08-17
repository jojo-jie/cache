package main

import (
	"fmt"
	"sync"
)

func main() {
	//runtime.GOMAXPROCS(8)
	/*trace.Start(os.Stderr)
	defer trace.Stop()*/
	urls := []string{"0.0.0.0:1121", "0.0.0.0:1122", "0.0.0.0:1123"}
	wg := sync.WaitGroup{}
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Println(url)
		}(url)
	}
	wg.Wait()
}

//链表前驱后继节点
//空link 节点插入删除顺序--> 引入哨兵节点
//检查链表代码是否正确的边界条件
//如果链表==nil
//如果链表只包含一个空节点
//如果只包含2哥空节点
//处理头尾节点是否可以正常工作 如何能够打破次元壁垒完成
//缓存空间换时间
