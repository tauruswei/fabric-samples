package main

import (
	"fmt"
	"strconv"
)

/**
 * @Author: fengxiaoxiao /13156050650@163.com
 * @Desc:
 * @Version: 1.0.0
 * @Date: 2021/12/3 5:24 下午
 */

func main() {
	fmt.Println(strconv.Atoi(""))
	u, _ := strconv.ParseUint("42", 10, 64)
	fmt.Println(u)

	var bytes []byte
	if string(bytes) == "" {
		fmt.Println("success")
	}

}
