package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dwang288/snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
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
	// inject our users model (db) into our application struct
	users *models.UserModel
	// add a template cache for parsed templates so we don't have to keep reparsing
	templateCache map[string]*template.Template
	// add formDecoder for automatically pulling out post body data
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Initialize a new sessionManager, set it to use our DB as the backing store
	// Set session TTL to 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Set Secure attribute on session cookies to indicate that this session cookie
	// should only be sent by a user's browser when a HTTPS connection is being used
	sessionManager.Cookie.Secure = true

	// Initialize new application struct with dependencies
	// Inject initialized snippets DB pool, initialized users DB pool,
	// template cache, and form decoder
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Set TLS settings to select the non default elliptic curve we want to use
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Specify and initialize a http.Server so we can use our custom errorLog.
	// Otherwise we could just use the http.ListenAndServe shortcut function.
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Add timeouts to connections
		// Protects agains slow client attacks, dropped connections clientside, etc
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// Use the ListenAndServeTLS() function on our custom http.Server
	// to start a new web server over HTTPS
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
