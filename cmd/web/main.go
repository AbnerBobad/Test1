package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"time"

	"github.com/AbnerBobad/final_project/internal/data"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
)

type application struct {
	addr *string

	products *data.ProductModel
	users    *data.UserModel

	logger        *slog.Logger
	templateCache map[string]*template.Template
	session       *sessions.Session
	tlsConfig     *tls.Config //new change
}

func main() {
	addr := flag.String("addr", "", "HTTP network address")
	dsn := flag.String("dsn", "", "PostgreSQL DSN")
	secret := flag.String("secret", "sjkda2+sd3ds+sdf3+asc3+sdf42+sld", "Secret key")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database connection pool established")
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()
	//session integration NEW
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	//tls configuration ECDHE || NEW CHANGES
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	app := &application{
		addr:     addr,
		products: &data.ProductModel{DB: db},
		users:    &data.UserModel{DB: db},
		// category:      &data.CategoryModel{DB: db},
		logger:        logger,
		templateCache: templateCache,
		session:       session,
		tlsConfig:     tlsConfig, //new changes
	}
	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
