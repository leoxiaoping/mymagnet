// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"test/md4"
	"test/tracker"
	"time"
)

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	mw := new(MyMainWindow)
	//var openAction *walk.Action
	//var outTE, outTE2 *walk.TextEdit

	ic, err2 := walk.NewIconFromResource("3") //修改过的NewIconFromResource方法，感谢谷歌上的老外
	if err2 != nil {
		log.Fatal(err2)
	}

	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		MyIcon:   ic,
		Title:    "磁性链接生成器 -- link51.net",
		MinSize:  Size{400, 440},
		Size:     Size{400, 440},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				MinSize: Size{0, 50},
				Text:    "打开本地种子文件(*.torrent)……",
				OnClicked: func() {
					mw.openAction_Triggered()
				},
			},
			TextEdit{AssignTo: &mw.outTE, ReadOnly: true},
			PushButton{
				MinSize: Size{0, 30},
				Text:    "复制生成的磁性链接",
				OnClicked: func() {
					walk.Clipboard().SetText(mw.outTE.Text())
				},
			},
			TextEdit{AssignTo: &mw.outTE2, ReadOnly: true},
			PushButton{
				MinSize: Size{0, 30},
				Text:    "复制生成的电驴地址",
				OnClicked: func() {
					walk.Clipboard().SetText(mw.outTE2.Text())
				},
			},

			TextEdit{AssignTo: &mw.logView, ReadOnly: true},

			PushButton{
				MinSize: Size{0, 30},
				Text:    "关于软件",
				OnClicked: func() {
					mw.aboutAction_Triggered()
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}
	mw.logView.AppendText(fmt.Sprintf("%s : %s \n", time.Now().Format("15:04:05"), "就绪"))
	mw.MainWindow.Run()
}

type MyMainWindow struct {
	*walk.MainWindow

	outTE        *walk.TextEdit
	outTE2       *walk.TextEdit
	logView      *walk.TextEdit
	prevFilePath string
}

func (mw *MyMainWindow) openAction_Triggered() {
	if err := mw.openImage(); err != nil {
		log.Print(err)
	}
}

func (mw *MyMainWindow) openImage() error {

	dlg := new(walk.FileDialog)

	dlg.FilePath = mw.prevFilePath
	dlg.Filter = "Files (*.torrent)|*.torrent"
	dlg.Title = "选择文件……"

	if ok, err := dlg.ShowOpen(mw); err != nil {
		return err
	} else if !ok {
		return nil
	}

	mw.prevFilePath = dlg.FilePath

	mw.logView.AppendText(fmt.Sprintf("%s : 打开%s \n", time.Now().Format("15:04:05"), mw.prevFilePath))

	fmt.Printf("路径：%s", mw.prevFilePath)
	fileName := filepath.Base(mw.prevFilePath)

	infoHash := tracker.FindInfoHash(mw.prevFilePath)
	margicLink := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", infoHash, url.QueryEscape(fileName))
	mw.outTE.SetText(margicLink) //设置磁性链接	margicLink := fmt.Sprintf("magnet:?xt=urn:btih:%s&dn=%s", infoHash, url.QueryEscape(fileName))

	f, _ := os.Open(mw.prevFilePath)
	defer f.Close()

	fInfo, _ := f.Stat()
	fileSize := fInfo.Size()

	filemd4, _ := file2md4(mw.prevFilePath)

	ed2kLink := fmt.Sprintf("ed2k://|file|%s|%d|%s|", url.QueryEscape(fileName), fileSize, filemd4)
	mw.outTE2.SetText(ed2kLink) //设置电驴链接

	return nil
}

func (mw *MyMainWindow) aboutAction_Triggered() {
	walk.MsgBox(mw, "关于链接生成器(win32平台)", "链接无忧 www.link51.com \n软件可以随意复制 \n Email：admin@link51.com", walk.MsgBoxIconInformation)
}

func file2sha1(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha1.New()
	_, erro := io.Copy(h, file)
	if erro != nil {
		return "", erro
	}
	out := fmt.Sprintf("%x", h.Sum(nil))
	return out, nil
}

func file2md4(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := md4.New()
	_, erro := io.Copy(h, file)
	if erro != nil {
		return "", erro
	}
	out := fmt.Sprintf("%x", h.Sum(nil))
	return out, nil
}
