package main

import (
    "fmt"
	"io"
	"net/http"
	"os"
    "path/filepath"
	"strconv"
)

func serveSaveStudent(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) 
	if err != nil {
		http.Redirect(w, r, "/?error=Max file size 10 MB", http.StatusSeeOther)
		return
	}

	idStr := r.FormValue("studentID")
	name := r.FormValue("studentName")
	cgpaStr := r.FormValue("studentCGPA")
	interest := r.FormValue("careerInterest")

	if idStr == "" || name == "" || cgpaStr == "" || interest == "" {
		http.Redirect(w, r, "/?error=All input fields are required", http.StatusSeeOther)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Redirect(w, r, "/?error=Student ID must be a number", http.StatusSeeOther)
		return
	}

	cgpa, err := strconv.ParseFloat(cgpaStr, 64)
	if err != nil {
		http.Redirect(w, r, "/?error=CGPA must be a number", http.StatusSeeOther)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to process file upload", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", id, ext)
	filePath := filepath.Join("uploaded_images", filename)

	newFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Failed to copy file data", http.StatusInternalServerError)
		return
	}

	imageURL := "/" + filePath

	save_student_chan := make(chan bool) 

	go func() {
		if _, exists := students[id]; !exists {
			student := Student{
				ID:            	id,
				Name:          	name,
				CGPA:          	cgpa,
				CareerInterest: interest,
				ImageURL:      	imageURL,
			}
	
			students[id] = student

			save_student_chan <- !exists
		}

		close(save_student_chan)
	}()

	done := <- save_student_chan

	if done {
		http.Redirect(w, r, "/?success="+fmt.Sprintf("Student with ID %d added successfully", id), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/?error="+fmt.Sprintf("Student with ID %d already exists", id), http.StatusSeeOther)
	}
}