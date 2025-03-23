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

package login

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"resonance/common"

	"github.com/yohcop/openid-go"
)

func Login(w http.ResponseWriter, r *http.Request) {
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

	//url, err := openid.RedirectURL("https://steamcommunity.com/openid", fmt.Sprintf("https://%s/login/callback?ticket=%x", r.Host, ticket), fmt.Sprintf("https://%s", r.Host))
	url, err := openid.BuildRedirectURL("https://steamcommunity.com/openid/login", "http://specs.openid.net/auth/2.0/identifier_select", "http://specs.openid.net/auth/2.0/identifier_select", fmt.Sprintf("https://%s/login/callback?ticket=%x", r.Host, ticket), fmt.Sprintf("https://%s", r.Host))
	if err != nil {
		common.WriteError(w, r, fmt.Sprintf("failed to create login redirect: %s", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}
