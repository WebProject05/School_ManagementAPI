package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"strconv"
	"strings"
	"sync"
)

var teachers = make(map[int]models.Teacher)
var mutex = &sync.Mutex{}
var nextID = 1

// Initialize some dummy data
// This init() function does not have to be called
// The Go runtime automatically runs this init() as this is a reserved keyword
func init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Physics",
	}
	nextID++

	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "9B",
		Subject:   "Mathematics",
	}
	nextID++

	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Michael",
		LastName:  "Doe",
		Class:     "11C",
		Subject:   "Chemistry",
	}
	nextID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println("Requested ID:", idStr)

	var teacherList []models.Teacher

	if idStr == "" {
		// No ID provided: Search by query parameters
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}
	} else {
		// ID provided: Convert the string ID to an integer
		id, err := strconv.Atoi(idStr)
		if err != nil {
			// If they typed something like "/teachers/apple", Atoi will fail.
			http.Error(w, "Invalid teacher ID. Must be a number.", http.StatusBadRequest)
			return
		}

		// Now we can safely compare 'id' (int) with 'teacher.ID' (int)
		for _, teacher := range teachers {
			if teacher.ID == id {
				teacherList = append(teacherList, teacher)
				break
			}
		}
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func postTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	// This will take the raw data and convert it to go structs
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newTecher := range newTeachers {
		newTecher.ID = nextID
		teachers[nextID] = newTecher
		addedTeachers[i] = newTecher
		nextID++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)

}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// A function that handles the get method route
		getTeachersHandler(w, r)

	case http.MethodPost:
		postTeacherHandler(w, r)

	case http.MethodPut:
		w.Write([]byte("Update (PUT) teacher"))
		return

	case http.MethodPatch:
		w.Write([]byte("Partial Update (PATCH) teacher"))
		return

	case http.MethodDelete:
		w.Write([]byte("Delete (DELETE) teacher"))
		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.Write([]byte("Hello from the teachers route!"))
}
