package deleteOldMetadata

import (
	"git_test/src/lib/es"
	"log"
)

const MIN_VERSION_COUNT = 6

func main() {
	buckets, e := es.SearchVersionStatus(MIN_VERSION_COUNT + 1)
	if e != nil {
		log.Println(e)
		return
	}
	for i := range buckets {
		bucket := buckets[i]
		for v := 0; v < bucket.Doc_count-MIN_VERSION_COUNT; i++ {
			es.DelMetadata(bucket.Key, v+int(bucket.Min_version.Value))
		}
	}
}
