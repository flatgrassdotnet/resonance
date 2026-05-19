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

package pages

import (
	"embed"
	"html/template"
	"io/fs"
)

type MainData struct {
	Header string
	Body   string
	Footer string
}

var (
	//go:embed templates
	templates      embed.FS
	TemplatesFS, _ = fs.Sub(templates, "templates")

	//go:embed assets
	assets      embed.FS
	AssetsFS, _ = fs.Sub(assets, "assets")

	Main = template.Must(template.New("main.html").ParseFS(TemplatesFS, "main.html"))
)
