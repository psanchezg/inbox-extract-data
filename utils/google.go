package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getToken(ctx context.Context, scope ...string) (*oauth2.Config, *oauth2.Token) {
	// If you don't have a client_secret.json, go here (as of July 2017):
	// https://auth0.com/docs/connections/social/google
	// https://web.archive.org/web/20170708123613/https://auth0.com/docs/connections/social/google

	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	/////////////////////////////////////////////////////////////////////////////
	// SCOPE https://developers.google.com/gmail/api/auth/scopes
	/////////////////////////////////////////////////////////////////////////////
	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	// If modifying this code to access other features of gmail, you may need
	// to change scope.
	if !slices.Contains(scope, gmail.GmailReadonlyScope) {
		scope = append(scope, gmail.GmailReadonlyScope)
	}
	if !slices.Contains(scope, sheets.SpreadsheetsScope) {
		scope = append(scope, sheets.SpreadsheetsScope)
	}
	if !slices.Contains(scope, drive.DriveReadonlyScope) {
		scope = append(scope, drive.DriveReadonlyScope)
	}
	config, err := google.ConfigFromJSON(b, scope...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	cacheFile, err := newTokenizer()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	token, err := tokenFromFile(cacheFile)
	if err != nil {
		token = getTokenFromWeb(config)
		saveToken(cacheFile, token)
	}
	return config, token
}

// ConnectToGmailService uses a Context, config and 1 or more scopes to retrieve a
// Token then generate a Client. It uses the client to connect to the gmail api
// service and returns a *gmail.Service.
func ConnectToGmailService(ctx context.Context, scope ...string) *gmail.Service {
	config, token := getToken(ctx, scope...)
	srv, err := gmail.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}
	return srv
}

// ConnectToSheetsService uses a Context, config and 1 or more scopes to retrieve a
// Token then generate a Client. It uses the client to connect to the gmail api
// service and returns a *gmail.Service.
func ConnectToSheetsService(ctx context.Context, scope ...string) *sheets.Service {
	config, token := getToken(ctx, scope...)
	srv, err := sheets.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to retrieve sheets Client %v", err)
	}
	return srv
}

// ConnectToDriveService uses a Context, config and 1 or more scopes to retrieve a
// Token then generate a Client. It uses the client to connect to the gmail api
// service and returns a *gmail.Service.
func ConnectToDriveService(ctx context.Context, scope ...string) *drive.Service {
	config, token := getToken(ctx, scope...)
	srv, err := drive.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}
	return srv
}

// newTokenizer returns a new token and generates credential file path and
// returns the generated credential path/filename along with any errors.
func newTokenizer() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir, url.QueryEscape("google.json")),
		err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read errors encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// getTokenFromWeb uses Config to request a Token. It returns the retrieved
// Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// saveToken uses a file path to create a file and store the token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// columnIndexToLetter convierte un índice de columna a la notación de letra de columna de Sheets (1 -> A, 2 -> B, ..., 27 -> AA, ...)
func ColumnIndexToLetter(index int) string {
	result := ""
	for index > 0 {
		index--
		result = string(rune('A'+(index%26))) + result
		index /= 26
	}
	return result
}

// columnToIndex convierte la letra de la columna en su índice correspondiente
func ColumnToIndex(column string) int {
	column = strings.ToUpper(column) // Asegurar que la columna esté en mayúsculas
	index := 0
	for i := 0; i < len(column); i++ {
		index = index*26 + int(column[i]-'A'+1)
	}
	return index - 1
}

// SplitCellReference separa una referencia de celda en sus componentes de columna y fila
func SplitCellReference(cell string) (string, int) {
	var column string
	var rowStr string

	for _, char := range cell {
		if unicode.IsLetter(char) {
			column += string(char)
		} else if unicode.IsDigit(char) {
			rowStr += string(char)
		}
	}

	row, err := strconv.Atoi(rowStr)
	if err != nil {
		// Manejar el error según sea necesario
		fmt.Println("Error al convertir la fila a un número:", err)
		return "", 0
	}

	return column, row
}
