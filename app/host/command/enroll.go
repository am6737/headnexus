package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/common/code"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
)

const (
	// lifetimeSeconds 24 hours
	lifetimeSeconds = 86400
)

func (h *HostHandler) GenerateEnrollCode(ctx context.Context, cmd *host.GenerateEnrollCode) (*host.GenerateEnrollCodeResponse, error) {
	r, err := h.repo.GetEnrollHost(ctx, &entity.GetEnrollHost{
		HostID: cmd.HostID,
	})
	if err != nil {
		h.logger.Error("failed to get host", err)
		return nil, err
	}

	if r.EnrollAt > 0 {
		return nil, errors.New("host already enrolled")
	}

	ecode := code.GenEnrollCode()

	if err := h.repo.EnrollHost(ctx, &entity.EnrollHost{
		HostID:              cmd.HostID,
		Code:                ecode,
		EnrollCodeExpiredAt: ctime.CurrentTimestampMillis() + lifetimeSeconds,
	}); err != nil {
		h.logger.Errorf("failed to enroll host: %v", err)
		return nil, err
	}

	return &host.GenerateEnrollCodeResponse{
		Code:            ecode,
		LifetimeSeconds: lifetimeSeconds,
	}, nil
}

func (h *HostHandler) EnrollCodeCheck(ctx context.Context, cmd *host.GenerateEnrollCodeCheck) (*host.GenerateEnrollCodeCheckResponse, error) {
	r, err := h.repo.GetEnrollHost(ctx, &entity.GetEnrollHost{
		HostID: cmd.HostID,
	})
	if err != nil {
		h.logger.Error("failed to get host", err)
		return nil, err
	}

	if r.EnrollAt > 0 {
		return &host.GenerateEnrollCodeCheckResponse{Exists: false}, nil
	}

	return &host.GenerateEnrollCodeCheckResponse{Exists: true}, nil
}

func (h *HostHandler) EnrollHost(ctx context.Context, cmd *host.EnrollHost) (*host.EnrollHostResponse, error) {
	r, err := h.repo.GetEnrollHost(ctx, &entity.GetEnrollHost{
		HostID: cmd.HostID,
	})
	if err != nil {
		h.logger.Error("failed to get host", err)
		return nil, err
	}

	if r.EnrollAt > 0 {
		return nil, errors.New("host already enrolled")
	}

	fmt.Println("EnrollCodeExpiredAt => ", ctime.FormatTimestamp(r.EnrollCodeExpiredAt))
	fmt.Println("ctime.CurrentTimestampMillis() => ", ctime.FormatTimestamp(ctime.CurrentTimestampMillis()))

	if r.EnrollCodeExpiredAt < ctime.CurrentTimestampMillis() {
		return nil, errors.New("host enroll code expired")
	}

	ecode := code.GenEnrollCode()
	enrollAt := ctime.CurrentTimestampMillis()

	if err := h.repo.EnrollHost(ctx, &entity.EnrollHost{
		HostID:   cmd.HostID,
		Code:     ecode,
		EnrollAt: enrollAt,
	}); err != nil {
		h.logger.Errorf("failed to enroll host: %v", err)
		return nil, err
	}

	ho, err := h.repo.HostOnline(ctx, &entity.HostOnline{
		ID: cmd.HostID,
	})
	if err != nil {
		h.logger.Errorf("failed to host online: %v", err)
		return nil, err
	}

	var lastSeenAt string
	var online bool
	if ho.Status == entity.Online {
		lastSeenAt = "Online"
		online = true
	} else {
		lastSeenAt = ctime.FormatTimeSince(ho.LastSeenAt)
	}

	return &host.EnrollHostResponse{
		EnrollAt:   enrollAt,
		LastSeenAt: lastSeenAt,
		Online:     online,
	}, nil
}
