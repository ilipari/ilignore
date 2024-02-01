/*
Copyright © 2024 Ignazio Lipari iglipari@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"ilipari/ilignore/cmd"
	"log/slog"
	"os"
)

func main() {
	var programLevel = new(slog.LevelVar) // Info by default
	programLevel.Set(slog.LevelWarn)
	// h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	// slog.SetDefault(slog.New(h))
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel}))
	slog.SetDefault(logger)

	cmd.Execute()
}
