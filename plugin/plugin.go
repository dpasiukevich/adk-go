package plugin

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/genai"
)

type Config struct {
	OnUserMessageCallback OnUserMessageCallback

	BeforeRunCallback BeforeRunCallback
	AfterRunCallback  AfterRunCallback

	BeforeAgentCallback agent.BeforeAgentCallback
	AfterAgentCallback  agent.AfterAgentCallback

	BeforeModelCallback llmagent.BeforeModelCallback
	AfterModelCallback  llmagent.AfterModelCallback

	BeforeToolCallback llmagent.BeforeToolCallback
	AfterToolCallback  llmagent.AfterToolCallback

	CloseFunc func() error
}

func New(cfg Config) (*Plugin, error) {
	return &Plugin{}, nil
}

type Plugin struct{}

type OnUserMessageCallback func(agent.InvocationContext, *genai.Content) (*genai.Content, error)

type BeforeRunCallback func(agent.InvocationContext) (*genai.Content, error)

type AfterRunCallback func(agent.InvocationContext, *genai.Content)
