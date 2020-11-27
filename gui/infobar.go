package gui

import (
	"log"
	"net"

	"github.com/lxn/walk"

	"webdav/config"
	"webdav/webdav"
)

var (
	mw                          *walk.MainWindow
	ni                          *walk.NotifyIcon
	err                         error
	icon                        *walk.Icon
	exit, reloadDrv, reloadConf *walk.Action

	listeners     = make([]net.Listener, 0)
	otherListener = make(map[string]net.Listener, 0)
)

func GuiRun() {
	icon, err = walk.Resources.Icon("MAIN_ICO")
	if err != nil {
		log.Fatal(err)
	}

	// mainWindow
	mw, err = walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}
	// notifyIcon
	ni, err = walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("WebDav Server [by Nolva]"); err != nil {
		log.Fatal(err)
	}
	_ = ni.SetVisible(true)
	// menus

	setMenu()
	// run
	mw.Run()
}

func setMenu() {
	_ = ni.ContextMenu().Actions().Clear()

	others := GenCheckable("其他服务器", func(checked bool) {
		if checked {
			for _, s := range config.GlobalConf.Servers {
				listener := webdav.StartServer(&s)
				listeners = append(listeners, listener)
			}
		} else {
			for _, l := range listeners {
				_ = l.Close()
			}
			listeners = []net.Listener{}
		}
	}, len(listeners) != 0)

	exit = GenMenuItem("退出", func() {
		_ = ni.Dispose()
		mw.Dispose()
		walk.App().Exit(0)
	})

	reloadDrv = GenMenuItem("重新查询盘符", setMenu)
	reloadConf = GenMenuItem("重新读取设置(test)", config.ReloadConfig)
	_ = reloadConf.SetEnabled(false)

	_ = ni.ContextMenu().Actions().Add(others)
	if len(config.GlobalConf.Servers) == 0 {
		_ = others.SetEnabled(false)
	}

	for _, it := range []string{
		"A", "B", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	} {
		dir := it + ":\\"
		if checkDir(dir) {
			func(it string) {
				_ = ni.ContextMenu().Actions().Add(GenCheckable(it,
					func(checked bool) {
						if checked {
							s := config.GlobalConf.Default
							s.Scope = dir
							otherListener[it] = webdav.StartServer(&s)
						} else {
							if otherListener[it] != nil {
								_ = otherListener[it].Close()
								delete(otherListener, it)
							}
						}
					}, otherListener[it] != nil))
			}(it)
		}
	}

	_ = ni.ContextMenu().Actions().Add(walk.NewSeparatorAction())
	_ = ni.ContextMenu().Actions().Add(reloadDrv)
	_ = ni.ContextMenu().Actions().Add(reloadConf)
	_ = ni.ContextMenu().Actions().Add(walk.NewSeparatorAction())
	_ = ni.ContextMenu().Actions().Add(exit)

	//TODO 追加单击事件，显示当前运行的server的配置
}
