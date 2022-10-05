package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
	// csv "github.com/Angling-Zainuddin/scrap-analysis/csv_lib"
)

func main() {
	fmt.Println("starting...")
	db, err := OpenDbConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range MapFileName {
		for _, v := range v {
			err := checkFileExistance(v["file_name"])
			if err != nil {
				fmt.Printf("failed to check file : %s", err)
				break
			}

			readAndWriteCsv(db, v["file_name"], v["phone_index"], v["name_index"])

		}
	}

	fmt.Println("process finished...")
}

func checkFileExistance(fileName string) error {
	_, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return nil
}

func readAndWriteCsv(db *gorm.DB, fileName string, phoneIndex string, nameIndex string) {
	fmt.Printf("reading and writing file with name = %s \n", fileName)
	csvReader, csvFile, err := OpenCsvFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvFile.Close()
	row, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	outputFileName := fmt.Sprintf("analyzed_%s", fileName)
	csvFileSave, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFileSave.Close()

	w := csv.NewWriter(csvFileSave)
	defer w.Flush()

	var data [][]string
	for index, each := range row {
		if index == 0 {
			rowtitle := each
			rowtitle = append(rowtitle, "crm_id", "crm_phone", "crm_whatsapp", "crm_name")
			data = append(data, rowtitle)
			continue
		}

		if index <= 300 {
			fmt.Println("skipped data index number = ", index)
			continue
		}

		phoneIndexInt, err := strconv.Atoi(phoneIndex)
		if err != nil {
			fmt.Printf("failed to convert str phone index %s", err)
		}
		nameIndexInt, err := strconv.Atoi(nameIndex)
		if err != nil {
			fmt.Printf("failed to convert str name index %s", err)
		}
		phonezero, err := formatPhoneZero(each[phoneIndexInt])
		if err != nil {
			fmt.Printf("error formating phone zero : %s \n", err)
			fmt.Println("phone =", each[phoneIndexInt])
			fmt.Println("name =", each[nameIndexInt])
			fmt.Println("file name = ", fileName)
		}

		phone, err := formatPhone(each[phoneIndexInt])
		if err != nil {
			fmt.Printf("error formating phone : %s \n", err)
			fmt.Println("phone =", each[phoneIndexInt])
			fmt.Println("name =", each[nameIndexInt])
			fmt.Println("file name = ", fileName)
		}

		customer := Customers{}
		if err == nil && phone != 0 {
			customer, err = GetCustomerByPhone(db, phone, phonezero)
			if err != nil {
				log.Fatalf("failed to get customer: %s", err)
				break
			}
		}

		dataappend := each
		emptyCustomer := Customers{}
		if customer != emptyCustomer {
			dataappend = append(dataappend, strconv.Itoa(customer.ID), customer.Phone, customer.WhatsApp, customer.Name)
		}

		data = append(data, dataappend)

	}
	w.WriteAll(data)
}

func formatPhoneZero(phone string) (int, error) {
	str := trimLeftChar(phone)
	var regex, err = regexp.Compile(`^(62)`)
	res := regex.ReplaceAllString(str, "0")
	withouhthypen := strings.Replace(res, "-", "", -1)
	intPhone, err := strconv.Atoi(withouhthypen)
	if err != nil {
		return 0, err
	}
	return intPhone, nil
}

func formatPhone(phone string) (int, error) {
	str := trimLeftChar(phone)
	withouhthypen := strings.Replace(str, "-", "", -1)
	intPhone, err := strconv.Atoi(withouhthypen)
	if err != nil {
		return 0, err
	}
	return intPhone, nil
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}
