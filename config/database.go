package config

import (
  "log"
  r "gopkg.in/gorethink/gorethink.v4"
)

func GetSession() *r.Session {
  session, err := r.Connect(r.ConnectOpts{
    Address: "localhost:28015",
    Database: "test",
	  Username: "admin",
	  Password: "asdqwe123",
  })

  if err != nil {
    log.Fatalln(err.Error())
  }

  return session
}