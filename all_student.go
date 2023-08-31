package main

import (
    "fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
)

func serveAllStudent(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/all_students.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	studentSlice := mapToSlice(students)

	sort.Slice(studentSlice, func(i, j int) bool {
		return studentSlice[i].ID < studentSlice[j].ID
	})

	perPage := 3 
    pageStr := r.URL.Query().Get("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    startIndex := (page - 1) * perPage
    endIndex := startIndex + perPage
    if endIndex > len(studentSlice) {
        endIndex = len(studentSlice)
    }

	data := studentSlice[startIndex:endIndex]

    prevPage := page - 1
    if prevPage < 1 {
        prevPage = 0
    }

    nextPage := page + 1
    if endIndex >= len(students) {
        nextPage = 0
    }

	templateData := struct {
        Students []Student
        PrevPage int
        NextPage int
    }{
        Students: data,
        PrevPage: prevPage,
        NextPage: nextPage,
    }

	if err := tmpl.Execute(w, templateData); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}