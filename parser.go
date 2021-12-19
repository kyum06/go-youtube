package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

type YOUTUBE struct{}

type playerResponseData struct {
	StreamingData struct {
		ExpiresInSeconds string   `json:"expiresInSeconds"`
		Formats          []Format `json:"formats"`
		AdaptiveFormats  []Format `json:"adaptiveFormats"`
		DashManifestURL  string   `json:"dashManifestUrl"`
		HlsManifestURL   string   `json:"hlsManifestUrl"`
	} `json:"streamingData"`
	VideoDetails struct {
		VideoID          string   `json:"videoId"`
		Title            string   `json:"title"`
		LengthSeconds    string   `json:"lengthSeconds"`
		Keywords         []string `json:"keywords"`
		ChannelID        string   `json:"channelId"`
		IsOwnerViewing   bool     `json:"isOwnerViewing"`
		ShortDescription string   `json:"shortDescription"`
		IsCrawlable      bool     `json:"isCrawlable"`
		Thumbnail        struct {
			Thumbnails []Thumbnail `json:"thumbnails"`
		} `json:"thumbnail"`
		AverageRating     float64 `json:"averageRating"`
		AllowRatings      bool    `json:"allowRatings"`
		ViewCount         string  `json:"viewCount"`
		Author            string  `json:"author"`
		IsPrivate         bool    `json:"isPrivate"`
		IsUnpluggedCorpus bool    `json:"isUnpluggedCorpus"`
		IsLiveContent     bool    `json:"isLiveContent"`
	} `json:"videoDetails"`
}

type Format struct {
	ItagNo           int    `json:"itag"`
	URL              string `json:"url"`
	MimeType         string `json:"mimeType"`
	Quality          string `json:"quality"`
	Cipher           string `json:"signatureCipher"`
	Bitrate          int    `json:"bitrate"`
	FPS              int    `json:"fps"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	LastModified     string `json:"lastModified"`
	ContentLength    int64  `json:"contentLength,string"`
	QualityLabel     string `json:"qualityLabel"`
	ProjectionType   string `json:"projectionType"`
	AverageBitrate   int    `json:"averageBitrate"`
	AudioQuality     string `json:"audioQuality"`
	ApproxDurationMs string `json:"approxDurationMs"`
	AudioSampleRate  string `json:"audioSampleRate"`
	AudioChannels    int    `json:"audioChannels"`
}

type Thumbnail struct {
	URL    string
	Width  uint
	Height uint
}

var (
	URL, _   = regexp.Compile(`https?://(youtu.be/|(www\.)?youtube.com/watch\?v=)[\-_\w]+`)
	BODY, _  = regexp.Compile(`{"res.*}}}`)
	JSURL, _ = regexp.Compile(`/s/p[^"]*`)
)

func (y YOUTUBE) Get(url string) (string, error) {

	if !URL.MatchString(url) {
		return "", errors.New("")
	}

	res, _ := http.Get(url)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	bytes := BODY.Find(body)

	var context playerResponseData
	json.Unmarshal(bytes, &context)

	formats := append(
		context.StreamingData.AdaptiveFormats,
		context.StreamingData.Formats...,
	)

	var format Format
	for _, v := range formats {
		if v.ItagNo == 251 {
			format = v
			break
		}
	}

	if format.ItagNo == 0 {
		return "", errors.New("")
	}

	jsURL := JSURL.Find(body)
	y.Decipher(&format, jsURL)

	return format.URL, nil
}

// https://www.youtube.com/watch?v=M8HtFw_GVPk
