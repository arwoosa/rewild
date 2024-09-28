package repository

import (
	"errors"
	"fmt"
	"net/http"
	"oosa_rewild/internal/helpers"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type LinkRepository struct{}
type LinkRequest struct {
	Url string `json:"url" validate:"required"`
}

type HTMLMeta struct {
	Url           string
	Title         string
	Description   string
	OGTitle       string
	OGDescription string
	OGImage       string
	//OGAuthor      string
	//OGPublisher   string
	OGSiteName string
}

func (r LinkRepository) Query(c *gin.Context) {
	var payload LinkRequest
	validateError := helpers.Validate(c, &payload)
	if validateError != nil {
		return
	}

	meta, err := r.GetMeta(c, payload.Url)
	if err != nil {
		helpers.ResponseBadRequestError(c, err.Error())
	} else {
		c.JSON(200, meta)
	}
}

func (r LinkRepository) GetMeta(c *gin.Context, Url string) (HTMLMeta, error) {
	var HTMLMeta HTMLMeta

	schemeMatch, _ := regexp.MatchString("^(http|https)://", Url)
	if !schemeMatch {
		Url = "http://" + Url
	}

	response, err := http.Get(Url)
	HTMLMeta.Url = Url

	if err != nil {
		return HTMLMeta, err
	} else {
		defer response.Body.Close()
		if response.StatusCode != 200 {
			errMessage := fmt.Sprintf("status code error: %d %s", response.StatusCode, response.Status)
			return HTMLMeta, errors.New(errMessage)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			helpers.ResponseBadRequestError(c, err.Error())
			return HTMLMeta, err
		}

		description, descriptionExists := doc.Find("meta[name=\"description\"]").First().Attr("content")
		if descriptionExists {
			HTMLMeta.Description = description
		}

		ogTitle, ogTitleExists := doc.Find("meta[property=\"og:title\"]").First().Attr("content")
		if ogTitleExists {
			HTMLMeta.OGTitle = ogTitle
		}

		ogDescription, ogDescriptionExists := doc.Find("meta[property=\"og:description\"]").First().Attr("content")
		if ogDescriptionExists {
			HTMLMeta.OGDescription = ogDescription
		}

		ogImage, ogImageExists := doc.Find("meta[property=\"og:image\"]").First().Attr("content")
		if ogImageExists {
			HTMLMeta.OGImage = ogImage
		}

		ogSitename, ogSitenameExists := doc.Find("meta[property=\"og:site_name\"]").First().Attr("content")
		if ogSitenameExists {
			HTMLMeta.OGSiteName = ogSitename
		}

		title := doc.Find("title").Text()
		HTMLMeta.Title = title
	}
	return HTMLMeta, nil
}
