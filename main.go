package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
)

type Sku struct{
	Price    string
	Stock    string
	SkuId    string
	OverSold bool
}

type MyMainWindow struct{
	*walk.MainWindow
	name   *walk.LineEdit
	id     *walk.LineEdit
	skuId  *walk.LineEdit
	add    *walk.PushButton
	table  *walk.TableView
	info *stupidModel
}

func main () {
	mw:=&MyMainWindow{info:newStupidModel()}
	if err:=(MainWindow{
		AssignTo:&mw.MainWindow,
		Title: "DailyCheckPrices",
		Size:Size{Width: 550,Height: 400},
		Layout: VBox{},
		Children: []Widget{
			Composite{
				Layout:HBox{},
				Children:[]Widget{
					Label{Text:"商品名: "},
					LineEdit{ AssignTo: &mw.name},
					Label{Text:"ID: "},
					LineEdit{AssignTo:&mw.id},
					Label{Text:"skuID: "},
					LineEdit{AssignTo:&mw.skuId},
					PushButton{
						AssignTo:&mw.add,
						Text:"加入",
					},
				},
			},
			TableView{
				AssignTo:&mw.table,
				Columns:[]TableViewColumn{
					{Title:"商品名"},
					{Title:"ID"},
					{Title:"skuID"},
					{Title:"价格"},
					{Title:"最低价"},
				},
				Model:mw.info,
				StyleCell:mw.info.colorSet,
			},
		},
	}.Create());err!=nil{
		log.Fatal(err)
	}

	mw.init()
	mw.add.Clicked().Attach(func(){
		go mw.addObject()
	})
	mw.info.PublishRowsReset()
	mw.Run()
	defer mw.writeObjects()
}

func (mw *MyMainWindow) init(){
	mw.objectInit()
}

func newStupidModel ()*stupidModel{
	s:=new(stupidModel)
	return s
}
