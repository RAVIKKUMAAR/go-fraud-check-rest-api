// @Title  ccf_os_mapper.go
// @Description  Mapper Implementaion
// @Depend  mappermd.go
// @Author  Liu Tie
// @Update  Liu Tie  2021/08/26  Move init logic to mappermd.go
package main

import (
	"encoding/json"
	"strconv"
	"time"
	"strings"
	"math"
)


//Used to save input JSON string of a transaction
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

// @title    Mapper_Init
// @description   Init mapper depend data structs
// @auth     Liu Tie
// @param    None
// @return   None
// @Update   Liu Tie 2021/08/26   Move this logic to the mappermd.go
func Mapper_Init(){}

// @title    Mapper_Gen
// @description   Mappint the JSON string of a transaction into a string that ccf tfs can read
// @auth     Liu Tie
// @param    input_data_json    string     JSON string of a 
// @return   string  ccf tfs requried string format
// @update
func Mapper_Gen(input_data_json string) string {
	//
	//Offset list , 9 fields,
	const offset_Merchant_State int = 0
	const offset_Zip int = 25
	const offset_Merchant_Name int = 75
	const offset_Merchant_City int = 125
	const offset_MCC int = 167
	const offset_Use_Chip int = 191
	const offset_Errors int = 194
	const offset_Year_Month_Day_Time int = 218
	const offset_Amount int = 219
	//
	//-------------------------------------
	//2 Run Time actions
	//------------------------------------
	//Load from JSON string into struct
	var input_data_refine Input_Data
	//
	json.Unmarshal([]byte(input_data_json), &input_data_refine)
	//fmt.Println(input_data_refine)
	//
    mapper_code_list := [...]float64{
        //-------------
        //@Merchant State , 3 digital , max index = 224, code len = 25
        //-------------
		//Merchant State_x0_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant State_x1_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant State_x2_0-2
        0.0, 0.0, 0.0,
        //Merchant State_x3_0
        1.0, 
        //Merchant State_x4_0
        1.0,
        //-------------
        //@Zip, 5 digitals ,  code len = 50
        //-------------
        //Zip_x0_0.0
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Zip_x1_0.0
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Zip_x2_0.0
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Zip_x3_0.0
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Zip_x4_0.0
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //-------------
        //@Merchant Name, 5 digitals ,  code len = 50
        //-------------
        //Merchant Name_x0_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant Name_x1_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant Name_x2_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant Name_x3_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant Name_x4_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //-------------
        //@Merchant City, 5 digitals ,  code len = 42
        //-------------
        //Merchant City_x0_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant City_x1_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant City_x2_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant City_x3_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //Merchant City_x4_0-1
        0.0, 0.0,
        //-------------
        //@MCC, 3 digital , max index = 224, code len = 24
        //-------------
        //MCC_x0_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //MCC_x1_0-9
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        //MCC_x2_0-1
        0.0, 0.0,
        //MCC_x3_0
        1.0, 
        //MCC_x4_0
        1.0,   
        //-------------
        //@Use Chip, 'Chip Transaction'/'Online Transaction'/'Swipe Transaction' , code len = 3
        //-------------
        //Use Chip_C/O/S
        0.0, 0.0, 0.0, 
        //-------------
        //@Errors?, code len = 24
        //-------------
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0,
        0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 
        0.0, 0.0, 0.0, 0.0,
        //-------------
        //@Year_Month_Day_Time,  code len = 1
        //-------------
        0.0,
        //-------------
        //@Amount,  code len = 1
        //-------------
        0.0}
    //
    //------------------------------------------------------------------------------------------------
    // 2.1 Merchant State
    //-------------------------------------------------------------------------------------------------
	V := Merchant_State_dict[input_data_refine.Merchant_State]
    for index := 0; index < 3; index++ {
        mapper_code_list[index*10 + (V % 10 )] = 1.0
        V = V / 10
	}
	//
	//--------------------------------------------------------------------------------------------------
    // 2.2 Zip
    //---------------------------------------------------------------------------------------------------
    V = int(input_data_refine.Zip)
    for index := 0; index < 5; index++ {
        mapper_code_list[offset_Zip + index*10 + (V % 10 )] = 1.0
        V = V / 10
	}
	//
	//--------------------------------------------------------------------------------------------------
    // 2.3 Merchant Name
    //---------------------------------------------------------------------------------------------------
    V = Merchant_Name_dict[input_data_refine.Merchant_Name]
    for index := 0; index < 5; index++ {
        mapper_code_list[offset_Merchant_Name + index*10 + (V % 10 )] = 1.0
        V = V / 10
	}
	//
	//--------------------------------------------------------------------------------------------------
    // 2.4 Merchant City
    //---------------------------------------------------------------------------------------------------
	V = Merchant_City_dict[input_data_refine.Merchant_City]
    for index := 0; index < 5; index++ {
        mapper_code_list[offset_Merchant_City + index*10 + (V % 10 )] = 1.0
		V  =  V  / 10 
	}
	//
	//--------------------------------------------------------------------------------------------------
    // 2.5 MCC
    //---------------------------------------------------------------------------------------------------
    V = MCC_dict[input_data_refine.MCC]
    for index := 0; index < 2; index++ {
        mapper_code_list[offset_MCC + index*10 + (V % 10 )] = 1.0
        V = V / 10
	}
    mapper_code_list[offset_MCC + 2*10 + (V % 10 )] = 1.0
    V = V / 10
	//
	//--------------------------------------------------------------------------------------------------
    // 2.6 Use Chip
    //---------------------------------------------------------------------------------------------------
    V = Use_Chip_dict[input_data_refine.Use_Chip]
    mapper_code_list[offset_Use_Chip + V] = 1.0
    //--------------------------------------------------------------------------------------------------
    // 2.7 Errors?
    //---------------------------------------------------------------------------------------------------
    V = Errors_dict[input_data_refine.Errors]
    mapper_code_list[offset_Errors + V] = 1.0
	//
	//--------------------------------------------------------------------------------------------------
    // 2.8 Year_Month_Day_Time
    //---------------------------------------------------------------------------------------------------
	//fmt.Println(input_data_refine.Time)
	hours,_ := strconv.Atoi((strings.Split(input_data_refine.Time,":")[0]))
	minutes,_ := strconv.Atoi((strings.Split(input_data_refine.Time,":")[1]))
	t:= time.Date(
		input_data_refine.Year, 
		time.Month(input_data_refine.Month), 
		input_data_refine.Day, 
		hours,
		minutes,
		0,0, time.UTC)
	//
	X := float64(t.UnixNano()/1e9)
    mapper_code_list[offset_Year_Month_Day_Time] = (X - Mapper_time_l) / (Mapper_time_h - Mapper_time_l)
	//
	//--------------------------------------------------------------------------------------------------
    // 2.9 Amount
    //---------------------------------------------------------------------------------------------------
    //print(input_data_refine.Amount[1:])
	X,_ = strconv.ParseFloat(input_data_refine.Amount[1:], 64)
    X = math.Log(math.Max(1,X))
    mapper_code_list[offset_Amount] = X / Mapper_amount
	//
	b,_ := json.Marshal(mapper_code_list)
	//
	return string(b)
	//
}
