package inputreader

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// InputReader is something that does stuff
type InputReader struct {
	Reader     io.Reader
	Writer     io.Writer
	buffReader *bufio.Reader
}

// InputOptions are the options for input
type InputOptions struct {
	Default string
}

// New returns and new InputReader
func New(reader io.Reader, writer io.Writer) *InputReader {
	return &InputReader{
		Reader:     reader,
		Writer:     writer,
		buffReader: bufio.NewReader(reader),
	}
}

func (i *InputReader) read() (string, error) {
	text, err := i.buffReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	// convert CRLF to LF
	text = strings.Trim(text, "\r\n")

	return text, nil
}

// Ask collects input from the user
func (i *InputReader) Ask(query string) (string, error) {
	fmt.Fprintf(i.Writer, "\n%s\n", query)

	return i.read()
}

// Select asks the user to select an option from the list
func (i *InputReader) Select(query string, list []string, opts *InputOptions) (string, error) {
	// Find default index which opts.Default indicates
	defaultIndex := -1
	defaultVal := opts.Default
	if defaultVal != "" {
		for i, item := range list {
			if item == defaultVal {
				defaultIndex = i
			}
		}

		// DefaultVal is set but doesn't exist in list
		if defaultIndex == -1 {
			// This error message is not for user
			// Should be found while development
			return "", fmt.Errorf("opt.Default is specified but item does not exist in list")
		}
	}

	// Construct the query & display it to user
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n\n", query))
	for i, item := range list {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}

	buf.WriteString("\n")
	fmt.Fprintf(i.Writer, buf.String())

	// resultStr and resultErr are return val of this function
	var resultStr string
	var resultErr error
	for {

		// Construct the asking line to input
		var buf bytes.Buffer
		buf.WriteString("Enter a number")

		// Add default val if provided
		if defaultIndex >= 0 {
			buf.WriteString(fmt.Sprintf(" (Default is %d)", defaultIndex+1))
		}

		buf.WriteString(": ")
		fmt.Fprintf(i.Writer, buf.String())

		// Read user input from reader.
		line, err := i.read()
		if err != nil {
			resultErr = err
			break
		}

		// line is empty but default is provided returns it
		if line == "" && defaultIndex >= 0 {
			resultStr = list[defaultIndex]
			break
		}

		if line == "" {
			fmt.Fprintf(i.Writer, "Input must not be empty. Answer by a number.\n\n")
			continue
		}

		// Convert user input string to int val
		n, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(i.Writer, "%q is not a valid input. Answer by a number.\n\n", line)
			continue
		}

		// Check answer is in range of list
		if n < 1 || len(list) < n {
			fmt.Fprintf(i.Writer, "%q is not a valid choice. Choose a number from 1 to %d.\n\n",
				line, len(list))
			continue
		}

		// Reach here means it gets ideal input.
		resultStr = list[n-1]
		break
	}

	// Insert the new line for next output
	fmt.Fprintf(i.Writer, "\n")

	return resultStr, resultErr
}
