package main

import (
    "fmt"
	"html/template"
	"net/http"
	"strconv"
)

func serveSearchStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	searchIDStr := r.URL.Query().Get("searchStudentID")
	if searchIDStr == "" {
		http.Redirect(w, r, "/?error=Student ID is required", http.StatusSeeOther)
		return
	}

	searchID, err := strconv.Atoi(searchIDStr)
	if err != nil {
		http.Redirect(w, r, "/?error=Student ID must be a number", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/student_details.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	search_student_chan := make(chan Student)

	go func () {
		student, exists := students[searchID]
		if exists {
			search_student_chan <- student
		}

		close(search_student_chan)
	}()

	student, exists := <-search_student_chan
	if !exists {
		http.Redirect(w, r, "/?notfound="+fmt.Sprintf("Student with ID %d not found", searchID), http.StatusSeeOther)
	} else {
		tmpl.Execute(w, student)
	}
}