package ui

import (
	"docker-tray/internal/app/logger"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type PassWindow struct {
	logger.Logger
	window   *gtk.Window
	password string
	cancel   bool
}

func NewPassWindow() *PassWindow {
	passWindow := PassWindow{}

	// passWindow.init()
	return &passWindow
}

// func (w *PassWindow) init() {
// 	gtk.Init(nil)
// 	w.window = w.createPassWindow()
// 	w.window.ShowAll()
// 	w.window.SetKeepAbove(true)

// 	gtk.Main()
// }

func (w *PassWindow) GetPass() (string, bool) {

	gtk.Init(nil)
	w.window = w.createPassWindow()
	w.window.ShowAll()
	w.window.SetKeepAbove(true)

	gtk.Main()

	return w.password, w.cancel
}

func (w *PassWindow) createPassWindow() *gtk.Window {
	var err error

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		w.LogFatal("Unable to create window:", err)
	}

	win.SetTitle("Password")
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetDefaultSize(200, 100)
	win.SetBorderWidth(10)
	win.SetResizable(false)

	win.Connect("destroy", func() {
		win.Destroy()
	})

	passwordLabel, err := gtk.LabelNew("Password: ")
	if err != nil {
		w.LogFatal("Unable to create password label:", err)
	}
	passwordLabel.SetWidthChars(12)

	messageLabel, err := gtk.LabelNew("Wrong password")
	if err != nil {
		w.LogFatal("Unable to create message label:", err)
	}
	messageLabel.SetWidthChars(5)
	messageLabel.SetNoShowAll(true)

	passwordEntry, err := gtk.EntryNew()
	if err != nil {
		w.LogFatal("Unable to create entry:", err)
	}
	passwordEntry.SetVisibility(false)

	okButton, err := gtk.ButtonNewWithLabel("OK")
	if err != nil {
		w.LogFatal("Unable to create OK button:", err)
	}
	cancelButton, err := gtk.ButtonNewWithLabel("Cancel")
	if err != nil {
		w.LogFatal("Unable to create Cancel button:", err)
	}

	okActivateFunc := func() {
		w.LogInfo("Нажата кнопка: Ok")
		pass, err := passwordEntry.GetText()
		if err != nil {
			w.LogError("Ошибка при получении текста", err)
		} else {
			w.password = strings.TrimSpace(pass)
		}

		messageLabel.Show()
		passwordEntry.SetText("")
		w.Close()
	}

	okButton.Connect("clicked", okActivateFunc)

	passwordEntry.Connect("activate", okActivateFunc)

	cancelButton.Connect("clicked", func() {
		w.LogInfo("Нажата кнопка: Cancel")
		w.cancel = true
		w.Close()
	})

	hBox1, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		w.LogError("Ошибка при создании элемента", err)
	}
	hBox2, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		w.LogError("Ошибка при создании элемента", err)
	}
	vBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		w.LogError("Ошибка при создании элемента", err)
	}

	hBox1.PackStart(passwordLabel, true, false, 5)
	hBox1.PackStart(passwordEntry, true, false, 5)

	hBox2.PackStart(cancelButton, true, false, 5)
	hBox2.PackStart(okButton, true, false, 5)

	vBox.PackStart(hBox1, false, false, 5)
	vBox.PackStart(messageLabel, false, false, 0)
	vBox.PackStart(hBox2, false, false, 5)

	win.Add(vBox)

	return win
}

// func (w *PassWindow) Refresh() (string, bool) {
// 	return "", true
// }

func (w *PassWindow) Close() {
	w.window.Close()
	w.window.Destroy()
	gtk.MainQuit()
}
