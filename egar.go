// get egarat go version
package main

import (
	b "bufio"
	by "bytes"
	fl "flag"
	f "fmt"
	te "text/template"
	i "io"
	u "net/url"
	o "os"
	r "regexp"
	s "strings"
	t "time"
	l "log"
)

const (
	blankCommentLine = `^[\s\t]*#?`
)

// globals

// Egar holding data for each renter
// which read from csv, elec, water,
// and fix files and divided by ratio
type Egar struct {
	Mansion           int
	Name              string
	ElecPlate         string
	EgarValue         float64
	ElecFare          float64
	//ElecSellemFare    float64
	//ElecBasementFare  float64
	WaterFare         float64
	FixesFare         float64
	Message string
	StairElecFare     float64
	BaseElecFare      float64
	StairCleaningFare float64
	Mobile            string
}

//Total ccalculate total egar value
func (e Egar) Total() float64 {
	return e.EgarValue +
	e.ElecFare +
	//e.ElecSellemFare +
	//e.ElecBasementFare +
	e.BaseElecFare +
	e.StairElecFare +
	e.StairCleaningFare +
	e.WaterFare +
	e.FixesFare
}

//TODO: to be capitalized and commented
var (
	//t0
	t0                        = t.Now()
	TimeDue                   = t.Now()
	dataFile                  = "egar.csv"
	waterdir                  = "water"
	elecdir                   = "elec"
	fixdir                    = "fix"
	messagedir = "message"
	prefix                    = ""
	// datakeys csv 1st raw keys
	dataKeys                  = []string{}
	// map datakeys to column order
	keysMap                   = map[string]int{}
	// csv file read to slice of column to value
	dataMap                   = []map[string]string{}
	grossWaterFare    float64 = 0.0
	grossFixFare      float64 = 0.0
	elecfares                 = map[string]float64{}
	notes                     = []string{}
	message = ""
	sep                       = string(o.PathSeparator)
	egarat                    = []Egar{}
	waterFactor       float64 = 7.5
	factor            float64 = 7.0
	stairElecPlate            = "790"
	basementElecPlate         = "789"
	//stairFare cleaning fare
	stairFare         float64 = 50
)

func IsEmptyOrCommentLine(line string)bool {
	result,err:=r.Match(blankCommentLine,[]byte(line))
	if err != nil {
		l.Panic(err)
	}
	return result
}

func init(){
	//setup logger
	mainLog,err:=o.Create("log")
	if err != nil{
		panic(err)
	}
	l.Default().SetOutput(mainLog)
}
// main app function


func GetEgarat() {
	//welcome
	l.Println("Welcome to ", o.Args, " ", t0)

	//finding and formating date path infix
	yearDue, monthDue, _ := TimeDue.Date()
	stringDue := f.Sprintf("%4d%[3]c%02[2]d", yearDue, monthDue, o.PathSeparator)
	l.Println(stringDue)

	//reading electricity file
	elecReader, err := o.Open(elecdir + sep + stringDue)
	if err != nil {
		panic(err)
	}
	//defer elecReader.Close()
	elecScanner := b.NewScanner(elecReader)

	for elecScanner.Scan() {
		//if IsEmptyOrCommentLine( elecScanner.Text()){
		//	continue
		//}

		l.Println(elecScanner.Text())
		fields := s.Fields(elecScanner.Text())
		//if len(fields)<2{
		//	continue
		//}
		ff := 0.0
		f.Sscan(fields[1], &ff)
		elecfares[fields[0]] = ff
	}
	l.Println("ef", elecfares)

	//reading water file
	waterReader, err := o.Open(waterdir + sep + stringDue)
	if err != nil {
		notes = append(notes, err.Error())
	}
	f.Fscan(waterReader, &grossWaterFare)
//	f.Println("water:", grossWaterFare)

	//reading fixation file
	fixReader, err := o.Open(fixdir + sep + stringDue)
	if err != nil {
		notes = append(notes, err.Error())
	}
	f.Fscan(fixReader, &grossFixFare)
	//f.Println("fix:", grossFixFare)

	//reading message file
	message=func()string{ 
		r,e:=o.ReadFile(messagedir + sep + stringDue);
		if e!=nil{panic(e)};
		return string(r)}()
	//messageReader, err := o.Open(messagedir + sep + stringDue)
	//if err != nil {
	//	notes = append(notes, err.Error())
	//}
	//f.Fscanln(messageReader, &message)
	//f.Println("fix:", grossFixFare)


	//reading main data
	dataReader, err := o.Open(dataFile)
	if err != nil {
		panic(err)
	}
	dataScanner := b.NewScanner(dataReader)

	dataScanner.Scan()
	for emptyOrComment, err := r.MatchString(blankCommentLine, dataScanner.Text()); emptyOrComment && err != nil; {
		if err != nil {
			panic(err)
		}
		dataScanner.Scan()
	}
	for i, key := range s.Fields(dataScanner.Text()) {
		dataKeys = append(dataKeys, key)
		keysMap[key] = i
	}
	//f.Println(dataKeys, keysMap)
	for dataScanner.Scan() {
		//f.Println(dataScanner.Text())
		fields := s.Fields(dataScanner.Text())
		if len(fields) == 0 {
			continue
		}
		egar := Egar{}
		f.Sscan(fields[keysMap["mansion"]], &(egar.Mansion))
		f.Sscan(fields[keysMap["name"]], &egar.Name)
		f.Sscan(fields[keysMap["elec"]], &egar.ElecPlate)
		f.Sscan(fields[keysMap["egar"]], &egar.EgarValue)
		f.Sscan(fields[keysMap["mobile"]], &egar.Mobile)
		egar.Message=message
		egarat = append(egarat, egar)
	}
	//f.Printf("%#v\n", egarat)
	for i, _ := range egarat {
		egarat[i].ElecFare = elecfares[egarat[i].ElecPlate]
		egarat[i].WaterFare = grossWaterFare / waterFactor
		egarat[i].FixesFare = grossFixFare / factor
		egarat[i].StairElecFare = elecfares[stairElecPlate] / factor
		egarat[i].BaseElecFare = elecfares[basementElecPlate] / factor
		egarat[i].StairCleaningFare = stairFare
		//f.Printf("%+v\n", egarat[i])
		//f.Println(egarat[i].Total())
	}
	_ = `
=======start====
{{$f:="%.2f"}}
{{/*range $_,$eg:=.*/}}
===========================
name:		Mr. {{$eg.Name}}
mansion:	-- {{$eg.Mansion}}
plate:		-- {{$eg.ElecPlate}}
rent:		{{$eg.EgarValue|printf $f}}
elec:		{{$eg.ElecFare|printf $f}}
water:		{{$eg.WaterFare|printf $f}}
stair:		{{$eg.StairElecFare|printf $f}}
base:		{{$eg.BaseElecFare|printf $f}}
clean:		{{$eg.StairCleaningFare|printf $f}}
service:	{{$eg.FixesFare|printf $f}}

--------------------------
total:	 	{{.Total|printf $f}}
==============================
{{/*end*/}}
=======end====

`

	temp := `
{{$f:="%.2f"}}
==================
name:		Mr. {{.Name}}
mansion:		-- {{.Mansion}}
plate:			-- {{.ElecPlate}}
rent:		{{.EgarValue|printf $f}}
elec:		{{.ElecFare|printf $f}}
water:		{{.WaterFare|printf $f}}
stair:		{{.StairElecFare|printf $f}}
base:		{{.BaseElecFare|printf $f}}
clean:		{{.StairCleaningFare|printf $f}}
service:**	{{.FixesFare|printf $f}}
---------------------
total:	 	{{.Total|printf $f}}
=====================
notes:
**{{.Message}}
{{/*end*/}}
`
	stam := te.New("main")
	stam = stam.Funcs(te.FuncMap{"total": Egar.Total})
	te.Must(stam.Parse(temp))
	for _, eg := range egarat {

		buff := by.Buffer{}
		// buff:=by.NewBuffer()
		stam.Execute(i.Writer(&buff), eg)
		// out:=i.ReadAll(buff)
		// out:=buff.String()
		out, _ := i.ReadAll(i.Reader(&buff))
		f.Printf("%s\nhttps://wa.me/%s/?text=%s", string(out), eg.Mobile, u.PathEscape(string(out)))
		//f.Printf("%s\n<a href=\"https://wa.me/%s/?text=%s\">Send WhatsApp --></a>", string(out), eg.Mobile, u.PathEscape(string(out)))

	}
}

func main() {
	//Year Monyh egar
	YearMonthEgar := t.Now()
	year := 0
	fl.IntVar(&year, "yeat", 2020, "year of egar")
	fl.Parse()
	f.Println("egar", year, YearMonthEgar)
	GetEgarat()
l.Printf("kk%f\n",egarat[len(egarat)-1].StairElecFare)
	f.Println("Done")
}
