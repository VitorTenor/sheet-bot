package configuration

import (
	"context"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func BuildGoogleSrv(ctx context.Context, config *ApplicationConfig) (*sheets.Service, error) {
	jwtConfig := &jwt.Config{
		Email:      config.Google.ClientEmail,
		PrivateKey: []byte(strings.Replace(config.Google.PrivateKey, "\\n", "\n", -1)),
		Scopes:     []string{config.Google.ApiUrl},
		TokenURL:   google.JWTTokenURL,
	}
	googleClient := jwtConfig.Client(ctx)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(googleClient))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
