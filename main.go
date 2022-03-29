package main

import (
	"github.com/prasetyanurangga/snaptify_api/spotify"
	"github.com/prasetyanurangga/snaptify_api/image_label"
	simplejson "github.com/bitly/go-simplejson"
	convert "github.com/benpate/convert"
    "math"
    "errors"
    "github.com/gin-gonic/gin"
    "os"
    "strings"
)

var errUnexpectedType = errors.New("Non-numeric type could not be converted to float")

func getFloatSwitchOnly(unk interface{}) (float64, error) {
    switch i := unk.(type) {
    case float64:
        return i, nil
    case float32:
        return float64(i), nil
    case int64:
        return float64(i), nil
    case int32:
        return float64(i), nil
    case int:
        return float64(i), nil
    case uint64:
        return float64(i), nil
    case uint32:
        return float64(i), nil
    case uint:
        return float64(i), nil
    default:
        return math.NaN(), errUnexpectedType
    }
}

func getLabel(url string) (string, bool){
	success := false
	m := map[string]interface{}{
	  "url": url,
	}
	apiKeyEnv := getEnv("API_KEY_IMAGE_LABELING")

	imgLabel := imageLabel.New(apiKeyEnv)


	response, err  := imgLabel.Get(m)
	
	biggestValue := 0.0;
	biggestKey := "";
	if err == nil {

		json, _ := simplejson.NewJson(response)

		for key, element := range json.MustMap() {
			feetFloat := convert.Float(element)
			if feetFloat > biggestValue {
				biggestValue = feetFloat
				biggestKey = key
			}
	    }
		success = true
	} else {
		success = false
	}


	return biggestKey, success
}

type track struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Artist     string  `json:"artist"`
    Image    string  `json:"image"`
    URL    	 string  `json:"url"`
}

func getTrackSpotify(keyword string) ([]track){

	clientIDEnv := getEnv("CLIENT_ID_SPOTIFY")
	clientSecretEnv := getEnv("CLIENT_SECRET_SPOTIFY")
	spot := spotify.New(clientIDEnv, clientSecretEnv)

	var trackList []track = []track{}
	authorized, _ := spot.Authorize()
	if authorized {

		// If we ere able to authorize then Get a simple album
		response, _ := spot.Get("search?q=%s&type=track", nil, keyword)

		// Parse response to a JSON Object and
		// get the album's name
		json, _ := simplejson.NewJson(response)
		tracks, existsTrack := json.CheckGet("tracks")
		itemsTrack, existsItem := tracks.CheckGet("items")


		if existsTrack && existsItem {
			items := itemsTrack.MustArray()
			for _, itemT := range items {
				item, _ := itemT.(map[string]interface{})
				externalUrl := item["external_urls"].(map[string]interface{})
				album := item["album"].(map[string]interface{})
				images := album["images"].([]interface{})
				artist := album["artists"].([]interface{})
				thumbnail := images[1].(map[string]interface{})
				var artistText []string

				for _, itemArtistRaw := range artist {
					itemArtist := itemArtistRaw.(map[string]interface{})
					artistText = append(artistText, itemArtist["name"].(string))
				}
		        trackList = append(trackList, track{
		            ID: item["id"].(string), 
		            Name: item["name"].(string), 
		            Artist: strings.Join(artistText, ","), 
		            URL: externalUrl["spotify"].(string), 
		            Image: thumbnail["url"].(string),
		        })
		    }
		}

	}


	return trackList
}


type RequestTrack struct{
    Url string `json:"url"`
}

type RequestTrackKeyword struct{
    Keyword string `json:"keyword"`
}

func getEnv(key string) string {

  return os.Getenv(key)
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
	    
// 	ip := c.ClientIP()

//         if ip != getEnv("WHITE_LIST_IP") {
//         	c.AbortWithStatus(500)
//             return
//         }
	    

        c.Next()
    }
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
    router.POST("/get_by_keyword", func(c *gin.Context) {
//         var requestTrackKeyword RequestTrackKeyword
//         c.BindJSON(&requestTrackKeyword)
//         tracks := getTrackSpotify(requestTrackKeyword.Keyword)
	    ip := c.ClientIP()
	    localIp := getEnv("WHITE_LIST_IP")
        c.JSON(200, gin.H{"data" : ip + " = " + localIp}) // Your custom response here
    })

    router.POST("/get_by_image", func(c *gin.Context) {
        var requestTrack RequestTrack
        c.BindJSON(&requestTrack)
        label, success := getLabel(requestTrack.Url)
        var tracks []track
        if success {
			tracks = getTrackSpotify(label)
        }
        c.JSON(200, gin.H{"data" : tracks}) // Your custom response here
    })

    router.Run()

}
