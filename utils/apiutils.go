package utils

import (
	"fmt"
	"io"
)

func ApplyApiObject(url string, obj interface{}) {
	resp, err := PostWithJson(url, obj)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(string(data))
}
