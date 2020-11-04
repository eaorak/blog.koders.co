package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Entry
	apiKeyFlag = flag.String("apiKey", "", "dev.to API Key")
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logger = logrus.WithFields(logrus.Fields{
		"service": "blog-update",
	})
}

func main() {
	flag.Parse()
	articleMap := make(map[int]*BlogArticle)
	jsonFile, err := os.Open("articles.json")
	if err != nil {
		logger.Fatalf("Unable to read file: %v", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var articles []*BlogArticle

	json.Unmarshal(byteValue, &articles)

	for _, a := range articles {
		articleMap[a.ID] = a
	}
	as, err := getArticles()
	if err != nil {
		logger.Fatalf("Unable to fetch articles: %v", err)
	}
	for _, a := range as {
		articleMap[a.ID] = &BlogArticle{
			ID:          a.ID,
			Title:       a.Title,
			Description: a.Description,
			CoverImage:  a.CoverImage,
			Slug:        a.Slug,
			URL:         a.URL,
			PublishedAt: a.PublishedAt.Format("2006-01-02"),
		}
	}

	var newArticles []*BlogArticle

	for _, a := range articleMap {
		newArticles = append(newArticles, a)
	}
	file, _ := json.MarshalIndent(newArticles, "", " ")

	err = ioutil.WriteFile("articles.json", file, 0644)
	if err != nil {
		logger.Fatalf("Unable to write articles: %v", err)
	}

	cleanPosts()

	t := template.Must(template.New("t1").Parse(`---
title: "{{.Title}}"
excerpt: "{{.Description}}"
date: "2020-08-13"
coverImage: "{{.CoverImage}}"
author:
  name: Ender
  picture: "assets/blog/authors/ender.jpg"
ogImage:
  url: "{{.CoverImage}}"
---

[Read more]({{.URL}})
`))
	for _, a := range newArticles {
		f, err := os.Create(fmt.Sprintf("_posts/%s.md", a.Slug))
		if err != nil {
			logger.Fatalf("Unable to write file: %v", err)
		}
		err = t.Execute(f, a)
		if err != nil {
			logger.Fatalf("Unable execute template: %v", err)
		}
	}
}

func cleanPosts() error {
	d, err := os.Open("_posts")
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join("_posts", name))
		if err != nil {
			return err
		}
	}
	return nil
}

func getArticles() ([]*Article, error) {
	url := "https://dev.to/api/articles?top=1&per_page=5"

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("api-key", *apiKeyFlag)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, err
	}

	var articles []*Article
	err = json.Unmarshal(body, &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

type Article struct {
	TypeOf               string      `json:"type_of"`
	ID                   int         `json:"id"`
	Title                string      `json:"title"`
	Description          string      `json:"description"`
	CoverImage           string      `json:"cover_image"`
	ReadablePublishDate  string      `json:"readable_publish_date"`
	SocialImage          string      `json:"social_image"`
	TagList              []string    `json:"tag_list"`
	Tags                 string      `json:"tags"`
	Slug                 string      `json:"slug"`
	Path                 string      `json:"path"`
	URL                  string      `json:"url"`
	CanonicalURL         string      `json:"canonical_url"`
	CommentsCount        int         `json:"comments_count"`
	PublicReactionsCount int         `json:"public_reactions_count"`
	CollectionID         interface{} `json:"collection_id"`
	CreatedAt            time.Time   `json:"created_at"`
	EditedAt             time.Time   `json:"edited_at"`
	CrosspostedAt        interface{} `json:"crossposted_at"`
	PublishedAt          time.Time   `json:"published_at"`
	LastCommentAt        time.Time   `json:"last_comment_at"`
	PublishedTimestamp   time.Time   `json:"published_timestamp"`
	User                 struct {
		Name            string `json:"name"`
		Username        string `json:"username"`
		TwitterUsername string `json:"twitter_username"`
		GithubUsername  string `json:"github_username"`
		WebsiteURL      string `json:"website_url"`
		ProfileImage    string `json:"profile_image"`
		ProfileImage90  string `json:"profile_image_90"`
	} `json:"user"`
	Organization struct {
		Name           string `json:"name"`
		Username       string `json:"username"`
		Slug           string `json:"slug"`
		ProfileImage   string `json:"profile_image"`
		ProfileImage90 string `json:"profile_image_90"`
	} `json:"organization"`
}

type BlogArticle struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	PublishedAt string `json:"published_at"`
}
