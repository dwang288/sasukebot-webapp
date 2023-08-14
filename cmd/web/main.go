package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/dwang288/snippetbox/internal/models"

	"github.com/go-playground/form/v4"
	// Adding this so 1) go mod tidy doesn't remove the package and
	// 2) the init() function of the package runs and registers itself with the
	// database/sql package. This is standard procedure with most the go sql drivers.
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// inject our model (db) into our application struct
	snippets *models.SnippetModel
	// add a template cache for parsed templates so we don't have to keep reparsing
	templateCache map[string]*template.Template
	// add formDecoder for automatically pulling out post body data
	formDecoder *form.Decoder
}

func main() {

	// Define a command line flag with the name addr and a default value
	addr := flag.String("addr", ":4000", "HTTP network address")
	// Flag for the mySQL DSN string
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	// Parse value stored in flag and assign to addr. Without parsing, addr will always
	// be set to the default value. Will panic if errors occur during parsing
	flag.Parse()

	// New INFO level logger with output destination, message string prefix, and flags
	// to indicate the additional information to include (local date and time) joined with |
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// New ERROR level logger with all the INFO level logger information + filename/line number
	// of error logged and logged to stderr
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open a db based on the passed in dsn string, defer pool close
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	// Initialize new application struct with dependencies
	// Inject initialized DB, template cache, and form decoder
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// Specify and initialize a http.Server so we can use our custom errorLog.
	// Otherwise we could just use the http.ListenAndServe shortcut function.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	// Use the ListenAndServe() function on our custom http.Server
	// to start a new web server.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB wraps sql.Open() and return a sql.DB connection pool
func openDB(dsn string) (*sql.DB, error) {
	// sql.Open initializes the pool and estabilshes it for future use, but does
	// not actually open any connections. DB connections are opened lazily when needed.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Opens an actual connection to the DB to check that we can connect
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
