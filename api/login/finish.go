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

package login

import (
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"resonance/common"
	"resonance/db"
)

func Finish(w http.ResponseWriter, r *http.Request) {
	ticket, err := hex.DecodeString(r.URL.Query().Get("ticket"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to decode ticket: %s", err), http.StatusBadRequest)
		return
	}
	if len(ticket) != 16 {
		common.WriteError(w, r, "invalid ticket", http.StatusBadRequest)
		return
	}

	token, err := db.TokenFromTicket(ticket)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			common.WriteError(w, r, "login not complete", http.StatusUnauthorized)
			return
		}

		common.WriteError(w, r, fmt.Sprintf("failed to get token: %s", err), http.StatusInternalServerError)
		return
	}

	be := base64.NewEncoder(base64.StdEncoding, w)
	defer be.Close()

	be.Write(token)
}
