package command

import "flag"
import log "github.com/sirupsen/logrus"

type Argument struct {
	name      string
	usage     string
	mandatory bool
}

type Command struct {
	name      string
	arguments []Argument
}

func (command *Command) ExtractAndValidateArguments() map[string]*string {
	values := make(map[string]*string)
	for _, arg := range command.arguments {
		values[arg.name] = flag.String(arg.name, "", arg.usage)
	}
	flag.Parse()
	isValid, validationErrorMessage := command.validateArguments(values)
	if isValid {
		return values
	} else {
		log.Error(validationErrorMessage)
		return nil
	}
}

func (command *Command) validateArguments(argValues map[string]*string) (bool, string) {
	for _, arg := range command.arguments {
		if arg.mandatory && *argValues[arg.name] == "" {
			return false, "Please provide value for " + arg.name
		}
	}
	return true, ""
}

type Action struct {
	command        Command
	argumentValues map[string]string
}

type CommandActionProcessor struct {
	commands []Command
}

func NewCommandActionProcessor() *CommandActionProcessor {
	createCommand := Command{
		name: "CREATE",
		arguments: []Argument{
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
			Argument{
				name:      "acousticLibraryID",
				usage:     "Acoustic Library ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
		},
	}

	updateCommand := Command{
		name: "UPDATE",
		arguments: []Argument{
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
			Argument{
				name:      "acousticLibraryID",
				usage:     "Acoustic Library ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
		},
	}

	readCommand := Command{
		name: "READ",
		arguments: []Argument{
			Argument{
				name:      "feedLocation",
				usage:     "File path of the feed",
				mandatory: true,
			},
			Argument{
				name:      "configLocation",
				usage:     "File path of the config",
				mandatory: true,
			},
			Argument{
				name:      "acousticLibraryID",
				usage:     "Acoustic Library ID",
				mandatory: true,
			},
			Argument{
				name:      "contentTypeID",
				usage:     "Content Type ID",
				mandatory: true,
			},
		},
	}

	return &CommandActionProcessor{
		commands: []Command{createCommand, updateCommand, readCommand},
	}
}
