package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Структура для передачі даних в шаблон FuelScreen
type FuelData struct {
	Title         string
	DefaultValues map[string]string
	Result        string
	Error         string
}

// Структура для передачі даних в шаблон GasScreen
type GasData struct {
	Text string
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/fuel/coal", coalHandler)
	http.HandleFunc("/fuel/mazut", mazutHandler)
	http.HandleFunc("/fuel/gas", gasHandler)
	http.HandleFunc("/calculate", calculateHandler)

	fmt.Println("Запуск веб-сервера на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Калькулятор викидів</title>
</head>
<body>
    <h1>Оберіть тип палива</h1>
    <ul>
        <li><a href="/fuel/coal">Вугілля</a></li>
        <li><a href="/fuel/mazut">Мазут</a></li>
        <li><a href="/fuel/gas">Природний газ</a></li>
    </ul>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func coalHandler(w http.ResponseWriter, r *http.Request) {
	fuelScreenHandler(w, r, "Вугілля", map[string]string{
		"Q_ir":      "20.47",
		"a_vin":     "0.8",
		"Ar":        "25.20",
		"Gamma_vin": "1.5",
		"n_зу":      "0.985",
		"k_tvS":     "0",
		"B":         "1096363",
	})
}

func mazutHandler(w http.ResponseWriter, r *http.Request) {
	fuelScreenHandler(w, r, "Мазут", map[string]string{
		"Q_ir":      "39.48",
		"a_vin":     "1",
		"Ar":        "0.15",
		"Gamma_vin": "0",
		"n_зу":      "0.985",
		"k_tvS":     "0",
		"B":         "70945",
	})
}

func gasHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("gasScreen").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>Природний газ</title>
</head>
<body>
    <h1>Природний газ</h1>
    <p>При спалюванні природного газу тверді частинки відсутні.
    Отже, показник емісії твердих частинок k_тв при спалюванні природного газу дорівнює нулю,
    і валовий викид твердих частинок при спалюванні природного газу також буде нульовим:
    E_тв=0т</p>
    <p><a href="/">Назад до вибору палива</a></p>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := GasData{
		Text: `При спалюванні природного газу тверді частинки відсутні.
		Отже, показник емісії твердих частинок k_тв при спалюванні природного газу дорівнює нулю,
		і валовий викид твердих частинок при спалюванні природного газу також буде нульовим:
		E_тв=0т`,
	}
	tmpl.Execute(w, data)
}

func fuelScreenHandler(w http.ResponseWriter, r *http.Request, title string, defaultValues map[string]string) {
	tmpl, err := template.New("fuelScreen").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Title}}</h1>
	{{if .Error}}
		<p style="color: red;">{{.Error}}</p>
	{{end}}
    <form action="/calculate" method="post">
        <input type="hidden" name="fuelType" value="{{.Title}}">
        <label for="qir">Q_ir:</label>
        <input type="text" id="qir" name="Q_ir" value="{{.DefaultValues.Q_ir}}"><br><br>

        <label for="avin">a_вин:</label>
        <input type="text" id="avin" name="a_вин" value="{{.DefaultValues.a_вин}}"><br><br>

        <label for="ar">Ar:</label>
        <input type="text" id="ar" name="Ar" value="{{.DefaultValues.Ar}}"><br><br>

        <label for="gvin">Gamma_vin:</label>
        <input type="text" id="gvin" name="Gamma_vin" value="{{.DefaultValues.Gamma_vin}}"><br><br>

        <label for="nzu">n_зу:</label>
        <input type="text" id="nzu" name="n_зу" value="{{.DefaultValues.n_зу}}"><br><br>

        <label for="ktvs">k_твS:</label>
        <input type="text" id="ktvs" name="k_твS" value="{{.DefaultValues.k_твS}}"><br><br>

        <label for="b">B:</label>
        <input type="text" id="b" name="B" value="{{.DefaultValues.B}}"><br><br>

        <button type="submit">Розрахувати</button>
    </form>

    {{if .Result}}
    <div>
        <h2>Результат розрахунку:</h2>
        <pre>{{.Result}}</pre>
    </div>
    {{end}}
	<p><a href="/">Назад до вибору палива</a></p>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := FuelData{
		Title:         title,
		DefaultValues: defaultValues,
	}
	tmpl.Execute(w, data)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Помилка обробки форми: "+err.Error(), http.StatusBadRequest)
		return
	}

	fuelType := r.FormValue("fuelType")
	qIr := r.FormValue("Q_ir")
	aVin := r.FormValue("a_вин")
	ar := r.FormValue("Ar")
	gVin := r.FormValue("Gamma_vin")
	nZu := r.FormValue("n_зу")
	ktvS := r.FormValue("k_твS")
	b := r.FormValue("B")

	ktv := calculateKtv(qIr, aVin, ar, gVin, nZu, ktvS)
	etv := calculateEtv(qIr, ktv, b)

	resultText := fmt.Sprintf(
		"Показник емісії твердих частинок при спалюванні становитиме: k_тв = %.2f г/ГДж;\n"+
			"Валовий викид при спалюванні становитиме: E_тв = %.2f т;", ktv, etv,
	)

	var data FuelData
	if fuelType == "Вугілля" {
		data = FuelData{
			Title: "Вугілля",
			DefaultValues: map[string]string{
				"Q_ir":      qIr,
				"a_вин":     aVin,
				"Ar":        ar,
				"Gamma_vin": gVin,
				"n_зу":      nZu,
				"k_твS":     ktvS,
				"B":         b,
			},
			Result: resultText,
		}
	} else if fuelType == "Мазут" {
		data = FuelData{
			Title: "Мазут",
			DefaultValues: map[string]string{
				"Q_ir":      qIr,
				"a_вин":     aVin,
				"Ar":        ar,
				"Gamma_vin": gVin,
				"n_зу":      nZu,
				"k_твS":     ktvS,
				"B":         b,
			},
			Result: resultText,
		}
	}

	tmpl, err := template.New("fuelScreen").Parse(`
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
</head>
<body>
    <h1>{{.Title}}</h1>
	{{if .Error}}
		<p style="color: red;">{{.Error}}</p>
	{{end}}
    <form action="/calculate" method="post">
        <input type="hidden" name="fuelType" value="{{.Title}}">
        <label for="qir">Q_ir:</label>
        <input type="text" id="qir" name="Q_ir" value="{{.DefaultValues.Q_ir}}"><br><br>

        <label for="avin">a_вин:</label>
        <input type="text" id="avin" name="a_вин" value="{{.DefaultValues.a_вин}}"><br><br>

        <label for="ar">Ar:</label>
        <input type="text" id="ar" name="Ar" value="{{.DefaultValues.Ar}}"><br><br>

        <label for="gvin">Gamma_vin:</label>
        <input type="text" id="gvin" name="Gamma_vin" value="{{.DefaultValues.Gamma_vin}}"><br><br>

        <label for="nzu">n_зу:</label>
        <input type="text" id="nzu" name="n_зу" value="{{.DefaultValues.n_зу}}"><br><br>

        <label for="ktvs">k_твS:</label>
        <input type="text" id="ktvs" name="k_твS" value="{{.DefaultValues.k_твS}}"><br><br>

        <label for="b">B:</label>
        <input type="text" id="b" name="B" value="{{.DefaultValues.B}}"><br><br>

        <button type="submit">Розрахувати</button>
    </form>

    {{if .Result}}
    <div>
        <h2>Результат розрахунку:</h2>
        <pre>{{.Result}}</pre>
    </div>
    {{end}}
	<p><a href="/">Назад до вибору палива</a></p>
</body>
</html>
`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func calculateKtv(
	qIrStr string,
	aVin string,
	arStr string,
	gVinStr string,
	nZu string,
	ktvS string,
) float64 {
	qRi, err := strconv.ParseFloat(strings.ReplaceAll(qIrStr, ",", "."), 64)
	if err != nil {
		return 0.0
	}
	aVinFloat, err := strconv.ParseFloat(strings.ReplaceAll(aVin, ",", "."), 64)
	if err != nil {
		return 0.0
	}
	arFloat, err := strconv.ParseFloat(strings.ReplaceAll(arStr, ",", "."), 64)
	if err != nil {
		return 0.0
	}
	gVinFloat, err := strconv.ParseFloat(strings.ReplaceAll(gVinStr, ",", "."), 64)
	if err != nil {
		return 0.0
	}
	nZuFloat, err := strconv.ParseFloat(strings.ReplaceAll(nZu, ",", "."), 64)
	if err != nil {
		return 0.0
	}
	ktvSFloat, err := strconv.ParseFloat(strings.ReplaceAll(ktvS, ",", "."), 64)
	if err != nil {
		return 0.0
	}

	return (1000000.0/qRi)*(aVinFloat*arFloat/(100.0-gVinFloat))*(1.0-nZuFloat) + ktvSFloat
}

func calculateEtv(qIrStr string, ktv float64, b string) float64 {
	qRiFloat, err := strconv.ParseFloat(strings.ReplaceAll(qIrStr, ",", "."), 64)
	if err != nil {
		return 0.0
	}
	bFloat, err := strconv.ParseFloat(strings.ReplaceAll(b, ",", "."), 64)
	if err != nil {
		return 0.0
	}

	return 1e-6 * ktv * qRiFloat * bFloat
}
