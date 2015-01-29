package api

import (
	"fmt"
)

const (
	originalUrl    = "http://farm%v.staticflickr.com/%v/%v_%v_o.%v"
	thumbnailUrl   = "http://farm%v.staticflickr.com/%v/%v_%v_%v.jpg"
	thumbnailSizeL = "q"
	thumbnailSizeS = "s"
)

/*
{
    "stat"      : "fail",
    "code"      : "97",
    "message"   : "Missing signature"
}
*/
type FlickrResult struct {
	Status string `json:"stat"`
	Error  string `json:"message"`
}

/*
PhotoSearchResult is the result of a flickr.photos.search API call.

{
    "photos":
    {
        "page":1,
        "pages":583,
        "perpage":1,
        "total":"583",
        "photo":
        [
            {
                "id":"14691360159",
                "owner":"48475357@N00",
                "secret":"5ff0a05549",
                "server":"3916",
                "farm":4,
                "title":"",
                "ispublic":1,
                "isfriend":0,
                "isfamily":0,
                "originalsecret":"d09d2d858b",
                "originalformat":"jpg"
            }
        ]
    },
    "stat":"ok"
}
*/
type PhotoSearchResult struct {
	FlickrResult
	Photos PhotoSearchPage
}
type PhotoSearchPage struct {
	// Current page number
	Page int
	// Total number of pages
	Pages  int
	Photos []PhotoSummary `json:"photo"`
}

type PhotoSummary struct {
	PhotoCommon
	Owner     string
	JsonTitle string `json:"title"`
	IsPublic  int    `json:"ispublic"`
	IsFriend  int    `json:"isfriend"`
	IsFamily  int    `json:"isfamily"`
}

type PhotoCommon struct {
	JsonId         string `json:"id"`
	Secret         string
	Server         string
	Farm           int
	OriginalSecret string `json:"originalsecret"`
	OriginalFormat string `json:"originalformat"`
}

/*
PhotoInfoResult is the result of a flickr.photos.getInfo API call.

{
    "photo":
    {
        "id":"14691360159",
        "secret":"5ff0a05549",
        "server":"3916",
        "farm":4,
        "dateuploaded":"1407690835",
        "isfavorite":0,
        "license":"0",
        "safety_level":"0",
        "rotation":0,
        "originalsecret":"d09d2d858b",
        "originalformat":"jpg",
        "owner":
        {
            "nsid":"48475357@N00",
            "username":"pete-t",
            "realname":"Peter Thompson",
            "location":"Vienna, Austria",
            "iconserver":"23",
            "iconfarm":1,
            "path_alias":"petert"
        },
        "title":
        {
            "_content":""
        },
        "description":
        {
            "_content":"Olympus digital camera"
        },
        "visibility":{"ispublic":1,"isfriend":0,"isfamily":0},
        "dates":{"posted":"1407690835","taken":"2013-08-15 16:07:38","takengranularity":"0","takenunknown":0,"lastupdate":"1407701318"},
        "views":"95",
        "editability":{"cancomment":0,"canaddmeta":0},
        "publiceditability":{"cancomment":1,"canaddmeta":0},
        "usage":{"candownload":1,"canblog":0,"canprint":0,"canshare":1},
        "comments":{"_content":"1"},
        "notes":{"note":[]},
        "people":{"haspeople":0},
        "tags":{"tag":[]},
        "urls":{"url":[{"type":"photopage","_content":"https:\/\/www.flickr.com\/photos\/petert\/14691360159\/"}]},
        "media":"photo"
    },
    "stat":"ok"
}
*/
type PhotoInfoResult struct {
	FlickrResult
	Photo PhotoInfo `json:"photo"`
}

type PhotoInfo struct {
	PhotoCommon
	JsonTitle       contentString `json:"title"`
	JsonDescription contentString `json:"description"`
}

type contentString struct {
	Content string `json:"_content"`
}

type Photoer interface {
	Id() string
	// Url to large format thumbnail (150 x 150)
	LargeThumbnailUrl() string
	// Url to small format thumbnail (75 x 75)
	SmallThumbnailUrl() string
	// URL of original-size photo
	OriginalUrl() string
	// The photo's title
	Title() string
}

type FullPhotoer interface {
	Photoer
	Description() string
}

func (p PhotoCommon) Id() string {
	return p.JsonId
}

func (p PhotoCommon) LargeThumbnailUrl() string {
	// "http://farm"+photo.farm+".staticflickr.com/"+photo.server+"/"+photo.id+"_"+photo.secret+"_q.jpg"
	return fmt.Sprintf(thumbnailUrl, p.Farm, p.Server, p.JsonId, p.Secret, thumbnailSizeL)
}

func (p PhotoCommon) SmallThumbnailUrl() string {
	// "http://farm"+photo.farm+".staticflickr.com/"+photo.server+"/"+photo.id+"_"+photo.secret+"_q.jpg"
	return fmt.Sprintf(thumbnailUrl, p.Farm, p.Server, p.JsonId, p.Secret, thumbnailSizeS)
}

func (p PhotoCommon) OriginalUrl() string {
	// 'http://farm'+photo.farm+'.staticflickr.com/'+photo.server+'/'+photo.id+'_'+photo.originalSecret+'_o.'+photo.originalFormat
	return fmt.Sprintf(originalUrl, p.Farm, p.Server, p.JsonId, p.OriginalSecret, p.OriginalFormat)
}

func (p PhotoSummary) Title() string {
	return p.Title()
}

func (p PhotoInfo) Title() string {
	return p.JsonTitle.Content
}

func (p PhotoInfo) Description() string {
	return p.JsonDescription.Content
}
