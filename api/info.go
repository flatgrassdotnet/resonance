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

package api

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"resonance/common"
	"resonance/db"
	"strings"
	"time"
)

type InfoResponse struct {
	NoteCooldown   int `json:"note_cooldown"`
	ReportCooldown int `json:"report_cooldown"`
}

func Info(w http.ResponseWriter, r *http.Request) {
	// token
	token, err := base64.StdEncoding.DecodeString(r.Header.Get("Authorization"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to decode token: %s", err), http.StatusInternalServerError)
		return
	}

	// steamid
	steamid, err := db.SteamIDFromToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.WriteError(w, r, "invalid token", http.StatusUnauthorized)
			return
		}

		common.WriteError(w, r, fmt.Sprintf("failed to get steamid: %s", err), http.StatusInternalServerError)
		return
	}

	// note cooldown
	mapname := strings.TrimSpace(r.URL.Query().Get("map"))
	if !common.IsValidMap(mapname) {
		common.WriteError(w, r, "invalid map", http.StatusBadRequest)
		return
	}

	latest, err := db.LatestNoteTime(steamid, mapname)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get latest note time: %s", err), http.StatusInternalServerError)
		return
	}

	notes, err := db.GetNoteCountByUserMap(steamid, mapname)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get note count: %s", err), http.StatusInternalServerError)
		return
	}

	var ir InfoResponse

	ir.NoteCooldown = max(0, int(time.Until(latest.Add(time.Minute*time.Duration(notes))).Seconds()))

	// report cooldown
	latest, err = db.LatestReportTime(steamid)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get latest report time: %s", err), http.StatusInternalServerError)
		return
	}

	ir.ReportCooldown = max(0, int(time.Until(latest.Add(time.Minute*5)).Seconds()))

	w.Header().Set("Content-Type", "text/json")
	err = json.NewEncoder(w).Encode(ir)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to write json: %s", err), http.StatusInternalServerError)
		return
	}
}
