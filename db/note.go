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

package db

import (
	"time"
)

type Note struct {
	Map      string     `json:"map"`
	ID       int        `json:"id"`
	SteamID  string     `json:"steamid,omitempty"`
	Comment  string     `json:"comment"`
	Position [3]float64 `json:"position"`
	Created  time.Time  `json:"created"`
}

func InsertNote(note Note) (int, error) {
	pos, err := toVectorJSON(note.Position)
	if err != nil {
		return 0, err
	}

	r, err := conn.Exec("INSERT INTO notes (steamid, map, position, comment) VALUES (?, ?, VEC_FromText(?), ?)", note.SteamID, note.Map, pos, note.Comment)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func DeleteNote(id int) error {
	_, err := conn.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func LatestNoteTime(steamid string) (time.Time, error) {
	var created time.Time
	err := conn.QueryRow("SELECT COALESCE(MAX(created), FROM_UNIXTIME(946702800)) FROM notes WHERE steamid = ?", steamid).Scan(&created)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return created, nil
}

func GetNotes(filter string, value string) ([]Note, error) {
	var args []any
	query := "SELECT id, steamid, map, VEC_ToText(position), comment, created FROM notes"

	switch filter {
	case "steamid":
		query += " WHERE steamid = ?"
		args = append(args, value)
	case "map":
		query += " WHERE map = ?"
		args = append(args, value)
	}

	r, err := conn.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var notes []Note

	for r.Next() {
		var n Note
		var pos string
		err = r.Scan(&n.ID, &n.SteamID, &n.Map, &pos, &n.Comment, &n.Created)
		if err != nil {
			return nil, err
		}

		n.Position, err = fromVectorJSON(pos)
		if err != nil {
			return nil, err
		}

		notes = append(notes, n)
	}

	return notes, nil
}

func GetNoteOwner(id int) (string, error) {
	var steamid string
	err := conn.QueryRow("SELECT steamid FROM notes WHERE id = ?").Scan(&steamid)
	if err != nil {
		return "", err
	}

	return steamid, nil
}

func HasNoteWithinDistance(mapname string, position [3]float64, distance int) (bool, error) {
	pos, err := toVectorJSON(position)
	if err != nil {
		return false, err
	}

	var count int
	err = conn.QueryRow("SELECT COUNT(*) FROM notes WHERE map = ? AND VEC_DISTANCE_EUCLIDEAN(VEC_FromText(?), position) < ?", mapname, pos, distance).Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}
