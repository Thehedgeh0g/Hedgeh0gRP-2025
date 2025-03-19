package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	amqpadapter "valuator/pkg/Infrastructure/amqp"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"valuator/pkg/Infrastructure/repository"
	"valuator/pkg/app/service"
)

var redisClient *redis.Client
var amqpConn *amqp.Connection
var amqpChannel *amqp.Channel

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
	var err error
	amqpConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer func(amqpConn *amqp.Connection) {
		err2 := amqpConn.Close()
		if err2 != nil {
			panic(err2)
		}
	}(amqpConn)
}

func main() {
	var err error

	amqpChannel, err = amqpConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a amqpChannel:", err)
	}
	defer func(channel *amqp.Channel) {
		err2 := channel.Close()
		if err2 != nil {
			panic(err2)
		}
	}(amqpChannel)

	value := os.Getenv("PORT")
	if value == "" {
		value = "8082"
	}
	fmt.Println("Listening on port " + value)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/api/text", textHandler)
	err = http.ListenAndServe(":"+value, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port " + value)
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		amqpDispatcher := amqpadapter.NewAMQPDispatcher(amqpChannel, "text")
		textRepository := repository.NewTextRepository(redisClient)
		textService := service.NewTextService(textRepository, amqpDispatcher)
		// Получаем текст из формы
		data := r.FormValue("text")

		hash, err := textService.EvaluateText(data)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			log.Println(err.Error())
			return
		}

		// Отправляем на страницу summary с результатами
		http.Redirect(w, r, fmt.Sprintf("/summary?id=%s", hash), http.StatusPermanentRedirect)
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
	hash := r.URL.Query().Get("id")
	textRepository := repository.NewTextRepository(redisClient)
	text, err := textRepository.FindByHash(hash)
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

func aboutHandler(w http.ResponseWriter, _ *http.Request) {
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

func textHandler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Query().Get("id")
	textRepository := repository.NewTextRepository(redisClient)
	text, err := textRepository.FindByHash(hash)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}

	_, err = w.Write([]byte(text.GetText()))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		log.Println(err.Error())
		return
	}
}
