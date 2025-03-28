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
	"context"
	"time"
)

type Note struct {
	Map      string     `json:"map,omitempty"`
	ID       int        `json:"id"`
	Author   string     `json:"author,omitempty"`
	Admin    bool       `json:"admin,omitempty"`
	Comment  string     `json:"comment"`
	Position [3]float64 `json:"position"`
	Created  time.Time  `json:"created"`
}

func InsertNote(ctx context.Context, note Note) (int, error) {
	pos, err := toVectorJSON(note.Position)
	if err != nil {
		return 0, err
	}

	r, err := conn.ExecContext(ctx, "INSERT INTO notes (author, map, position, comment) VALUES (?, ?, VEC_FromText(?), ?)", note.Author, note.Map, pos, note.Comment)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func DeleteNote(ctx context.Context, id int) error {
	_, err := conn.ExecContext(ctx, "DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func LatestNoteTime(ctx context.Context, steamid string, mapname string) (time.Time, error) {
	var created time.Time
	err := conn.QueryRowContext(ctx, "SELECT COALESCE(MAX(created), FROM_UNIXTIME(946702800)) FROM notes WHERE author = ? AND map = ?", steamid, mapname).Scan(&created)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return created, nil
}

func GetNotes(ctx context.Context, filter string, value string) ([]Note, error) {
	var args []any
	query := "SELECT n.id, n.author, u.admin, n.map, VEC_ToText(n.position), n.comment, n.created FROM notes n JOIN users u ON n.author = u.steamid"

	switch filter {
	case "steamid":
		query += " WHERE n.author = ?"
		args = append(args, value)
	case "map":
		query += " WHERE n.map = ?"
		args = append(args, value)
	}

	r, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	var notes []Note

	for r.Next() {
		var n Note
		var pos string
		err = r.Scan(&n.ID, &n.Author, &n.Admin, &n.Map, &pos, &n.Comment, &n.Created)
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

func GetNoteOwner(ctx context.Context, id int) (string, error) {
	var steamid string
	err := conn.QueryRowContext(ctx, "SELECT author FROM notes WHERE id = ?", id).Scan(&steamid)
	if err != nil {
		return "", err
	}

	return steamid, nil
}

func HasNoteWithinDistance(ctx context.Context, mapname string, position [3]float64, distance int) (bool, error) {
	pos, err := toVectorJSON(position)
	if err != nil {
		return false, err
	}

	var count int
	err = conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM notes WHERE map = ? AND VEC_DISTANCE_EUCLIDEAN(VEC_FromText(?), position) < ?", mapname, pos, distance).Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}
