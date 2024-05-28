package command

import (
	"context"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type CreateEnrollHost struct {
	//HostID     string
	EnrollCode string
	//UserID     string
}

type CreateEnrollHostResponse struct {
	Online     bool
	EnrollAt   int64
	LastSeenAt string
	HostID     string
	Config     string
}

type CreateEnrollHostHandler decorator.CommandHandler[*CreateEnrollHost, *CreateEnrollHostResponse]

func NewCreateEnrollHostHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) CreateEnrollHostHandler {
	return &createEnrollHostHandler{
		logger: logger,
		repos:  repos,
	}
}

type createEnrollHostHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *createEnrollHostHandler) Handle(ctx context.Context, cmd *CreateEnrollHost) (*CreateEnrollHostResponse, error) {
	host, err := h.repos.HostRepo.GetHostByEnrollCode(ctx, cmd.EnrollCode)
	if err != nil {
		return nil, err
	}

	r, err := h.repos.HostRepo.GetEnrollHost(ctx, &entity.GetEnrollHost{
		HostID: host.ID,
	})
	if err != nil {
		h.logger.Error("failed to get host", err)
		return nil, err
	}

	if host.EnrollAt > 0 {
		return nil, errors.New("host already enrolled")
	}

	if r.EnrollCodeExpiredAt < ctime.CurrentTimestampMillis() {
		return nil, errors.New("host enroll code expired")
	}

	if r.Code != cmd.EnrollCode {
		return nil, errors.New("主机注册码不存在或已过期")
	}

	enrollAt := ctime.CurrentTimestampMillis()

	if err := h.repos.HostRepo.EnrollHost(ctx, &entity.EnrollHost{
		HostID:   host.ID,
		Code:     cmd.EnrollCode,
		EnrollAt: enrollAt,
	}); err != nil {
		h.logger.Errorf("failed to enroll host: %v", err)
		return nil, err
	}

	ho, err := h.repos.HostRepo.HostOnline(ctx, &entity.HostOnline{
		ID: host.ID,
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

	if host.Role == "none" {
		lighthouses, err := h.repos.HostRepo.Find(ctx, &entity.HostFindOptions{
			UserID:    host.Owner,
			NetworkID: host.NetworkID,
			Role:      "lighthouse",
		})
		if err != nil {
			h.logger.WithError(err).Error("failed to find lighthouse host")
			return nil, err
		}
		host.Config.StaticHostMap = genLighthousesStaticMap(lighthouses)
		for _, lighthouse := range lighthouses {
			host.Config.Lighthouse.Hosts = append(host.Config.Lighthouse.Hosts, lighthouse.IPAddress)
		}
	}

	marshal, err := yaml.Marshal(&host.Config)
	if err != nil {
		h.logger.WithError(err).Error("failed to marshal host config")
		return nil, err
	}

	return &CreateEnrollHostResponse{
		EnrollAt:   enrollAt,
		LastSeenAt: lastSeenAt,
		Online:     online,
		HostID:     host.ID,
		Config:     string(marshal),
	}, nil
}

func genLighthousesStaticMap(lighthouses []*entity.Host) map[string][]string {
	lighthousesStaticMap := make(map[string][]string)
	for _, lighthouse := range lighthouses {
		lighthousesStaticMap[lighthouse.IPAddress] = []string{fmt.Sprintf("%s:%d", lighthouse.PublicIP, lighthouse.Port)}
	}
	return lighthousesStaticMap
}
