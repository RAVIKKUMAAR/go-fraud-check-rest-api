// @Title  ccf_os_api_svr_rc.go
// @Description 
//  Expose REST API to front-end CICS applicaiton
//  Recevie a new transaction as input
//  Combind with the history data stored in Redis
//  Call TFS ccf model get predict values
//  Return this predict to REST API caller
//  For easy testing, input can be only a index, API can find the transaction body in test_220_100k_os.csv
//  REST API port 8080
// @Depend 
//  test_220_100k_os.csv
//  mapper.so
//  Redis Server
//  Tensorflow serving Server
// @Author  Liu Tie
// @Update  Liu Tie  2021/08/26  Sepreated from mapper for easy reading 
package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"strconv"
	"time"
	"github.com/gomodule/redigo/redis"
	"strings"
	"bytes"
	"encoding/csv"
	"os"
	"bufio"
	"io"
	"log"
	"plugin"
)

//REST API Input JSON 
type Tran struct {
	Tx_index  int `json:"tx_index"`
}
//TFS Return Predict JSON
type Predict struct {
	Predictions [7][1][1]float64 `json:"predictions"`
}

//Create Test Case Dict
//Load test case from test_220_100k_os.csv
//REST API can get pick test case based on the input tx_index
var test_case_dict = make(map[int]string)
func Create_test_case_dict(csv_row_dict map[int]string) {
	csvFile, _ := os.Open("test_220_100k_os.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	defer csvFile.Close()
	//Remove Header
	line, error := reader.Read()
	//Dict of Test Case data
	//csv_row_dict := make(map[int]string)
	//Read the CSV Body
	for {
		line, error = reader.Read()
		//
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		//
		type Input_Data struct {
			//Empty='missing_value'
			Merchant_State string
			//Empty=0.0
			Zip float64
			Merchant_Name string
			Merchant_City string
			MCC int
			Use_Chip string
			//Empty='missing_value'
			Errors string
			Year int
			Month int
			Day int
			Time string
			Amount string
		}
		//
		var csv_row Input_Data
		csv_row_index,_ := 	strconv.Atoi(line[0])
		csv_row.Year,_ 	= 	strconv.Atoi(line[3])
		csv_row.Month,_ = 	strconv.Atoi(line[4])
		csv_row.Day,_ 	= 	strconv.Atoi(line[5])
		csv_row.Time 	= 	line[6]
		csv_row.Amount 	= 	line[7]
		csv_row.Use_Chip = 	line[8]
		csv_row.Merchant_Name 	= 	line[9]
		csv_row.Merchant_City 	= 	line[10]
		//Empty='missing_value'
		if line[11] == "" {
			csv_row.Merchant_State = "missing_value"
		}else{
			csv_row.Merchant_State = line[11]
		}
		//Empty=0.0
		if line[12] == "" {
			csv_row.Zip	= 0.0
		}else{
			csv_row.Zip,_ = strconv.ParseFloat(line[12], 64)
		}
		csv_row.MCC,_ 	= 	strconv.Atoi(line[13])
		//Empty='missing_value'
		if line[14] == "" {
			csv_row.Errors	= "missing_value"
		}else{
			csv_row.Errors = line[14]
		}
		//
		//fmt.Println(csv_row)
		//
		csv_row_json_byte, _ := json.Marshal(csv_row)
		//fmt.Println(csv_row_index)
		//fmt.Println(string(csv_row_json))
		//
		csv_row_dict[csv_row_index] = string(csv_row_json_byte)
		//
		// fmt.Println(csv_row_dict[csv_row_index])
		//
		// Sample JSON Output
		// {
		// 	"Merchant_State": "missing_value",
		// 	"Zip": 0,
		// 	"Merchant_Name": "-7421093378627544099",
		// 	"Merchant_City": " ONLINE",
		// 	"MCC": 5311,
		// 	"Use_Chip": "Online Transaction",
		// 	"Errors": "missing_value",
		// 	"Year": 2002,
		// 	"Month": 12,
		// 	"Day": 10,
		// 	"Time": "06:10",
		// 	"Amount": "$128.41"
		// }
		//
		//break
		//
	}
	//
	fmt.Println(len(csv_row_dict))
	//
}

//Create Mapper funcation
var Mapper_Gen_fun = Create_Mapper()
func Create_Mapper() func(string) string {
	//
	mapper, err := plugin.Open("mapper.so")
    if err != nil {
        panic(err)
    }
	
	//func ()
	Mapper_Init, err := mapper.Lookup("Mapper_Init")
	if err != nil {
		panic(err)
	}
	//
	//func Mapper_Gen(input_data_json string) string
	Mapper_Gen, err := mapper.Lookup("Mapper_Gen")
	if err != nil {
		panic(err)
	}
	//Mapper_Init_fun()
	Mapper_Init.(func())()
	//
	return Mapper_Gen.(func(string) string)
	//
}

//Create Redis Pool
var pool = newPool()
func newPool() *redis.Pool {
	//
	REDISADD 	= os.Getenv("REDISADD")
	fmt.Println("Redis URL=>",REDISADD)
	//
	redis_pool := &redis.Pool{
		MaxIdle: 80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			//c, err := redis.Dial("tcp", "9.47.86.127:6379")
			c, err := redis.Dial("tcp", REDISADD)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	//
	client_0 := redis_pool.Get()
	client_1 := redis_pool.Get()
	client_2 := redis_pool.Get()
	client_3 := redis_pool.Get()
	//
	client_0.Close()
	client_1.Close()
	client_2.Close()
	client_3.Close()
	//
	return redis_pool
}

//Call TFS
func Call_TFS(input_test_case_json string) Predict{
	//
	TFSADD	= os.Getenv("TFSADD")
	url   := fmt.Sprintf("http://%v/v1/models/model:predict",TFSADD)
	//
	//url := "http://localhost:8501/v1/models/model:predict"
    fmt.Println("URL:>", url)
    // var jsonStr = []byte(input_test_case_json)
	// post_req,_ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    // post_req.Header.Set("Content-Type", "application/json")
    // http_client := &http.Client{}
	//start := time.Now()
    // post_resp,post_err := http_client.Do(post_req)
    // if post_err != nil {
    //     panic(post_err)
    // }
	
    // defer post_resp.Body.Close()
	//
	client := &http.Client{}
    post_resp, err := client.Post(url,  "application/json", bytes.NewBuffer([]byte(input_test_case_json)))
    if err != nil {
        panic(err)
    }
    defer post_resp.Body.Close()
	//
    body,_ := ioutil.ReadAll(post_resp.Body)
	predict_str := string(body)
    //fmt.Print(predict_str,"\n")
	//
	var predict_ret Predict
	json.Unmarshal([]byte(predict_str), &predict_ret)
	fmt.Println(predict_ret)
	//
	return predict_ret
}

//Read and Save Env Var
var REDISADD string //Reids Server IP:Port
var TFSADD string   //Tensorflow Server IP:Port

//
func main() {
	//----------------------------------------
	// 1 Creat Test case Dict 
	//----------------------------------------
	Create_test_case_dict(test_case_dict)

	//----------------------------------------
	// 2 Handle REST API Request 
	//----------------------------------------
	CCFHanlder := http.HandlerFunc(CCF)
	http.Handle("/ccf_inference", CCFHanlder)
	fmt.Println("Backedn Server Started ready for POST URL http://127.0.0.1:8080/ccf_inference")
	http.ListenAndServe(":8080", nil)
}

//Handle REST API Request
func CCF(w http.ResponseWriter, r *http.Request) {
	//----------------------------------------------------
	// 2.1 Read JSON parameter from REST API request
	//----------------------------------------------------
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		errorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var e Tran
	//
	decoder := json.NewDecoder(r.Body)
	//
	err := decoder.Decode(&e)
	//
	if err != nil {
		//
		errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		//
		return
	}
	fmt.Println(e)
	//----------------------------------------------------
	// 2.2 Get Test Case and Mapping
	//----------------------------------------------------
	tx_index := e.Tx_index
	start_Mapper  := time.Now()
	//fmt.Println("Test Case Str =",tx_index,test_case_dict[tx_index])
	mapped_json_str := Mapper_Gen_fun(test_case_dict[tx_index])
	elapsed_Mapper := time.Since(start_Mapper)
	fmt.Println("\nMapper elapsed=",elapsed_Mapper)
	//fmt.Println("Mapped Str    =",mapped_json_str)
	//----------------------------------------------------
	// 2.3 Read JSON String DMS1~7 from redis (296.67µs)
	//----------------------------------------------------
	//
	fmt.Println(tx_index)
	client := pool.Get()
	defer client.Close()
	//
	start_Redis := time.Now()
	value, err 	:= client.Do("GET", strconv.Itoa(tx_index))
	//
	redis_str  := string(value.([]byte))
	//
	redis_str_index := strings.IndexAny( redis_str, "]], [[")
	redis_str_new_history := redis_str[redis_str_index + 5 : len(redis_str)-2]
	mapped_json_str = `{"instances": [[` + redis_str_new_history + ", [" + mapped_json_str + "]]}"
	elapsed_Redis := time.Since(start_Redis)
	fmt.Println("\nRedis elapsed=",elapsed_Redis)
	//
	if err != nil {
		panic(err)
	}
	//
	//----------------------------------------
    // 3 Call TFS REST API with the JSON
    //----------------------------------------
	start_TFS := time.Now()
	predict_ret := Call_TFS(redis_str)
	elapsed_TFS := time.Since(start_TFS)
	fmt.Println("\nTFS elapsed=",elapsed_TFS)
	//
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := make(map[string]string)
	//
	resp["TX       Index"] = strconv.Itoa(e.Tx_index)
	resp["Predict  Value"] = strconv.FormatFloat(predict_ret.Predictions[6][0][0], 'E', -1, 64)
	resp["Elapsed Mapper"] = fmt.Sprintf("%v",elapsed_Mapper)
	resp["Elapsed  Redis"] = fmt.Sprintf("%v",elapsed_Redis)
	resp["Elapsed    TFS"] = fmt.Sprintf("%v",elapsed_TFS)
	//
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
	//
	//
	return
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}
