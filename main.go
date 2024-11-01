package main

import (
   "database/sql"
   "fmt"
   "log"
   "net/http"

   // 使用标准库的 sqlite3 驱动
   _ "modernc.org/sqlite"
)

func main() {
   http.HandleFunc("/user", getUser)
   log.Fatal(http.ListenAndServe(":8080", nil))
}

func getUser(w http.ResponseWriter, r *http.Request) {
   username := r.URL.Query().Get("username")

   // 这里存在 SQL 注入漏洞
   query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)

   db, err := sql.Open("sqlite", "example.db")
   if err != nil {
       log.Fatal(err)
   }
   defer db.Close()

   rows, err := db.Query(query)
   if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
       return
   }
   defer rows.Close()

   for rows.Next() {
       var id int
       var name string
       err := rows.Scan(&id, &name)
       if err != nil {
           http.Error(w, err.Error(), http.StatusInternalServerError)
           return
       }
       fmt.Fprintf(w, "ID: %d, Username: %s\n", id, name)
   }
}
