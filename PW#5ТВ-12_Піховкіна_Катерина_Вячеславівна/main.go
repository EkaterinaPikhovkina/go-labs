// main.go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/", mainScreenHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainScreenHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("main_screen").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Головний екран</title>
</head>
<body>
    <h1>Оберіть завдання</h1>
    <button onclick="location.href='/task1'" type="button">Завдання 1</button><br><br>
    <button onclick="location.href='/task2'" type="button">Завдання 2</button>
</body>
</html>
`))
	tmpl.Execute(w, nil)
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	defaults := map[string]string{
		"w1":     "0.01",
		"w2":     "0.07",
		"w3":     "0.015",
		"w4":     "0.02",
		"w5":     "0.18",
		"tв1":    "30",
		"tв2":    "10",
		"tв3":    "100",
		"tв4":    "15",
		"tв5":    "2",
		"kп_max": "43",
	}

	data := struct {
		Defaults map[string]string
		Result   string
	}{
		Defaults: defaults,
		Result:   "",
	}

	if r.Method == "POST" {
		w1 := parseFloat(r.FormValue("w1"))
		w2 := parseFloat(r.FormValue("w2"))
		w3 := parseFloat(r.FormValue("w3"))
		w4 := parseFloat(r.FormValue("w4"))
		w5 := parseFloat(r.FormValue("w5"))

		tv1 := parseFloat(r.FormValue("tв1"))
		tv2 := parseFloat(r.FormValue("tв2"))
		tv3 := parseFloat(r.FormValue("tв3"))
		tv4 := parseFloat(r.FormValue("tв4"))
		tv5 := parseFloat(r.FormValue("tв5"))

		kpmax := parseFloat(r.FormValue("kп_max"))

		wos := calculateWos(w1, w2, w3, w4, w5)
		wds := calculateWds(w1, tv1, w2, tv2, w3, tv3, w4, tv4, w5, tv5, kpmax)

		if wos > wds {
			data.Result = "Надійність двоколової системи електропередачі є вищою ніж одноколової."
		} else {
			data.Result = "Надійність одноколової системи електропередачі є вищою ніж двоколової."
		}
	}

	tmpl := template.Must(template.New("task1").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Завдання 1</title>
</head>
<body>
    <h1>Завдання 1</h1>
    <form method="POST" action="/task1">
        <div>
            <label for="w1">w1:</label>
            <input type="text" id="w1" name="w1" value="{{.Defaults.w1}}">
            <label for="tв1">tв1:</label>
            <input type="text" id="tв1" name="tв1" value="{{.Defaults.tв1}}">
        </div>
        <div>
            <label for="w2">w2:</label>
            <input type="text" id="w2" name="w2" value="{{.Defaults.w2}}">
            <label for="tв2">tв2:</label>
            <input type="text" id="tв2" name="tв2" value="{{.Defaults.tв2}}">
        </div>
        <div>
            <label for="w3">w3:</label>
            <input type="text" id="w3" name="w3" value="{{.Defaults.w3}}">
            <label for="tв3">tв3:</label>
            <input type="text" id="tв3" name="tв3" value="{{.Defaults.tв3}}">
        </div>
        <div>
            <label for="w4">w4:</label>
            <input type="text" id="w4" name="w4" value="{{.Defaults.w4}}">
            <label for="tв4">tв4:</label>
            <input type="text" id="tв4" name="tв4" value="{{.Defaults.tв4}}">
        </div>
        <div>
            <label for="w5">w5:</label>
            <input type="text" id="w5" name="w5" value="{{.Defaults.w5}}">
            <label for="tв5">tв5:</label>
            <input type="text" id="tв5" name="tв5" value="{{.Defaults.tв5}}">
        </div>
        <div>
            <label for="kп_max">kп.max:</label>
            <input type="text" id="kп_max" name="kп_max" value="{{.Defaults.kп_max}}">
        </div>
        <button type="submit">Розрахувати</button>
    </form>
    {{if .Result}}
    <div>
        <h2>Результат:</h2>
        <p>{{.Result}}</p>
    </div>
    {{end}}
</body>
</html>
`))
	tmpl.Execute(w, data)
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	defaults := map[string]string{
		"w":      "0.01",
		"tв":     "0.045",
		"kп":     "0.004",
		"Pм":     "5120",
		"Tм":     "6451",
		"Зпер_а": "23.6",
		"Зпер_п": "17.6",
	}

	data := struct {
		Defaults map[string]string
		Result   string
	}{
		Defaults: defaults,
		Result:   "",
	}

	if r.Method == "POST" {
		w := parseFloat(r.FormValue("w"))
		tv := parseFloat(r.FormValue("tв"))
		kp := parseFloat(r.FormValue("kп"))
		pm := parseFloat(r.FormValue("Pм"))
		tm := parseFloat(r.FormValue("Tм"))
		z_per_a := parseFloat(r.FormValue("Зпер_а"))
		z_per_p := parseFloat(r.FormValue("Зпер_п"))

		mzper := calculateMZper(w, tv, kp, pm, tm, z_per_a, z_per_p)
		data.Result = fmt.Sprintf("Математичне сподівання збитків від переривання електропостачання М(Зпер) = %.2f грн.", mzper)
	}

	tmpl := template.Must(template.New("task2").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Завдання 2</title>
</head>
<body>
    <h1>Завдання 2</h1>
    <form method="POST" action="/task2">
        <div>
            <label for="w">w:</label>
            <input type="text" id="w" name="w" value="{{.Defaults.w}}">
            <label for="tв">tв:</label>
            <input type="text" id="tв" name="tв" value="{{.Defaults.tв}}">
        </div>
        <div>
            <label for="kп">kп:</label>
            <input type="text" id="kп" name="kп" value="{{.Defaults.kп}}">
        </div>
        <div>
            <label for="Pм">Pм:</label>
            <input type="text" id="Pм" name="Pм" value="{{.Defaults.Pм}}">
            <label for="Tм">Tм:</label>
            <input type="text" id="Tм" name="Tм" value="{{.Defaults.Tм}}">
        </div>
        <div>
            <label for="Зпер_а">Зпер.а:</label>
            <input type="text" id="Зпер_а" name="Зпер_а" value="{{.Defaults.Зпер_а}}">
        </div>
        <div>
            <label for="Зпер_п">Зпер.п:</label>
            <input type="text" id="Зпер_п" name="Зпер_п" value="{{.Defaults.Зпер_п}}">
        </div>
        <button type="submit">Розрахувати</button>
    </form>
    {{if .Result}}
    <div>
        <h2>Результат:</h2>
        <p>{{.Result}}</p>
    </div>
    {{end}}
</body>
</html>
`))
	tmpl.Execute(w, data)
}

func calculateWos(w1, w2, w3, w4, w5 float64) float64 {
	return w1 + w2 + w3 + w4 + w5
}

func calculateWds(w1, tv1, w2, tv2, w3, tv3, w4, tv4, w5, tv5, kpmax float64) float64 {
	wos := calculateWos(w1, w2, w3, w4, w5)
	tvos := (w1*tv1 + w2*tv2 + w3*tv3 + w4*tv4 + w5*tv5) / wos
	kaos := (wos * tvos) / 8760
	kpos := (1.2 * kpmax) / 8760
	wdk := 2 * wos * (kaos + kpos)
	wds := wdk + 0.02
	return wds
}

func calculateMZper(w, tv, kp, pm, tm, z_per_a, z_per_p float64) float64 {
	mw_neda := w * tv * pm * tm
	mw_nedp := kp * pm * tm
	mz_per := z_per_a*mw_neda + z_per_p*mw_nedp
	return mz_per
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64) // Ignore error for simplicity in this example
	return v
}
