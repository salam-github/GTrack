package create

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Artist struct {
	Id           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Relations    map[string][]string `json:"relations"`
}

type Relations struct {
	Id        int                 `json:"id"`
	Relations map[string][]string `json:"datesLocations"`
}

func getRelations() ([]Relations, error) {
	body, err := createBody("https://groupietrackers.herokuapp.com/api/relation")

	if err != nil {
		fmt.Println("getRelations error")
		return []Relations{}, err
	}

	var temp struct {
		Relations []Relations `json:"index"`
	}

	err = json.Unmarshal(body, &temp)
	if err != nil {
		fmt.Println("getRelations error")
		return []Relations{}, err
	}

	for i := range temp.Relations {
		relations := make(map[string][]string, 1)

		for key, dates := range temp.Relations[i].Relations {
			var modDates []string
			var modKey string

			modKey = strings.ReplaceAll(key, "-", ", ")
			modKey = strings.ReplaceAll(modKey, "_", " ")

			for _, date := range dates {
				modDates = append(modDates, strings.ReplaceAll(date, "-", "."))
			}

			relations[modKey] = modDates
			temp.Relations[i].Relations = relations
		}
	}

	return temp.Relations, err
}
func GetArtists() ([]Artist, error) {
	var artists []Artist
	body, err := createBody("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("GetArtist error")
		return []Artist{}, err
	}

	var temp []struct {
		Id           int      `json:"id"`
		Image        string   `json:"image"`
		Name         string   `json:"name"`
		Members      []string `json:"members"`
		CreationDate int      `json:"creationDate"`
		FirstAlbum   string   `json:"firstAlbum"`
		Relations    string   `json:"relations"`
	}

	err = json.Unmarshal(body, &temp)
	if err != nil {
		fmt.Println("GetArtist error")
		return []Artist{}, err
	}

	relations, err := getRelations()
	if err != nil {
		fmt.Println("GetArtist error")
		return []Artist{}, err
	}

	for v := range temp {
		for key, dates := range relations[v].Relations {
			var modified []string
			for _, date := range dates {
				modified = append(modified, strings.ReplaceAll(date, "-", "."))
			}
			relations[v].Relations[key] = modified
		}

		value := Artist{
			temp[v].Id,
			temp[v].Image,
			temp[v].Name,
			temp[v].Members,
			temp[v].CreationDate,
			strings.ReplaceAll(temp[v].FirstAlbum, "-", "."),
			relations[v].Relations,
		}

		artists = append(artists, value)
	}

	return artists, err
}
func createBody(link string) ([]byte, error) {
	var body []byte
	response, err := http.Get(link)
	if err != nil {
		fmt.Println("createBody error")
		return body, err
	}

	body, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("createBody error")
		return body, err
	}

	return body, err
}
