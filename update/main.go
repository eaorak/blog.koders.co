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
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const dailyPostCount = 5 // Daily post count fetched from dev.to API
const maxPostCount = 500 // Maximum post count. Old posts will be deleted if total posts exceedes this number.

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
	logger.Info("Articles loaded from json file: %d", len(articleMap))
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
			PublishDate: a.PublishedAt.Format("2006-01-02"),
			PublishedAt: a.PublishedAt,
		}
	}
	logger.Info("Articles fetched from dev.to API: %d", len(as))

	var newArticles []*BlogArticle

	for _, a := range articleMap {
		newArticles = append(newArticles, a)
	}
	sort.Slice(newArticles, func(i, j int) bool { return newArticles[i].PublishedAt.After(newArticles[j].PublishedAt) })

	if len(newArticles) > maxPostCount {
		newArticles = newArticles[:maxPostCount]
	}

	file, _ := json.MarshalIndent(newArticles, "", " ")

	err = ioutil.WriteFile("articles.json", file, 0644)
	if err != nil {
		logger.Fatalf("Unable to write articles: %v", err)
	}

	logger.Info("Articles writed to json file: %d", len(as))

	cleanPosts()

	logger.Info("_posts directory cleaned")

	t := template.Must(template.New("t1").Funcs(template.FuncMap{
		"html": func(value interface{}) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
	}).Parse(`---
title: '{{html .Title}}'
excerpt: '{{html .Description}}'
date: '{{.PublishDate}}'
coverImage: '{{.CoverImage}}'
author:
  name: Koders
  picture: "assets/blog/authors/koders.png"
ogImage:
  url: '{{.CoverImage}}'
---

{{html .Description}}

[Read more]({{.URL}})
`))
	for _, a := range newArticles {
		f, err := os.Create(fmt.Sprintf("_posts/%s.md", a.Slug))
		if err != nil {
			logger.Fatalf("Unable to write file: %v", err)
		}
		a.Title = strings.ReplaceAll(a.Title, "'", "''")
		a.Description = strings.ReplaceAll(a.Description, "'", "''")
		err = t.Execute(f, a)
		if err != nil {
			logger.Fatalf("Unable execute template: %v", err)
		}
	}
	logger.Info("Blog posts created")
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
	url := fmt.Sprintf("https://dev.to/api/articles?top=1&per_page=%d", dailyPostCount)

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
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CoverImage  string    `json:"cover_image"`
	Slug        string    `json:"slug"`
	URL         string    `json:"url"`
	PublishDate string    `json:"publish_date"`
	PublishedAt time.Time `json:"published_at"`
}
