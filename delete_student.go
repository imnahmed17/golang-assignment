package main

import (
    "fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func serveDeleteStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deleteIDStr := r.FormValue("deleteStudentID")
	if deleteIDStr == "" {
		http.Redirect(w, r, "/?error=Student ID is required", http.StatusSeeOther)
		return
	}

	deleteID, err := strconv.Atoi(deleteIDStr)
	if err != nil {
		http.Redirect(w, r, "/?error=Student ID must be a number", http.StatusSeeOther)
		return
	}

	delete_student_chan := make(chan bool)

	go func () {
		student, exists := students[deleteID] 
		if exists {
			err = os.Remove(strings.TrimLeft(student.ImageURL, "/"))
			if err != nil {
				http.Error(w, "Failed to remove file", http.StatusInternalServerError)
				return
			}

			delete(students, deleteID)
		}

		delete_student_chan <- exists
	}()

	done := <- delete_student_chan

	if done {
		http.Redirect(w, r, "/?success="+fmt.Sprintf("Student with ID %d deleted successfully", deleteID), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/?notfound="+fmt.Sprintf("Student with ID %d not found", deleteID), http.StatusSeeOther)
	}
}