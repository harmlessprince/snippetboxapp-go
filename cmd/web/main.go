package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/harmlessprince/snippetboxapp/pkg/models"
	"github.com/harmlessprince/snippetboxapp/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type contextKey string

var contextKeyUser = contextKey("user")

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	snippetModel interface {
		Insert(string, string, string) (int, error)
		Get(int) (*models.Snippet, error)
		Latest() ([]*models.Snippet, error)
	}
	userModel interface {
		Insert(string, string, string) error
		Authenticate(string, string) (int, error)
		Get(int) (*models.User, error)
	}
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func main() {

	addr := flag.String("addr", ":4000", "Http Network Address/Port")
	dsn := flag.String("dsn", "root:password@/snippetbox?parseTime=true", "MYSQL Database Connection string")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret for session manager")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	infoLog.Println("Database connection successful")
	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Fatal(err)
	}

	app := &application{
		errorLog:      errLog,
		infoLog:       infoLog,
		snippetModel:  &mysql.SnippetModel{DB: db},
		userModel:     &mysql.UserModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}

	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: errLog,
		//Connection Timeout Settings
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on port %s", *addr)
	err = server.ListenAndServe()
	if err != nil {
		errLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
