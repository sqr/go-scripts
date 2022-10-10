package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

func main() {

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the zone ID
	zoneId, err := api.ZoneIDByName("arturobracero.com")
	if err != nil {
		log.Fatal(err)
	}

	currentIp := getCurrentIp()

	args := os.Args[1:]
	for _, r := range args {
		recordIp, recordId, err := getRecordIP(r, api, zoneId)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Current IP is:", currentIp)
		fmt.Println("Record:", r)
		fmt.Println("Cloudflare points to:", recordIp)
		fmt.Println("ID of the record is:", recordId)
		if recordIp != currentIp {
			fmt.Println(updateRecord(r, api, currentIp, recordId, zoneId))
		}
		if recordIp == currentIp {
			fmt.Println("Record is already updated.")
		}
	}
}

func updateRecord(recordName string, api *cloudflare.API, currentIp string, recordId string, zoneId string) string {
	err := api.UpdateDNSRecord(context.Background(), zoneId, recordId, cloudflare.DNSRecord{Content: currentIp})
	if err != nil {
		log.Fatal(err)
	}
	return "Record updated"
}

func getRecordIP(name string, api *cloudflare.API, id string) (string, string, error) {
	recs, err := api.DNSRecords(context.Background(), id, cloudflare.DNSRecord{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	if len(recs) == 0 {
		return "", "", errors.New("Empty array")
	}
	recordId := recs[0].ID
	recordIp := recs[0].Content
	return recordIp, recordId, nil
}

func getCurrentIp() string {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	sb = strings.TrimSuffix(sb, "\n")
	return sb
}
