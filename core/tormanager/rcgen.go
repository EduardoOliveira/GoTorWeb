package tormanager

import (
	"html/template"
	"os"
)

func (tm *TorManager) genRC() error {

	t, err := template.ParseFiles("./templates/torrc.tmpl")
	if err != nil {
		return err
	}

	f, err := os.Create("./go-torrc")
	if err != nil {
		return err
	}

	err = t.Execute(f, tm)
	if err != nil {
		return err
	}
	return nil
}
