package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// main is the entry point of the program.
func main() {
    // Check for correct number of arguments and if the -o flag is present.
    outputOption := ""
    if len(os.Args) < 4 {
        printUsage()
        os.Exit(1)
    }
    if len(os.Args) == 5 && strings.HasPrefix(os.Args[4], "-o=") {
        outputOption = strings.TrimPrefix(os.Args[4], "-o=")
        if outputOption != "json" && outputOption != "start" && outputOption != "end" {
            printUsage()
            os.Exit(1)
        }
    }

    formatID := os.Args[1] // First argument: format identifier.
    startTimeArg := os.Args[2] // Second argument: start time.
    endTimeArg := os.Args[3] // Third argument: end time.

    // Determine the format based on the format identifier.
    var format string
    switch formatID {
    case "1":
        format = "2006-01-02T15:04:05" // ISO 8601 format.
    case "2":
        format = "01-02-2006" // American format.
    case "3":
        format = "02-01-2006" // European format.
    case "4":
        format = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC 2822 format.
    case "5":
        // Unix Timestamp format - special handling required.
    case "6":
        // Custom strftime format - format is specified in startTimeArg.
        format = startTimeArg
        startTimeArg = endTimeArg // Shift argument positions for custom format.
    default:
        fmt.Println("Invalid format identifier")
        os.Exit(1)
    }

    // Parse start and end times.
    startTime, err := parseTimeArgument(startTimeArg)
    if err != nil {
        fmt.Printf("Error parsing start time: %s\n", err)
        os.Exit(1)
    }
    endTime, err := parseTimeArgument(endTimeArg)
    if err != nil {
        fmt.Printf("Error parsing end time: %s\n", err)
        os.Exit(1)
    }

    // Output the formatted date range based on the specified output option.
    if formatID == "5" {
        handleOutput(startTime.Unix(), endTime.Unix(), outputOption)
    } else {
        handleOutput(startTime.Format(format), endTime.Format(format), outputOption)
    }
}

// parseTimeArgument converts a time range argument into a time.Time.
func parseTimeArgument(arg string) (time.Time, error) {
    currentTime := time.Now()

    // Handle special keywords "now" and "today".
    if arg == "now" {
        return currentTime, nil
    } else if arg == "today" {
        // Return the start of the current day.
        return time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location()), nil
    } else if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "+") {
        // Handle relative time arguments.
        return parseRelativeTime(arg, currentTime)
    }

    // Return an error if the argument is not in a recognized format.
    return time.Time{}, fmt.Errorf("invalid time argument: %s", arg)
}

// parseRelativeTime handles relative time calculations.
func parseRelativeTime(arg string, referenceTime time.Time) (time.Time, error) {
    unit := arg[len(arg)-1:] // Extract the unit (m, h, d, w, M).
    number, err := strconv.Atoi(arg[1 : len(arg)-1]) // Extract the number part.
    if err != nil {
        return time.Time{}, fmt.Errorf("invalid time number: %s", arg)
    }
    if arg[0] == '-' {
        number = -number // Make the number negative for past dates.
    }

    // Calculate the time based on the unit and number.
    switch unit {
    case "m":
        return referenceTime.Add(time.Duration(number) * time.Minute), nil
    case "h":
        return referenceTime.Add(time.Duration(number) * time.Hour), nil
    case "d":
        return referenceTime.AddDate(0, 0, number), nil
    case "w":
        return referenceTime.AddDate(0, 0, 7*number), nil
    case "M":
        return referenceTime.AddDate(0, number, 0), nil
    default:
        return time.Time{}, fmt.Errorf("invalid time unit: %s", unit)
    }
}

// printUsage prints detailed usage instructions for the tool
func printUsage() {
	// Detailed instructions for how to use the tool, including format identifiers and examples
	fmt.Println("Usage: ./strftime <format_id> <start_time> <end_time>")
	fmt.Println("\nFormat Identifiers:")
	fmt.Println("  1: ISO 8601 (e.g., 2006-01-02T15:04:05)")
	fmt.Println("  2: American Format (e.g., 01-02-2006)")
	fmt.Println("  3: European Format (e.g., 02-01-2006)")
	fmt.Println("  4: RFC 2822 (e.g., Mon, 02 Jan 2006 15:04:05 -0700)")
	fmt.Println("  5: Unix Timestamp (seconds since Unix epoch)")
	fmt.Println("  6: Custom strftime format (specified as part of <start_time>)")

	fmt.Println("\nTime Arguments:")
	fmt.Println("  Start and end times can be specified in several ways:")
	fmt.Println("  - Relative times: Prefix with '+' or '-' followed by a number and a unit.")
	fmt.Println("    Units: 'm' (minutes), 'h' (hours), 'd' (days), 'w' (weeks), 'M' (months)")
	fmt.Println("    Examples: -1d (1 day ago), +1w (1 week in the future)")
	fmt.Println("  - Special keywords: 'now' (current moment) and 'today' (start of current day)")
	fmt.Println("  - For custom format (6), the format string should be the second argument.")

	fmt.Println("\nExamples:")
	fmt.Println("  ./strftime 1 -1d +1d     -> ISO 8601 format, from 1 day ago to 1 day in the future")
	fmt.Println("  ./strftime 3 today +1w   -> European format, from today to 1 week in the future")
	fmt.Println("  ./strftime 6 \"%Y/%m/%d %H:%M:%S\" -2h now -> Custom format, from 2 hours ago to now")

	fmt.Println("\nNote:")
	fmt.Println("  For the Unix Timestamp format (5), the time range will be output as two timestamps.")
}

// handleOutput prints the output based on the specified option.
func handleOutput(start interface{}, end interface{}, outputOption string) {
    switch outputOption {
    case "json":
        // Output in JSON format.
        output := map[string]interface{}{"start": start, "end": end}
        jsonData, err := json.Marshal(output)
        if err != nil {
            fmt.Printf("Error generating JSON output: %s\n", err)
            os.Exit(1)
        }
        fmt.Println(string(jsonData))
    case "start":
        // Output only the start date/time.
        fmt.Println(start)
    case "end":
        // Output only the end date/time.
        fmt.Println(end)
    default:
        // Default output: both start and end dates/times.
        fmt.Println("Start:", start, "\nEnd:", end)
    }
}