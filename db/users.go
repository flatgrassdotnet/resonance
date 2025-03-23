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
	"crypto/rand"
	"time"
)

type User struct {
	SteamID string    `json:"steamid"`
	Admin   bool      `json:"admin"`
	Muted   bool      `json:"muted"`
	Created time.Time `json:"created"`
}

func InsertUser(ctx context.Context, user User) error {
	_, err := conn.ExecContext(ctx, "INSERT IGNORE INTO users (steamid) VALUES (?)", user.SteamID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(ctx context.Context, steamid string) error {
	_, err := conn.ExecContext(ctx, "DELETE FROM users WHERE steamid = ?", steamid)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(ctx context.Context, steamid string) (User, error) {
	var u User
	err := conn.QueryRowContext(ctx, "SELECT steamid, admin, muted, created FROM users WHERE steamid = ?", steamid).Scan(&u.SteamID, &u.Admin, &u.Muted, &u.Created)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func InsertTicket(ctx context.Context, steamid string, ticket []byte) error {
	_, err := conn.ExecContext(ctx, "REPLACE INTO tickets (steamid, ticket) VALUES (?, ?)", steamid, ticket)
	if err != nil {
		return err
	}

	return nil
}

func TokenFromTicket(ctx context.Context, ticket []byte) ([]byte, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}

	var steamid string
	err = conn.QueryRowContext(ctx, "SELECT steamid FROM tickets WHERE ticket = ? AND created > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 5 MINUTE)", ticket).Scan(&steamid)
	if err != nil {
		return nil, err
	}

	_, err = conn.ExecContext(ctx, "REPLACE INTO sessions (steamid, token) VALUES (?, ?)", steamid, token)
	if err != nil {
		return nil, err
	}

	_, err = conn.ExecContext(ctx, "DELETE FROM tickets WHERE ticket = ?", ticket)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func SteamIDFromToken(ctx context.Context, token []byte) (string, error) {
	var steamid string
	err := conn.QueryRowContext(ctx, "SELECT steamid FROM sessions WHERE token = ? AND created > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 MONTH)", token).Scan(&steamid)
	if err != nil {
		return "", err
	}

	return steamid, nil
}
