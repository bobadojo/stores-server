package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand/v2"

	pb "github.com/bobadojo/go/pkg/stores/v1/storespb"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const defaultName = "world"

var (
	addr     = flag.String("addr", "127.0.0.1:50051", "Address of grpc server.")
	key      = flag.String("api-key", "", "API key.")
	token    = flag.String("token", "", "Authentication token.")
	keyfile  = flag.String("keyfile", "", "Path to a Google service account key file.")
	audience = flag.String("audience", "", "Audience.")
	insecure = flag.Bool("insecure", false, "Insecure connections.")
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

func run() error {
	flag.Parse()

	// Set up a connection to the server.
	var conn *grpc.ClientConn
	var err error
	if *insecure {
		conn, err = grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
	} else {
		p := x509.NewCertPool()
		p.AppendCertsFromPEM([]byte(roots))
		tlsConfig := &tls.Config{
			RootCAs:            p,
			InsecureSkipVerify: true,
		}

		conn, err = grpc.Dial(*addr, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		//conn, err := grpc.Dial(os.Getenv("hello.endpoints.agentio.cloud.goog:443"), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})))
		//conn, err := grpc.Dial("hello.endpoints.agentio.cloud.goog:443", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
	}
	defer conn.Close()
	c := pb.NewStoresClient(conn)

	if *keyfile != "" {
		log.Printf("Authenticating using Google service account key in %s", *keyfile)
		keyBytes, err := ioutil.ReadFile(*keyfile)
		if err != nil {
			log.Fatalf("Unable to read service account key file %s: %v", *keyfile, err)
		}

		tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, *audience)
		if err != nil {
			log.Fatalf("Error building JWT access token source: %v", err)
		}
		jwt, err := tokenSource.Token()
		if err != nil {
			return nil
		}
		*token = jwt.AccessToken
		// NOTE: the generated JWT token has a 1h TTL.
		// Make sure to refresh the token before it expires by calling TokenSource.Token() for each outgoing requests.
		// Calls to this particular implementation of TokenSource.Token() are cheap.
	}

	ctx := context.Background()
	if *key != "" {
		log.Printf("Using API key, len=%d", len(*key))
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", *key)
	}
	if *token != "" {
		log.Printf("Using authentication token: %s", *token)
		ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	minLat := 39.0 + float32(rand.Uint32()%100)/100.0
	minLng := -85.0 + float32(rand.Uint32()%100)/100.0

	deltaLat := float32(0.5)
	deltaLng := float32(0.5)

	req := &pb.FindStoresRequest{
		Bounds: &pb.BoundingBox{
			Max: &pb.Location{
				Latitude:  minLat + deltaLat,
				Longitude: minLng + deltaLng,
			},
			Min: &pb.Location{
				Latitude:  minLat,
				Longitude: minLng,
			},
		},
	}
	log.Printf("%+v", req)

	// Contact the server and print out its response.
	r, err := c.FindStores(ctx, req)
	if err != nil {
		return err
	}
	log.Printf("%d stores", r.Count)
	for i, s := range r.Stores {
		r2, err := c.GetStore(ctx, &pb.GetStoreRequest{
			Name: s.Name,
		})
		if err != nil {
			return err
		}

		b, err := json.Marshal(r2)
		if err != nil {
			return err
		}
		log.Printf("%d: %s", i, string(b))

	}
	return nil
}
