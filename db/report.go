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
	"context"
	"time"
)

type Report struct {
	ID      int       `json:"id"`
	NoteID  int       `json:"noteid"`
	SteamID string    `json:"steamid"`
	Reason  string    `json:"reason"`
	Created time.Time `json:"created"`
}

func InsertReport(ctx context.Context, report Report) (int, error) {
	r, err := conn.ExecContext(ctx, "INSERT INTO reports (noteid, steamid, reason) VALUES (?, ?, ?)", report.NoteID, report.SteamID, report.Reason)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func LatestReportTime(ctx context.Context, steamid string) (time.Time, error) {
	var created time.Time
	err := conn.QueryRowContext(ctx, "SELECT COALESCE(MAX(created), FROM_UNIXTIME(946702800)) FROM reports WHERE steamid = ?", steamid).Scan(&created)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return created, nil
}
