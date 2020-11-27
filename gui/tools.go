package gui

import (
	"log"
	"os"

	"github.com/lxn/walk"
)

func GenMenuItem(caption string, action walk.EventHandler) *walk.Action {
	act := walk.NewAction()
	_ = act.SetText(caption)
	act.Triggered().Attach(action)
	return act
}
func checkDir(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}
func GenCheckable(caption string, w func(bool), initial bool) *walk.Action {
	act := walk.NewAction()
	_ = act.SetText(caption)
	_ = act.SetCheckable(true)
	err = act.SetChecked(initial)
	if err != nil {
		log.Println(err.Error())
	}
	act.Triggered().Attach(func() {
		w(act.Checked())
	})
	return act
}
