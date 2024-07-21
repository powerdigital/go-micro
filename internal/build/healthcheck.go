package build

import (
	"context"
	"net/http"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const readinessEndpoint = "/health"

type healthcheckFn func(ctx context.Context) error

type healthcheck struct {
	checks []healthcheckFn
}

func (h *healthcheck) add(fn healthcheckFn) {
	h.checks = append(h.checks, fn)
}

func (h *healthcheck) do(ctx context.Context) error {
	var (
		errs error
		mu   sync.Mutex
		wg   sync.WaitGroup
	)

	for _, c := range h.checks {
		wg.Add(1)

		go func(fn healthcheckFn) {
			if err := fn(ctx); err != nil {
				mu.Lock()
				errs = multierror.Append(err)
				mu.Unlock()
			}

			wg.Done()
		}(c)
	}

	wg.Wait()

	return errors.Wrap(errs, "healthcheck")
}

func (h *healthcheck) handler(w http.ResponseWriter, req *http.Request) {
	if err := h.do(req.Context()); err != nil {
		zerolog.Ctx(req.Context()).Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := req.Body.Close(); err != nil {
		zerolog.Ctx(req.Context()).Err(err).Msg("close healthcheck request body")
	}
}
