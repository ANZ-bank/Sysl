package syslutil

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

//ReadCMDFlags reads command line flags from file specified by configPath.
func ReadCMDFlags(configPath string) ([]string, error) {
	data, ferr := ioutil.ReadFile(configPath)
	if ferr != nil {
		return nil, ferr
	}

	re := regexp.MustCompile(`(-+[^=\s\n]+)?((=?"[^=\n]+")|(=?[^=\s\n]+))?\s*`)
	flags := re.FindAllString(strings.TrimSpace(string(data)), -1)

	for i, flag := range flags {
		newFlag := strings.TrimSpace(flag)
		start := strings.Index(newFlag, "=")
		if newFlag[start+1] == '"' && newFlag[len(newFlag)-1] == '"' {
			// flag value matches pattern `=?"[^=\n"]+"`, it looks like --app-name="Test App"
			flags[i] = flag[0:start+1] + flag[start+2:len(newFlag)-1]
		} else {
			flags[i] = newFlag
		}
	}

	return flags, nil
}

//PopulateCMDFlagsFromFile reads command line flags from file specified by cmdArgs and this flag starts with @,
//like `sysl codegen @file`.
func PopulateCMDFlagsFromFile(cmdArgs []string) ([]string, error) {
	if len(cmdArgs) < 3 {
		return nil, fmt.Errorf("command arguments are not enough")
	}

	var fileFlag string
	for _, flag := range cmdArgs[2:] {
		if strings.HasPrefix(flag, "@") {
			fileFlag = flag
			break
		}
	}

	flags, err := ReadCMDFlags(strings.Replace(fileFlag, "@", "", 1))
	if err != nil {
		return nil, err
	}

	return append(cmdArgs[0:2], flags...), err
}
