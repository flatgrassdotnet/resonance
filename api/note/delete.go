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
)

func Delete(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		common.WriteError(w, r, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		common.WriteError(w, r, "invalid id", http.StatusBadRequest)
		return
	}

	// token
	token, err := base64.StdEncoding.DecodeString(r.Header.Get("Authorization"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to decode token: %s", err), http.StatusInternalServerError)
		return
	}

	// steamid
	steamid, err := db.SteamIDFromToken(r.Context(), token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.WriteError(w, r, "invalid token", http.StatusUnauthorized)
			return
		}

		common.WriteError(w, r, fmt.Sprintf("failed to get steamid: %s", err), http.StatusInternalServerError)
		return
	}

	// owner check
	owner, err := db.GetNoteOwner(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.WriteError(w, r, "note doesn't exist", http.StatusBadRequest)
			return
		}

		common.WriteError(w, r, fmt.Sprintf("failed to get note owner: %s", err), http.StatusInternalServerError)
		return
	}

	if steamid != owner {
		common.WriteError(w, r, "not note owner", http.StatusBadRequest)
		return
	}

	err = db.DeleteNote(r.Context(), id)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to delete note: %s", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
