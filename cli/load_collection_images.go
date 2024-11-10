package main

import (
	nft_proxy "github.com/alphabatem/nft-proxy"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	services "github.com/alphabatem/nft-proxy/services"
	"github.com/gagliardetto/solana-go"
	"log"
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

		// Fetch the asset metadata using the URI
		assetMetadata, err := l.fetchMetadata(assetURI)
		if err != nil {
			log.Printf("Error fetching metadata from URI %s: %v", assetURI, err)
			continue
		}
		// Pass the fetched metadata to the fileDataWorker
		l.fileDataIn <- assetMetadata
	}
}

//Downloads required files & passes to `mediaWorker`
func (l *collectionLoader) fileDataWorker(assetURI token_metadata.Data) {
	for assetMetadata := range l.fileDataIn {
		log.Printf("Downloading file: %s", assetMetadata.URI)

		// Simulate file download (replace with actual download logic)
		// Assuming we download and create media file metadata
		media := &nft_proxy.Media{
			MediaUri: assetMetadata.ImageUri,
		}

		// Pass the media data to mediaWorker
		l.mediaIn <- media
	}	
}

//Stores media data down to SQL
func (l *collectionLoader) mediaWorker() {
	var solImgSvc services.SolanaImageService

    for m := range l.mediaIn {
        // Error handling for database saving
        err := 
        if err != nil {
            log.Printf("Error saving media: %v", err)
        }
        log.Printf("M: %s", m.MediaUri)
    }
}