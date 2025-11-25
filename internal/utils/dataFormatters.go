package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// dataToJson converts a slice of maps containing data into a JSON string
func DataToJson(data []map[string]interface{}) (string, error) {
	enco, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(enco), nil
}

// dataToCSV converts a slice of maps into a CSV string
func DataToCSV(data []map[string]interface{}) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("no data to convert")
	}

	// Extract headers from the first map
	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}
	// Map iteration order is intentionally randomized, so we use sorting for consistency
	// See https://go.dev/blog/maps#iteration-order
	sort.Strings(headers)

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write headers
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write headers: %w", err)
	}

	// Write rows
	for _, row := range data {
		record := make([]string, len(headers))
		for i, header := range headers {
			if val, ok := row[header]; ok && val != nil {
				switch v := val.(type) {
				case string:
					record[i] = v
				case fmt.Stringer:
					record[i] = v.String()
				case int, int8, int16, int32, int64:
					record[i] = fmt.Sprintf("%d", v)
				case float32, float64:
					record[i] = strconv.FormatFloat(reflect.ValueOf(v).Float(), 'f', -1, 64)
				case bool:
					record[i] = strconv.FormatBool(v)
				default:
					record[i] = fmt.Sprintf("%v", v)
				}
			}
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("csv write error: %w", err)
	}

	return buf.String(), nil
}

// dataHTMLTable converts a slice of maps into an HTML table string.
// - rows: each map represents one table row (key -> cell value).
// - Columns are the union of all map keys, sorted alphabetically for deterministic output.
func DataToHTMLTable(rows []map[string]interface{}) (string, error) {
	var b strings.Builder

	// start table
	b.WriteString("<table>")

	// no rows - return error
	if len(rows) == 0 {
		return "", fmt.Errorf("no data to convert")
	}

	// collect all column names
	colSet := make(map[string]struct{})
	for _, r := range rows {
		for k := range r {
			colSet[k] = struct{}{}
		}
	}

	cols := make([]string, 0, len(colSet))
	for k := range colSet {
		cols = append(cols, k)
	}
	sort.Strings(cols) // deterministic order

	// header
	b.WriteString("<thead><tr>")
	for _, c := range cols {
		b.WriteString("<th>")
		b.WriteString(html.EscapeString(c))
		b.WriteString("</th>")
	}
	b.WriteString("</tr></thead>")

	// body
	b.WriteString("<tbody>")
	for _, r := range rows {
		b.WriteString("<tr>")
		for _, c := range cols {
			b.WriteString("<td>")
			val, ok := r[c]
			if ok && val != nil {
				// convert value to string and escape HTML
				b.WriteString(html.EscapeString(fmt.Sprint(val)))
			}
			// if missing or nil -> empty cell
			b.WriteString("</td>")
		}
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody>")

	// end table
	b.WriteString("</table>")

	return b.String(), nil
}
