package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PluginConfig struct {
	PluginName        string
	PluginDescription string
	PluginVersion     string
	PluginFileName    string
	AuthorName        string
	AuthorURI         string
	TextDomain        string
	DomainPath        string
	MainClassName     string
	MainFunctionName  string
	Namespace         string
	PluginPrefix      string
	ConstantPrefix    string
	OriginalName      string
}

func CollectPluginConfig(pluginName string) (*PluginConfig, error) {
	reader := bufio.NewReader(os.Stdin)

	// Convert plugin name to various formats for technical identifiers
	className := strings.Title(strings.ReplaceAll(pluginName, "-", ""))
	functionName := strings.ReplaceAll(pluginName, "-", "_")
	prefix := strings.ToLower(strings.ReplaceAll(pluginName, "-", ""))
	constantPrefix := strings.ToUpper(strings.ReplaceAll(pluginName, "-", ""))

	config := &PluginConfig{
		PluginFileName:   pluginName + ".php",
		TextDomain:       pluginName,
		DomainPath:       "/languages",
		MainClassName:    className,
		MainFunctionName: functionName + "_init",
		Namespace:        className,
		PluginPrefix:     prefix,
		ConstantPrefix:   constantPrefix,
		OriginalName:     pluginName,
	}

	fmt.Print("Plugin Name [" + strings.Title(strings.ReplaceAll(pluginName, "-", " ")) + "]: ")
	headerName, _ := reader.ReadString('\n')
	headerName = strings.TrimSpace(headerName)
	if headerName == "" {
		headerName = strings.Title(strings.ReplaceAll(pluginName, "-", " "))
	}
	config.PluginName = headerName

	fmt.Print("Plugin Description: ")
	description, _ := reader.ReadString('\n')
	config.PluginDescription = strings.TrimSpace(description)

	fmt.Print("Plugin Version [1.0.0]: ")
	version, _ := reader.ReadString('\n')
	version = strings.TrimSpace(version)
	if version == "" {
		version = "1.0.0"
	}
	config.PluginVersion = version

	fmt.Print("Plugin File Name [" + config.PluginFileName + "]: ")
	fileName, _ := reader.ReadString('\n')
	fileName = strings.TrimSpace(fileName)
	if fileName != "" {
		config.PluginFileName = fileName
	}

	fmt.Print("Author Name: ")
	author, _ := reader.ReadString('\n')
	config.AuthorName = strings.TrimSpace(author)

	fmt.Print("Author URI: ")
	uri, _ := reader.ReadString('\n')
	config.AuthorURI = strings.TrimSpace(uri)

	fmt.Print("Text Domain [" + config.TextDomain + "]: ")
	textDomain, _ := reader.ReadString('\n')
	textDomain = strings.TrimSpace(textDomain)
	if textDomain != "" {
		config.TextDomain = textDomain
	}

	fmt.Print("Domain Path [" + config.DomainPath + "]: ")
	domainPath, _ := reader.ReadString('\n')
	domainPath = strings.TrimSpace(domainPath)
	if domainPath != "" {
		config.DomainPath = domainPath
	}

	fmt.Print("Main Class Name [" + config.MainClassName + "]: ")
	mainClass, _ := reader.ReadString('\n')
	mainClass = strings.TrimSpace(mainClass)
	if mainClass != "" {
		config.MainClassName = mainClass
	}

	fmt.Print("Main Function Name [" + config.MainFunctionName + "]: ")
	mainFunction, _ := reader.ReadString('\n')
	mainFunction = strings.TrimSpace(mainFunction)
	if mainFunction != "" {
		config.MainFunctionName = mainFunction
	}

	fmt.Print("Namespace [" + config.Namespace + "]: ")
	namespace, _ := reader.ReadString('\n')
	namespace = strings.TrimSpace(namespace)
	if namespace != "" {
		config.Namespace = namespace
	}

	fmt.Print("Plugin Prefix [" + config.PluginPrefix + "]: ")
	pluginPrefix, _ := reader.ReadString('\n')
	pluginPrefix = strings.TrimSpace(pluginPrefix)
	if pluginPrefix != "" {
		config.PluginPrefix = pluginPrefix
	}

	fmt.Print("Constant Prefix [" + config.ConstantPrefix + "]: ")
	constantPrefixInput, _ := reader.ReadString('\n')
	constantPrefixInput = strings.TrimSpace(constantPrefixInput)
	if constantPrefixInput != "" {
		config.ConstantPrefix = constantPrefixInput
	}

	return config, nil
}
