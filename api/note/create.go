/*
	resonance - echoes across all your favorite maps
	Copyright (C) 2025  patapancakes <patapancakes@pagefault.games>

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

package note

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"resonance/common"
	"resonance/db"
	"strconv"
	"strings"
	"time"
)

func Create(w http.ResponseWriter, r *http.Request) {
	// token
	token, err := base64.StdEncoding.DecodeString(r.Header.Get("Authorization"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to decode token: %s", err), http.StatusInternalServerError)
		return
	}

	var note db.Note

	// steamid
	note.Author, err = db.SteamIDFromToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.WriteError(w, r, "invalid token", http.StatusUnauthorized)
			return
		}

		common.WriteError(w, r, fmt.Sprintf("failed to get steamid: %s", err), http.StatusInternalServerError)
		return
	}

	u, err := db.GetUser(note.Author)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get user: %s", err), http.StatusInternalServerError)
		return
	}
	if u.Muted {
		common.WriteError(w, r, "muted", http.StatusForbidden)
		return
	}

	// map
	note.Map = strings.TrimSpace(r.PostFormValue("map"))
	if !common.IsValidMap(note.Map) {
		common.WriteError(w, r, "invalid map", http.StatusBadRequest)
		return
	}

	// cooldown
	latest, err := db.LatestNoteTime(note.Author, note.Map)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get latest post time: %s", err), http.StatusInternalServerError)
		return
	}

	notes, err := db.GetNoteCountByUserMap(note.Author, note.Map)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get note count: %s", err), http.StatusInternalServerError)
		return
	}

	if latest.Add(time.Minute * time.Duration(notes)).After(time.Now().UTC()) {
		common.WriteError(w, r, "rate limited", http.StatusTooManyRequests)
		return
	}

	// comment
	note.Comment = strings.TrimSpace(r.PostFormValue("comment"))
	if note.Comment == "" || len(note.Comment) > 250 {
		common.WriteError(w, r, "invalid comment", http.StatusBadRequest)
		return
	}

	// position
	posStr := strings.Split(r.PostFormValue("pos"), ",")
	if len(posStr) != len(note.Position) {
		common.WriteError(w, r, "invalid pos", http.StatusBadRequest)
		return
	}

	for i, p := range posStr {
		note.Position[i], err = strconv.ParseFloat(p, 64)
		if err != nil {
			common.WriteError(w, r, "invalid pos", http.StatusBadRequest)
			return
		}
	}

	// proximity check
	hasNote, err := db.HasNoteWithinDistance(note.Map, note.Position, 70)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to determine note proximity: %s", err), http.StatusInternalServerError)
		return
	}
	if hasNote {
		common.WriteError(w, r, "distance requirement not met", http.StatusBadRequest)
		return
	}

	id, err := db.InsertNote(note)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to insert note: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, id)
}
