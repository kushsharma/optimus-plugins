package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/odpf/optimus/models"
	"github.com/odpf/optimus/plugin"
	"github.com/odpf/optimus/plugin/base"
)

var (
	Name = "neo"

	// will be injected by build system
	Version = "latest"
	Image   = "ghcr.io/kushsharma/optimus-task-neo-executor"

	_ models.CommandLineMod = &Neo{}
)

type Neo struct {
	log hclog.Logger
}

// GetSchema provides basic task details
func (n *Neo) PluginInfo() (*models.PluginInfoResponse, error) {
	return &models.PluginInfoResponse{
		Name:          Name,
		Description:   "Near earth object tracker",
		PluginType:    models.PluginTypeTask,
		PluginMods:    []models.PluginMod{models.ModTypeCLI},
		PluginVersion: Version,
		APIVersion:    []string{strconv.Itoa(base.ProtocolVersion)},

		// docker image that will be executed as the actual transformation task
		Image: fmt.Sprintf("%s:%s", Image, Version),
		// this is where the secret required by docker container will be mounted
		SecretPath: "/tmp/key.json",
	}, nil
}

// GetQuestions provides questions asked via optimus cli when a new job spec
// is requested to be created
func (n *Neo) GetQuestions(ctx context.Context, req models.GetQuestionsRequest) (*models.GetQuestionsResponse, error) {
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
	return &models.GetQuestionsResponse{
		Questions: tQues,
	}, nil
}

// ValidateQuestion validate how stupid user is
// Each question config in GetQuestions will send a validation request
func (n *Neo) ValidateQuestion(ctx context.Context, req models.ValidateQuestionRequest) (*models.ValidateQuestionResponse, error) {
	var err error
	switch req.Answer.Question.Name {
	case "Start":
		err = func(ans interface{}) error {
			d, ok := ans.(string)
			if !ok || d == "" {
				return errors.New("not a valid string")
			}
			// can choose to parse here for a valid date but we will use optimus
			// macros here {{.DSTART}} instead of actual dates
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
			// macros here {{.DEND}} instead of actual dates
			// _, err := time.Parse(time.RFC3339, d)
			// return err
			return nil
		}(req.Answer.Value)
	}
	if err != nil {
		return &models.ValidateQuestionResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}
	return &models.ValidateQuestionResponse{
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

// DefaultConfig are a set of key value pair which will be embedded in job
// specification. These configs can be requested by the docker container before
// execution
// PluginConfig Value can contain valid go templates and they will be parsed at
// runtime
func (n *Neo) DefaultConfig(ctx context.Context, request models.DefaultConfigRequest) (*models.DefaultConfigResponse, error) {
	start, _ := findAnswerByName("Start", request.Answers)
	end, _ := findAnswerByName("End", request.Answers)

	conf := []models.PluginConfig{
		{
			Name:  "RANGE_START",
			Value: start.Value,
		},
		{
			Name:  "RANGE_END",
			Value: end.Value,
		},
	}
	return &models.DefaultConfigResponse{
		Config: conf,
	}, nil
}

// DefaultAssets are a set of files which will be embedded in job
// specification in assets folder. These configs can be requested by the
// docker container before execution.
func (n *Neo) DefaultAssets(ctx context.Context, _ models.DefaultAssetsRequest) (*models.DefaultAssetsResponse, error) {
	return &models.DefaultAssetsResponse{}, nil
}

// override the compilation behaviour of assets - if needed
func (n *Neo) CompileAssets(ctx context.Context, req models.CompileAssetsRequest) (*models.CompileAssetsResponse, error) {
	return &models.CompileAssetsResponse{
		Assets: req.Assets,
	}, nil
}

func main() {
	plugin.Serve(func(log hclog.Logger) interface{} {
		return &Neo{
			log: log,
		}
	})
}
