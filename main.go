package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"github.com/gorilla/handlers"
	//"time"
)

func main() {
	// Set the router as the default one shipped with Gin

	router := gin.Default()
	hasLoaded := false

	// Serve frontend static login page file

	router.StaticFS("/hi/", http.Dir("./login/build"))

	router.GET("/dashReq", func(c *gin.Context) {

		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		if lastname == "hello" {

			router.StaticFS("/dash", http.Dir("./login/build1"))

		}

	})

	router.POST("/", func(c *gin.Context) {

		//fmt.Printf("BODY: %s", c.Request.Header)

		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			//handle read response error
		}

		//result := decrypt([]byte(string(body)))

		if string(body) == "17b84c8330f86af407ec45cd1ac3e9bc183d38d3c13d64ff06fbd699ccb3c69e" && hasLoaded == false {

			router.StaticFS("/dash", http.Dir("./login/build1"))
			hasLoaded = true

		}

		fmt.Printf("%s\n", string(body))

		//fmt.Printf("%s\n", string(result))

	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "HEAD", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://localhost"
		},
		MaxAge: 12 * time.Hour,
	}))

	//router.Use(cors.Default())

	router.Run(":5000")

	//log.Fatal(http.ListenAndServe(":5000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))

}

/* func decrypt(data []byte) []byte {

	key := []byte("000102030405060708090a0b0c0d0e0f")
	nonce1 := []byte("101112131415161718191a1b1c1d1e1d")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//nonce := make([]byte, gcm.NonceSize())

	//fmt.Printf("THIS IS NONCE SIZE: %d", len(nonce))

	plaintext, err := gcm.Open(nil, nonce1[:12], data, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
} */

func queryParams(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.URL.Query() {
		fmt.Printf("%s: %s\n", k, v)
	}
}
