package main

import (
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	"valuator/pkg/Infrastructure/repository"
	"valuator/pkg/app/service"
)

var redisClient *redis.Client

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Адрес Redis
	})
}

func main() {
	value := os.Getenv("PORT")
	if value == "" {
		value = "8082"
	}
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/about", aboutHandler)
	err := http.ListenAndServe(":"+value, nil)
	if err != nil {
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		textRepository := repository.NewTextRepository(redisClient)
		textService := service.NewTextService(textRepository)
		// Получаем текст из формы
		data := r.FormValue("text")

		textID, err := textService.ProcessText(data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		// Отправляем на страницу summary с результатами
		http.Redirect(w, r, fmt.Sprintf("/summary?id=%s", textID.String()), http.StatusPermanentRedirect)
		return
	}

	ts, err := template.ParseFiles("pages/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	textRepository := repository.NewTextRepository(redisClient)
	text, err := textRepository.FindByID(uuid.MustParse(id))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	tmpl, _ := template.ParseFiles("pages/summary.html")
	err = tmpl.Execute(w, map[string]interface{}{
		"Rank":       text.GetRank(),
		"Similarity": boolToInt(text.GetSimilarity()),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("pages/about.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
	err = tmpl.Execute(w, map[string]interface{}{
		"Name":  "Ilya Lezhnin",
		"Email": "Flipdoge87@gmail.com",
	})
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}
