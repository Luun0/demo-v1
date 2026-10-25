// Package migrations embeds the SQL schema files.
package migrations

import "embed"

// FS holds *.sql migrations applied in lexical order.
//
//go:embed *.sql
var FS embed.FS
