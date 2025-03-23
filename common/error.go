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

package common

import (
	"net/http"
	"resonance/db"
)

func WriteError(w http.ResponseWriter, r *http.Request, error string, code int) error {
	http.Error(w, error, code)

	var e db.Error
	e.Endpoint = r.URL.Path
	e.Error = error

	err := db.InsertError(r.Context(), e)
	if err != nil {
		return err
	}

	return nil
}
