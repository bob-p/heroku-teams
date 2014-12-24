package main

import(
  "github.com/mholt/binding"
  "github.com/bgentry/heroku-go"
  "time"
)

// User model
type User struct {
  ID int
  Email string
  CreatedAt time.Time `db:"created_at"`
}

type Users []User

func (u *User) FieldMap() binding.FieldMap {
  return binding.FieldMap{
    &u.Email: "email",
  }
}

func(u User) AddToHeroku(client heroku.Client) {
  apps, _ := client.AppList(nil)
  for _, app := range apps {
    client.CollaboratorCreate(app.Id, u.Email, nil)
  }
}

func(u User) RemoveFromHeroku(client heroku.Client) {
  apps, _ := client.AppList(nil)
  for _, app := range apps {
    client.CollaboratorDelete(app.Id, u.Email)
  }
}
