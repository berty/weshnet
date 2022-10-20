package tinder

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"berty.tech/berty/v2/go/internal/logutil"
)

const defaultTTL = time.Hour

type AdvertiseOptions struct {
	Filters map[string]struct{}
}
type AdvertiseOption func(opts *AdvertiseOptions) error

func (o *AdvertiseOptions) apply(opts ...AdvertiseOption) error {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return fmt.Errorf("uanble to apply option: %w", err)
		}
	}

	return nil
}

func StartAdvertisesFilterDrivers(drivers ...string) AdvertiseOption {
	return func(opts *AdvertiseOptions) error {
		opts.Filters = map[string]struct{}{}
		for _, driver := range drivers {
			opts.Filters[driver] = struct{}{}
		}

		return nil
	}
}

// Register advertise topic on each of his driver
func (s *Service) StartAdvertises(ctx context.Context, topic string, opts ...AdvertiseOption) error {
	if len(s.drivers) == 0 {
		return fmt.Errorf("no driver available to advertise")
	}

	var aopts AdvertiseOptions
	if err := aopts.apply(opts...); err != nil {
		return fmt.Errorf("failed to advertise: %w", err)
	}

	for _, driver := range s.drivers {
		if aopts.Filters != nil {
			// skip filter driver
			if _, filter := aopts.Filters[driver.Name()]; filter {
				continue
			}
		}

		// start background job
		go func(driver IDriver) {
			if err := s.advertise(ctx, driver, topic); err != nil {
				s.logger.Debug("advertise ended", zap.Error(err))
			}
		}(driver)
	}

	return nil
}

func (s *Service) advertise(ctx context.Context, d IDriver, topic string) error {
	for {
		currentAddrs := s.networkNotify.GetLastUpdatedAddrs(ctx)

		now := time.Now()
		ttl, err := d.Advertise(ctx, topic)
		took := time.Since(now)

		var deadline time.Duration
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// retry in 30 seconds
			deadline = time.Second * 30
		} else {
			if ttl == 0 {
				ttl = defaultTTL
			}
			deadline = 4 * ttl / 5
		}

		s.logger.Debug("advertise",
			zap.String("driver", d.Name()),
			logutil.PrivateString("topic", topic),
			zap.Duration("ttl", ttl),
			zap.Duration("took", took),
			zap.Duration("next", deadline),
			zap.Error(err),
		)

		waitctx, cancel := context.WithTimeout(ctx, deadline)
		// wait for network update or waitctx expire
		_, ok := s.networkNotify.WaitForUpdate(waitctx, currentAddrs)
		cancel()

		// filter addrs

		if !ok {
			select {
			// check for parent ctx
			case <-ctx.Done():
				// main context has expire, stop
				return ctx.Err()
			default:
				// wait context has expire, continue
			}
		}
	}
}
