package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
)

// NewApplicationDefault creates a new default application.
func NewApplicationSql(addr string) (a *ApplicationSql) {
	// default config
	defaultRouter := chi.NewRouter()
	defaultAddr := ":8080"
	if addr != "" {
		defaultAddr = addr
	}

	a = &ApplicationSql{
		rt:   defaultRouter,
		addr: defaultAddr,
		db:   connectDatabase(),
	}
	return
}

// ApplicationDefault is the default application.
type ApplicationSql struct {
	// rt is the router.
	rt *chi.Mux
	// addr is the address to listen.
	addr string

	db *sql.DB
}

// TearDown tears down the application.
func (a *ApplicationSql) TearDown() (err error) {
	return
}

// SetUp sets up the application.
func (a *ApplicationSql) SetUp() (err error) {
	// dependencies
	// - repository
	rp := repository.NewRepositoryProductMySql(a.db)
	// - handler
	hd := handler.NewHandlerProduct(rp)

	// router
	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)
	// - endpoints
	a.rt.Route("/products", func(r chi.Router) {
		// GET /products/{id}
		r.Get("/{id}", hd.GetById())
		// POST /products
		r.Post("/", hd.Create())
		// PUT /products/{id}
		r.Put("/{id}", hd.UpdateOrCreate())
		// PATCH /products/{id}
		r.Patch("/{id}", hd.Update())
		// DELETE /products/{id}
		r.Delete("/{id}", hd.Delete())
	})

	return
}

// Run runs the application.
func (a *ApplicationSql) Run() (err error) {
	defer a.db.Close()
	err = http.ListenAndServe(a.addr, a.rt)
	return
}

func connectDatabase() (db *sql.DB) {
	// config
	config := mysql.Config{
		User:      "user1",
		Passwd:    "secret_password",
		Addr:      "localhost:3306",
		Net:       "tcp",
		DBName:    "my_db",
		ParseTime: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		panic(err.Error())
	}

	return

}
