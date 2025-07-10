package asdf

import (
	"fmt"
	"strings"

	"github.com/bitrise-io/bitrise/v2/log"
	"github.com/bitrise-io/toolprovider/provider"
)

func (a *AsdfToolProvider) InstallToolPlugin(toolName string) error {
	plugin := provider.GetToolPlugin(toolName)
	if plugin == nil || (plugin.PluginName == "" && plugin.PluginURL == "") {
		return nil
	}

	installed, err := a.pluginAlreadyInstalled(*plugin)
	if err != nil {
		log.Warnf("Failed to check if plugin is already installed: %v", err)
	}
	if installed {
		log.Debugf("Tool plugin %s is already installed, skipping installation.", toolName)
		return nil
	}

	pluginAddArgs := []string{"add"}
	if plugin.PluginName != "" {
		pluginAddArgs = append(pluginAddArgs, plugin.PluginName)
	}
	if plugin.PluginURL != "" {
		pluginAddArgs = append(pluginAddArgs, plugin.PluginURL)
	}

	_, err = a.ExecEnv.RunAsdfPlugin(pluginAddArgs...)
	if err != nil {
		return err
	}

	installed, err = a.pluginAlreadyInstalled(*plugin)
	if err != nil {
		return fmt.Errorf("failed to check if plugin was installed successfully: %w", err)
	}
	if !installed {
		return fmt.Errorf("failed to install tool plugin %s", toolName)
	}

	return nil
}

func (a *AsdfToolProvider) pluginAlreadyInstalled(plugin provider.ToolPlugin) (bool, error) {
	pluginListArgs := []string{"list", "--urls"}
	out, err := a.ExecEnv.RunAsdfPlugin(pluginListArgs...)
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if (plugin.PluginName == "" || strings.HasPrefix(line, plugin.PluginName)) &&
			(plugin.PluginURL == "" || strings.Contains(line, plugin.PluginURL)) {
			return true, nil
		}
	}

	return false, nil
}
