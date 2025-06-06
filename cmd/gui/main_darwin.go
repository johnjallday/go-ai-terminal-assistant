//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"strings"

	"go-ai-terminal-assistant/agents"
	chat "go-ai-terminal-assistant/internal/chat"
	"go-ai-terminal-assistant/models"
	router "go-ai-terminal-assistant/router"
	convstorage "go-ai-terminal-assistant/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	fstorage "fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	gohook "github.com/robotn/gohook"
)

func main() {
	a := app.NewWithID("go-ai-terminal-assistant.gui")
	w := a.NewWindow("AI Chat")
	w.Resize(fyne.NewSize(600, 400))

	session := chat.NewSession()

	historyText := widget.NewMultiLineEntry()
	historyText.Wrapping = fyne.TextWrapWord
	scroll := container.NewVScroll(historyText)
	scroll.SetMinSize(fyne.NewSize(600, 350))

	// Status bar shows currently enabled agents
	statusBar := widget.NewLabel("")
	updateStatusBar := func() {
		var names []string
		for _, reg := range session.Router.ListAllAgents() {
			if reg.Enabled {
				names = append(names, reg.Agent.GetName())
			}
		}
		if len(names) > 0 {
			statusBar.SetText("Enabled Agents: " + strings.Join(names, ", "))
		} else {
			statusBar.SetText("Enabled Agents: none")
		}
	}
	updateStatusBar()

	var reloadAgents func()
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Type a message or /command and press Enter or /help")
	entry.OnSubmitted = func(text string) {
		if text == "" {
			return
		}
		fyne.Do(func() {
			addMessage(historyText, "You", text)
			scroll.ScrollToBottom()
			entry.Disable()
		})
		go func(input string) {
			if strings.HasPrefix(input, "/") {
				fyne.Do(func() {
					handleCommand(input, session, historyText, scroll, w, updateStatusBar)
					if reloadAgents != nil {
						reloadAgents()
					}
				})
			} else {
				response, err := session.ProcessMessage(input)
				fyne.Do(func() {
					if err != nil {
						addMessage(historyText, "Error", err.Error())
					} else {
						addMessage(historyText, "Assistant", response)
					}
					scroll.ScrollToBottom()
				})
			}
			fyne.Do(func() {
				entry.Enable()
				entry.SetText("")
				w.Canvas().Focus(entry)
			})
		}(text)
	}

	// Build Chat and Agents tabs
	chatTab := container.NewTabItem("Chat",
		container.NewBorder(nil, entry, nil, nil, scroll),
	)
	agentsTabContent, reloadAgents := makeAgentsTab(session.Router, updateStatusBar, w)
	agentsTab := container.NewTabItem("Agents", agentsTabContent)
	tabs := container.NewAppTabs(chatTab, agentsTab)
	w.SetContent(container.NewBorder(nil, statusBar, nil, nil, tabs))

	evChan := gohook.Start()
	defer gohook.End()
	gohook.Register(gohook.KeyDown, []string{"space", "command", "shift"}, func(e gohook.Event) {
		fyne.Do(func() {
			w.Show()
			w.RequestFocus()
			w.Canvas().Focus(entry)
		})
	})
	go gohook.Process(evChan)

	w.ShowAndRun()
}

func addMessage(historyText *widget.Entry, sender, message string) {
	historyText.SetText(historyText.Text + fmt.Sprintf("%s: %s\n", sender, message))
}

// makeAgentsTab creates the Agents management UI, listing all agents
// with controls to enable/disable, solo/unsolo, and (if supported) launch tools.
// It returns the UI component and a reload function to refresh the list.
func makeAgentsTab(router *router.AgentRouter, updateStatus func(), parent fyne.Window) (fyne.CanvasObject, func()) {
	list := container.NewVBox()
	scroll := container.NewVScroll(list)
	scroll.SetMinSize(fyne.NewSize(600, 350))

	var reload func()
	reload = func() {
		list.RemoveAll()
		// List all agent registrations
		for _, reg := range router.ListAllAgents() {
			name := reg.Agent.GetName()
			desc := reg.Agent.GetDescription()
			tags := strings.Join(reg.Tags, ", ")

			// Enabled checkbox (set before hooking change to avoid recursion)
			enabledCheck := widget.NewCheck("Enabled", nil)
			enabledCheck.SetChecked(reg.Enabled)
			enabledCheck.OnChanged = func(on bool) {
				router.EnableAgent(name, on)
				reload()
			}

			// Solo button to isolate this agent
			soloBtn := widget.NewButton("Solo", func() {
				router.SoloAgent(name)
				reload()
			})

			// Agent info and tags
			aLabel := widget.NewLabel(fmt.Sprintf("%s (P%d): %s", name, reg.Priority, desc))
			tLabel := widget.NewLabel(fmt.Sprintf("Tags: %s", tags))
			ctrls := container.NewHBox(enabledCheck, soloBtn)

			// Vertical box for each agent entry
			entry := container.NewVBox(aLabel, tLabel, ctrls)

			// If this agent provides tools (e.g. Reaper scripts), list them
			if tp, ok := reg.Agent.(agents.ToolProvider); ok {
				tools := tp.Tools()
				if len(tools) > 0 {
					var toolNames []string
					for _, tool := range tools {
						toolNames = append(toolNames, tool.Name)
					}
					selectEntry := widget.NewSelectEntry(toolNames)
					agent := reg.Agent
					agentName := name
					// normalize removes underscores and dots, lowercases for fuzzy matching
					normalize := func(s string) string {
						s = strings.ToLower(s)
						s = strings.ReplaceAll(s, "_", "")
						s = strings.ReplaceAll(s, ".", "")
						return s
					}
					launchBtn := widget.NewButton("Launch", func() {
						input := strings.TrimSpace(selectEntry.Text)
						if input == "" {
							dialog.ShowInformation("Select Tool", "Please select a tool to launch", parent)
							return
						}
						var chosen string
						inputNorm := normalize(input)
						for _, tn := range toolNames {
							if strings.EqualFold(tn, input) || strings.Contains(normalize(tn), inputNorm) {
								chosen = tn
								break
							}
						}
						if chosen == "" {
							dialog.ShowError(fmt.Errorf("tool not found: %s", input), parent)
							return
						}
						resp, err := agent.Handle(chosen, nil, "")
						if err != nil {
							dialog.ShowError(err, parent)
						} else {
							dialog.ShowInformation(fmt.Sprintf("%s Tool", agentName), resp, parent)
						}
					})
					entry.Add(widget.NewLabel("Tools:"))
					entry.Add(container.NewHBox(selectEntry, launchBtn))
				}
			}

			list.Add(entry)
		}
		list.Refresh()
		updateStatus()
	}

	// Unsolo All button
	unsoloBtn := widget.NewButton("Unsolo All", func() {
		router.UnsoloAgents()
		reload()
	})

	// Initial load
	reload()

	// Layout: Unsolo button at top, scrollable list below
	return container.NewVBox(unsoloBtn, scroll), reload
}

// handleCommand processes slash-commands similar to CLI
func handleCommand(input string, session *chat.ChatSession, historyText *widget.Entry, scroll *container.Scroll, parent fyne.Window, updateStatus func()) {
	parts := strings.Fields(input)
	cmd := parts[0]
	args := parts[1:]
	switch cmd {
	case "/clear":
		historyText.SetText("")
		scroll.ScrollToBottom()
		return
	case "/help", "/h":
		var lines []string
		for _, c := range chat.SlashCommands {
			lines = append(lines, fmt.Sprintf("%s: %s", c.Command, c.Description))
		}
		addMessage(historyText, "System", strings.Join(lines, "\n"))
		scroll.ScrollToBottom()
	case "/model":
		var options []string
		for _, m := range models.AvailableModels {
			options = append(options, fmt.Sprintf("%s - %s", m.Name, m.DisplayName))
		}
		sel := widget.NewSelect(options, func(selected string) {
			if selected == "" {
				return
			}
			name := strings.Fields(selected)[0]
			session.SelectedModel = name
			session.ModelDisplayName = models.GetModelDisplayName(name)
			addMessage(historyText, "System", fmt.Sprintf("Switched to model: %s", session.ModelDisplayName))
			scroll.ScrollToBottom()
		})
		dialog.NewCustom("Select Model", "Cancel", sel, parent).Show()
	case "/agents":
		var lines []string
		for _, ag := range session.Router.ListAgents() {
			lines = append(lines, fmt.Sprintf("%s - %s", ag.GetName(), ag.GetDescription()))
		}
		addMessage(historyText, "System", strings.Join(lines, "\n"))
		scroll.ScrollToBottom()
	case "/tools":
		var toolLines []string
		toolLines = append(toolLines, "🛠️ Available tools:")
		for _, ag := range session.Router.ListAgents() {
			if tp, ok := ag.(agents.ToolProvider); ok {
				for _, tool := range tp.Tools() {
					toolLines = append(toolLines, fmt.Sprintf(" - %s: %s (%s agent)", tool.Name, tool.Description, ag.GetName()))
				}
			}
		}
		if len(toolLines) == 1 {
			addMessage(historyText, "System", "No tools available")
		} else {
			addMessage(historyText, "System", strings.Join(toolLines, "\n"))
		}
		scroll.ScrollToBottom()
	case "/status":
		addMessage(historyText, "System", session.Router.GetAgentStatus())
		scroll.ScrollToBottom()
	case "/enable":
		if len(args) < 1 {
			addMessage(historyText, "Error", "Usage: /enable <agent>")
		} else if session.Router.EnableAgent(args[0], true) {
			addMessage(historyText, "System", fmt.Sprintf("%s Agent enabled", args[0]))
		} else {
			addMessage(historyText, "Error", fmt.Sprintf("Agent '%s' not found", args[0]))
		}
		scroll.ScrollToBottom()
	case "/disable":
		if len(args) < 1 {
			addMessage(historyText, "Error", "Usage: /disable <agent>")
		} else if session.Router.EnableAgent(args[0], false) {
			addMessage(historyText, "System", fmt.Sprintf("%s Agent disabled", args[0]))
		} else {
			addMessage(historyText, "Error", fmt.Sprintf("Agent '%s' not found", args[0]))
		}
		scroll.ScrollToBottom()
	case "/solo":
		if len(args) < 1 {
			addMessage(historyText, "Error", "Usage: /solo <agent>")
		} else if session.Router.SoloAgent(args[0]) {
			addMessage(historyText, "System", fmt.Sprintf("Solo mode: Only %s Agent is enabled", args[0]))
			addMessage(historyText, "System", "Use '/unsolo' to re-enable all agents")
			// Print available tools for the soloed agent, if any
			for _, ag := range session.Router.ListAgents() {
				if tp, ok := ag.(agents.ToolProvider); ok {
					tools := tp.Tools()
					if len(tools) > 0 {
						var lines []string
						lines = append(lines, "Available tools:")
						for _, tool := range tools {
							lines = append(lines, fmt.Sprintf(" - %s: %s", tool.Name, tool.Description))
						}
						addMessage(historyText, "System", strings.Join(lines, "\n"))
					}
				}
			}
		} else {
			addMessage(historyText, "Error", fmt.Sprintf("Agent '%s' not found", args[0]))
		}
		scroll.ScrollToBottom()
	case "/unsolo":
		session.Router.UnsoloAgents()
		addMessage(historyText, "System", "All agents re-enabled")
		scroll.ScrollToBottom()
	case "/tag":
		if len(args) < 1 {
			addMessage(historyText, "Error", "Usage: /tag <tag>")
		} else {
			agents := session.Router.GetAgentsByTag(args[0])
			if len(agents) == 0 {
				addMessage(historyText, "System", fmt.Sprintf("No agents found with tag '%s'", args[0]))
			} else {
				var lines []string
				for _, ag := range agents {
					lines = append(lines, fmt.Sprintf("%s - %s", ag.GetName(), ag.GetDescription()))
				}
				addMessage(historyText, "System", strings.Join(lines, "\n"))
			}
		}
		scroll.ScrollToBottom()
	case "/config":
		cfg := session.Factory.GetConfig()
		var lines []string
		lines = append(lines, fmt.Sprintf("Example Agents Enabled: %v", cfg.EnableExampleAgents))
		lines = append(lines, fmt.Sprintf("Code Review Agent: %v", cfg.EnableCodeReview))
		lines = append(lines, fmt.Sprintf("Data Analysis Agent: %v", cfg.EnableDataAnalysis))
		if cfg.WeatherAPIKey != "" {
			lines = append(lines, "Weather API Key: Configured")
		} else {
			lines = append(lines, "Weather API Key: Not set")
		}
		for agent, pr := range cfg.CustomAgentPriority {
			lines = append(lines, fmt.Sprintf("%s priority: %d", agent, pr))
		}
		addMessage(historyText, "System", strings.Join(lines, "\n"))
		scroll.ScrollToBottom()
	case "/store":
		if session.LastPrompt != "" && session.LastResponse != "" {
			if err := convstorage.StoreOpenAIResponse(session.LastPrompt, session.LastResponse, session.SelectedModel); err != nil {
				addMessage(historyText, "Error", fmt.Sprintf("Error storing response: %v", err))
			} else {
				addMessage(historyText, "System", "Response saved to file")
			}
		} else {
			addMessage(historyText, "Error", "No previous conversation to store")
		}
		scroll.ScrollToBottom()
	case "/load":
		fileDialog := dialog.NewFileOpen(func(rc fyne.URIReadCloser, err error) {
			if err != nil || rc == nil {
				return
			}
			defer rc.Close()
			path := rc.URI().Path()
			modelName, prompt, resp, err := convstorage.LoadConversation(path)
			if err != nil {
				addMessage(historyText, "Error", fmt.Sprintf("Error loading conversation: %v", err))
			} else {
				addMessage(historyText, "System", fmt.Sprintf("Loaded conversation from %s", path))
				addMessage(historyText, "System", fmt.Sprintf("Model: %s\nUser: %s\nAssistant: %s", modelName, prompt, resp))
				session.LastPrompt = prompt
				session.LastResponse = resp
				session.SelectedModel = modelName
				session.ModelDisplayName = models.GetModelDisplayName(modelName)
				addMessage(historyText, "System", fmt.Sprintf("Switched to model: %s", session.ModelDisplayName))
			}
			scroll.ScrollToBottom()
		}, parent)
		fileDialog.SetFilter(fstorage.NewExtensionFileFilter([]string{".txt"}))
		fileDialog.Show()
	case "/list":
		files, err := convstorage.ListConversationFiles()
		if err != nil {
			addMessage(historyText, "Error", fmt.Sprintf("Error listing conversations: %v", err))
		} else if len(files) == 0 {
			addMessage(historyText, "System", "No saved conversations found")
		} else {
			var lines []string
			for _, f := range files {
				lines = append(lines, f.DisplayName)
			}
			addMessage(historyText, "System", strings.Join(lines, "\n"))
		}
		scroll.ScrollToBottom()
	default:
		addMessage(historyText, "Error", fmt.Sprintf("Unknown command: %s", cmd))
		scroll.ScrollToBottom()
	}
	updateStatus()
}
