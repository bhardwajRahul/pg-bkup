// Package utils /
package utils

/*
MIT License

Copyright (c) 2023 Jonas Kaninda

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

import "os"

const RestoreExample = "restore --dbname database --file db_20231219_022941.sql.gz\n" +
	"restore --dbname database --storage s3 --path /custom-path --file db_20231219_022941.sql.gz"
const BackupExample = "backup --dbname database --disable-compression\n" +
	"backup --dbname database --storage s3 --path /custom-path --disable-compression"

const MainExample = "backup --dbname database --disable-compression\n" +
	"backup --dbname database --storage s3 --path /custom-path\n" +
	"restore --dbname database --file db_20231219_022941.sql.gz"

var Version string

func VERSION(def string) string {
	build := os.Getenv("VERSION")
	if build == "" {
		return def
	}
	return build
}
func FullVersion() string {
	ver := Version
	if b := VERSION(""); b != "" {
		return b
	}
	return ver
}
