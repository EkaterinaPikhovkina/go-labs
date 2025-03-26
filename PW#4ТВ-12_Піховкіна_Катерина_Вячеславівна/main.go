package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Templates
var mainTemplate = template.Must(template.New("main").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Веб-калькулятор</title>
</head>
<body>
    <h1>Оберіть завдання</h1>
    <button onclick="location.href='/task1'">Завдання 1</button><br><br>
    <button onclick="location.href='/task2'">Завдання 2</button><br><br>
    <button onclick="location.href='/task3'">Завдання 3</button>
</body>
</html>
`))

var task1Template = template.Must(template.New("task1").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Завдання 1</title>
</head>
<body>
    <h1>Завдання 1</h1>
    <form action="/calculate_task1" method="post">
        <label for="I_к">I_к:</label>
        <input type="text" id="I_к" name="I_к" value="{{.Defaults.I_к}}"><br><br>
        <label for="t_ф">t_ф:</label>
        <input type="text" id="t_ф" name="t_ф" value="{{.Defaults.t_ф}}"><br><br>
        <label for="S_м">S_м:</label>
        <input type="text" id="S_м" name="S_м" value="{{.Defaults.S_м}}"><br><br>
        <label for="T_м">T_м:</label>
        <input type="text" id="T_м" name="T_м" value="{{.Defaults.T_м}}"><br><br>
        <label for="j_ек">j_ек:</label>
        <input type="text" id="j_ек" name="j_ек" value="{{.Defaults.j_ек}}"><br><br>
        <label for="U_ном">U_ном:</label>
        <input type="text" id="U_ном" name="U_ном" value="{{.Defaults.U_ном}}"><br><br>
        <label for="C_т">C_т:</label>
        <input type="text" id="C_т" name="C_т" value="{{.Defaults.C_т}}"><br><br>

        <button type="submit">Розрахувати</button>
    </form>
    {{if .Result}}
    <div>
        <h2>Результат:</h2>
        <p>{{.Result}}</p>
    </div>
    {{end}}
    <br><button onclick="location.href='/'">На головну</button>
</body>
</html>
`))

var task2Template = template.Must(template.New("task2").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Завдання 2</title>
</head>
<body>
    <h1>Завдання 2</h1>
    <form action="/calculate_task2" method="post">
        <label for="U_с_н">U_с.н:</label>
        <input type="text" id="U_с_н" name="U_с_н" value="{{.Defaults.U_с_н}}"><br><br>
        <label for="S_к">S_к:</label>
        <input type="text" id="S_к" name="S_к" value="{{.Defaults.S_к}}"><br><br>
        <label for="U_к_percent">U_к%:</label>
        <input type="text" id="U_к_percent" name="U_к_percent" value="{{.Defaults.U_к_percent}}"><br><br>
        <label for="S_ном_т">S_ном.т:</label>
        <input type="text" id="S_ном_т" name="S_ном_т" value="{{.Defaults.S_ном_т}}"><br><br>

        <button type="submit">Розрахувати</button>
    </form>
    {{if .Result}}
    <div>
        <h2>Результат:</h2>
        <p>{{.Result}}</p>
    </div>
    {{end}}
    <br><button onclick="location.href='/'">На головну</button>
</body>
</html>
`))

var task3Template = template.Must(template.New("task3").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Завдання 3</title>
</head>
<body>
    <h1>Завдання 3</h1>
    <form action="/calculate_task3" method="post">
        <label for="U_к_max">U_к.max:</label>
        <input type="text" id="U_к_max" name="U_к_max" value="{{.Defaults.U_к_max}}"><br><br>
        <label for="U_в_н">U_в.н:</label>
        <input type="text" id="U_в_н" name="U_в_н" value="{{.Defaults.U_в_н}}"><br><br>
        <label for="U_н_н">U_н.н:</label>
        <input type="text" id="U_н_н" name="U_н_н" value="{{.Defaults.U_н_н}}"><br><br>
        <label for="R_с_н">R_с.н:</label>
        <input type="text" id="R_с_н" name="R_с_н" value="{{.Defaults.R_с_н}}"><br><br>
        <label for="X_с_н">X_с.н:</label>
        <input type="text" id="X_с_н" name="X_с_н" value="{{.Defaults.X_с_н}}"><br><br>
        <label for="R_с_min">R_с.min:</label>
        <input type="text" id="R_с_min" name="R_с_min" value="{{.Defaults.R_с_min}}"><br><br>
        <label for="X_с_min">X_с.min:</label>
        <input type="text" id="X_с_min" name="X_с_min" value="{{.Defaults.X_с_min}}"><br><br>
        <label for="l_л">l_л:</label>
        <input type="text" id="l_л" name="l_л" value="{{.Defaults.l_л}}"><br><br>
        <label for="R_0">R_0:</label>
        <input type="text" id="R_0" name="R_0" value="{{.Defaults.R_0}}"><br><br>
        <label for="X_0">X_0:</label>
        <input type="text" id="X_0" name="X_0" value="{{.Defaults.X_0}}"><br><br>

        <button type="submit">Розрахувати</button>
    </form>
    {{if .Result}}
    <div>
        <h2>Результат:</h2>
        <p>{{.Result}}</p>
    </div>
    {{end}}
    <br><button onclick="location.href='/'">На головну</button>
</body>
</html>
`))

// Data structures to pass to templates
type TaskData struct {
	Defaults map[string]string
	Result   string
	Title    string
}

// Handlers
func mainHandler(w http.ResponseWriter, r *http.Request) {
	mainTemplate.Execute(w, nil)
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	defaults := map[string]string{
		"I_к":   "2.5",
		"t_ф":   "2.5",
		"S_м":   "1300",
		"T_м":   "4000",
		"j_ек":  "1.4",
		"U_ном": "10",
		"C_т":   "92",
	}
	data := TaskData{Defaults: defaults, Title: "Завдання 1"}
	task1Template.Execute(w, data)
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	defaults := map[string]string{
		"U_с_н":       "10.5",
		"S_к":         "200",
		"U_к_percent": "10.5",
		"S_ном_т":     "6.3",
	}
	data := TaskData{Defaults: defaults, Title: "Завдання 2"}
	task2Template.Execute(w, data)
}

func task3Handler(w http.ResponseWriter, r *http.Request) {
	defaults := map[string]string{
		"U_к_max": "11.1",
		"U_в_н":   "115",
		"U_н_н":   "11",
		"S_ном_т": "6.3",
		"R_с_н":   "10.65",
		"X_с_н":   "24.02",
		"R_с_min": "34.88",
		"X_с_min": "65.68",
		"l_л":     "12.37",
		"R_0":     "0.64",
		"X_0":     "0.363",
	}
	data := TaskData{Defaults: defaults, Title: "Завдання 3"}
	task3Template.Execute(w, data)
}

func calculateTask1Handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	iKStr := r.FormValue("I_к")
	tFStr := r.FormValue("t_ф")
	sMStr := r.FormValue("S_м")
	// tMStr := r.FormValue("T_м") // Removed unused tMStr
	jEkStr := r.FormValue("j_ек")
	uNomStr := r.FormValue("U_ном")
	cTStr := r.FormValue("C_т")

	iK, _ := strconv.ParseFloat(iKStr, 64)
	tF, _ := strconv.ParseFloat(tFStr, 64)
	sM, _ := strconv.ParseFloat(sMStr, 64)
	// tM, _ := strconv.ParseFloat(tMStr, 64) // Removed unused tM
	jEk, _ := strconv.ParseFloat(jEkStr, 64)
	uNom, _ := strconv.ParseFloat(uNomStr, 64)
	cT, _ := strconv.ParseFloat(cTStr, 64)

	s := calculateS(sM, uNom, jEk)
	smin := calculateSmin(iK, tF, cT)

	var result string
	if s > smin {
		result = "Вибираємо кабель ААБ 10 3×25 з допустимим струмом I_доп=90 А. " +
			"Однак за термічною стійкістю до дії струмів КЗ s≤s_min," +
			"що вимагає збільшення перерізу жил кабелю до 50 〖мм〗^2."
	} else {
		result = "Вибираємо кабель ААБ 10 3×25 з допустимим струмом I_доп=90 А."
	}

	defaults := map[string]string{
		"I_к":   iKStr,
		"t_ф":   tFStr,
		"S_м":   sMStr,
		"T_м":   r.FormValue("T_м"), // Keep T_м in defaults for template rendering
		"j_ек":  jEkStr,
		"U_ном": uNomStr,
		"C_т":   cTStr,
	}
	data := TaskData{Defaults: defaults, Result: result, Title: "Завдання 1"}
	task1Template.Execute(w, data)
}

func calculateTask2Handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	uSnStr := r.FormValue("U_с_н")
	skStr := r.FormValue("S_к")
	ukPercentageStr := r.FormValue("U_к_percent")
	sNomTStr := r.FormValue("S_ном_т")

	uSn, _ := strconv.ParseFloat(uSnStr, 64)
	sk, _ := strconv.ParseFloat(skStr, 64)
	ukPercentage, _ := strconv.ParseFloat(ukPercentageStr, 64)
	sNomT, _ := strconv.ParseFloat(sNomTStr, 64)

	ip0 := calculateIp0(uSn, sk, ukPercentage, sNomT)
	result := fmt.Sprintf("Початкове діюче значення струму трифазного КЗ дорівнює %.2f кА.", ip0)

	defaults := map[string]string{
		"U_с_н":       uSnStr,
		"S_к":         skStr,
		"U_к_percent": ukPercentageStr,
		"S_ном_т":     sNomTStr,
	}
	data := TaskData{Defaults: defaults, Result: result, Title: "Завдання 2"}
	task2Template.Execute(w, data)
}

func calculateTask3Handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	ukMaxStr := r.FormValue("U_к_max")
	uvnStr := r.FormValue("U_в_н")
	unnStr := r.FormValue("U_н_н")
	rSnStr := r.FormValue("R_с_н")
	xSnStr := r.FormValue("X_с_н")
	rSminStr := r.FormValue("R_с_min")
	xSminStr := r.FormValue("X_с_min")
	llStr := r.FormValue("l_л")
	r0Str := r.FormValue("R_0")
	x0Str := r.FormValue("X_0")

	ukMax, _ := strconv.ParseFloat(ukMaxStr, 64)
	uvn, _ := strconv.ParseFloat(uvnStr, 64)
	unn, _ := strconv.ParseFloat(unnStr, 64)
	rSn, _ := strconv.ParseFloat(rSnStr, 64)
	xSn, _ := strconv.ParseFloat(xSnStr, 64)
	rSmin, _ := strconv.ParseFloat(rSminStr, 64)
	xSmin, _ := strconv.ParseFloat(xSminStr, 64)
	ll, _ := strconv.ParseFloat(llStr, 64)
	r0, _ := strconv.ParseFloat(r0Str, 64)
	x0, _ := strconv.ParseFloat(x0Str, 64)

	iln := calculateIln(ukMax, uvn, unn, rSn, xSn, rSmin, xSmin, ll, r0, x0)
	result := fmt.Sprintf(
		"Струми трифазного та двофазного КЗ в точці 10 в нормальному та мінімальному режимах:\n"+
			"I_(л.н)^((3) )=%.2f А\n"+
			"I_(л.н)^((2) )=%.2f А\n"+
			"I_(л.н.min)^((3) )=%.2f А\n"+
			"I_(л.н.min)^((2) )=%.2f А\n"+
			"Аварійний режим на станції не передбачений, оскільки вона живить споживачів 3 категорії.",
		iln[0], iln[1], iln[2], iln[3],
	)

	defaults := map[string]string{
		"U_к_max": ukMaxStr,
		"U_в_н":   uvnStr,
		"U_н_н":   unnStr,
		"R_с_н":   rSnStr,
		"X_с_н":   xSnStr,
		"R_с_min": rSminStr,
		"X_с_min": xSminStr,
		"l_л":     llStr,
		"R_0":     r0Str,
		"X_0":     x0Str,
	}
	data := TaskData{Defaults: defaults, Result: result, Title: "Завдання 3"}
	task3Template.Execute(w, data)
}

// Calculation functions (Go versions of Kotlin functions)
func calculateS(sM float64, uNom float64, jEk float64) float64 {
	iM := (sM / 2) / (math.Sqrt(3.0) * uNom)
	return iM / jEk
}

func calculateSmin(iK float64, tF float64, cT float64) float64 {
	return (iK * math.Sqrt(tF)) / cT
}

func calculateIp0(uSn float64, sk float64, ukPercentage float64, sNomT float64) float64 {
	xC := (uSn * uSn) / sk
	xT := (ukPercentage / 100) * (uSn * uSn) / sNomT
	xTotal := xC + xT
	return uSn / (math.Sqrt(3.0) * xTotal)
}

func calculateIln(ukMax float64, uvn float64, unn float64, rSn float64, xSn float64, rSmin float64,
	xSmin float64, ll float64, r0 float64, x0 float64) []float64 {
	xt := ukMax / 100 * (uvn * uvn) / 6.3
	xSh := xSn + xt
	xShMin := xSmin + xt

	kPr := (unn * unn) / (uvn * uvn)

	rShN := rSn * kPr
	xShN := xSh * kPr

	rL := ll * r0
	xL := ll * x0

	rSigmaN := rL + rShN
	xSigmaN := xL + xShN
	zSigmaN := math.Sqrt(rSigmaN*rSigmaN + xSigmaN*xSigmaN)

	iLn3 := (unn * 1000) / (math.Sqrt(3.0) * zSigmaN)
	iLn2 := iLn3 * math.Sqrt(3.0) / 2

	rSigmaMin := rL + (rSmin * kPr)
	xSigmaMin := xL + (xShMin * kPr)
	zSigmaMin := math.Sqrt(rSigmaMin*rSigmaMin + xSigmaMin*xSigmaMin)

	iLnMin3 := (unn * 1000) / (math.Sqrt(3.0) * zSigmaMin)
	iLnMin2 := iLnMin3 * math.Sqrt(3.0) / 2

	return []float64{iLn3, iLn2, iLnMin3, iLnMin2}
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)
	http.HandleFunc("/task3", task3Handler)
	http.HandleFunc("/calculate_task1", calculateTask1Handler)
	http.HandleFunc("/calculate_task2", calculateTask2Handler)
	http.HandleFunc("/calculate_task3", calculateTask3Handler)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
