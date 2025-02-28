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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"resonance/common"
	"resonance/db"
)

type ViewResponse struct {
	Notes []db.Note `json:"notes"`
}

func View(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("map") {
		common.WriteError(w, r, "missing map", http.StatusBadRequest)
		return
	}

	var vr ViewResponse
	var err error

	vr.Notes, err = db.GetNotes("map", r.URL.Query().Get("map"))
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		common.WriteError(w, r, fmt.Sprintf("failed to read notes for map: %s", err), http.StatusInternalServerError)
		return
	}

	// strip steamid
	for i := range vr.Notes {
		vr.Notes[i].Author = ""
	}

	w.Header().Set("Content-Type", "text/json")
	err = json.NewEncoder(w).Encode(vr)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to write json: %s", err), http.StatusInternalServerError)
		return
	}
}
