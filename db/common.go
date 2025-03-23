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

package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var conn *sql.DB

func Init(username string, password string, address string, name string) error {
	var err error
	conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, password, address, name))
	if err != nil {
		return err
	}

	return nil
}

func fromVectorJSON(vector string) ([3]float64, error) {
	var out [3]float64
	err := json.Unmarshal([]byte(vector), &out)
	if err != nil {
		return [3]float64{}, err
	}

	return out, nil
}

func toVectorJSON(vector [3]float64) (string, error) {
	out, err := json.Marshal(vector)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
