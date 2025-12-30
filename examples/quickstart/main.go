// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package provides a quickstart ADK agent.
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/genai"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/internal/context/gocontext"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/plugin"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	logPlugin, err := plugin.New(plugin.Config{
		BeforeAgentCallback: func(ctx agent.CallbackContext) (*genai.Content, error) {
			log.Printf("BeforeAgentCallback invoked for agent: %s", ctx.AgentName())
			return nil, nil // no-op, keeping original content
		},

		AfterRunCallback: func(ic agent.InvocationContext, c *genai.Content) {
			log.Printf("AfterRunCallback invoked for agent: %s", ic.Agent().Name())
		},
	})

runner.New(runner.Config{
	PluginConfig: runner.PluginConfig{
		Plugins: []plugin.Plugin{
			logPlugin,
		},
	},
})

	a, err := llmagent.New(llmagent.Config{
		Name:        "weather_time_agent",
		Model:       model,
		Description: "Agent to answer questions about the time and weather in a city.",
		Instruction: "Your SOLE purpose is to answer questions about the current time and weather in a specific city. You MUST refuse to answer any questions unrelated to time or weather.",
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				log.Printf("Starting agent invocation: %s", ctx.InvocationID())

				goctx, ok := ctx.(gocontext.Holder)
				if !ok {
					log.Printf("Context does not implement gocontext.Holder")
				}

				if err := goctx.SetGoContext(context.WithValue(goctx.GoContext(), "key", "value")); err != nil {
					log.Printf("Failed to set context value: %v", err)
				}

				log.Printf("Set context value: %s", goctx.GoContext().Value("key"))

				return nil, nil
			},
		},
		AfterAgentCallbacks: []agent.AfterAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				goctx, ok := ctx.(gocontext.Holder)
				if !ok {
					log.Printf("Context does not implement gocontext.Holder")
				}

				log.Printf("GOT_VALUE: %s", goctx.GoContext().Value("key"))

				return nil, nil
			},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
