package index

import (
	"github.com/atthakorn/search-engine/internal/config"
	"github.com/blevesearch/bleve"
	"log"
	"github.com/blevesearch/bleve/mapping"
	"github.com/atthakorn/search-engine/internal/blevex/lang/th"
	"io/ioutil"
	"path/filepath"
	"github.com/atthakorn/search-engine/internal/crawler"

	"strings"
	"time"
	"os"
)



func Index() {

	index := newIndex()
	defer index.Close()

	benchmark := benchmark(index, indexing)
	benchmark()

}





func setupDataDirectory() {

	//destroy any outdated data
	os.RemoveAll(config.IndexPath)

	//create data placeholder
	os.Mkdir(config.IndexPath, 0755)
}



func newIndex() bleve.Index {

	setupDataDirectory()

	indexPath := config.IndexPath

	index, err := bleve.Open(indexPath)

	log.Printf("Creating new index...")

	indexMapping := buildIndexMapping()
	index, err = bleve.New(indexPath, indexMapping)

	if err != nil {
		log.Printf("Terminate indexer, cannot create index at %s", indexPath)
		return nil
	}

	return index
}

func buildIndexMapping() (mapping.IndexMapping) {

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = th.AnalyzerName

	return indexMapping
}



func benchmark(index bleve.Index, indexing func(index bleve.Index) (int, error)) func(){

	return func() {

		startTime := time.Now()

		count, err := indexing(index)
		if err != nil {
			log.Printf("Indexing Error: %v", err)

		}
		indexDuration := time.Since(startTime)
		indexDurationSeconds := float64(indexDuration) / float64(time.Second)
		timePerDocument := float64(indexDuration) / float64(count)
		log.Printf("Indexed %d documents, in %.2fs (average %.2f ms/document)", count, indexDurationSeconds, timePerDocument/float64(time.Millisecond))
	}
}

func indexing(index bleve.Index) (count int, err error) {

	//total number of indexed documents
	count = 0
	dataPath := config.DataPath
	entries, err := ioutil.ReadDir(dataPath)

	if err != nil {
		log.Printf("Terminate indexer, cannot load data entries at %s", dataPath)
		return 0, err
	}

	//bulk by indexing 100 documents at a time
	batch := index.NewBatch()
	batchSize := 50
	batchCount:= 0
	for _, entry := range entries {

		//skip entry if it is directory, or not json file
		if entry.IsDir() ||  !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}


		file := filepath.Join(dataPath, entry.Name())
		json, err := crawler.LoadString(file)

		if err != nil {
			log.Printf("Fail to load crawler data file: %s", file)
			return 0, err
		}

		//unmarshal json
		var datas []crawler.Data
		err = crawler.Unmarshal(json, &datas)

		if err != nil {
			log.Printf("Fail to unmarshal data from json: %s", file)
			return 0, err
		}

		for _, data := range datas {
			count++
			batch.Index(data.URL, &Data{
				URL:   data.URL,
				Title: data.Title,
				Body:  strings.Join(data.Texts, " · "),
			})

			batchCount++
			if batchCount >= batchSize {
				err = index.Batch(batch)

				if err != nil {
					log.Printf("Bulk indexing error: %v", err)
					return 0, err
				}

				log.Printf("Documents already indexed: %d", count)
				batch = index.NewBatch()
				batchCount = 0
			}

		}
	}

	// flush the last batch
	if batchCount > 0 {
		err = index.Batch(batch)
		if err != nil {
			log.Printf("Bulk indexing error: %v", err)
			return 0, err
		}
	}

	return count, nil
}


