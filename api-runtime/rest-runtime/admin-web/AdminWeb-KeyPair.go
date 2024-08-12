// The CB-Spider is a sub-Framework of the Cloud-Barista Multi-Cloud Project.
// The CB-Spider Mission is to connect all the clouds with a single interface.
//
//      * Cloud-Barista: https://github.com/cloud-barista
//
// by CB-Spider Team, 2024.08.

package adminweb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cr "github.com/cloud-barista/cb-spider/api-runtime/common-runtime"
	"github.com/labstack/echo/v4"
)

// KeyPairInfo structure
type KeyPairInfo struct {
	IId          IID
	Fingerprint  string
	TagList      []Tag
	KeyValueList []KeyValue
}

type IID struct {
	NameId   string
	SystemId string
}

type Tag struct {
	Key   string
	Value string
}

type KeyValue struct {
	Key   string
	Value string
}

// Function to fetch KeyPairs
func fetchKeyPairs(connConfig string) ([]*KeyPairInfo, error) {
	resBody, err := getResourceList_with_Connection_JsonByte(connConfig, "keypair")
	if err != nil {
		return nil, fmt.Errorf("error fetching KeyPairs: %v", err)
	}

	var info struct {
		ResultList []*KeyPairInfo `json:"keypair"`
	}
	if err := json.Unmarshal(resBody, &info); err != nil {
		return nil, fmt.Errorf("error decoding KeyPairs: %v", err)
	}

	sort.Slice(info.ResultList, func(i, j int) bool {
		return info.ResultList[i].IId.NameId < info.ResultList[j].IId.NameId
	})

	return info.ResultList, nil
}

// Handler function to render the KeyPair management page
func KeyPairManagement(c echo.Context) error {
	connConfig := c.Param("ConnectConfig")
	if connConfig == "region not set" {
		htmlStr := `
			<html>
			<head>
			    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
				<style>
				th {
				  border: 1px solid lightgray;
				}
				td {
				  border: 1px solid lightgray;
				  border-radius: 4px;
				}
				</style>
			    <script type="text/javascript">
				alert(connConfig)
			    </script>
			</head>
			<body>
				<br>
				<br>
				<label style="font-size:24px;color:#606262;">&nbsp;&nbsp;&nbsp;Please select a Connection Configuration! (MENU: 2.CONNECTION)</label>	
			</body>
		`

		return c.HTML(http.StatusOK, htmlStr)
	}

	regionName, err := getRegionName(connConfig)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	keyPairs, err := fetchKeyPairs(connConfig)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	data := struct {
		ConnectionConfig string
		RegionName       string
		KeyPairs         []*KeyPairInfo
	}{
		ConnectionConfig: connConfig,
		RegionName:       regionName,
		KeyPairs:         keyPairs,
	}

	templatePath := filepath.Join(os.Getenv("CBSPIDER_ROOT"), "/api-runtime/rest-runtime/admin-web/html/keypair.html")
	tmpl, err := template.New("keypair.html").Funcs(template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}).ParseFiles(templatePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error loading template: " + err.Error()})
	}

	return tmpl.Execute(c.Response().Writer, data)
}

// Function to create a KeyPair
func CreateKeyPair(c echo.Context) error {
	connConfig := c.QueryParam("ConnectionName")
	keyPairName := c.Param("Name")

	url := fmt.Sprintf("http://%s:%s/spider/keypair", cr.ServiceIPorName, cr.ServicePort)
	reqBody := fmt.Sprintf(`{"ConnectionName": "%s", "ReqInfo": {"Name": "%s"}}`, connConfig, keyPairName)
	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create KeyPair"})
	}

	var keyPair KeyPairInfo
	if err := json.NewDecoder(resp.Body).Decode(&keyPair); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"Result": "true", "PrivateKey": keyPair.KeyValueList[0].Value})
}

// Function to delete a KeyPair
func DeleteKeyPair(c echo.Context) error {
	connConfig := c.QueryParam("ConnectionName")
	keyPairName := c.Param("Name")

	url := fmt.Sprintf("http://%s:%s/spider/keypair/%s?connectionName=%s", cr.ServiceIPorName, cr.ServicePort, keyPairName, connConfig)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete KeyPair"})
	}

	return c.JSON(http.StatusOK, map[string]string{"Result": "true"})
}
