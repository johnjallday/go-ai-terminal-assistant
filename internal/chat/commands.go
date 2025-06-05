package chat

// CommandInfo represents a slash-command and its description.
type CommandInfo struct {
	Command     string
	Description string
}

// SlashCommands lists all available slash commands and their descriptions.
var SlashCommands = []CommandInfo{
	{"/help, /h", "Show this help message"},
	{"/model", "Change the model interactively"},
	{"/agents", "List enabled agents"},
	{"/tools", "List all available tools"},
	{"/status", "Show agent status"},
	{"/enable <agent>", "Enable an agent"},
	{"/disable <agent>", "Disable an agent"},
	{"/solo <agent>", "Solo an agent"},
	{"/unsolo", "Re-enable all agents"},
	{"/tag <tag>", "List agents by tag"},
	{"/config", "Show agent configuration"},
	{"/store", "Store last conversation to file"},
	{"/load", "Load a conversation from file"},
	{"/list", "List saved conversations"},
}
