/*
	resonance - echoes across all your favorite maps
	Copyright (C) 2025  Pancakes <patapancakes@pagefault.games>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"resonance/api"
	"resonance/api/login"
	"resonance/api/note"
	"resonance/db"
)

func main() {
	// webserver related
	host := flag.String("host", "", "host to listen on")
	port := flag.Int("port", 80, "port to listen on")

	// database related
	dbuser := flag.String("dbuser", "resonance", "database server user name")
	dbpass := flag.String("dbpass", "", "database server user password")
	dbproto := flag.String("dbproto", "tcp", "database connection protocol")
	dbaddr := flag.String("dbaddr", "127.0.0.1", "database server address")
	dbname := flag.String("dbname", "resonance", "database name")

	flag.Parse()

	// set up database
	err := db.Init(*dbuser, *dbpass, *dbproto, *dbaddr, *dbname)
	if err != nil {
		log.Fatalf("failed to open database connection: %s", err)
	}

	// set up webserver

	// files
	http.Handle("GET /img/", http.FileServer(http.Dir("data")))

	// login
	http.HandleFunc("GET /login", login.Login)
	http.HandleFunc("GET /login/callback", login.Callback)
	http.HandleFunc("GET /login/finish", login.Finish)

	// misc
	http.HandleFunc("GET /stats", api.Stats)
	http.HandleFunc("GET /info", api.Info)
	http.HandleFunc("GET /nuke", api.Nuke)

	// note
	http.HandleFunc("GET /note/view", note.View)
	http.HandleFunc("GET /note/mine", note.Mine)
	http.HandleFunc("GET /note/delete", note.Delete)
	http.HandleFunc("POST /note/create", note.Create)
	http.HandleFunc("POST /note/report", note.Report)

	log.Printf("Starting web server on %s:%d", *host, *port)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		log.Fatalf("error in web server: %s", err)
	}
}
