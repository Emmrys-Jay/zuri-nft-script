package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Record models each row in the csv file
type Record struct {
	Format           string      `json:"format" default:"CHIP-0007"`
	Name             string      `json:"name"`
	Description      string      `json:"description"`
	MintingTool      string      `json:"minting_tool"`
	SensitiveContent bool        `json:"sensitive_content"`
	SeriesNumber     int         `json:"series_number"`
	SeriesTotal      int         `json:"series_total"`
	Attributes       []Attribute `json:"attributes"`
	UUID             string      `json:"uuid"`
	Hash             string      `json:"sha256,omitempty"`
}

// Attribute models each attribute from the csv file
type Attribute struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

func main() {
	csvFile := flag.String("csv", "hngi9-csv-file.csv", "a csv file containing hng files")
	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatalf("error: Could not open file %v\n", *csvFile)
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error: Could not read file %v\n", *csvFile)
	}

	_ = convertRecords(&records, *csvFile)

	writeCSV(*csvFile, &records)
}

// convertRecords performs all basic logic required to parse a csv file, generate its JSON,
// and output an updated csv file in a step by step order
func convertRecords(r *[][]string, filename string) *[]Record {
	records := make([]Record, 0)
	var mintingTool string

	// name of directory where all JSON files will be stored
	dir := "nft-jsons"

	err := createJsonDir(dir)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatalf("error: Could not create dir %v, got error %v", dir, err)
		}
	}

	for k := 0; k < len(*r); k++ {
		if k == 0 {
			(*r)[k] = append((*r)[k], "SHA256")
			continue
		}

		record := Record{
			Format:           "CHIP-0007",
			Name:             (*r)[k][2],
			Description:      (*r)[k][3],
			SensitiveContent: false,
			SeriesTotal:      420,
			Attributes: []Attribute{
				{
					TraitType: "gender",
					Value:     (*r)[k][4],
				},
			},
			UUID: (*r)[k][6],
		}

		attrbs := strings.Split((*r)[k][5], ",")
		if len(attrbs) > 1 && attrbs[0] != "" {
			for _, v := range attrbs {
				i := strings.Split(v, ":")
				if len(i) > 1 {
					attrb := Attribute{
						TraitType: strings.ToLower(i[0]),
						Value:     strings.ToLower(i[1]),
					}

					record.Attributes = append(record.Attributes, attrb)
				}
			}
		}

		if isTeamName(&record) {
			mintingTool = (*r)[k][0]
		}

		if !isValid(&record) {
			continue
		}

		sn, err := strconv.Atoi((*r)[k][0])
		if err != nil {
			continue
		}
		record.SeriesNumber = sn
		record.MintingTool = mintingTool

		rs, err := json.MarshalIndent(record, "", " ")
		if err != nil {
			log.Fatalf("error: Could not marshal json file")
		}

		sha, err := generateJSONFileSHA256((*r)[k][1], dir, rs)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		record.Hash = string(sha)
		(*r)[k] = append((*r)[k], record.Hash)

		records = append(records, record)
	}

	return &records
}

// isValid checks that a record has some required fields
func isValid(r *Record) bool {
	return r.Description != "" && r.Name != "" && r.UUID != ""
}

// isTeamName checks a record/row in the csv file to check if the row contains
// only the team name
func isTeamName(record *Record) bool {
	return record.Name == "" && record.Description == "" && len(record.Attributes) == 1 && record.UUID == ""
}

// writeCSV creates a new csv file and writes the new data into it
func writeCSV(filename string, records *[][]string) {
	filename = strings.TrimSuffix(filename, ".csv")

	file, err := os.Create(filename + ".output.csv")
	if err != nil {
		log.Fatal("error: Could not create output csv file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.WriteAll(*records)
}

// createJsonDir creates the directory where all JSON files will be stores
func createJsonDir(name string) error {
	err := os.Mkdir(name, 0750)
	if err != nil {
		return err
	}

	return nil
}

// generateJSONFileSHA256 generates a json file for a JSON data,
// and returns the SHA256 hash of the input json
func generateJSONFileSHA256(filename string, dirName string, json []byte) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := filepath.Join(wd, dirName, filename+".json")
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(json))
	if err != nil {
		return "", err
	}

	sha := sha256.Sum256(json)
	return fmt.Sprintf("%x", sha), nil
}
