// package main

// import (
// 	"net/http"

// 	"github.com/labstack/echo/v4"
// )

// func main() {
// 	e := echo.New()

// 	// route sederhana
// 	e.GET("/", func(c echo.Context) error {
// 		return c.String(http.StatusOK, "Hello from Echo + Go + Docker!")
// 	})

// 	// listen di port 8080
// 	e.Logger.Fatal(e.Start(":8080"))
// }

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "http://192.168.18.72:8080/app/api/v1/login-provider"
	body := []byte(`{"email":"surya.apandi28@gmail.com","password":"supersecret123"}`)

	client := &http.Client{}
	success := 0
	overLimit := 0

	for i := 1; i <= 1500; i++ {
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(i, "Error ->", err)
			continue
		}

		if resp.StatusCode == 429 || resp.StatusCode == 500 {
			overLimit++
			fmt.Printf("%d OVER LIMIT -> %d\n", i, resp.StatusCode)

			// reset hitungan ke awal
			fmt.Println("⚠️ Terkena limit, reset counter ke awal...")
			success = 0
			overLimit = 0
			i = 0 // supaya loop mulai lagi dari awal (karena akan di-increment jadi 1)
		} else {
			success++
			fmt.Printf("%d Success -> %d\n", i, resp.StatusCode)
		}

		// baca body (optional)
		_, _ = io.ReadAll(resp.Body)
		resp.Body.Close()

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("==== Test selesai ====")
	fmt.Println("Success:", success)
	fmt.Println("Over limit:", overLimit)
	fmt.Println("Total request yang dilakukan:", success+overLimit)
}
