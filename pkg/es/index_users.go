package es

var indexUsers = `
{
	"settings": {
        "index": {
            "number_of_shards": 3,
            "number_of_replicas": 2
        }
    },
	"mappings": {
        "document": {
		  "properties": {
		    "crawled_at": {
		      "format": "EEE MMM dd HH:mm:ss Z YYYY",
		      "type": "date"
		    },
		    "id": {
		      "type": "long"
		    },
		    "id_str": {
		      "type": "text"
		    },
		    "name": {
		      "type": "keyword",
		      "index": true
		    },
		    "screen_name": {
		      "type": "keyword",
		      "index": true
		    },
		    "description": {
		      "type": "text",
		      "index": true
		    },
		    "followers_count": {
		      "type": "long"
		    },
		    "friends_count": {
		      "type": "long"
		    },
		    "listed_count": {
		      "type": "long"
		    },
		    "created_at": {
		      "format": "EEE MMM dd HH:mm:ss Z YYYY",
		      "type": "date"
		    },
		    "favourites_count": {
		      "type": "long"
		    },
		    "statuses_count": {
		      "type": "long"
		    },
		    "profile_image_url": {
		      "type": "text"
		    },
		    "profile_banner_url": {
		      "type": "text"
		    },
		    "user_class": {
		      "type": "keyword"
		    },
		    "user_class_score": {
		      "type": "long"
		    }
		  }
		}
    }
}
`
