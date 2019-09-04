package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

const (
	config = ".objects"
	corruptInfo="objects list corrupt"
)

var (
	order = binary.LittleEndian
)

type objectInfo struct{
	objectInfoFixPart
	name string
}

type objectInfoFixPart struct {
	Id          int64
	SkuId       int64
	Price       float64
	LowestPrice float64
	Len         int64
}

func (oi *objectInfo)Equal(ei *objectInfo)bool{
	return oi.Id ==ei.Id && oi.SkuId ==ei.SkuId
}

func (oi *objectInfo)getPrise() bool {
	res, err := http.Get(fmt.Sprintf("https://item.taobao.com/item.htm?id=%d",oi.Id))
	if err != nil {
		return false
	}
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false
	}
	reg1 := regexp.MustCompile(`skuMap\s+:\s*\{.*\}`)
	regResultTmp := reg1.Find(s)
	reg2 := regexp.MustCompile(`"[0-9;:]*":{[^}]*}`)
	regResult := reg2.FindAll(regResultTmp, -1)
	for _, skuString := range regResult {
		n := bytes.IndexByte(skuString, '{')
		if n == -1 {
			continue
		}
		skuString = skuString[n:]
		var sku Sku
		if err := json.Unmarshal(skuString, &sku); err != nil {
			continue
		}
		checkId, err := strconv.Atoi(sku.SkuId)
		if err != nil {
			continue
		}
		if int64(checkId) == oi.SkuId {
			if ans, err := strconv.ParseFloat(sku.Price, 64); err == nil {
				oi.Price =ans
				if ans<oi.LowestPrice {
					oi.LowestPrice =oi.Price
				}
				return true
			}
			return false
		}
	}
	return false
}

func (mw *MyMainWindow) objectInit(){
	s,err:=ioutil.ReadFile(config)
	b:=bytes.NewReader(s)
	if err!=nil{
		log.Printf("read objects fail")
		return
	}
	for b.Len()!=0{
		if oi:=readObject(b);oi!=nil{
			mw.info.items=append(mw.info.items,oi)
		}else{
			return
		}
	}
	return
}

func (mw *MyMainWindow) addObject(){
	defer func(){
		_=mw.name.SetText("")
		_=mw.id.SetText("")
		_=mw.skuId.SetText("")
	}()
	oi:=new(objectInfo)
	oi.name=mw.name.Text()
	oi.Len =int64(len(oi.name))
	if i,err:=strconv.Atoi(mw.id.Text());err!=nil{
		return
	}else{
		oi.Id =int64(i)
	}
	if i,err:=strconv.Atoi(mw.skuId.Text());err!=nil{
		return
	}else{
		oi.SkuId =int64(i)
	}
	if !mw.isNew(oi){
		return
	}
	oi.Price =math.MaxFloat64
	oi.LowestPrice =math.MaxFloat64
	mw.info.items=append(mw.info.items,oi)
	mw.updatePrices()
}

func (mw *MyMainWindow)isNew(test *objectInfo)bool{
	for _,oi:=range mw.info.items{
		if test.Equal(oi){
			return false
		}
	}
	return true
}

func (mw *MyMainWindow)updatePrices(){
	for i,_:=range mw.info.items{
		if !mw.info.items[i].getPrise(){
			mw.info.items[i].Price =-1.0
		}
	}
	mw.info.PublishRowsReset()
}

func (m *stupidModel)RowCount() int{
	return len(m.items)
}

type stupidModel struct {
	walk.TableModelBase
	items []*objectInfo
}

func (m *stupidModel)Value(row,col int)interface{}{
	oi:=m.items[row]
	switch col{
	case 0:
		return oi.name
	case 1:
		return oi.Id
	case 2:
		return oi.SkuId
	case 3:
		return oi.Price
	case 4:
		return oi.LowestPrice
	}
	log.Fatalf("unexpected column")
	return nil
}

func readObject (r io.Reader)*objectInfo{
	oi:=new(objectInfo)
	if err:=binary.Read(r,order,&oi.objectInfoFixPart);err!=nil{
		return nil
	}
	s:=make([]byte,oi.Len)
	if n,err:=r.Read(s);err!=nil||n!=int(oi.Len){
		return nil
	}
	oi.name=string(s)
	return oi
}

func (mw *MyMainWindow) writeObjects() {
	file,err:=os.OpenFile(config,os.O_WRONLY|os.O_CREATE|os.O_TRUNC,0666)
	if err!=nil{
		log.Println("create ",config," fail")
		return
	}
	for _,oi:=range mw.info.items{
		if !oi.store(file){
			log.Println("store objects fail")
		}
	}
}

func (oi objectInfo) store(w io.Writer) bool {
	if err:=binary.Write(w,order,oi.objectInfoFixPart);err!=nil{
		return false
	}
	n,err:=w.Write([]byte(oi.name))
	return err==nil && n==int(oi.Len)
}

func (m *stupidModel) colorSet(style *walk.CellStyle){
	oi:= m.items[style.Row()]
	if oi.Price==oi.LowestPrice{
		style.BackgroundColor=walk.RGB(159, 215, 255)
	}
}