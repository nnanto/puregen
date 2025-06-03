package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nnanto/puregen/helper"
	"github.com/spf13/cobra"
)

var outputFile string

var creatorCmd = &cobra.Command{
	Use:   "creator-mode",
	Short: "Interactive mode to create and edit YAML IDL files",
	Long:  `Creator mode provides an interactive interface to create and edit YAML IDL files that adhere to the SchemaYAML structure.`,
	Run: func(cmd *cobra.Command, args []string) {
		if outputFile == "" {
			fmt.Fprintf(os.Stderr, "Error: output file is required\n")
			os.Exit(1)
		}

		creator := helper.NewCreator(outputFile)
		if err := creator.LoadOrCreate(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading/creating file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Creator mode for file: %s\n", outputFile)
		if creator.FileExists() {
			fmt.Println("Loaded existing file")
		} else {
			fmt.Println("Created new file")
		}

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Println("\nWhat would you like to do?")
			fmt.Println("1. Create/Edit a message")
			fmt.Println("2. Create/Edit a service")
			fmt.Println("3. Show current schema")
			fmt.Println("4. Exit")
			fmt.Print("Enter your choice (1-4): ")

			choice, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				continue
			}

			choice = strings.TrimSpace(choice)
			switch choice {
			case "1":
				if err := handleMessageCreation(creator, reader); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			case "2":
				if err := handleServiceCreation(creator, reader); err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			case "3":
				creator.ShowSchema()
			case "4", "exit", "quit":
				fmt.Println("Goodbye!")
				return
			default:
				fmt.Println("Invalid choice. Please enter 1-4.")
			}
		}
	},
}

func handleMessageCreation(creator *helper.Creator, reader *bufio.Reader) error {
	fmt.Print("Enter message name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("message name cannot be empty")
	}

	fmt.Print("Enter message description (optional): ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	description = strings.TrimSpace(description)

	message := creator.CreateMessage(name, description)

	// Add fields
	for {
		fmt.Print("Add a field? (y/n): ")
		addField, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		addField = strings.TrimSpace(strings.ToLower(addField))

		if addField != "y" && addField != "yes" {
			break
		}

		fmt.Print("Field name: ")
		fieldName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fieldName = strings.TrimSpace(fieldName)

		if fieldName == "" {
			fmt.Println("Field name cannot be empty")
			continue
		}

		fmt.Print("Field type: ")
		fieldType, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fieldType = strings.TrimSpace(fieldType)

		if fieldType == "" {
			fmt.Println("Field type cannot be empty")
			continue
		}

		fmt.Print("Field description (optional): ")
		fieldDescription, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fieldDescription = strings.TrimSpace(fieldDescription)

		fmt.Print("Is required? (y/n): ")
		requiredStr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		requiredStr = strings.TrimSpace(strings.ToLower(requiredStr))
		required := requiredStr == "y" || requiredStr == "yes"

		fmt.Print("Is repeated? (y/n): ")
		repeatedStr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		repeatedStr = strings.TrimSpace(strings.ToLower(repeatedStr))
		repeated := repeatedStr == "y" || repeatedStr == "yes"

		field := creator.CreateField(fieldType, fieldDescription, required, repeated)
		message.Fields[fieldName] = field

		// Validate field type
		if !creator.IsValidType(fieldType) {
			fmt.Printf("Warning: Type '%s' is not a primitive type and is not defined as a message\n", fieldType)
		}
	}

	if err := creator.AddMessage(name, message); err != nil {
		return err
	}

	if err := creator.Save(); err != nil {
		return err
	}

	fmt.Printf("Message '%s' added successfully!\n", name)
	return nil
}

func handleServiceCreation(creator *helper.Creator, reader *bufio.Reader) error {
	fmt.Print("Enter service name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("service name cannot be empty")
	}

	fmt.Print("Enter service description (optional): ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	description = strings.TrimSpace(description)

	service := creator.CreateService(description)

	// Add methods
	for {
		fmt.Print("Add a method? (y/n): ")
		addMethod, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		addMethod = strings.TrimSpace(strings.ToLower(addMethod))

		if addMethod != "y" && addMethod != "yes" {
			break
		}

		fmt.Print("Method name: ")
		methodName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		methodName = strings.TrimSpace(methodName)

		if methodName == "" {
			fmt.Println("Method name cannot be empty")
			continue
		}

		fmt.Print("Method description (optional): ")
		methodDescription, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		methodDescription = strings.TrimSpace(methodDescription)

		fmt.Print("Input message type: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Input type cannot be empty")
			continue
		}

		fmt.Print("Output message type: ")
		output, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		output = strings.TrimSpace(output)

		if output == "" {
			fmt.Println("Output type cannot be empty")
			continue
		}

		fmt.Print("Is streaming? (y/n): ")
		streamingStr, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		streamingStr = strings.TrimSpace(strings.ToLower(streamingStr))
		streaming := streamingStr == "y" || streamingStr == "yes"

		method := creator.CreateMethod(methodName, methodDescription, input, output, streaming)
		service.Methods[methodName] = method

		// Validate input and output types
		if !creator.IsValidType(input) {
			fmt.Printf("Warning: Input type '%s' is not a primitive type and is not defined as a message\n", input)
		}
		if !creator.IsValidType(output) {
			fmt.Printf("Warning: Output type '%s' is not a primitive type and is not defined as a message\n", output)
		}
	}

	if err := creator.AddService(name, service); err != nil {
		return err
	}

	if err := creator.Save(); err != nil {
		return err
	}

	fmt.Printf("Service '%s' added successfully!\n", name)
	return nil
}
