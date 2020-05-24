package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"text/template"
	"time"
)

var task1Templ = `
Задача №1.
Электролампы изготавливаются на трех заводах.
1-й завод производит {{.PH1}}% общего количества ламп, 2-й – {{.PH2}}%, а 3-й – остальную часть.
Продукция 1-го завода содержит {{.PAH1}}% бракованных ламп, 2-го – {{.PAH2}}%, 3-го – {{.PAH3}}%.
В магазин поступает продукция всех трех заводов. Купленная лампа оказалась с браком.
Какова вероятность того, что она произведена {{.Number}}-м заводом?
`

var task2Templ = `
Задача №2.
2 стрелка независимо один от другого стреляют по одной мишени, делая каждый по одному выстрелу.
Вероятность попадания в мишень для 1-го стрелка p1 = {{.P1}}; для 2-го стрелка p2 = {{.P2}}.
После стрельбы в мишени обнаружена одна пробоина.
Найти вероятность того, что эта пробоина принадлежит {{.Number}}-му стрелку.
`

var task3Templ = `
Задача №3.
Расследуются причины авиационной катастрофы, о которых можно сделать четыре гипотезы: H1, H2, H3, H4.
Согласно статистике P(H1)= {{.PH1}}; P(H2)= {{.PH2}}; P(H3)= {{.PH3}}; P(H4)= {{.PH4}}.
Обнаружено, что в ходе катастрофы произошло событие A={воспламенение горючего}.
Условные вероятности события A при гипотезах H1, H2, H3, H4 согласно той же статистике известны:
P(A|H1)= {{.PAH1}}; P(A|H2)= {{.PAH2}}; P(A|H3)= {{.PAH3}}; P(A|H4)= {{.PAH4}}.
Какая гипотеза наиболее вероятна до известия о событии A и после него?
`

var dtl1Templ = `{{template "task1" .}}
Подробное решение:
P(H1) = {{printf "%.2f" (divBy100 .PH1)}}
P(H2) = {{printf "%.2f" (divBy100 .PH2)}}
P(H3) = {{printf "%.2f" (divBy100 .PH3)}}

P(A|H1) = {{printf "%.2f" (divBy100 .PAH1)}}
P(A|H2) = {{printf "%.2f" (divBy100 .PAH2)}}
P(A|H3) = {{printf "%.2f" (divBy100 .PAH3)}}

P(A) = {{printf "%.5f" .PA}}

Ответ:
P(H{{.Number}}|A) = P(H{{.Number}}) * P(A|H{{.Number}}) / P(A) = {{printf "%.5f" .Answer}}
`

var dtl2Templ = `{{template "task2" .}}
Подробное решение:
H1 = {1-й стрелок попал, 2-й не попал}
H2 = {1-й стрелок не попал, 2-й попал}
H3 = {оба стрелка не попали}
H4 = {оба стрелка попали}

P(H1) = P1 * (1 - P2) = {{printf "%.5f" .PH1}}
P(H2) =	(1 - P1) * P2 = {{printf "%.5f" .PH2}}

P(A|H3), P(A|H4) = 0
P(A|H1), P(A|H2) = 1

P(A) = {{printf "%.5f" .PA}}

Ответ:
P(H{{.Number}}|A) = P(H{{.Number}}) * P(A|H{{.Number}}) / P(A) = {{printf "%.5f" .Answer}}
`

var dtl3Templ = `{{template "task3" .}}
Подробное решение:
P(A) = {{printf "%.5f" .PA}}

P(H1|A) = {{printf "%.5f" .PH1A}}
P(H2|A) = {{printf "%.5f" .PH2A}}
P(H3|A) = {{printf "%.5f" .PH3A}}
P(H4|A) = {{printf "%.5f" .PH4A}}

Ответ:
до известия -  max{P(Hi)} = P({{.MaxPH.Key}}) = {{printf "%.2f" .MaxPH.Value}}, Hi = {{.MaxPH.Key}}; 
после - max{P(Hi|A)} = P({{.MapPHA.Key}}|A) = {{printf "%.5f" .MapPHA.Value}}, Hi = {{.MapPHA.Key}}
`

var ans1Templ = `{{template "task1" .}}
Ответ: {{printf "%.5f" .Answer}}
`

var ans2Templ = `{{template "task2" .}}
Ответ: {{printf "%.5f" .Answer}}
`

var ans3Templ = `{{template "task3" .}}
Ответ: до известия - {{.MaxPH.Key}}, P({{.MaxPH.Key}}) = {{printf "%.2f" .MaxPH.Value}}; 
после - {{.MapPHA.Key}}, P({{.MapPHA.Key}}|A) = {{printf "%.5f" .MapPHA.Value}}
`

var taskFile, dtlFile, ansFile, errFile *os.File

type data1 struct {
	Number int
	PH1    int
	PH2    int
	PH3    int
	PAH1   int
	PAH2   int
	PAH3   int
	PA     float64
	Answer float64
}

type data2 struct {
	Number int
	P1     float64
	P2     float64
	PH1    float64
	PH2    float64
	PA     float64
	Answer float64
}

type data3 struct {
	PH1    float64
	PH2    float64
	PH3    float64
	PH4    float64
	PAH1   float64
	PAH2   float64
	PAH3   float64
	PAH4   float64
	PH1A   float64
	PH2A   float64
	PH3A   float64
	PH4A   float64
	PA     float64
	MaxPH  sortField
	MapPHA sortField
}

type sortField struct {
	Key   string
	Value float64
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			printErrorsToFile(r.(error))
		}
	}()

	os.Remove("ошибки.txt")
	rand.Seed(time.Now().UnixNano())
	taskNum := 0

	fmt.Println("Введите нужное количество вариантов: ")

	_, err := fmt.Scanln(&taskNum)

	if err != nil {
		printErrorsToFile(err)
	}

	taskFile, err = os.Create("задачи.txt")

	if err != nil {
		printErrorsToFile(err)
	}

	dtlFile, err = os.Create("задачи_с_подробным_решением.txt")

	if err != nil {
		printErrorsToFile(err)
	}

	ansFile, err = os.Create("задачи_с_ответами.txt")

	if err != nil {
		printErrorsToFile(err)
	}

	for i := 1; i <= taskNum; i++ {
		taskFile.WriteString(fmt.Sprintf("\r\nВариант %d\r\n", i))
		dtlFile.WriteString(fmt.Sprintf("\r\nВариант %d\r\n", i))
		ansFile.WriteString(fmt.Sprintf("\r\nВариант %d\r\n", i))

		err = task1()

		if err != nil {
			printErrorsToFile(err)
		}

		err = task2()

		if err != nil {
			printErrorsToFile(err)
		}

		err = task3()

		if err != nil {
			printErrorsToFile(err)
		}

		taskFile.WriteString("\r\n")
		dtlFile.WriteString("\r\n")
		ansFile.WriteString("\r\n")
	}
}

func task1() error {
	PH1 := rand.Intn(98) + 1
	PH2 := rand.Intn(99-PH1) + 1
	PH3 := 100 - PH1 - PH2
	Number := rand.Intn(3) + 1

	PAH1 := rand.Intn(48) + 1
	PAH2 := rand.Intn(49-PAH1) + 1
	PAH3 := 50 - PAH1 - PAH2

	PA := divBy100(PH1)*divBy100(PAH1) + divBy100(PH2)*divBy100(PAH2) + divBy100(PH3)*divBy100(PAH3)
	var Answer float64

	if Number == 1 {
		Answer = divBy100(PH1) * divBy100(PAH1) / PA
	} else if Number == 2 {
		Answer = divBy100(PH2) * divBy100(PAH2) / PA
	} else {
		Answer = divBy100(PH3) * divBy100(PAH3) / PA

	}

	d := data1{
		Number: Number,
		PH1:    PH1,
		PH2:    PH2,
		PH3:    PH3,
		PAH1:   PAH1,
		PAH2:   PAH2,
		PAH3:   PAH3,
		PA:     PA,
		Answer: Answer,
	}

	funcMap := template.FuncMap{
		"divBy100": divBy100,
	}

	task, err := template.New("task1").Funcs(funcMap).Parse(task1Templ)

	if err != nil {
		return err
	}

	dtl, err := task.New("dtl1").Parse(dtl1Templ)

	if err != nil {
		return err
	}

	ans, err := task.New("ans1").Parse(ans1Templ)

	if err != nil {
		return err
	}

	err = task.Execute(taskFile, d)

	if err != nil {
		return err
	}

	err = dtl.Execute(dtlFile, d)

	if err != nil {
		return err
	}

	err = ans.Execute(ansFile, d)

	if err != nil {
		return err
	}

	return nil
}

func task2() error {
	P1 := divBy100(rand.Intn(99) + 1)
	P2 := divBy100(rand.Intn(99) + 1)
	Number := rand.Intn(2) + 1

	PH1 := P1 * (1 - P2)
	PH2 := (1 - P1) * P2

	PA := PH1 + PH2
	var Answer float64

	if Number == 1 {
		Answer = PH1 / PA
	} else {
		Answer = PH2 / PA
	}

	d := data2{
		Number: Number,
		P1:     P1,
		P2:     P2,
		PH1:    PH1,
		PH2:    PH2,
		PA:     PA,
		Answer: Answer,
	}

	task, err := template.New("task2").Parse(task2Templ)

	if err != nil {
		return err
	}

	dtl, err := task.New("dtl2").Parse(dtl2Templ)

	if err != nil {
		return err
	}

	ans, err := task.New("ans2").Parse(ans2Templ)

	if err != nil {
		return err
	}

	err = task.Execute(taskFile, d)

	if err != nil {
		return err
	}

	err = dtl.Execute(dtlFile, d)

	if err != nil {
		return err
	}

	err = ans.Execute(ansFile, d)

	if err != nil {
		return err
	}

	return nil
}

func task3() error {
	PH1Int := rand.Intn(97) + 1
	PH2Int := rand.Intn(98-PH1Int) + 1
	PH3Int := rand.Intn(99-PH1Int-PH2Int) + 1
	PH4Int := 100 - PH1Int - PH2Int - PH3Int

	PH1 := divBy100(PH1Int)
	PH2 := divBy100(PH2Int)
	PH3 := divBy100(PH3Int)
	PH4 := divBy100(PH4Int)

	sorted := []sortField{
		{
			Key:   "H1",
			Value: PH1,
		},
		{
			Key:   "H2",
			Value: PH2,
		},
		{
			Key:   "H3",
			Value: PH3,
		},
		{
			Key:   "H4",
			Value: PH4,
		},
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	if sorted[0].Value == sorted[1].Value {
		sorted[0].Key = fmt.Sprintf("%s,%s", sorted[0].Key, sorted[1].Key)

		if sorted[1].Value == sorted[2].Value {
			sorted[0].Key = fmt.Sprintf("%s,%s", sorted[0].Key, sorted[2].Key)
		}
	}

	MaxPH := sorted[0]

	PAH1 := divBy100(rand.Intn(99) + 1)
	PAH2 := divBy100(rand.Intn(99) + 1)
	PAH3 := divBy100(rand.Intn(99) + 1)
	PAH4 := divBy100(rand.Intn(99) + 1)

	PA := PH1*PAH1 + PH2*PAH2 + PH3*PAH3 + PH4*PAH4

	PH1A := PH1 * PAH1 / PA
	PH2A := PH2 * PAH2 / PA
	PH3A := PH3 * PAH3 / PA
	PH4A := PH4 * PAH4 / PA

	sorted = []sortField{
		{
			Key:   "H1",
			Value: PH1A,
		},
		{
			Key:   "H2",
			Value: PH2A,
		},
		{
			Key:   "H3",
			Value: PH3A,
		},
		{
			Key:   "H4",
			Value: PH4A,
		},
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})

	if sorted[0].Value == sorted[1].Value {
		sorted[0].Key = fmt.Sprintf("%s,%s", sorted[0].Key, sorted[1].Key)

		if sorted[1].Value == sorted[2].Value {
			sorted[0].Key = fmt.Sprintf("%s,%s", sorted[0].Key, sorted[2].Key)

			if sorted[2].Value == sorted[3].Value {
				sorted[0].Key = fmt.Sprintf("%s,%s", sorted[0].Key, sorted[3].Key)
			}
		}
	}

	MaxPHA := sorted[0]

	d := data3{
		PH1:    PH1,
		PH2:    PH2,
		PH3:    PH3,
		PH4:    PH4,
		PAH1:   PAH1,
		PAH2:   PAH2,
		PAH3:   PAH3,
		PAH4:   PAH4,
		PA:     PA,
		PH1A:   PH1A,
		PH2A:   PH2A,
		PH3A:   PH3A,
		PH4A:   PH4A,
		MaxPH:  MaxPH,
		MapPHA: MaxPHA,
	}

	task, err := template.New("task3").Parse(task3Templ)

	if err != nil {
		return err
	}

	dtl, err := task.New("dtl3").Parse(dtl3Templ)

	if err != nil {
		return err
	}

	ans, err := task.New("ans3").Parse(ans3Templ)

	if err != nil {
		return err
	}

	err = task.Execute(taskFile, d)

	if err != nil {
		return err
	}

	err = dtl.Execute(dtlFile, d)

	if err != nil {
		return err
	}

	err = ans.Execute(ansFile, d)

	if err != nil {
		return err
	}

	return nil
}

func divBy100(n int) float64 {
	return float64(n) / 100
}

func printErrorsToFile(err error) {
	var unexpErr error
	errFile, unexpErr = os.Create("ошибки.txt")

	if unexpErr != nil {
		fmt.Fprintln(os.Stderr, unexpErr)
		os.Exit(1)
	}

	errStr := fmt.Sprintf("Программа завершилась с ошибкой:\r\n%s", err.Error())
	errFile.WriteString(errStr)
	os.Exit(1)
}
