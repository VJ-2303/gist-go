package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	maxBytes := 1_048_576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type at character %d", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be more than %d bytes", maxBytesError.Limit)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) readParamsID(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil {
		return 0, nil
	}
	return id, nil
}

func (app *application) generateHash(plainText string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (app *application) checkPassAndHash(plaintext string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
	if err != nil {
		return false
	}
	return true
}

func (app *application) GenerateURLSafeString(length int) (string, error) {
	// Calculate the number of random bytes needed. Each byte becomes about 1.33 characters in Base64,
	// so we calculate a slightly larger buffer to be safe and then trim.
	// For 8 chars, we need 6 random bytes. (6 * 8 / 6 = 8)
	byteLength := (length * 3) / 4

	// Create a byte slice to hold the random data.
	randomBytes := make([]byte, byteLength)

	// Read cryptographically secure random bytes into the slice.
	_, err := rand.Read(randomBytes)
	if err != nil {
		// If reading random bytes fails, it's a critical error.
		return "", fmt.Errorf("could not generate random bytes: %w", err)
	}

	// Encode the random bytes into a URL-safe Base64 string.
	// URLEncoding uses '-' and '_' instead of '+' and '/'
	encodedString := base64.URLEncoding.EncodeToString(randomBytes)

	// Base64 encoding can sometimes add padding ('='). We'll remove it.
	// For our calculation, it shouldn't happen, but it's good practice.
	// And we trim to the exact length requested.
	return encodedString[:length], nil
}
