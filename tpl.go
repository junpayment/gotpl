package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
	"encoding/base64"
)

// Reads a YAML document from the values_in stream, uses it as values
// for the tpl_files templates and writes the executed templates to
// the out stream.
func ExecuteTemplates(values_in io.Reader, out io.Writer, tpl_files ...string) error {
	buf := bytes.NewBuffer(nil)
	_, err := io.Copy(buf, values_in)
	if err != nil {
		return fmt.Errorf("Failed to read standard input: %v", err)
	}

	for _, tpl_file := range tpl_files {
		tpl, err := template.ParseFiles(tpl_file)
		if err != nil {
			return fmt.Errorf("Error parsing template(s): %v", err)
		}

		var values map[string]interface{}
		err = yaml.Unmarshal(buf.Bytes(), &values)
		if err != nil {
			return fmt.Errorf("Failed to parse standard input: %v", err)
		}

		err = convert2base64(&values)
		if err != nil {
			return fmt.Errorf("Failed to parse standard input: %v", err)
		}

		err = tpl.Execute(out, values)
		if err != nil {
			return fmt.Errorf("Failed to parse standard input: %v", err)
		}
	}

	return nil
}

// convert map to base64
func convert2base64(values *map[string]interface{}) (error) {

	if _, ok := (*values)["base64"]; ok {
		for k, v := range (*values)["base64"].(map[interface{}]interface{}) {
			(*values)["base64"].(map[interface{}]interface{})[k] =
				base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
	}

	return nil
}

func main() {
	err := ExecuteTemplates(os.Stdin, os.Stdout, os.Args[1:]...)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
