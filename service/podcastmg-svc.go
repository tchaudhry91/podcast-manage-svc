package service

import (
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
	"time"
)

// TokenClaims is a custom claims struct to issue JWT tokens
type TokenClaims struct {
	EmailID string `json:"email_id"`
	jwt.StandardClaims
}

// PodcastManageService is service to manage podcast rss-feeds
type PodcastManageService interface {
	CreateUser(ctx context.Context, emailID, password string) error
	GetUser(ctx context.Context, emailID string) (podcastmg.User, error)
	GetPodcastDetails(ctx context.Context, url string) (podcastmg.Podcast, error)
	Subscribe(ctx context.Context, emailID, podcastURL string) error
	GetUserSubscriptions(ctx context.Context, emailID string) ([]podcastmg.Podcast, error)
	GetSubscriptionDetails(ctx context.Context, emailID, podcastURL string) (podcastmg.Podcast, error)
	GetToken(ctx context.Context, emailID, password string) (string, error)
}

type podcastManageService struct {
	store              podcastmg.Store
	tokenSigningString string
}

// NewSQLStorePodcastManageService returns a pmg-svc backed by a SQL based DB Store
func NewSQLStorePodcastManageService(dialect, connectionString, tokenSigningString string) (PodcastManageService, error) {
	var svc podcastManageService
	store := podcastmg.NewDBStore(dialect, connectionString)
	err := store.Connect()
	if err != nil {
		return &svc, err
	}
	defer store.Close()
	err = store.Migrate()
	if err != nil {
		return &svc, err
	}

	svc = podcastManageService{
		store:              store,
		tokenSigningString: tokenSigningString,
	}
	return &svc, nil
}

// CreateUser registers a new user in the store
func (svc *podcastManageService) CreateUser(ctx context.Context, emailID string, password string) error {
	user, err := podcastmg.NewUser(emailID, password)
	if err != nil {
		return err
	}
	err = svc.store.Connect()
	if err != nil {
		return err
	}
	defer svc.store.Close()
	err = svc.store.CreateUser(&user)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns a user object if found in the store
func (svc *podcastManageService) GetUser(ctx context.Context, emailID string) (podcastmg.User, error) {
	var user podcastmg.User
	err := svc.store.Connect()
	if err != nil {
		return user, err
	}
	defer svc.store.Close()
	user, err = svc.store.GetUserByEmail(emailID)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GetPodcastDetails returns a parsed Podcast object for a given xml-url
func (svc *podcastManageService) GetPodcastDetails(ctx context.Context, url string) (podcastmg.Podcast, error) {
	return podcastmg.BuildPodcastFromURL(url)
}

// Subscribe adds a podcast subscription to a user and saves it in the database
func (svc *podcastManageService) Subscribe(ctx context.Context, emailID, podcastURL string) error {
	err := svc.store.Connect()
	if err != nil {
		return err
	}
	defer svc.store.Close()
	user, err := svc.store.GetUserByEmail(emailID)
	if err != nil {
		return err
	}
	podcast, err := podcastmg.BuildPodcastFromURL(podcastURL)
	if err != nil {
		return err
	}
	user.AddSubscription(podcast)
	err = svc.store.UpdateUser(&user)
	if err != nil {
		return err
	}
	return nil
}

// GetUserSubscriptions returns a list of podcasts that the user is subscribed to
func (svc *podcastManageService) GetUserSubscriptions(ctx context.Context, emailID string) ([]podcastmg.Podcast, error) {
	var subscriptions []podcastmg.Podcast
	err := svc.store.Connect()
	if err != nil {
		return subscriptions, err
	}
	defer svc.store.Close()
	user, err := svc.store.GetUserByEmail(emailID)
	if err != nil {
		return subscriptions, err
	}
	return user.GetSubscriptions(), nil
}

// GetSubscriptionDetails returns a populated podcast with items based on the user subscription
func (svc *podcastManageService) GetSubscriptionDetails(ctx context.Context, emailID, podcastURL string) (podcastmg.Podcast, error) {
	var podcast podcastmg.Podcast
	err := svc.store.Connect()
	if err != nil {
		return podcast, err
	}
	defer svc.store.Close()
	podcast, err = svc.store.GetPodcastBySubscription(emailID, podcastURL)
	if err != nil {
		return podcast, err
	}
	return podcast, nil
}

// GetToken returns a JWT token for service authorization
func (svc *podcastManageService) GetToken(ctx context.Context, emailID string, password string) (tokenString string, err error) {
	err = svc.store.Connect()
	if err != nil {
		return tokenString, err
	}
	user, err := svc.store.GetUserByEmail(emailID)
	if err != nil {
		return tokenString, err
	}
	err = user.ComparePassword(password)
	if err != nil {
		return tokenString, err
	}

	claims :=
		TokenClaims{
			emailID,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			},
		}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(svc.tokenSigningString))
	if err != nil {
		return tokenString, err
	}
	return tokenString, err
}
