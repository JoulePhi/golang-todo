package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ParseJSON(r *http.Request, v interface{}) error {
    if r.Body == nil {
        return fmt.Errorf("empty request body")
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(v)
}

func GetIDParam(r *http.Request) (int64, error) {
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) == 0 {
        return 0, fmt.Errorf("invalid URL path")
    }

    idStr := parts[len(parts)-1]
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        return 0, fmt.Errorf("invalid ID format")
    }

    return id, nil
}