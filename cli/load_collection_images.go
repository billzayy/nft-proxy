package main

import (
	nft_proxy "github.com/alphabatem/nft-proxy"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	services "github.com/alphabatem/nft-proxy/services"
	"github.com/gagliardetto/solana-go"
	"log"
	"io"
	"net/http"
	"os"
)

var solSvc SolanaService
var solImgSvc SolanaImageService

type collectionLoader struct {
	metaWorkerCount  int
	fileWorkerCount  int
	mediaWorkerCount int

	metaDataIn chan *token_metadata.Metadata
	fileDataIn chan *nft_proxy.NFTMetadataSimple
	mediaIn    chan *nft_proxy.Media
}

func main() {
	log.Printf("Loading collection images: %s", "TODO")

	l := collectionLoader{
		metaWorkerCount:  3,
		fileWorkerCount:  3,
		mediaWorkerCount: 1,
		metaDataIn:       make(chan *token_metadata.Metadata),
		fileDataIn:       make(chan *nft_proxy.NFTMetadataSimple),
		mediaIn:          make(chan *nft_proxy.Media),
	}

	err := solSvc.Start()

	if err != nil {
		return err
	}

	l.spawnWorkers()

	//TODO Get collection
	err := l.loadCollection()
	if err != nil {
		panic(err)
	}

	//TODO Fetch all the mints for that collection
	err := l.fetchAllMints()
	if err != nil {
		panic(err)
	}

	//TODO Fetch Mints/Hash List

	//TODO Batch into batches of 100
	//TODO Pass to metaDataIn<-

	//TODO Fetch all the metadata accounts for that collection
	err := l.fetchAccount()
	if err != nil {
		panic(err)
	}
	//TODO Fetch all images for the accounts
	//TODO Fetch Image
	//TODO Resize Image 500x500
	//TODO Fetch Media
}

func (l *collectionLoader) spawnWorkers() {
	for i := 0; i < l.metaWorkerCount; i++ {
		go l.metaDataWorker()
	}
	for i := 0; i < l.fileWorkerCount; i++ {
		go l.fileDataWorker()
	}
	for i := 0; i < l.mediaWorkerCount; i++ {
		go l.mediaWorker()
	}
}

func (l *collectionLoader) loadCollection() error {
	log.Fatal("Loading collection...")

	// Get block hash from rpc
	solHash, err := solSvc.RecentBlockhash()

	if err != nil{
		return err
	}

	data, _,err := solSvc.TokenData(solHash)

	if err != nil {
		return err
	}

	log.Printf("Collection Loaded: %v", data.Collection)

	for _, metadata := range data {
		l.metaDataIn <- &token_metadata.Metadata{
			Key:             metadata.Key,
			UpdateAuthority: metadata.UpdateAuthority,
			Mint:            metadata.Mint,
			Data:            metadata.Data,
			Collection: 	 metadata.Collection
		}
	}

	return nil
}

func (l *collectionLoader) fetchAllMints() error{
	var mints []solana.PublicKey

	if l.metaDataIn == nil {
		return err
	}

	for metadata := range l.metaDataIn{
		mints = append(mints, metadata.Mints)
	}

	log.Printf("Mints: %v",mints)

	return nil
}

func (l *collectionLoader) fetchAccount() error{
	for metadata := range l.metaDataIn {
        err := solImgSvc.FetchAccount(metadata.Key)
        if err != nil {
            return fmt.Errorf("failed to fetch account: %w", err)
        }
    }
    return nil
}

//Fetches the off-chain data from the on-chain account & passes to `fileDataWorker`
func (l *collectionLoader) metaDataWorker() {
	for metadata := range l.metaDataIn {
		assetURI := metadata.Data.Uri

		log.Printf("Fetching metadata from URI: %s", assetURI)

		resp, err := http.Get(assetURI)
		if err != nil {
			log.Fatalf("failed to fetch metadata from URI %s: %w", uri, err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %w", err)
		}

		var assetMetadata nft_proxy.NFTMetadataSimple

		err = json.Unmarshal(body, &metadata)
		if err != nil {
			log.Fatalf("failed to unmarshal JSON data: %w", err)
		}

		l.fileDataIn <- assetMetadata
	}
}

//Downloads required files & passes to `mediaWorker`
func (l *collectionLoader) fileDataWorker() {
	filePath := "./raw_cache/solona/"
	for assetMetadata := range l.fileDataIn {
		log.Printf("Downloading file: %s", assetMetadata.URI)

		response, err := http.Get(assetMetadata.URI)
		if err != nil {
			log.Fatalf("failed to download file: %v", err)
		}
		defer response.Body.Close()

		// Create the local file
		file, err := os.Create(filepath)
		if err != nil {
			log.Fatalf("failed to create file: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("failed to save file: %v", err)
		}

		media := &nft_proxy.Media{
			MediaUri: assetMetadata.ImageUri,
		}

		l.mediaIn <- media
	}	
}

//Stores media data down to SQL
func (l *collectionLoader) mediaWorker() {
	var db *services.SqliteService

	err := db.Start()
	if err != nil {
		log.Fatalf(err)
	}

	defer db.Shutdown()

    for m := range l.mediaIn {
        // Error handling for database saving
        _, err = db.Create(m.MediaURI)
        if err != nil {
            log.Printf("Error saving media: %v", err)
        }
        log.Printf("M: %s", m.MediaUri)
    }
}