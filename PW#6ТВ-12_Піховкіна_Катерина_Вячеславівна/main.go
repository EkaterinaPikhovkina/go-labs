package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// ResultObject для передачі результатів розрахунку
type ResultObject struct {
	GroupUtilizationRate       float64
	EffectiveNumberOfUnits     float64
	EstimatedActivePowerFactor float64
	EstimatedActiveLoad        float64
	EstimatedReactiveLoad      float64
	FullPower                  float64
	EstimatedGroupCurrent      float64
}

// seq генерує послідовність цілих чисел від 0 до n-1.
func seq(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

// CalculateHandler обробляє головну сторінку з формою
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Створюємо набір функцій для шаблону, включаючи "seq"
	funcMap := template.FuncMap{
		"seq": seq,
	}

	tmpl, err := template.New("calculate.html").Funcs(funcMap).ParseFiles("calculate.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Початкові значення для таблиці
	data := struct {
		NominalEfficiency   []float64
		LoadPowerFactor     []float64
		LoadVoltage         []float64
		Units               []float64
		NominalPower        []float64
		UtilizationFactor   []float64
		ReactivePowerFactor []float64
		Rows                int // Додаємо кількість рядків для шаблону
	}{
		NominalEfficiency:   []float64{0.92, 0.92, 0.92, 0.92, 0.92, 0.92, 0.92, 0.92},
		LoadPowerFactor:     []float64{0.9, 0.9, 0.9, 0.9, 0.9, 0.9, 0.9, 0.9},
		LoadVoltage:         []float64{0.38, 0.38, 0.38, 0.38, 0.38, 0.38, 0.38, 0.38},
		Units:               []float64{4.0, 2.0, 4.0, 1.0, 1.0, 1.0, 1.0, 2.0},
		NominalPower:        []float64{20.0, 14.0, 42.0, 36.0, 20.0, 40.0, 32.0, 20.0},
		UtilizationFactor:   []float64{0.15, 0.12, 0.15, 0.3, 0.5, 0.2, 0.2, 0.65},
		ReactivePowerFactor: []float64{1.33, 1.0, 1.33, 1.52, 0.75, 1.0, 1.0, 0.75},
		Rows:                8, // Передаємо кількість рядків
	}

	tmpl.Execute(w, data)
}

// CalculateResultHandler обробляє відправку форми та розрахунки
func CalculateResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	rows := 8

	nominalEfficiency := make([]float64, rows)
	loadPowerFactor := make([]float64, rows)
	loadVoltage := make([]float64, rows)
	units := make([]float64, rows)
	nominalPower := make([]float64, rows)
	utilizationFactor := make([]float64, rows)
	reactivePowerFactor := make([]float64, rows)

	for i := 0; i < rows; i++ {
		nominalEfficiency[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("nominalEfficiency[%d]", i)), 64)
		loadPowerFactor[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("loadPowerFactor[%d]", i)), 64)
		loadVoltage[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("loadVoltage[%d]", i)), 64)
		units[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("units[%d]", i)), 64)
		nominalPower[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("nominalPower[%d]", i)), 64)
		utilizationFactor[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("utilizationFactor[%d]", i)), 64)
		reactivePowerFactor[i], _ = strconv.ParseFloat(r.FormValue(fmt.Sprintf("reactivePowerFactor[%d]", i)), 64)
	}

	calculateValues := func() ResultObject {
		n := units
		pN := nominalPower
		kV := utilizationFactor
		tgPhi := reactivePowerFactor

		totalPower := 0.0
		for i := 0; i < len(n); i++ {
			totalPower += n[i] * pN[i]
		}

		totalWeightedPower := 0.0
		for i := 0; i < len(n); i++ {
			totalWeightedPower += n[i] * pN[i] * kV[i]
		}
		groupUtilizationRate := totalWeightedPower / totalPower

		powerSquared := math.Pow(totalPower, 2)
		sumPowerSquared := 0.0
		for i := 0; i < len(n); i++ {
			sumPowerSquared += math.Pow(n[i]*pN[i], 2)
		}
		effectiveNumberOfUnits := powerSquared / sumPowerSquared

		estimatedActivePowerFactor := 1.25
		estimatedActiveLoad := estimatedActivePowerFactor * totalWeightedPower
		estimatedReactiveLoad := groupUtilizationRate * totalPower * tgPhi[0] // Assuming tgPhi[0] is representative
		estimatedGroupCurrent := estimatedActiveLoad / loadVoltage[0]         // Assuming loadVoltage[0] is representative

		return ResultObject{
			GroupUtilizationRate:       groupUtilizationRate,
			EffectiveNumberOfUnits:     effectiveNumberOfUnits,
			EstimatedActivePowerFactor: estimatedActivePowerFactor,
			EstimatedActiveLoad:        estimatedActiveLoad,
			EstimatedReactiveLoad:      estimatedReactiveLoad,
			FullPower:                  math.Sqrt(math.Pow(estimatedActiveLoad, 2) + math.Pow(estimatedReactiveLoad, 2)),
			EstimatedGroupCurrent:      estimatedGroupCurrent,
		}
	}

	results := calculateValues()

	resultTmpl, err := template.ParseFiles("result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultTmpl.Execute(w, results)
}

func main() {
	http.HandleFunc("/", CalculateHandler)
	http.HandleFunc("/calculate", CalculateResultHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
