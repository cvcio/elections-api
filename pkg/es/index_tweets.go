package es

var indexTweets = `
{
	"settings": {
        "index": {
            "number_of_shards": 3,
            "number_of_replicas": 2,
            "mapping": {
                "total_fields": {
                    "limit": 3600
                }
            }
        }
    },
	"mappings": {
        "document": {
		  "properties": {
		    "created_at": {
		      "format": "EEE MMM dd HH:mm:ss Z YYYY",
		      "type": "date"
		    },
		    "text": {
		      "type": "text",
		      "index": true
		    },
		    "source": {
		      "type": "text",
		      "index": true
		    },
		    "in_reply": {
		      "type": "object",
		      "properties": {
		        "status": {
		          "type": "long"
		        },
		        "user_id": {
		          "type": "long"
		        },
		        "screen_name": {
		          "type": "text",
		          "index": true
		        }
		      }
		    },
		    "user": {
		      "type": "object",
		      "properties": {
		        "id": {
		          "type": "long"
		        },
		        "name": {
		          "type": "keyword",
		          "index": true
		        },
		        "screen_name": {
		          "type": "keyword",
		          "index": true
		        },
		        "location": {
		          "type": "text"
		        },
		        "description": {
		          "type": "text",
		          "index": true
		        },
		        "url": {
		          "type": "text"
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
		        "lang": {
		          "type": "text"
		        },
		        "profile_image_url": {
		          "type": "text"
		        },
		        "profile_banner_url": {
		          "type": "text"
		        }
		      }
		    },
		    "geo": {
		      "type": "geo_point",
		      "ignore_malformed": true
		    },
		    "place": {
		      "type": "object",
		      "properties": {
		        "id": {
		          "type": "text"
		        },
		        "url": {
		          "type": "text"
		        },
		        "place_type": {
		          "type": "text"
		        },
		        "name": {
		          "type": "text"
		        },
		        "full_name": {
		          "type": "text",
		          "index": true
		        },
		        "country_code": {
		          "type": "text"
		        },
		        "country": {
		          "type": "text"
		        }
		      }
		    },
		    "retweet_count": {
		      "type": "long"
		    },
		    "favorite_count": {
		      "type": "long"
		    },
		    "lang": {
		      "type": "text"
		    },
		    "hashtags": {
		      "properties": {
		        "text": {
		          "type": "text",
		          "index": true
		        }
		      }
		    },
		    "user_mentions": {
		      "properties": {
		        "screen_name": {
		          "type": "keyword",
		          "index": true
		        },
		        "name": {
		          "type": "keyword",
		          "index": true
		        },
		        "id": {
		          "type": "long"
		        }
		      }
		    },
		    "media": {
		      "properties": {
		        "id": {
		          "type": "long"
		        },
		        "media_url_https": {
		          "type": "text"
		        },
		        "expanded_url": {
		          "type": "text"
		        },
		        "type": {
		          "type": "text"
		        }
		      }
		    },
		    "retweeted_status": {
		      "type": "object",
		      "properties": {
		        "created_at": {
		          "format": "EEE MMM dd HH:mm:ss Z YYYY",
		          "type": "date"
		        },
		        "id": {
		          "type": "long"
		        },
		        "text": {
		          "type": "text",
		          "index": true
		        },
		        "source": {
		          "type": "text",
		          "index": true
		        },
		        "in_reply": {
		          "type": "object",
		          "properties": {
		            "status": {
		              "type": "long"
		            },
		            "user_id": {
		              "type": "long"
		            },
		            "screen_name": {
		              "type": "keyword"
		            }
		          }
		        },
		        "user": {
		          "type": "object",
		          "properties": {
		            "id": {
		              "type": "long"
		            },
		            "name": {
		              "type": "keyword",
		              "index": true
		            },
		            "screen_name": {
		              "type": "keyword",
		              "index": true
		            },
		            "location": {
		              "type": "text"
		            },
		            "description": {
		              "type": "text",
		              "index": true
		            },
		            "url": {
		              "type": "text"
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
		            "lang": {
		              "type": "text"
		            },
		            "profile_image_url": {
		              "type": "text"
		            },
		            "profile_banner_url": {
		              "type": "text"
		            }
		          }
		        }
		      }
		    },
		    "quoted_status": {
		      "type": "object",
		      "properties": {
		        "created_at": {
		          "format": "EEE MMM dd HH:mm:ss Z YYYY",
		          "type": "date"
		        },
		        "id": {
		          "type": "long"
		        },
		        "text": {
		          "type": "text",
		          "index": true
		        },
		        "source": {
		          "type": "text",
		          "index": true
		        },
		        "in_reply": {
		          "type": "object",
		          "properties": {
		            "status": {
		              "type": "long"
		            },
		            "user_id": {
		              "type": "long"
		            },
		            "screen_name": {
		              "type": "keyword",
		              "index": true
		            }
		          }
		        },
		        "user": {
		          "type": "object",
		          "properties": {
		            "id": {
		              "type": "long"
		            },
		            "name": {
		              "type": "keyword",
		              "index": true
		            },
		            "screen_name": {
		              "type": "keyword",
		              "index": true
		            },
		            "location": {
		              "type": "text"
		            },
		            "description": {
		              "type": "text",
		              "index": true
		            },
		            "url": {
		              "type": "text"
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
		            "lang": {
		              "type": "text"
		            },
		            "profile_image_url": {
		              "type": "text"
		            },
		            "profile_banner_url": {
		              "type": "text"
		            }
		          }
		        }
		      }
		    }
		  }
		}
    }
}`
