package service

import (
	"context"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
	"time"
)

var (
	// ErrDBConn indicates a failure to connect to the database
	ErrDBConn = errors.New("DB Connection Failed")

	// ErrUserCreate indicate a failure to create a user
	ErrUserCreate = errors.New("Failed to create User")

	// ErrUserFetch indicates a failure to fetch the user from the Datastore
	ErrUserFetch = errors.New("Failed to get user")

	//ErrPodcastBuild indicates a failure to build the podcast for the given URL
	ErrPodcastBuild = errors.New("Failed to build podcast from given URL")

	// ErrUserUpdate indicates a failure to save an updated user to the Datastore
	ErrUserUpdate = errors.New("Failed to save user update to Database")

	// ErrPodcastFetch indicates a failure to fetch user subscriptions from the Datastore
	ErrPodcastFetch = errors.New("Failed to fetch subscriptions")

	// ErrInvalidPassword indicates a failure to match password
	ErrInvalidPassword = errors.New("Invalid password provided")
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
	logger             log.Logger
	tokenSigningString string
}

// NewSQLStorePodcastManageService returns a pmg-svc backed by a SQL based DB Store
func NewSQLStorePodcastManageService(dialect, connectionString, tokenSigningString string, logger log.Logger) (PodcastManageService, error) {
	var svc podcastManageService
	store := podcastmg.NewDBStore(dialect, connectionString)
	err := store.Connect()
	if err != nil {
		logger.Log("err", err)
		return &svc, ErrDBConn
	}
	defer store.Close()
	err = store.Migrate()
	if err != nil {
		logger.Log("err", err)
		return &svc, errors.New("Error Migrating DB Structure")
	}

	svc = podcastManageService{
		store:              store,
		tokenSigningString: tokenSigningString,
		logger:             logger,
	}
	return &svc, nil
}

// CreateUser registers a new user in the store
func (svc *podcastManageService) CreateUser(ctx context.Context, emailID string, password string) error {
	user, err := podcastmg.NewUser(emailID, password)
	if err != nil {
		svc.logger.Log("err", err)
		return ErrUserCreate
	}
	err = svc.store.Connect()
	if err != nil {
		svc.logger.Log("err", err)
		return ErrDBConn
	}
	defer svc.store.Close()
	err = svc.store.CreateUser(&user)
	if err != nil {
		svc.logger.Log("err", err)
		return ErrUserCreate
	}
	return nil
}

// GetUser returns a user object if found in the store
func (svc *podcastManageService) GetUser(ctx context.Context, emailID string) (podcastmg.User, error) {
	var user podcastmg.User
	err := svc.store.Connect()
	if err != nil {
		svc.logger.Log("err", err)
		return user, ErrDBConn
	}
	defer svc.store.Close()
	user, err = svc.store.GetUserByEmail(emailID)
	if err != nil {
		svc.logger.Log("err", err)
		return user, ErrUserFetch
	}
	return user, nil
}

// GetPodcastDetails returns a parsed Podcast object for a given xml-url
func (svc *podcastManageService) GetPodcastDetails(ctx context.Context, url string) (podcastmg.Podcast, error) {
	podcast, err := podcastmg.BuildPodcastFromURL(url)
	if err != nil {
		svc.logger.Log("err", err)
		return podcast, ErrPodcastBuild
	}
	return podcast, nil
}

// Subscribe adds a podcast subscription to a user and saves it in the database
func (svc *podcastManageService) Subscribe(ctx context.Context, emailID, podcastURL string) error {
	err := svc.store.Connect()
	if err != nil {
		svc.logger.Log("err", err)
		return ErrDBConn
	}
	defer svc.store.Close()
	user, err := svc.store.GetUserByEmail(emailID)
	if err != nil {
		svc.logger.Log("err", err)
		return ErrUserFetch
	}
	podcast, err := podcastmg.BuildPodcastFromURL(podcastURL)
	if err != nil {
		svc.logger.Log("err", err)
		return ErrPodcastBuild
	}
	user.AddSubscription(podcast)
	err = svc.store.UpdateUser(&user)
	if err != nil {
		svc.logger.Log("err", err)
		return ErrUserUpdate
	}
	return nil
}

// GetUserSubscriptions returns a list of podcasts that the user is subscribed to
func (svc *podcastManageService) GetUserSubscriptions(ctx context.Context, emailID string) ([]podcastmg.Podcast, error) {
	var subscriptions []podcastmg.Podcast
	err := svc.store.Connect()
	if err != nil {
		svc.logger.Log("err", err)
		return subscriptions, ErrDBConn
	}
	defer svc.store.Close()
	user, err := svc.store.GetUserByEmail(emailID)
	if err != nil {
		svc.logger.Log("err", err)
		return subscriptions, ErrUserFetch
	}
	return user.GetSubscriptions(), nil
}

// GetSubscriptionDetails returns a populated podcast with items based on the user subscription
func (svc *podcastManageService) GetSubscriptionDetails(ctx context.Context, emailID, podcastURL string) (podcastmg.Podcast, error) {
	var podcast podcastmg.Podcast
	err := svc.store.Connect()
	if err != nil {
		svc.logger.Log("err", err)
		return podcast, ErrDBConn
	}
	defer svc.store.Close()
	podcast, err = svc.store.GetPodcastBySubscription(emailID, podcastURL)
	if err != nil {
		svc.logger.Log("err", err)
		return podcast, ErrPodcastFetch
	}
	return podcast, nil
}

// GetToken returns a JWT token for service authorization
func (svc *podcastManageService) GetToken(ctx context.Context, emailID string, password string) (tokenString string, err error) {
	err = svc.store.Connect()
	if err != nil {
		svc.logger.Log("err", err)
		return tokenString, ErrDBConn
	}
	user, err := svc.store.GetUserByEmail(emailID)
	if err != nil {
		svc.logger.Log("err", err)
		return tokenString, ErrUserFetch
	}
	err = user.ComparePassword(password)
	if err != nil {
		svc.logger.Log("err", err)
		return tokenString, ErrInvalidPassword
	}

	claims := TokenClaims{
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
