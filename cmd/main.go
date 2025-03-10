package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/icemint0828/imgedit"
	"github.com/icemint0828/imgedit/internal/app"
)

const (
	Version = "1.2.0"
)

func main() {
	flag.Bool(app.OptionVertical.Name, false, app.OptionVertical.Usage)
	flag.Uint(app.OptionWidth.Name, 0, app.OptionWidth.Usage)
	flag.Uint(app.OptionHeight.Name, 0, app.OptionHeight.Usage)
	flag.Float64(app.OptionRatio.Name, 0, app.OptionRatio.Usage)
	flag.Uint(app.OptionLeft.Name, 0, app.OptionLeft.Usage)
	flag.Uint(app.OptionTop.Name, 0, app.OptionTop.Usage)
	flag.CommandLine.Usage = usage
	permuteArgs(os.Args[1:])
	flag.Parse()

	// validation for flag and args
	args := flag.Args()

	if len(args) != 2 {
		exitOnError(errors.New("argument is missing"))
	}
	subCommandName, imagePath := args[0], args[1]
	subCommand := app.SupportedSubCommands.FindSubCommand(subCommandName)
	if subCommand == nil {
		exitOnError(errors.New(fmt.Sprintf("%s is not supported for subcommand", subCommandName)))
	}

	if !subCommand.ValidOption() {
		exitOnError(errors.New(fmt.Sprintf("%s is not valid for option", subCommandName)))
	}

	if !exists(imagePath) {
		exitOnError(errors.New(fmt.Sprintf("file does not exist : %s", imagePath)))
	}

	// run application
	err := app.NewApp(subCommand, imagePath).Run()
	if err != nil {
		exitOnError(err)
	}
}

func exitOnError(err error) {
	fmt.Println(err.Error())
	fmt.Println("see: imgedit -help")
	os.Exit(1)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func usage() {
	commandName := path.Base(os.Args[0])
	fmt.Printf("%s version %s\n\n", commandName, Version)
	fmt.Printf("Try running %s like:\n", commandName)
	fmt.Printf("%s <sub command> <image path> -<option> | for example:\n\n", commandName)
	fmt.Printf("%s reverse test.png -vertical\n", commandName)
	fmt.Printf("%s resize test.png -width 500 -height 500\n\n", commandName)
	fmt.Printf("[sub command]\n")
	for _, subCommand := range app.SupportedSubCommands {
		fmt.Printf("\n  %s : %s\n", subCommand.Name, subCommand.Usage)
		if len(subCommand.RequiredOptions) > 0 {
			fmt.Printf("    (required options)\n")
			for _, option := range subCommand.RequiredOptions {
				fmt.Printf("      -%s : %s\n", option.Name, option.Usage)
			}
		}
		if len(subCommand.OptionalOptions) > 0 {
			fmt.Printf("    (optional options)\n")
			for _, option := range subCommand.OptionalOptions {
				fmt.Printf("      -%s : %s\n", option.Name, option.Usage)
			}
		}
	}
	fmt.Printf("\n[supported extensions]\n")
	var supportedExtensions []string
	for _, e := range imgedit.SupportedExtensions {
		supportedExtensions = append(supportedExtensions, string(e))
	}
	fmt.Printf("    %s\n", strings.Join(supportedExtensions, "/"))
}

func permuteArgs(args []string) {
	var flagArgs []string
	var nonFlagArgs []string

	for i := 0; i < len(args); i++ {
		v := args[i]
		if v[0] == '-' {
			optionName := v[1:]
			switch optionName {
			case app.OptionHeight.Name, app.OptionWidth.Name, app.OptionRatio.Name, app.OptionLeft.Name, app.OptionTop.Name:
				/* out of index */
				if len(args) <= i+1 {
					exitOnError(errors.New(fmt.Sprintf("argument is missing for %s", v)))
				}
				/* the next flag has come */
				optionVal := args[i+1]
				if optionVal[0] == '-' {
					exitOnError(errors.New(fmt.Sprintf("argument is missing for %s", v)))
				}
				flagArgs = append(flagArgs, args[i:i+2]...)
				i++
			default:
				flagArgs = append(flagArgs, args[i])
			}
		} else {
			nonFlagArgs = append(nonFlagArgs, args[i])
		}
	}
	permutedArgs := append(flagArgs, nonFlagArgs...)

	/* replace args */
	for i := 0; i < len(args); i++ {
		args[i] = permutedArgs[i]
	}
}
