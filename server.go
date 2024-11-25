package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var sandBoxPort = "4000"
var gCurrSchema string = `""" no schema """`

func isPortAvailable(port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func HandleDb2GqlIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	portAvailable := isPortAvailable(sandBoxPort)
	tmpl, _ := template.New("html").Parse(htmlContent)
	data := struct {
		SandboxAvailable bool
		SandBoxPort      string
	}{
		SandboxAvailable: portAvailable,
		SandBoxPort:      sandBoxPort,
	}
	_ = tmpl.Execute(w, data)
}

func HandleFields(w http.ResponseWriter, r *http.Request) {
	tableName := r.URL.Query().Get("table")
	columns, err := getTableSchema(db, DbName, tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(columns)
}

func HandleTables(w http.ResponseWriter, r *http.Request) {
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

func HandleGenerateSchema(w http.ResponseWriter, r *http.Request) {
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
	gCurrSchema = schema // global
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(schema))
}

func HandlePreviewCurrentSchema(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte(gCurrSchema))
}
