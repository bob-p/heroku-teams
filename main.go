package main

import(
  "github.com/codegangsta/negroni"
  "github.com/gorilla/pat"
   _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
  "github.com/bgentry/heroku-go"
  "gopkg.in/unrolled/render.v1"
  "log"
  "net/http"
  "os"
)

var (
  port            = os.Getenv("PORT")
  db_opts         = os.Getenv("DATABASE_URL")
  heroku_username = os.Getenv("HEROKU_USERNAME")
  heroku_pw       = os.Getenv("HEROKU_PW")
)

type appContext struct {
  db *sqlx.DB
  render *render.Render
  herokuClient heroku.Client
}

func(appC appContext) GetUsers(users *[]User) error {
  return appC.db.Select(users, "SELECT * FROM users")
}

func(appC *appContext) homepageHandler(w http.ResponseWriter, req *http.Request) {
  users := []User{}
  appC.GetUsers(&users)

  appC.render.HTML(w, http.StatusOK, "index", users)
}

func(appC *appContext) usersHandler(w http.ResponseWriter, req *http.Request) {
  user := new(User)
  errs := binding.Bind(req, user)
  if errs.Handle(w) {
    return
  }
  createUser := "INSERT INTO users (email, created_at) VALUES ($1, $2)"
  appC.db.MustExec(createUser, user.Email, time.Now().Local())

  //Add to heroku apps
  user.AddToHeroku(appC.herokuClient)

  http.Redirect(w, req, "/", 301)
}

func(appC * appContext) deleteUsersHandler(w http.ResponseWriter, req *http.Request) {
  user := User{}
  err := appC.db.Get(&user, "SELECT * FROM users WHERE id = $1", req.URL.Query().Get(":id"))

  if (err != nil) {
    log.Print(err)
    appC.render.HTML(w, http.StatusNotFound, "404", nil)
  } else {
    //Remove from heroku apps
    user.RemoveFromHeroku(appC.herokuClient)

    removeUser := "DELETE FROM users where id = $1"
    appC.db.MustExec(removeUser, user.Id)

    http.Redirect(w, req, "/", 301)
  }
}

func main() {
  db, err := sqlx.Connect("postgres", db_opts)
  if err != nil {
    log.Fatalln(err)
  }

  render := render.New(render.Options{
    Layout: "layout",
    Extensions:    []string{".html"},
  })

  herokuClient := heroku.Client{Username: heroku_username, Password: heroku_pw}

  appC := &appContext{db: db, render: render, herokuClient: herokuClient}

  pat := pat.New()
  pat.Post("/users/{id}/delete", appC.deleteUsersHandler) //Delete user
  pat.Post("/users", appC.usersHandler) //Add new user
  pat.Get("/", appC.homepageHandler) //List all users

  n := negroni.Classic()
  n.UseHandler(pat)
  n.Run(":" + port)
}
