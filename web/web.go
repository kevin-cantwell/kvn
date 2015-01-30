package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	jsn "github.com/timehop/goth/json"
)

func main() {
	rand.Seed(time.Now().Unix())

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/image", GifHandler)
	http.Handle("/", r)
	log.Println("http://localhost:" + os.Getenv("PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func writeError(response http.ResponseWriter, err error, msg string) {
	response.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		log.Println("ERROR:", err.Error(), msg)
	} else {
		log.Println("ERROR:", msg)
	}
	response.Write([]byte(msg))
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}
	t.Execute(response, nil)
}

func GifHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	q := request.Form["q"]
	if len(q) < 1 {
		writeError(response, errors.New("No query specified"), "No query specified")
		return
	}

	query := q[0]

	apiKey := os.Getenv("GIPHY_API_KEY")
	if apiKey == "" {
		writeError(response, errors.New("y u no set GIPHY_API_KEY???"), "An unknown error occured")
		return
	}

	fmt.Println(q)
	resp, err := http.Get("http://api.giphy.com/v1/gifs/search?q=" + url.QueryEscape(query) + "&api_key=" + apiKey)
	if err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}

	var giphyResp jsn.Data
	if err = json.Unmarshal(body, &giphyResp); err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}

	var data []interface{}
	if data, err = giphyResp.Array("data"); err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}
	if len(data) == 0 {
		writeError(response, errors.New("no can haz"), "No images could be found for your query :(")
		return
	}

	urls := make([]string, len(data))

	for i, img := range data {
		var image map[string]interface{}
		var ok bool
		if image, ok = img.(map[string]interface{}); !ok {
			writeError(response, err, "An unknown error occured")
			return
		}
		if urls[i], err = jsn.Data(image).String("images.original.url"); err != nil {
			writeError(response, err, "An unknown error occured")
			return
		}
	}

	t, err := template.ParseFiles("gif.html")
	if err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}

	n := rand.Intn(len(urls))
	p := struct{ Urls []string }{Urls: []string{urls[n]}}
	t.Execute(response, &p)
}
