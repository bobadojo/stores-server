package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/bobadojo/go/pkg/stores/v1/storespb"
	"github.com/tidwall/rtree"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type storesServer struct {
	stores []*storespb.Store
	tree   rtree.RTree

	storespb.UnimplementedStoresServer
}

func NewStoresServer() (*storesServer, error) {
	stores, err := readStores("data/stores.csv")
	if err != nil {
		return nil, err
	}
	tree := buildRTree(stores)
	return &storesServer{stores: stores, tree: tree}, nil
}

func readStores(filename string) ([]*storespb.Store, error) {
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
			// United States Post Office,401,North Ashley Drive,Tampa,FL,33602,27.9473718,-82.4595414

			Name:  idFromField(fields[0]),
			Type:  "office",
			Title: fields[0],
			Location: &storespb.Location{
				Latitude:  atof(fields[6]),
				Longitude: atof(fields[7]),
			},
			Address: &storespb.Address{
				Street:     fields[2],
				City:       fields[3],
				State:      fields[4],
				ZipCode:    atoi(fields[5]),
				RegionCode: "us",
			},
		}
		stores = append(stores, store)
	}
	return stores, nil
}

func idFromField(s string) string {
	hash := md5.Sum([]byte(s))
	return "stores/" + hex.EncodeToString(hash[:])[0:8]
}

func atoi(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
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

func buildRTree(stores []*storespb.Store) rtree.RTree {
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

func (s *storesServer) FindStores(ctx context.Context, request *storespb.FindStoresRequest) (*storespb.FindStoresResponse, error) {
	min := [2]float64{
		float64(request.Bounds.Min.Latitude),
		float64(request.Bounds.Min.Longitude),
	}
	max := [2]float64{
		float64(request.Bounds.Max.Latitude),
		float64(request.Bounds.Max.Longitude),
	}
	limit := int(request.Limit)
	if limit == 0 {
		limit = 50
	} else if limit > 50 {
		limit = 50
	}
	var stores []*storespb.Store
	s.tree.Search(min, max, func(min, max [2]float64, data interface{}) bool {
		stores = append(stores, data.(*storespb.Store))
		return len(stores) < limit
	})
	response := &storespb.FindStoresResponse{
		Stores: stores,
		Count:  int32(len(stores)),
	}
	return response, nil
}

func (s *storesServer) ListStores(ctx context.Context, request *storespb.ListStoresRequest) (*storespb.ListStoresResponse, error) {
	limit := int(request.PageSize)
	if limit == 0 {
		limit = 50
	} else if limit > 50 {
		limit = 50
	}
	offset := 0
	if request.PageToken != "" {
		b, err := base64.RawURLEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "bad page token: %+v", err)
		}
		offset, err = strconv.Atoi(string(b))
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "bad page token: %+v", err)
		}
	}
	var stores []*storespb.Store
	for i := offset; i < min(offset+limit, len(s.stores)); i++ {
		stores = append(stores, s.stores[i])
	}
	nextPageToken := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", offset+limit)))
	response := &storespb.ListStoresResponse{
		Stores:        stores,
		NextPageToken: nextPageToken,
	}
	return response, nil
}

func (s *storesServer) GetStore(ctx context.Context, request *storespb.GetStoreRequest) (*storespb.Store, error) {
	for i := 0; i < len(s.stores); i++ {
		if s.stores[i].Name == request.Name {
			return s.stores[i], nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "%q not found", request.Name)
}
