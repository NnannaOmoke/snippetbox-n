package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox-n/internal/models"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger //do we have to use pointers?
	userModel      *models.UserModel
	snippetModel   *models.SnippetModel //has to be a pointer because it contains a db context, which we don't want copying
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func parseArgs() (string, string) {
	port := flag.String("port", ":4000", "Port binded to by HTTPS server")
	dsn := flag.String("dsn", "web:komboyagi.2006Y@/snippetbox?parseTime=true", "The means of connecting to your database")
	flag.Parse()
	return *port, *dsn
}

func main() {
	infoLog := log.New(os.Stdout, "INFO:/t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR:/t", log.Ldate|log.Ltime|log.Lshortfile)

	port, dsn := parseArgs()

	db, dErr := openDB(dsn)
	if dErr != nil {
		errLog.Fatal(dErr)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	application := &Application{
		errLog,
		infoLog,
		&models.UserModel{DB: db},
		&models.SnippetModel{DB: db},
		templateCache,
		formDecoder,
		sessionManager,
	}

	tlsConfig := tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:         port,
		ErrorLog:     errLog,
		Handler:      application.routes(),
		TLSConfig:    &tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	infoLog.Printf("Starting Server on port: %s", port)

	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errLog.Fatal(err)
}
