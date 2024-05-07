package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/network/entity"
)

func (h *NetworkHandler) Update(ctx context.Context, cmd *network.UpdateNetwork) (*entity.Network, error) {
	if cmd.Name == "" {
		return nil, errors.New("name is required")
	}

	updatedNetwork, err := h.repo.Update(ctx, &entity.Network{
		ID:   cmd.ID,
		Name: cmd.Name,
	})
	if err != nil {
		h.logger.WithField("cmd", cmd).Error("error updating network")
		return nil, err
	}

	return updatedNetwork, nil
}
