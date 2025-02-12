package configs

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	CredentialsFile = "credentials.json"
	GoogleApiUrl    = "https://www.googleapis.com/auth/spreadsheets"
)

type Credentials struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

func loadCredentials() (*Credentials, error) {
	file, err := os.ReadFile(CredentialsFile)
	if err != nil {
		return nil, err
	}
	var credentials Credentials
	err = json.Unmarshal(file, &credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func BuildGoogleSrv(ctx context.Context) (*sheets.Service, error) {
	credentials, err := loadCredentials()
	if err != nil {
		return nil, err
	}
	return createJWTClient(ctx, credentials)
}

func createJWTClient(ctx context.Context, credentials *Credentials) (*sheets.Service, error) {
	config := &jwt.Config{
		Email:      credentials.ClientEmail,
		PrivateKey: []byte(strings.Replace(credentials.PrivateKey, "\\n", "\n", -1)),
		Scopes:     []string{GoogleApiUrl},
		TokenURL:   google.JWTTokenURL,
	}
	googleClient := config.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
