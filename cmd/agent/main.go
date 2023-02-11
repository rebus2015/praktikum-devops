package main

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"runtime"
	"time"
	//"net/http"
	//"net/url"
	//"os"
)

type gauge float64
type counter int64

type metrics struct {
	alloc       gauge
	buckhashsys gauge
	pollcount   counter
	randomvalue gauge
}

func (m *metrics) Update() {

	var memoryStat runtime.MemStats

	runtime.ReadMemStats(&memoryStat)

	(*m).alloc = gauge(memoryStat.Alloc)
	(*m).buckhashsys = gauge(memoryStat.BuckHashSys)
	(*m).pollcount++
	(*m).randomvalue = gauge(rand.Float64())
}

func Uint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf) // Always succeeds, no need to check error
	return binary.LittleEndian.Uint64(buf)
}

const pollinterval time.Duration = 2 * time.Second
const reportintelval time.Duration = 10 * time.Second
//const endpoint string := "http://localhost:8080/"

func main() {

	m := metrics{}
	updticker := time.NewTicker(pollinterval)
	sndticker := time.NewTicker(reportintelval)
   
	//client := &http.Client{}
	
	defer updticker.Stop()
	defer sndticker.Stop()

	

	for {
		select {
		case t := <-updticker.C:
			{
				m.Update()
				fmt.Printf("%v %v", t, m)
				fmt.Println("")
			}
		case s := <-sndticker.C:
			{

				fmt.Printf("%v Send Statistic", s)
				fmt.Println("")
			}
		}
	}


	// адрес сервиса (как его писать, расскажем в следующем уроке)
	

	// приглашение в консоли

	// // конструируем HTTP-клиент
	// client := &http.Client{}
	// // конструируем запрос
	// // запрос методом POST должен, кроме заголовков, содержать тело
	// // тело должно быть источником потокового чтения io.Reader
	// // в большинстве случаев отлично подходит bytes.Buffer
	// request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(data.Encode()))
	// if err != nil {
	//     fmt.Println(err)
	//     os.Exit(1)
	// }
	// // в заголовках запроса сообщаем, что данные кодированы стандартной URL-схемой
	// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	// // отправляем запрос и получаем ответ
	// response, err := client.Do(request)
	// if err != nil {
	//     fmt.Println(err)
	//     os.Exit(1)
	// }
	// // печатаем код ответа
	// fmt.Println("Статус-код ", response.Status)
	// defer response.Body.Close()
	// // читаем поток из тела ответа
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	//     fmt.Println(err)
	//     os.Exit(1)
	// }
	// // и печатаем его
	// fmt.Println(string(body))
}
