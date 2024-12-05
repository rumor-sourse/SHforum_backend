package logic

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

/*
压力测试，增加10000个帖子
*/
func TestCreatePost(t *testing.T) {
	url := "http://127.0.0.1:8088/api/v1/post/create"
	method := "POST"

	for i := 0; i < 10000; i++ {
		payload := strings.NewReader(fmt.Sprintf(`{
		"title": "stress-test-%d",
		"content": "stress test content",
		"community_id": 4}`, i+11))

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1NDQyNjY0NzYwNDEwNTIxNiwidXNlcm5hbWUiOiJ3Y2s0IiwiaXNzIjoibXktcHJvamVjdCIsImV4cCI6MTczMjc3NzEyNH0.uyRgXK2YgFa-HoDWKa-ErleYJ4MLi1r2sHydQrJWtcE")
		req.Header.Add("Accept", "*/*")
		req.Header.Add("Host", "127.0.0.1:8088")
		req.Header.Add("Connection", "keep-alive")

		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		res.Body.Close()
	}

}
