package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	jsn "github.com/timehop/goth/json"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/image", GifHandler)
	http.Handle("/", r)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		panic(err)
	}
}

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	t.Execute(response, nil)
}

func GifHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	q := request.Form["q"]
	if len(q) < 1 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("no query specified"))
		return
	}

	query := q[0]

	apiKey := os.Getenv("GIPHY_API_KEY")
	if apiKey == "" {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("y u no set GIPHY_API_KEY???"))
		return
	}

	fmt.Println(q)
	resp, err := http.Get("http://api.giphy.com/v1/gifs/search?q=" + url.QueryEscape(query) + "&api_key=" + apiKey)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	var giphyResp jsn.Data
	if err = json.Unmarshal(body, &giphyResp); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	var data []interface{}
	if data, err = giphyResp.Array("data"); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	if len(data) == 0 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("no images found"))
		return
	}

	var image map[string]interface{}
	var ok bool
	if image, ok = data[rand.Intn(len(data))].(map[string]interface{}); !ok {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("image item is not an object"))
		return
	}

	var url string
	if url, err = jsn.Data(image).String("images.original.url"); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}

	t, err := template.ParseFiles("gif.html")
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(err.Error()))
		return
	}
	p := struct{ GifUrl string }{GifUrl: url}
	t.Execute(response, &p)
}

// func ClockHandler(response http.ResponseWriter, request *http.Request) {
// 	vars := mux.Vars(request)
// 	tzCode := vars["tzCode"]
// 	t, err := template.ParseFiles("clock.html")
// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte(err.Error()))
// 		return
// 	}
// 	p := struct{ Timezone string }{Timezone: tzCode}
// 	t.Execute(response, &p)
// }
