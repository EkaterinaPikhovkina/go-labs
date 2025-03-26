package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type InputData struct {
	PC     string
	Sigma  string
	B      string
	Delta  string
	Result string
}

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/calculate", handleCalculate)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("form.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	defaults := map[string]string{
		"PC":    "5",
		"Sigma": "0.25",
		"B":     "7",
		"Delta": "5",
	}

	data := InputData{
		PC:     defaults["PC"],
		Sigma:  defaults["Sigma"],
		B:      defaults["B"],
		Delta:  defaults["Delta"],
		Result: "",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println(err)
	}
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pcStr := r.FormValue("pc")
	sigmaStr := r.FormValue("sigma")
	bStr := r.FormValue("b")
	deltaStr := r.FormValue("delta")

	profit := calculateProfit(pcStr, sigmaStr, bStr, deltaStr)
	result := fmt.Sprintf("Відповідь: Прибуток П = %.2f тис. грн.", profit)

	tmpl, err := template.ParseFiles("form.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := InputData{
		PC:     pcStr,
		Sigma:  sigmaStr,
		B:      bStr,
		Delta:  deltaStr,
		Result: result,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		log.Println(err)
	}
}

func calculateProfit(pcStr, sigmaStr, bStr, deltaStr string) float64 {
	pc, err := strconv.ParseFloat(pcStr, 64)
	if err != nil {
		return 0.0
	}
	sigma, err := strconv.ParseFloat(sigmaStr, 64)
	if err != nil {
		return 0.0
	}
	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return 0.0
	}
	delta, err := strconv.ParseFloat(deltaStr, 64)
	if err != nil {
		return 0.0
	}

	pMin := pc * (1 - delta/100)
	pMax := pc * (1 + delta/100)

	deltaW := integrateProbabilityDensity(pc, sigma, pMin, pMax, 0.01)

	w1 := pc * 24 * deltaW
	p := w1 * b

	w2 := pc * 24 * (1 - deltaW)
	penalty := w2 * b

	profit := p - penalty

	return profit
}

func probabilityDensity(p, pC, sigma float64) float64 {
	coefficient := 1 / (sigma * math.Sqrt(2*math.Pi))
	exponent := -math.Pow(p-pC, 2) / (2 * math.Pow(sigma, 2))
	return coefficient * math.Exp(exponent)
}

func integrateProbabilityDensity(pC, sigma, lowerBound, upperBound float64, stepSize float64) float64 {
	if stepSize == 0 {
		stepSize = 0.01
	}
	sum := 0.0
	p := lowerBound
	for p < upperBound {
		sum += probabilityDensity(p, pC, sigma) * stepSize
		p += stepSize
	}
	return sum
}
