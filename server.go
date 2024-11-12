package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
)

func renderHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, htmlContent)
}

func handleFields(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	columns, err := getTableSchema(db, DbName, tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(columns)
}

func handleTables(w http.ResponseWriter, r *http.Request) {
	query := `SHOW TABLES`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tables = append(tables, table)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tables)

}

func handleGenerateSchema(w http.ResponseWriter, r *http.Request) {
	var selectedFields map[string]map[string][]string
	if err := json.NewDecoder(r.Body).Decode(&selectedFields); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	var flatSelectedFields = []string{}
	all := selectedFields["selectedFields"]
	for _, fieldOfTable := range all {
		flatSelectedFields = append(flatSelectedFields, fieldOfTable...)
	}
	schema := generateGraphQLSchemaWithCommentsWithFields(DbName, flatSelectedFields)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(schema))
}
