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

import "context"

func GetUserCount(ctx context.Context) (int, error) {
	var count int
	err := conn.QueryRowContext(ctx, "SELECT COUNT(DISTINCT author) FROM notes").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetNoteCount(ctx context.Context) (int, error) {
	var count int
	err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM notes").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetNoteCountByUserMap(ctx context.Context, steamid string, mapname string) (int, error) {
	var count int
	err := conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM notes WHERE author = ? AND map = ?", steamid, mapname).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetMapCount(ctx context.Context) (int, error) {
	var count int
	err := conn.QueryRowContext(ctx, "SELECT COUNT(DISTINCT map) FROM notes").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetMaps(ctx context.Context) (map[string]int, error) {
	r, err := conn.QueryContext(ctx, "SELECT map, COUNT(map) FROM notes GROUP BY map")
	if err != nil {
		return nil, err
	}

	maps := make(map[string]int)

	for r.Next() {
		var mapname string
		var count int
		err = r.Scan(&mapname, &count)
		if err != nil {
			return nil, err
		}

		maps[mapname] = count
	}

	return maps, nil
}
