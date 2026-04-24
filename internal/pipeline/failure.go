package pipeline

import (
	"context"
	"fmt"

	"shipyard/internal/db"
	"shipyard/internal/docker"
	"shipyard/internal/events"
)

type FailureHandler struct {
	DB     *db.DB
	Bus    *events.Bus
	Docker *docker.Client
}

func (h *FailureHandler) Handle(ctx context.Context, evt events.PipelineEvent) error {
	id := evt.DeploymentID

	dep, err := h.DB.GetDeployment(ctx, id)
	if err != nil {
		return fmt.Errorf("get deployment: %w", err)
	}

	errMsg := evt.ErrorMsg
	if errMsg == "" {
		errMsg = "pipeline failed"
	}

	h.DB.AppendLogLine(ctx, id, "stderr", fmt.Sprintf("[error] %s", errMsg))

	dep.Status = db.StatusFailed
	h.DB.UpdateDeployment(ctx, dep)

	// Best-effort container cleanup
	if dep.ContainerID != "" {
		h.Docker.Stop(ctx, dep.ContainerID)
		h.Docker.Remove(ctx, dep.ContainerID)
	}

	return nil
}
