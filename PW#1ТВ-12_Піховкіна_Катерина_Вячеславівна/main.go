package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Структура для даних калькулятора 1
type Calculator1Data struct {
	Result string
	HP     string
	CP     string
	SP     string
	NP     string
	OP     string
	WP     string
	AP     string
}

// Структура для даних калькулятора 2
type Calculator2Data struct {
	Result string
	CG     string
	HG     string
	OG     string
	SG     string
	AG     string
	WG     string
	VG     string
	Qidaf  string
}

func main() {
	// Обробник для головної сторінки з вибором завдання
	http.HandleFunc("/", indexHandler)

	// Обробники для калькуляторів
	http.HandleFunc("/task1", calculator1Handler)
	http.HandleFunc("/task1/calculate", calculate1Handler)
	http.HandleFunc("/task2", calculator2Handler)
	http.HandleFunc("/task2/calculate", calculate2Handler)

	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Обробник для головної сторінки (index.html)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

// Обробник для відображення форми калькулятора 1 (calculator1.html)
func calculator1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("calculator1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := Calculator1Data{} // Початкові пусті дані
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

// Обробник для виконання обчислень калькулятора 1
func calculate1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не дозволений", http.StatusMethodNotAllowed)
		return
	}

	// Отримуємо значення з форми
	hpStr := r.FormValue("hp")
	cpStr := r.FormValue("cp")
	spStr := r.FormValue("sp")
	npStr := r.FormValue("np")
	opStr := r.FormValue("op")
	wpStr := r.FormValue("wp")
	apStr := r.FormValue("ap")

	// Перетворюємо рядки в числа
	hp, err := strconv.ParseFloat(hpStr, 64)
	cp, err := strconv.ParseFloat(cpStr, 64)
	sp, err := strconv.ParseFloat(spStr, 64)
	np, err := strconv.ParseFloat(npStr, 64)
	op, err := strconv.ParseFloat(opStr, 64)
	wp, err := strconv.ParseFloat(wpStr, 64)
	ap, err := strconv.ParseFloat(apStr, 64)

	if err != nil {
		renderError1(w, "Невірний формат числа", hpStr, cpStr, spStr, npStr, opStr, wpStr, apStr)
		return
	}

	// 1. Розрахунок коефіцієнтів переходу
	krs := 100 / (100 - wp)
	krg := 100 / (100 - wp - ap)

	// 2. Розрахунок складу сухої маси палива
	hS := hp * krs
	cS := cp * krs
	sS := sp * krs
	nS := np * krs
	oS := op * krs
	aS := ap * krs

	// 3. Розрахунок складу горючої маси палива
	hG := hp * krg
	cG := cp * krg
	sG := sp * krg
	nG := np * krg
	oG := op * krg

	// 4. Розрахунок нижчої теплоти згорання для робочої маси
	qrp := 339*cp + 1030*hp - 108.8*(op-sp) - 25*wp

	// 5. Перерахунок теплоти на суху масу
	qch := (qrp + 0.025*wp) * (100 / (100 - wp))

	// 6. Перерахунок теплоти на горючу масу
	qgh := (qrp + 0.025*wp) * (100 / (100 - wp - ap))

	// Форматуємо результат у рядок з округленням до 2 знаків
	resultText := fmt.Sprintf(`
		Krs: %.2f
		Krg: %.2f
		HS: %.2f
		CS: %.2f
		SS: %.2f
		NS: %.2f
		OS: %.2f
		AS: %.2f
		HG: %.2f
		CG: %.2f
		SG: %.2f
		NG: %.2f
		OG: %.2f
		Qrp (МДж/кг): %.2f
		Qch (МДж/кг): %.2f
		Qgh (МДж/кг): %.2f
	`, krs, krg, hS, cS, sS, nS, oS, aS, hG, cG, sG, nG, oG, qrp/1000, qch/1000, qgh/1000)

	// Створюємо шаблон та передаємо дані з результатом
	tmpl, err := template.ParseFiles("calculator1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Calculator1Data{
		Result: resultText,
		HP:     hpStr,
		CP:     cpStr,
		SP:     spStr,
		NP:     npStr,
		OP:     opStr,
		WP:     wpStr,
		AP:     apStr,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

// Функція для відображення помилок калькулятора 1
func renderError1(w http.ResponseWriter, errorMessage, hp, cp, sp, np, op, wp, ap string) {
	tmpl, err := template.ParseFiles("calculator1.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Calculator1Data{
		Result: errorMessage,
		HP:     hp,
		CP:     cp,
		SP:     sp,
		NP:     np,
		OP:     op,
		WP:     wp,
		AP:     ap,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

// Обробник для відображення форми калькулятора 2 (calculator2.html)
func calculator2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("calculator2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := Calculator2Data{} // Початкові пусті дані
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

// Обробник для виконання обчислень калькулятора 2
func calculate2Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не дозволений", http.StatusMethodNotAllowed)
		return
	}

	// Отримуємо значення з форми
	cgStr := r.FormValue("cg")
	hgStr := r.FormValue("hg")
	ogStr := r.FormValue("og")
	sgStr := r.FormValue("sg")
	agStr := r.FormValue("ag")
	wgStr := r.FormValue("wg")
	vgStr := r.FormValue("vg")
	qidafStr := r.FormValue("qidaf")

	// Перетворюємо рядки в числа
	cg, err := strconv.ParseFloat(cgStr, 64)
	hg, err := strconv.ParseFloat(hgStr, 64)
	og, err := strconv.ParseFloat(ogStr, 64)
	sg, err := strconv.ParseFloat(sgStr, 64)
	ag, err := strconv.ParseFloat(agStr, 64)
	wg, err := strconv.ParseFloat(wgStr, 64)
	vg, err := strconv.ParseFloat(vgStr, 64)
	qidaf, err := strconv.ParseFloat(qidafStr, 64)

	if err != nil {
		renderError2(w, "Невірний формат числа", cgStr, hgStr, ogStr, sgStr, agStr, wgStr, vgStr, qidafStr)
		return
	}

	cR := cg * (100 - wg - ag) / 100
	hR := hg * (100 - wg - ag) / 100
	oR := og * (100 - (wg / 10) - (ag / 10)) / 100
	sR := sg * (100 - wg - ag) / 100
	aR := ag * (100 - wg) / 100
	vR := vg * (100 - wg) / 100

	// Розрахунок нижчої теплоти згоряння
	qR := qidaf*(100-wg-aR)/100 - 0.025*wg

	// Форматуємо результат у рядок з округленням до 2 знаків
	resultText := fmt.Sprintf(`
		HR: %.2f
		CR: %.2f
		SR: %.2f
		OR: %.2f
		VR: %.2f
		AR: %.2f
		QR (МДж/кг): %.2f
	`, hR, cR, sR, oR, vR, aR, qR)

	// Створюємо шаблон та передаємо дані з результатом
	tmpl, err := template.ParseFiles("calculator2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Calculator2Data{
		Result: resultText,
		CG:     cgStr,
		HG:     hgStr,
		OG:     ogStr,
		SG:     sgStr,
		AG:     agStr,
		WG:     wgStr,
		VG:     vgStr,
		Qidaf:  qidafStr,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}

// Функція для відображення помилок калькулятора 2
func renderError2(w http.ResponseWriter, errorMessage, cg, hg, og, sg, ag, wg, vg, qidaf string) {
	tmpl, err := template.ParseFiles("calculator2.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Calculator2Data{
		Result: errorMessage,
		CG:     cg,
		HG:     hg,
		OG:     og,
		SG:     sg,
		AG:     ag,
		WG:     wg,
		VG:     vg,
		Qidaf:  qidaf,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		return
	}
}
