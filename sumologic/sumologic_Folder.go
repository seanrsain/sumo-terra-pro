package sumologic

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

//READ
func (s *Client) GetFolder(id string) (*Folder, error) {
	log.Println("####Begin GetFolder####")

	url := fmt.Sprintf("v2/content/folders/%s", id)
	log.Printf("Folder read url: %s", url)

	//Execute the folder read request
	rawFolder, _, err := s.Get(url)

	//If there was an error, exit here and return it
	if err != nil {
		return nil, err
	}

	//Parse a Folder struct from the response
	var folder Folder
	err = json.Unmarshal(rawFolder, &folder)

	//Exit here if there was an error parsing the json
	if err != nil {
		return nil, err
	}

	log.Println("####End GetFolder####")
	return &folder, nil
}

//DELETE
func (s *Client) DeleteFolder(id string) error {
	log.Println("####Begin DeleteFolder####")

	log.Printf("Deleting Folder Id: %s", id)

	url := fmt.Sprintf("v2/content/%s/delete", id)
	log.Printf("Folder delete url: %s", url)

	//start the deletion job
	rawJID, err := s.Delete(url)

	var jid JobId

	//Parse the response for the JobId
	err = json.Unmarshal(rawJID, &jid)

	//Exit here if there was an error parsing the json
	if err != nil {
		return err
	}

	url = fmt.Sprintf("v2/content/%s/delete/%s/status", id, jid.ID)
	log.Printf("Folder delete job status url: %s", url)

	//Ensure the job has completed before proceeding
	for {

		log.Println("====Start Job Status Check====")

		//check the status of the job
		rawStatus, _, err := s.Get(url)

		//If there were no errors during the request, proceed
		if err != nil {
			return err
		}

		// Parse the Job Status message
		var status Status
		err = json.Unmarshal(rawStatus, &status)

		//Exit here if there was an error parsing the json
		if err != nil {
			return err
		}

		log.Printf("Job Status: %s", status.Status)
		log.Printf("Job Message: %s", status.StatusMessage)
		log.Println("Job Errors:")
		log.Println(status.Errors)
		log.Println("====End Job Status Check====")

		if status.Status == "Success" {
			break
		}

		log.Printf("Sleeping for 1 second before retrying Job Status check...")
		time.Sleep(1 * time.Second)
	}

	log.Println("####End DeleteFolder####")
	return err
}

//CREATE
func (s *Client) CreateFolder(folder Folder) (string, error) {
	log.Println("####Begin CreateFolder####")

	url := "v2/content/folders"
	log.Printf("Create folder url: %s", url)

	//Initiate folder creation
	responseData, err := s.Post(url, folder)

	//Exit if there was an error during the request
	if err != nil {
		return "", err
	}

	//Parse Response
	var folderResponse Folder
	err = json.Unmarshal(responseData, &folderResponse)

	//Catch parsing errors
	if err != nil {
		return "", err
	}

	log.Printf("New folder ID is: %s", folderResponse.ID)
	return folderResponse.ID, nil
}

//FOLDER UPDATE

func (s *Client) UpdateFolder(folder Folder) error {
	log.Println("####Begin folder update####")

	url := fmt.Sprintf("v2/content/folders/%s", folder.ID)
	log.Printf("Update folder job status url: %s", url)

	_, err := s.Put(url, folder)

	log.Println("####End folder update####")
	return err
}

//Retrieve PersonalFolder for current credentials
func (s *Client) getPersonalFolder() (*Folder, error) {
	log.Println("####Begin loading Personal Folder####")
	url := "v2/content/folders/personal"
	rawFolder, _, err := s.Get(url)
	if err != nil {
		return nil, err
	}

	var personalFolder Folder
	err = json.Unmarshal(rawFolder, &personalFolder)

	if err != nil {
		return nil, err
	}

	log.Printf("PF Id: %s", personalFolder.ID)
	log.Printf("Name: %s", personalFolder.Name)
	log.Printf("Description: %s", personalFolder.Description)

	log.Println("####End loading Personal Folder####")
	return &personalFolder, nil
}
