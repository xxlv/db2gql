package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

//go:embed index.html
var htmlContent string

var schemaCache map[string]map[string][]Column = make(map[string]map[string][]Column)

var typeMap = map[string]string{
	"varchar":    "String",
	"text":       "String",
	"char":       "String",
	"int":        "Int",
	"tinyint":    "Int",
	"smallint":   "Int",
	"mediumint":  "Int",
	"bigint":     "Int",
	"decimal":    "Float",
	"float":      "Float",
	"double":     "Float",
	"boolean":    "Boolean",
	"bool":       "Boolean",
	"tinyint(1)": "Boolean",
	"date":       "DateTime",
	"time":       "DateTime",
	"datetime":   "DateTime",
	"timestamp":  "DateTime",
	"json":       "JSON",
	"blob":       "String",
	"binary":     "String",
	"varbinary":  "String",
}

type Descriptor struct {
	Name       string
	Type       string
	Default    string
	IsList     bool
	IsRequired bool
	Comment    string
}

func getTableSchema(db *sql.DB, dbname string, tableName string) ([]Column, error) {
	if cached, ok := schemaCache[dbname]; ok {
		if v, ok := cached[tableName]; ok {
			return v, nil
		}
	} else {
		schemaCache[dbname] = map[string][]Column{}
	}

	query := `SELECT
    COLUMN_NAME AS 'Field',
    COLUMN_TYPE AS 'Type',
    IS_NULLABLE AS 'Null',
    COLUMN_KEY AS 'Key',
    COLUMN_DEFAULT AS 'Default',
    EXTRA AS 'Extra',
    COLUMN_COMMENT AS 'Comment'
	FROM
		INFORMATION_SCHEMA.COLUMNS
	WHERE
    TABLE_SCHEMA = '$DB'
    AND TABLE_NAME = '$TABLE_NAME';`

	query = strings.Replace(query, "$DB", dbname, 1)
	query = strings.Replace(query, "$TABLE_NAME", tableName, 1)
	rows, err := db.Query(query)
	log.Default().Println("prepare load schema from table ", tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var column Column
		if err := rows.Scan(&column.Name, &column.Type, &column.Null, &column.Key, &column.Default, &column.Extra, &column.Comment); err != nil {
			return nil, err
		}
		column.Table = tableName
		columns = append(columns, column)
	}

	// save cache
	colschema := schemaCache[dbname]
	colschema[tableName] = columns
	schemaCache[dbname] = colschema

	return columns, nil
}

type Column struct {
	Table   string
	Name    string
	Type    string
	Null    string
	Key     string
	Default sql.NullString
	Extra   string
	Comment string
}

func generateGraphQLSchemaWithCommentsWithFields(dbname string, crossTablefields []string) string {
	var allCols []Column
	var types map[string]any = make(map[string]any)
	for _, corrdinate := range crossTablefields {
		tableField := strings.Split(corrdinate, ".")
		tableName := tableField[0]
		cols, _ := getTableSchema(db, dbname, tableName)
		types[tableName] = nil
		for _, col := range cols {
			if col.Name == tableField[1] {
				rawcomment := col.Comment
				if rawcomment == "" {
					rawcomment = genComment(col)
				}
				col.Comment = rawcomment
				allCols = append(allCols, col)
			}
		}
	}
	schemaGenerator := &SchemaGenerator{
		Name:       asTypeNameFromKeys(types),
		RawColumns: allCols,
	}
	return schemaGenerator.Gen()
}

func genComment(col Column) string {
	return fmt.Sprintf("Field: %s", col.Name)
}

func mapMySQLTypeToGraphQL(mysqlType string, dbnull string) string {
	nullable := dbnull == "YES"
	graphqlType := "String"
	for prefix, gqlType := range typeMap {
		if strings.HasPrefix(mysqlType, prefix) {
			graphqlType = gqlType
			break
		}
	}
	if !nullable {
		graphqlType += "!"
	}

	return graphqlType
}

var defaultDb = "aitask"
var db *sql.DB
var dbUser = flag.String("dbuser", "root", "Database username")
var dbPassword = flag.String("dbpassword", "123456", "Database password")
var dbName = flag.String("dbname", defaultDb, "Database name")
var dbHost = flag.String("dbhost", "127.0.0.1", "Database host")
var dbPort = flag.String("dbport", "3306", "Database port")

var serverPort = flag.String("port", "8080", "Server port")
var DbName string // for global use ...

func main() {
	flag.Parse()
	DbName = *dbName

	if *dbName == "" || len(*dbName) <= 0 {
		println(DbName)
		flag.PrintDefaults()
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *dbUser, *dbPassword, *dbHost, *dbPort, *dbName)
	log.Println(dsn)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/", renderHTML)
	http.HandleFunc("/tables", handleTables)
	http.HandleFunc("/fields", handleFields)
	http.HandleFunc("/generateSchema", handleGenerateSchema)
	http.HandleFunc("/previewSchema", previewCurrentSchema)
	fmt.Println("Server is running on http://localhost:" + *serverPort)
	log.Fatal(http.ListenAndServe(":"+*serverPort, nil))
}
