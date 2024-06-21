package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/bobadojo/go/pkg/stores/v1/storespb"
	"github.com/tidwall/rtree"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

func ReadStores(filename string) ([]*storespb.Store, error) {
	var stores []*storespb.Store
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	reader.Read()
	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		store := &storespb.Store{
			Name:  fields[0],
			Type:  fields[1],
			Title: fields[2],
			Location: &storespb.Location{
				Latitude:  atof(fields[10]),
				Longitude: atof(fields[11]),
			},
			Address: &storespb.Address{
				Street:     fields[3],
				City:       fields[4],
				State:      fields[5],
				ZipCode:    atoi(fields[6]),
				RegionCode: fields[13],
				County:     fields[14],
			},
			StoreHours: fields[9],
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func atoi(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return int32(i)
}

func atof(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Fatalf("%s", err)
	}
	return float32(f)
}

func BuildRTree(stores []*storespb.Store) rtree.RTree {
	var tree rtree.RTree
	for _, store := range stores {
		point := [2]float64{
			float64(store.Location.Latitude),
			float64(store.Location.Longitude),
		}
		tree.Insert(point, point, store)
	}
	return tree
}

func run() error {
	stores, err := ReadStores("data/us_post_offices_01-27-23-sample.csv")
	if err != nil {
		return err
	}
	tree := BuildRTree(stores)

	point := [2]float64{
		38.08973633943899,
		-96.8860452622175,
	}
	min := point
	max := point
	epsilon := 0.001
	min[0] -= epsilon
	min[1] -= epsilon
	max[0] += epsilon
	max[1] += epsilon

	tree.Search(min, max, func(min, max [2]float64, data interface{}) bool {
		log.Printf("%+v", data)
		return true
	})

	return nil
}
