package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func get_metar(icao_code string) string {
	url := "https://metar.vatsim.net/" + icao_code
	response, _ := http.Get(url)
	defer response.Body.Close()
	byteArray, _ := ioutil.ReadAll(response.Body)
	metar := string(byteArray)
	return (metar)
}

func find_wind(splited []string) string {
	if splited[2] != "AUTO" {
		return (splited[2])
	} else {
		return (splited[3])
	}
}

func find_vis(splited []string) string {
	if splited[2] != "AUTO" && len(splited[3]) <= 4 {
		return (splited[3])
	} else {
		if splited[3] == "CAVOK" {
			return ("CVOK")
		} else {
			return (splited[4])
		}
	}
}

func find_qnh(splited []string) string {
	for i := 3; i < len(splited); i++ {
		if strings.HasPrefix(splited[i], "Q") && len(splited[i]) == 5 {
			return (splited[i][1:])
		}
	}
	return ("----")
}

func find_alt(splited []string) string {
	for i := 3; i < len(splited); i++ {
		if strings.HasPrefix(splited[i], "A") && len(splited[i]) == 5 {
			return (splited[i][1:])
		}
	}
	for i := 3; i < len(splited); i++ {
		if strings.HasPrefix(splited[i], "Q") && len(splited[i]) == 5 {
			qnh := splited[i][1:]
			qnh_float, _ := strconv.ParseFloat(qnh, 64)
			return (strconv.FormatFloat(qnh_float/0.3386, 'f', 0, 64))
		}
	}
	return ("----")
}

func find_temp(splited []string) string {
	for i := 3; i < len(splited); i++ {
		if (strings.HasPrefix(splited[i], "Q") || strings.HasPrefix(splited[i], "A")) && len(splited[i]) == 5 {
			index := strings.Index(splited[i-1], "/")
			if strings.HasPrefix(splited[i-1][:index], "M") {
				return (splited[i-1][:index])
			} else {
				return (" " + splited[i-1][:index])
			}

		}
	}
	return ("--")
}

func is_imc(splited []string) (bool, error) {
	vis, err := strconv.Atoi(find_vis(splited))
	if err != nil {
		return false, errors.New("この空港はIMC/VMCの判定に対応していません")
	}
	if vis < 5000 { // 視程が5KM未満の場合はIMC
		return true, nil
	}
	for i := 3; i < len(splited); i++ { // 雲底が1000ft未満の場合はIMC
		if strings.HasPrefix(splited[i], "BKN") || strings.HasPrefix(splited[i], "OVC") {
			ceiling, err := strconv.Atoi(splited[i][3:])
			if err != nil {
				return false, errors.New("この空港はIMC/VMCの判定に対応していません")
			}
			if ceiling < 10 {
				return true, nil
			}
		}
	}
	return false, nil
}

func imcvmc(splited []string) string {
	isimc, err := is_imc(splited)
	if err != nil {
		return ("Err")
	}
	if isimc {
		return ("I")
	} else {
		return ("V")
	}
}

func main() {
	argv := os.Args
	argc := len(argv)
	if argc == 1 {
		fmt.Fprintln(os.Stderr, "usage -> wx RJTT RJFF ...")
		os.Exit(1)
	}
	for {
		for i := 0; i < argc-1; i++ {
			metar := get_metar(argv[i+1])
			//fmt.Println("\n" + metar)
			splited := strings.Split(metar, " ")
			airport_code := splited[0]
			metar_time := splited[1]
			wind := find_wind(splited)
			vis := find_vis(splited)
			qnh := find_qnh(splited)
			tmp := find_temp(splited)
			alt := find_alt(splited)
			cond := imcvmc(splited)
			is_imc(splited)
			fmt.Printf("%s %s %s\t%s %s %s/%s %s\n",
				airport_code, metar_time, wind, vis, tmp, qnh, alt, cond)
		}
		time.Sleep(5 * time.Minute)
		fmt.Print("\033[" + strconv.Itoa(argc-1) + "A")
	}
}
