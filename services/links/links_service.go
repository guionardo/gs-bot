package links

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

type (
	LinksService struct {
		repository *LinksRepository
		logger     *logrus.Entry
	}
)

func (s *LinksService) FetchLink(url string) (model *LinksModel, err error) {
	url = strings.Trim(url, " ")
	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == 404 {
		err = fmt.Errorf("Erro ao buscar link: %s - %s", url, resp.Status)
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	meta := extractHTMLMeta(resp.Body)
	model = &LinksModel{
		URL:            url,
		Title:          meta.Title,
		Image:          meta.Image,
		Description:    meta.Description,
		SiteName:       meta.SiteName,
		LastChecked:    time.Now(),
		LastStatusCode: resp.StatusCode,
		LastStatus:     resp.Status,
	}
	return model, nil
}

func getTitleFromHTML(bodyReader io.Reader) (title string, err error) {
	tkn := html.NewTokenizer(bodyReader)

	var isTitle bool

	for {

		tt := tkn.Next()

		switch {
		case tt == html.ErrorToken:
			return

		case tt == html.StartTagToken:

			t := tkn.Token()

			isTitle = t.Data == "title"

		case tt == html.TextToken:

			t := tkn.Token()

			if isTitle {

				return t.Data, nil
				isTitle = false
			}
		}
	}
	return "", nil
}

// GetUnreaden retorna os links não lidos há mais de 24h limitados a 5 por chat, em ordem de criação
func (lsv *LinksService) GetUnreaden() (map[int64][]*LinksModel, error) {
	allLinks, err := lsv.repository.GetUnreaden(time.Duration(24 * time.Hour))
	if err != nil {
		return nil, err
	}
	// Separa os links por chat
	var linksByChat = make(map[int64][]*LinksModel)
	for _, link := range allLinks {
		if len(linksByChat[link.ChatID]) == 0 {
			linksByChat[link.ChatID] = make([]*LinksModel, 0)
		}
		if len(linksByChat[link.ChatID]) < 5 {
			linksByChat[link.ChatID] = append(linksByChat[link.ChatID], link)
		}
	}

	return linksByChat, nil
}
