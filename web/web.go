package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/kevin-cantwell/kvn/docgifs"
	jsn "github.com/timehop/goth/json"
)

func main() {
	go docgifs.PeriodicallyRefresh()

	rand.Seed(time.Now().Unix())

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/image", GifHandler)
	r.HandleFunc("/slimemold", SlimeMoldHandler)
	r.HandleFunc("/slimemold/{asset}", SlimeMoldAssetHandler)
	r.HandleFunc("/docgif", DocGifHandler)
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

func SlimeMoldHandler(response http.ResponseWriter, request *http.Request) {
	t, err := template.ParseFiles("slimemold/slime_mold.html")
	if err != nil {
		writeError(response, err, "An unknown error occured")
		return
	}
	t.Execute(response, nil)
}

func SlimeMoldAssetHandler(response http.ResponseWriter, request *http.Request) {
	asset := mux.Vars(request)["asset"]
	http.ServeFile(response, request, filepath.Join("slimemold", asset))
}

func DocGifHandler(response http.ResponseWriter, request *http.Request) {
	page := docgifs.CurrentPage()
	fmt.Println("docgif:", `"doc gif `+page.SearchText+`"`, page.GiphyURL)

	t, err := template.ParseFiles("docgifs/docgifs.html")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(response, page)
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
