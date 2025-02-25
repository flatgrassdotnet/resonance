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
	"encoding/hex"
	"fmt"
	"net/http"
	"resonance/common"
	"resonance/db"
	"strings"

	"github.com/yohcop/openid-go"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	// ticket
	ticket, err := hex.DecodeString(r.URL.Query().Get("ticket"))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to decode ticket: %s", err), http.StatusBadRequest)
		return
	}
	if len(ticket) != 16 {
		common.WriteError(w, r, "invalid ticket", http.StatusBadRequest)
		return
	}

	id, err := openid.Verify(fmt.Sprintf("https://%s%s", r.Host, r.URL), discoveryCache, nonceStore)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("error while verifying login callback: %s", err), http.StatusInternalServerError)
		return
	}

	var u db.User

	u.SteamID = strings.TrimPrefix(id, "https://steamcommunity.com/openid/id/")

	err = db.InsertUser(u)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to insert user: %s", err), http.StatusInternalServerError)
		return
	}

	err = db.InsertTicket(u.SteamID, ticket)
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to insert link: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/authok.html", http.StatusSeeOther)
}
