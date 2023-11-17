package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	argv := os.Args
	argc := len(argv)
	for {
		for i := 0; i < argc-1; i++ {
			url := "https://metar.vatsim.net/" + argv[i+1]
			response, _ := http.Get(url)
			defer response.Body.Close()
			byteArray, _ := ioutil.ReadAll(response.Body)
			result := string(byteArray)
			airport_code := result[0:4]
			metar_time := result[7:12]
			wind := strings.Split(result[(strings.Index(result[4:], "KT")-4):(strings.Index(result[4:], "KT")+6)], " ")[len(strings.Split(result[(strings.Index(result[4:], "KT")-4):(strings.Index(result[4:], "KT")+6)], " "))-1]
			var vis string
			//fmt.Println(result)
			//fmt.Println(result[(strings.Index(result[4:], "KT") + 10) : (strings.Index(result[4:], "KT")+10)+1])
			if result[(strings.Index(result[4:], "KT")+10):(strings.Index(result[4:], "KT")+10)+1] != "V" {
				vis = result[strings.Index(result[4:], "KT")+7 : strings.Index(result[4:], "KT")+11]
			} else {
				vis = result[strings.Index(result[4:], "KT")+7+8 : strings.Index(result[4:], "KT")+11+8]
			}
			qnh := result[strings.Index(result[4:], "Q")+5 : strings.Index(result[4:], "Q")+9]
			alt_f, _ := strconv.ParseFloat(qnh, 64)
			alt_f = alt_f / 0.3386
			alt := fmt.Sprintf("%.0f", alt_f)
			fmt.Print(airport_code + " ")
			fmt.Print(metar_time + " ")
			fmt.Print(wind + "\t")
			fmt.Print(vis + " ")
			fmt.Print(qnh + "/" + alt + "\t")
			print("\n")
		}
		time.Sleep(5 * time.Minute)
		fmt.Print("\033[" + strconv.Itoa(argc-1) + "A")
	}
}
