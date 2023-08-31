package main

import (
    "fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
)

type Student struct {
	ID            	int     `json:"id"`
	Name          	string  `json:"name"`
	CGPA          	float64 `json:"cgpa"`
	CareerInterest	string  `json:"career_interest"`
	ImageURL      	string  `json:"image_url"`
}

var students = make(map[int]Student)

func mapToSlice(m map[int]Student) []Student {
    result := make([]Student, 0, len(m))
    for _, student := range m {
        result = append(result, student)
    }
    return result
}

func serveNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/error_page.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
        SuccessMessage string
        NotFoundMessage string
        ErrorMessage   string
    }{
        SuccessMessage: r.URL.Query().Get("success"),
        NotFoundMessage: r.URL.Query().Get("notfound"),
        ErrorMessage: r.URL.Query().Get("error"),
    }

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// fmt.Println("Stored Student Data:")
	// for id, student := range students {
	// 	fmt.Printf("ID: %d, Name: %s, CGPA: %.2f, Career Interest: %s, Image URL: %s\n", id, student.Name, student.CGPA, student.CareerInterest, student.ImageURL)
	// }
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", serveIndex)
	router.HandleFunc("/submit", serveSaveStudent)
	router.HandleFunc("/all-student", serveAllStudent)
	router.HandleFunc("/search-student", serveSearchStudent)
	router.HandleFunc("/delete-student", serveDeleteStudent)
	router.NotFoundHandler=http.HandlerFunc(serveNotFound)
	http.Handle("/", router)

	os.MkdirAll("uploaded_images", os.ModePerm)

	fs := http.FileServer(http.Dir("uploaded_images"))
    http.Handle("/uploaded_images/", http.StripPrefix("/uploaded_images/", fs))

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}