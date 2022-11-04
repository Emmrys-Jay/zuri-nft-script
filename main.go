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

// These constant variables represent the ordering of each fields in the input csv file
const (
	TEAM_NAMES = iota
	SERIES_NUMBER
	FILENAME
	NAME
	DESCRIPTION
	GENDER
	ATTRIBUTES
	UUID
)

type (
	// Record models each row in the csv file
	Record struct {
		Format           string      `json:"format" default:"CHIP-0007"`
		Name             string      `json:"name"`
		Description      string      `json:"description"`
		MintingTool      string      `json:"minting_tool"`
		SensitiveContent bool        `json:"sensitive_content"`
		SeriesNumber     int         `json:"series_number"`
		SeriesTotal      int         `json:"series_total"`
		Gender           string      `json:"gender"`
		Attributes       []Attribute `json:"attributes"`
		UUID             string      `json:"uuid"`
		Collections      Collection  `json:"collection"`
	}

	// Attribute models each attribute from the csv file
	Attribute struct {
		TraitType string `json:"trait_type"`
		Value     string `json:"value"`
	}

	// Collection models the collection json collection fiels
	Collection struct {
		Name       string       `json:"name"`
		ID         string       `json:"id"`
		Attributes []CAttribute `json:"attributes"`
	}

	// CAttribute models attribute objects under collection
	CAttribute struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}
)

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

	err = convertRecords(&records, *csvFile)
	if err != nil {
		log.Fatalf("convertRecords: %v", err)
	}

	err = writeCSV(*csvFile, &records)
	if err != nil {
		log.Fatalf("writeCSV: %v", err)
	}
}

// convertRecords performs all basic logic required to parse a csv file, generate its JSON,
// and modify a slice to form an updated csv file in a step by step order
func convertRecords(r *[][]string, filename string) error {
	records := make([]Record, 0)
	var mintingTool string

	// name of directory where all JSON files will be stored
	dir := "nft-jsons"

	if err := createJsonDir(dir); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fmt.Errorf("error: Could not create dir %v, got error %v", dir, err)
		}
	}

	for k := 0; k < len(*r); k++ {
		if k == 0 {
			(*r)[k] = append((*r)[k], "Hash")
			continue
		}

		record := Record{
			Format:           "CHIP-0007",
			Name:             strings.TrimSpace((*r)[k][NAME]),
			Description:      strings.TrimSpace((*r)[k][DESCRIPTION]),
			SensitiveContent: false,
			SeriesTotal:      420,
			Gender:           strings.TrimSpace((*r)[k][GENDER]),
			Attributes:       []Attribute{},
			UUID:             strings.TrimSpace((*r)[k][UUID]),
		}

		attrbs := strings.Split((*r)[k][ATTRIBUTES], ";")
		if len(attrbs) > 1 && attrbs[0] != "" {
			for _, v := range attrbs {
				i := strings.Split(v, ":")
				if len(i) > 1 {
					attrb := Attribute{
						TraitType: strings.TrimSpace(strings.ToLower(i[0])),
						Value:     strings.TrimSpace(strings.ToLower(i[1])),
					}

					record.Attributes = append(record.Attributes, attrb)
				}
			}
		}

		if hasTeamName((*r)[k]) {
			mintingTool = (*r)[k][TEAM_NAMES]
		}

		if !isValid(&record) {
			continue
		}

		sn, err := strconv.Atoi((*r)[k][SERIES_NUMBER])
		if err != nil {
			continue
		}
		record.SeriesNumber = sn
		record.MintingTool = mintingTool

		record.Collections = Collection{
			Name: "Zuri NFT Tickets for Free Lunch",
			ID:   "b774f676-c1d5-422e-beed-00ef5510c64d",
			Attributes: []CAttribute{
				{
					"description",
					"Rewards for accomplishments during HNGi9.",
				},
			},
		}

		rs, err := json.MarshalIndent(record, "", " ")
		if err != nil {
			return fmt.Errorf("error: Could not marshal json file")
		}

		sha, err := generateJSONFileSHA256((*r)[k][FILENAME], dir, rs)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}

		(*r)[k] = append((*r)[k], sha)

		records = append(records, record)
	}

	return nil
}

// isValid checks that a record has some required fields
func isValid(r *Record) bool {
	return r.Description != "" && r.Name != "" && r.UUID != "" && r.Gender != ""
}

// HasTeamName checks a record/row in the csv file to ascertain
// that the row contains a team name
func hasTeamName(record []string) bool {
	return record[TEAM_NAMES] != ""
}

// writeCSV creates a new csv file and writes the new data into it
func writeCSV(filename string, records *[][]string) error {
	filename = strings.TrimSuffix(filename, ".csv")

	file, err := os.Create(filename + ".output.csv")
	if err != nil {
		return fmt.Errorf("error: Could not create output csv file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(*records)
	if err != nil {
		return fmt.Errorf("error: Could not write into output csv file")
	}

	return nil
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
