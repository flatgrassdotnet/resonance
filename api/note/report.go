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

func Report(w http.ResponseWriter, r *http.Request) {
	// token
	token, err := base64.StdEncoding.DecodeString(r.Header.Get("Authorization"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to decode token: %s", err), http.StatusInternalServerError)
		return
	}

	var report db.Report

	// steamid
	report.SteamID, err = db.SteamIDFromToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.WriteError(w, r, "invalid token", http.StatusUnauthorized)
			return
		}

		common.WriteError(w, r, fmt.Sprintf("failed to get steamid: %s", err), http.StatusInternalServerError)
		return
	}

	// muted check
	u, err := db.GetUser(report.SteamID)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get user: %s", err), http.StatusInternalServerError)
		return
	}
	if u.Muted {
		common.WriteError(w, r, "muted", http.StatusForbidden)
		return
	}

	// cooldown
	latest, err := db.LatestReportTime(report.SteamID)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to get latest post time: %s", err), http.StatusInternalServerError)
		return
	}

	if latest.Add(time.Minute).After(time.Now().UTC()) {
		common.WriteError(w, r, "rate limited", http.StatusTooManyRequests)
		return
	}

	// note id
	report.NoteID, err = strconv.Atoi(r.PostFormValue("noteid"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("invalid noteid: %s", err), http.StatusBadRequest)
		return
	}

	// reason
	report.Reason = strings.TrimSpace(r.PostFormValue("reason"))
	if report.Reason == "" || len(report.Reason) > 250 {
		common.WriteError(w, r, "invalid reason", http.StatusBadRequest)
		return
	}

	_, err = db.InsertReport(report)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to insert report: %s", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
