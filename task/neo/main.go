package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/odpf/optimus/plugin"

	"github.com/odpf/optimus/models"
	"github.com/odpf/optimus/plugin/task"

	hplugin "github.com/hashicorp/go-plugin"
)

var (
	Name           = "neo"
	DatetimeFormat = "2006-01-02"
	Version        = "latest"
	Image          = "ghcr.io/kushsharma/optimus-task-neo-executer"
)

type Neo struct{}

// basic task details
func (n *Neo) GetTaskSchema(ctx context.Context, req models.GetTaskSchemaRequest) (models.GetTaskSchemaResponse, error) {
	return models.GetTaskSchemaResponse{
		Name:        Name,
		Description: "Near earth object tracker",

		// docker image that will be executed as the actual transformation task
		Image:      fmt.Sprintf("%s:%s", Image, Version),
		SecretPath: "/tmp/key.json",
	}, nil
}

// will be asked via opctl
func (n *Neo) GetTaskQuestions(ctx context.Context, req models.GetTaskQuestionsRequest) (models.GetTaskQuestionsResponse, error) {
	tQues := []models.PluginQuestion{
		{
			Name:   "Start",
			Prompt: "Date range start",
			Help:   "YYYY-MM-DD format",
		},
		{
			Name:   "End",
			Prompt: "Date range end",
			Help:   "YYYY-MM-DD format",
		},
	}
	return models.GetTaskQuestionsResponse{
		Questions: tQues,
	}, nil
}

// validate how stupid user is
func (n *Neo) ValidateTaskQuestion(ctx context.Context, req models.ValidateTaskQuestionRequest) (models.ValidateTaskQuestionResponse, error) {
	var err error
	switch req.Answer.Question.Name {
	case "Start":
		err = func(ans interface{}) error {
			d, ok := ans.(string)
			if !ok || d == "" {
				return errors.New("not a valid string")
			}
			// can choose to parse here for a valid date but we will use optimus
			// macros here instead of actual dates
			// _, err := time.Parse(time.RFC3339, d)
			// return err
			return nil
		}(req.Answer.Value)
	case "End":
		err = func(ans interface{}) error {
			d, ok := ans.(string)
			if !ok || d == "" {
				return errors.New("not a valid string")
			}
			// can choose to parse here for a valid date but we will use optimus
			// macros here instead of actual dates
			// _, err := time.Parse(time.RFC3339, d)
			// return err
			return nil
		}(req.Answer.Value)
	}
	if err != nil {
		return models.ValidateTaskQuestionResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return models.ValidateTaskQuestionResponse{
		Success: true,
	}, nil
}

func findAnswerByName(name string, answers []models.PluginAnswer) (models.PluginAnswer, bool) {
	for _, ans := range answers {
		if ans.Question.Name == name {
			return ans, true
		}
	}
	return models.PluginAnswer{}, false
}

// configs we will need from users
func (n *Neo) DefaultTaskConfig(ctx context.Context, request models.DefaultTaskConfigRequest) (models.DefaultTaskConfigResponse, error) {
	start, _ := findAnswerByName("Start", request.Answers)
	end, _ := findAnswerByName("End", request.Answers)

	conf := []models.TaskPluginConfig{
		{
			Name:  "RANGE_START",
			Value: start.Value,
		},
		{
			Name:  "RANGE_END",
			Value: end.Value,
		},
	}
	return models.DefaultTaskConfigResponse{
		Config: conf,
	}, nil
}

// none
func (n *Neo) DefaultTaskAssets(ctx context.Context, _ models.DefaultTaskAssetsRequest) (models.DefaultTaskAssetsResponse, error) {
	return models.DefaultTaskAssetsResponse{}, nil
}

// passthrough
func (n *Neo) CompileTaskAssets(ctx context.Context, req models.CompileTaskAssetsRequest) (models.CompileTaskAssetsResponse, error) {
	return models.CompileTaskAssetsResponse{
		Assets: req.Assets,
	}, nil
}

// a task should always have a destination, it could be endpoint, table, bucket, etc
// in our case it is actually nothing
func (n *Neo) GenerateTaskDestination(ctx context.Context, request models.GenerateTaskDestinationRequest) (models.GenerateTaskDestinationResponse, error) {
	return models.GenerateTaskDestinationResponse{
		Destination: "none",
	}, nil
}

// as this task doesn't need dependency resolution, just leaving this empty
func (n *Neo) GenerateTaskDependencies(ctx context.Context, request models.GenerateTaskDependenciesRequest) (response models.GenerateTaskDependenciesResponse, err error) {
	return response, nil
}

func main() {
	neo := &Neo{}

	var handshakeConfig = hplugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   plugin.MagicCookieKey,
		MagicCookieValue: plugin.MagicCookieValue,
	}
	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: map[string]hplugin.Plugin{
			plugin.TaskPluginName: &task.Plugin{Impl: neo},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: hplugin.DefaultGRPCServer,
	})
}
