package main

import (
    "database/sql"
    "flag"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "time"

    _ "github.com/mattn/go-sqlite3"
    "forum-mvp/internal/app"
)

func main() {
    addr := flag.String("addr", ":8080", "http listen address")
    dataDir := flag.String("data", "./data", "data directory for sqlite")
    tplDir := flag.String("templates", "./internal/web/templates", "templates dir")
    flag.Parse()

    if err := os.MkdirAll(*dataDir, 0755); err != nil { log.Fatal(err) }

    dbPath := filepath.Join(*dataDir, "forum.db")
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil { log.Fatal(err) }

    schema, err := os.ReadFile("internal/db/schema.sql")
    if err != nil { log.Fatal(err) }
    if _, err := db.Exec(string(schema)); err != nil { log.Fatal(err) }

    tpls := template.Must(template.ParseGlob(filepath.Join(*tplDir, "*.html")))

    a := &app.App{ DB: db, Templates: tpls, CookieName: "forum_session", SessionTTL: 7*24*time.Hour }

    // Routes
    http.HandleFunc("/", a.HandleIndex)
    http.HandleFunc("/register", a.HandleRegister)
    http.HandleFunc("/login", a.HandleLogin)
    http.HandleFunc("/logout", a.HandleLogout)
    http.HandleFunc("/post", a.HandleShowPost)
    http.HandleFunc("/post/new",  a.RequireAuth(a.HandleNewPost))
http.HandleFunc("/comment/new", a.RequireAuth(a.HandleNewComment))
http.HandleFunc("/like",        a.RequireAuth(a.HandleLike))

    log.Printf("listening on %s", *addr)
    log.Fatal(http.ListenAndServe(*addr, logRequest(http.DefaultServeMux)))
}

func logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        fmt.Printf("%s %s %s\n", r.Method, r.URL.Path, time.Since(start))
    })
}
