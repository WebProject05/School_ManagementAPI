package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
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

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	return validFields[field]
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Connect to DB (Note: Pass a shared *sql.DB instance from main() in production)
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Parse ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println("Requested ID:", idStr)

	// Always initialize as an empty slice (so JSON encodes as [] instead of null if empty)
	teacherList := make([]models.Teacher, 0)

	if idStr != "" && idStr != "teachers" {
		// --- CASE 1: Fetch single teacher by ID ---
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid teacher ID. Must be a number.", http.StatusBadRequest)
			return
		}

		var teacher models.Teacher
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?"

		err = db.QueryRow(query, id).Scan(
			&teacher.ID,
			&teacher.FirstName,
			&teacher.LastName,
			&teacher.Email,
			&teacher.Class,
			&teacher.Subject,
		)

		if err == sql.ErrNoRows {
			// No teacher found with that ID
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}

		teacherList = append(teacherList, teacher)

	} else {
		// --- CASE 2: Search teachers by query params or fetch all ---
		// firstName := r.URL.Query().Get("first_name")
		// lastName := r.URL.Query().Get("last_name")

		// Dynamically build SQL query and arguments
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []interface{}

		query, args = addFilters(r, query, args)

		// Application of sorting to the get handler of teachers
		query = addSorting(r, query)

		// if firstName != "" {
		// 	query += " AND first_name = ?"
		// 	args = append(args, firstName)
		// }
		// if lastName != "" {
		// 	query += " AND last_name = ?"
		// 	args = append(args, lastName)
		// }

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(
				&teacher.ID,
				&teacher.FirstName,
				&teacher.LastName,
				&teacher.Email,
				&teacher.Class,
				&teacher.Subject,
			)
			if err != nil {
				http.Error(w, "Error scanning row data", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}

		// Check for errors during iteration
		if err = rows.Err(); err != nil {
			http.Error(w, "Error iterating table rows", http.StatusInternalServerError)
			return
		}
	}

	// Build & return JSON response
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

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]

	if len(sortParams) > 0 {
		var validSortClauses []string // Store valid valid "field ASC/DESC" strings here

		for _, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}

			field := parts[0]
			order := strings.ToUpper(parts[1]) // Normalize to uppercase (ASC/DESC)

			// 1. Validate the database field
			if !isValidSortField(field) {
				continue
			}

			// 2. Validate the sorting direction strictly
			if order != "ASC" && order != "DESC" {
				continue
			}

			// Add to our valid list
			validSortClauses = append(validSortClauses, field+" "+order)
		}

		// 3. Only append ORDER BY if we actually have valid sorting parameters
		if len(validSortClauses) > 0 {
			query += " ORDER BY " + strings.Join(validSortClauses, ", ")
		}
	}
	return query
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ? "
			args = append(args, value)
		}
	}
	return query, args
}

func postTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// Note: Replace this local ConnectDb call with a shared *sql.DB instance passed from main()
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher
	// Decode raw request JSON
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil || len(newTeachers) == 0 {
		http.Error(w, "Invalid or empty request body", http.StatusBadRequest)
		return
	}

	// Begin database transaction for batch insert
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Will safely rollback if not committed

	stmt, err := tx.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Error preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		// FIXED: Included missing newTeacher.Email argument to match 5 placeholders (?)
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error writing to the database: "+err.Error(), http.StatusInternalServerError)
			return
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
			return
		}

		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	// Commit transaction if all inserts succeed
	if err := tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
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

// PUT method
func updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)

	// If the teacher witht he given query is not found
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to retrive data", http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updatedTeacher.FirstName,
		updatedTeacher.LastName,
		updatedTeacher.Email,
		updatedTeacher.Class,
		updatedTeacher.Subject,
		id, // The last '?' placeholder is for the WHERE clause
	)

	if err != nil {
		http.Error(w, "Error updating Teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)

}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// A function that handles the get method route
		getTeachersHandler(w, r)

	case http.MethodPost:
		postTeacherHandler(w, r)

	case http.MethodPut:
		// A function that handles Put method
		updateTeacherHandler(w, r)

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
