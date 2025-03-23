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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"resonance/common"
	"resonance/db"
)

type StatsResponse struct {
	UserCount int `json:"user_count"`
	NoteCount int `json:"note_count"`
	MapCount  int `json:"map_count"`

	Maps map[string]int `json:"maps"`
}

func Stats(w http.ResponseWriter, r *http.Request) {
	var sr StatsResponse
	var err error

	sr.UserCount, err = db.GetUserCount(r.Context())
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get user count: %s", err), http.StatusInternalServerError)
		return
	}

	sr.NoteCount, err = db.GetNoteCount(r.Context())
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get note count: %s", err), http.StatusInternalServerError)
		return
	}

	sr.MapCount, err = db.GetMapCount(r.Context())
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get map count: %s", err), http.StatusInternalServerError)
		return
	}

	sr.Maps, err = db.GetMaps(r.Context())
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get maps: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	err = json.NewEncoder(w).Encode(sr)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to write json: %s", err), http.StatusInternalServerError)
		return
	}
}
