package comic

import (
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Comic struct {
	Month      string
	Num        int
	Link       string
	Year	   string
	News       string
	SafeTitle  string
	Transcript string
	Alt        string
	Img        string
	Title      string
	Day        string
}

func GetComic(n int) (*Comic, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", n)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var comic Comic
	if err := json.NewDecoder(resp.Body).Decode(&comic); err != nil {
		return nil, err
	}

	return &comic, nil
}

func BuildIndex(client *redis.Client, n int) error {
	for i := 1; i <= n; i++ {
		hashName := fmt.Sprintf("comic_%d", i)
		if res, _ := client.Exists(hashName).Result(); res == 1 {
			log.Printf("comic:%d already built", i)
			continue
		}

		log.Printf("building inverted index for comic:%d", i)

		cm, err := GetComic(i)
		if err != nil {
			return fmt.Errorf("encounter error when building index for file %d: %s", i, err)
		}

		re, err := regexp.Compile("[^a-zA-Z0-9 '\n]+")
		if err != nil {
			return err
		}
		client.HSet(hashName, "transcript", cm.Transcript)
		client.HSet(hashName, "url", cm.Img)

		tmp := re.ReplaceAllString(cm.Transcript, "")
		transClean := strings.ReplaceAll(tmp, "\n", " ")
		transLower := strings.ToLower(transClean)

		words := strings.Split(transLower, " ")
		for _, w := range words {
			client.SAdd(w, hashName)
		}

	}

	return nil
}

func Search(client *redis.Client, keys...string) []string {
	var set mapset.Set

	for i, key := range keys {
		values, _ := client.SMembers(key).Result()
		if i == 0 {
			set = mapset.NewSet(stringToInterface(values)...)
		} else {
			curSet := mapset.NewSet(stringToInterface(values)...)
			set = set.Intersect(curSet)
			if set.Cardinality() == 0 {
				return nil
			}
		}
	}

	setSlice := set.ToSlice()
	var res []string

	for _, element := range setSlice {
		cm, _ := client.HGetAll(element.(string)).Result()
		comicJson, _ := json.Marshal(cm)
		res = append(res, string(comicJson))
	}

	return res
}

func stringToInterface(s []string) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}




