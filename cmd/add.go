package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//ElemType is the type of element to add
var ElemType string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please provide element name")
		}
		for k := range templates {
			if k == ElemType {
				return nil
			}
		}

		return fmt.Errorf("invalid type specified; must be in the keys of %v", templates)
	},
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		var found os.FileInfo = nil
		for _, file := range files {
			if file.Name() == "src" && file.IsDir() {
				found = file
				break
			}
		}
		if found == nil {
			log.Fatal("must be run in root dir of a react project")
		} else {
			createFiles(args[0])
		}
	},
}

func init() {
	addCmd.Flags().StringVarP(&ElemType, "type", "t", "component", "type of element to add")
	rootCmd.AddCommand(addCmd)
}

func createFiles(elementName string) {
	if toApply, ok := templates[ElemType]; ok {
		realName := elementName
		if toApply.capitalize {
			realName = strings.Title(elementName)
		}
		dir := fmt.Sprintf(toApply.rootDir, realName)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal("unable to create directory")
		}
		for _, t := range toApply.templates {
			realFileName := t.filename
			if strings.Count(t.filename, "%")-2*strings.Count(t.filename, "%%") > 0 {
				realFileName = fmt.Sprintf(t.filename, realName)
			}
			realContent := t.content
			if strings.Count(t.content, "%")-2*strings.Count(t.content, "%%") > 0 {
				realContent = fmt.Sprintf(t.content, realName)
			}
			err := ioutil.WriteFile(dir+"/"+realFileName, []byte(realContent), 0644)
			if err != nil {
				log.Printf("An error occured while creating file %v (%v)\n", realFileName, err)
			}
		}
	} else {

	}
}

var templates = map[string]template{
	"component": componentTemplate,
	"c":         componentTemplate,
	"reducer":   reducerTemplate,
	"r":         reducerTemplate,
	"action":    actionTemplate,
	"a":         actionTemplate,
}

type template struct {
	templates  []fileTemplate
	rootDir    string
	capitalize bool
}

type fileTemplate struct {
	filename string
	content  string
}

var componentTemplate = template{
	rootDir:    "src/components/%[1]s",
	capitalize: true,
	templates: []fileTemplate{
		{
			filename: "index.js",
			content: `import %[1]sContainer from './%[1]sContainer'
export default %[1]sContainer`,
		}, {
			filename: "%[1]s.js",
			content: `//@flow
import React, { Component } from 'react';

import './style.css'

type Props = {}
type State = {}
class %[1]s extends Component<Props, State> {
	render(){
		return (
			<div className="custom"/>
		)
	}
}
export default %[1]s;`,
		}, {
			filename: "%[1]sContainer.js",
			content: `// @flow
import { connect } from 'react-redux'
import %[1]s from './%[1]s'

const mapStateToProps = (state) => {
	return {}
}

const mapDispatchToProps = (dispatch) => {
	return {}
}

export default connect(mapStateToProps, mapDispatchToProps)(%[1]s)`,
		}, {
			filename: "style.css",
			content: `
.custom{
	background-color: lime;
}
`,
		},
	},
}

var reducerTemplate = template{
	rootDir:    "src/redux/reducers/%[1]s",
	capitalize: false,
	templates: []fileTemplate{
		{
			filename: "index.js",
			content: `import %[1]s from './%[1]s'
export default %[1]s
export * from './%[1]s'`,
		}, {
			filename: "%[1]s.js",
			content: `//@flow
export const MY_ACTION = 'MY_ACTION'
type stateType = {}
type actionType = { type: string }
const initState = {}
const %[1]s = (state: stateType = initState, action: actionType = { type: MY_ACTION }) => {
	switch (action.type) {
		case MY_ACTION:
			return { ...state, someExtraProp: "some data" }
		default:
			return state
	}
}
export default %[1]s`,
		},
	},
}

var actionTemplate = template{
	rootDir:    "src/redux/actions/%[1]s",
	capitalize: false,
	templates: []fileTemplate{
		{
			filename: "index.js",
			content:  `export * from './%[1]s'`,
		}, {
			filename: "%[1]s.js",
			content: `//@flow
import { STATE_CHANGE } from "../../reducers";


export const myAction = () => {
	return {
		type: STATE_CHANGE,
	}
}`,
		},
	},
}
