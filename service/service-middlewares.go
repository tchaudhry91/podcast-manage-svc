package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/tchaudhry91/podcast-manage-svc/podcastmg"
	"time"
)

type loggingMiddleware struct {
	logger log.Logger
	next   PodcastManageService
}

func MakeNewLoggingMiddleware(logger log.Logger, next PodcastManageService) PodcastManageService {
	return loggingMiddleware{
		logger,
		next,
	}
}

func (mw loggingMiddleware) CreateUser(ctx context.Context, emailID, password string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "CreateUser",
			"user", emailID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.CreateUser(ctx, emailID, password)
	return
}

func (mw loggingMiddleware) GetUser(ctx context.Context, emailID string) (user podcastmg.User, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUser",
			"user", emailID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	user, err = mw.next.GetUser(ctx, emailID)
	return
}

func (mw loggingMiddleware) GetPodcastDetails(ctx context.Context, url string) (podcast podcastmg.Podcast, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetPodcastDetails",
			"url", url,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	podcast, err = mw.next.GetPodcastDetails(ctx, url)
	return
}

func (mw loggingMiddleware) Subscribe(ctx context.Context, emailID, podcastURL string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Subscribe",
			"user", emailID,
			"url", podcastURL,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.Subscribe(ctx, emailID, podcastURL)
	return
}

func (mw loggingMiddleware) Unsubscribe(ctx context.Context, emailID, podcastURL string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Unsubscribe",
			"user", emailID,
			"url", podcastURL,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.Unsubscribe(ctx, emailID, podcastURL)
	return
}
func (mw loggingMiddleware) GetUserSubscriptions(ctx context.Context, emailID string) (subscriptions []podcastmg.Podcast, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetUserSubscriptions",
			"user", emailID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	subscriptions, err = mw.next.GetUserSubscriptions(ctx, emailID)
	return
}

func (mw loggingMiddleware) GetSubscriptionDetails(ctx context.Context, emailID, podcastURL string) (podcast podcastmg.Podcast, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetSubscriptionDetails",
			"user", emailID,
			"podcast", podcastURL,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	podcast, err = mw.next.GetSubscriptionDetails(ctx, emailID, podcastURL)
	return
}

func (mw loggingMiddleware) GetToken(ctx context.Context, emailID, password string) (token string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetToken",
			"user", emailID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	token, err = mw.next.GetToken(ctx, emailID, password)
	return
}
