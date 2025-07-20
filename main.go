package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	scope := flag.String("scope", "machine", "The scope to read variables from. Can be 'machine' or 'user'.")
	appendVar := flag.String("append", "", "The variable to append to.")
	removeVar := flag.String("remove", "", "The variable to remove from.")
	value := flag.String("value", "", "The value to append or remove.")
	separator := flag.String("separator", ";", "The separator to use.")
	flag.Parse()

	var key registry.Key
	var path string

	switch *scope {
	case "machine":
		key = registry.LOCAL_MACHINE
		path = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
	case "user":
		key = registry.CURRENT_USER
		path = `Environment`
	default:
		log.Fatalf("Invalid scope: %s. Please use 'machine' or 'user'.", *scope)
	}

	access := uint32(registry.QUERY_VALUE)
	if *appendVar != "" || *removeVar != "" {
		access = registry.ALL_ACCESS
	}

	k, err := registry.OpenKey(key, path, access)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	if *appendVar != "" {
		if *value == "" {
			log.Fatal("The -value flag must be provided when using -append.")
		}

		oldVal, _, err := k.GetStringValue(*appendVar)
		if err != nil && err != registry.ErrNotExist {
			log.Fatalf("Failed to read existing value for %s: %v", *appendVar, err)
		}

		newVal := oldVal
		if oldVal != "" {
			newVal += *separator
		}
		newVal += *value

		if err := k.SetStringValue(*appendVar, newVal); err != nil {
			log.Fatalf("Failed to set new value for %s: %v", *appendVar, err)
		}

		fmt.Printf("Successfully appended to %s\n", *appendVar)
		return
	}

	if *removeVar != "" {
		if *value == "" {
			log.Fatal("The -value flag must be provided when using -remove.")
		}

		oldVal, _, err := k.GetStringValue(*removeVar)
		if err != nil {
			log.Fatalf("Failed to read existing value for %s: %v", *removeVar, err)
		}

		parts := strings.Split(oldVal, *separator)
		var newParts []string
		for _, part := range parts {
			if part != *value {
				newParts = append(newParts, part)
			}
		}
		newVal := strings.Join(newParts, *separator)

		if err := k.SetStringValue(*removeVar, newVal); err != nil {
			log.Fatalf("Failed to set new value for %s: %v", *removeVar, err)
		}

		fmt.Printf("Successfully removed from %s\n", *removeVar)
		return
	}

	names, err := k.ReadValueNames(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range names {
		val, _, err := k.GetStringValue(name)
		if err != nil {
			log.Printf("failed to get string for %s: %v", name, err)
			continue
		}
		fmt.Printf("%s=%s\n", name, val)
	}
}
