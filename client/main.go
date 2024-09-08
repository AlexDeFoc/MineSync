package client

import (
	"bytes"
	"encoding/json"
	"io"
	"mineSync/core"
	"mineSync/globals"
	"net/http"
)

// FLOW
// STORE METADATA ON SERVER
// REQUEST CHANGED FILES LIST
// REQUEST REMOVED FILES LIST - remove them including empty folders
// REQUEST for each CHANGED FILE - ITS DATA

// Make url
var url string

func Run() error {
	// Set for errors
	globals.SetFileName("client.go")
	globals.SetFuncName("Run")

    // Get url
    url = "http://" + globals.NetworkURL + ".loca.lt" + "/"

    // Get metadata
    metadata, err := core.GetMetadata()
    if err != nil {
        globals.Error("Error getting metadata", err)
        return err
    }

    // Store metadata on server
    err = store(metadata)
    if err != nil {
        globals.Error("Error storing metadata")
        return err
    }

    // Req changed & removed from server
    changed, err := reqList("changed")
    if err != nil {
        globals.Error("Error getting changed")
        return err
    }

    removed, err := reqList("removed")
    if err != nil {
        globals.Error("Error removed changed")
        return err
    }

    for i := range changed {
        globals.Info("C/N:", changed[i])
    }
    for i := range removed {
        globals.Info("R:", removed[i])
    }

    // Req changed data
    for i := range changed {
        // Log start of sync
        globals.Info("Starting to sync:", changed[i])

        // Req file data
        data, err := reqChangedData(changed[i])
        if err != nil {
            globals.Error("Error getting data")
            return err
        }

        // Call write
        err = core.Write(changed[i], data)
        if err != nil {
            globals.Error("Error writing data")
            return err
        }

        globals.Info("W:", changed[i])
    }

    // Remove missing files on server
    err = core.Remove(removed)
    if err != nil {
        globals.Error("Error removing files")
        return err
    }

    if len(changed) == 0 {
        globals.Info("Nothing new")
    }
    if len(removed) == 0 {
        globals.Info("Nothing removed")
    }

    sendFinish()
    globals.Info("FINISHED SYNC")

    // End
    return nil
}

func sendFinish() error {
    // Make request
    req, err := http.NewRequest(http.MethodPut, url, nil)
    if err != nil {
        globals.Error("Error making request", err)
        return err
    }

    // Set header
    req.Header.Set("finish", "1")

    // Make def client
    client := http.Client{}

    // Make client DO req
    _, err = client.Do(req)
    if err != nil {
        globals.Error("Error doing req", err)
        return err
    }

    return nil
}

func reqChangedData(listEntry string) ([]byte, error) {
    // Make data
    var data []byte

    // Make req
    req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader([]byte(listEntry)))
    if err != nil {
        globals.Error("Error making req", err)
        return data, err
    }

    // Set req header
    req.Header.Set("data", "1")

    // Make def client
    client := http.Client{}

    // Make client DO req
    resp, err := client.Do(req)
    if err != nil {
        globals.Error("Error doing req", err)
        return data, err
    }

    // Decode data
    data, err = io.ReadAll(resp.Body)
    if err != nil {
        globals.Error("Error decoding data", err)
        return data, err
    }

    // End
    resp.Body.Close()
    return data, nil
}

func reqList(listType string) ([]string, error) {
    // Make list
    var list []string

    // Make req
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        globals.Error("Error creating req", err)
        return list, err
    }

    // Set req header
    req.Header.Set(listType, "1")

    // Make def client
    client := http.Client{}

    // Make client DO req
    resp, err := client.Do(req)
    if err != nil {
        globals.Error("Error doing req", err)
        return list, err
    }

    // Decode data
    err = json.NewDecoder(resp.Body).Decode(&list)
    if err != nil {
        globals.Error("Error decoding list", err)
        return list, err
    }

    // End
    resp.Body.Close()
    return list, nil
}

func store(metadata map[string][]byte) error {
    // Marshal the data
    jsonData, err := json.Marshal(metadata)
    if err != nil {
        globals.Error("Error encoding metadata", err)
        return err
    }

    // Create request
    req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(jsonData))
    if err != nil {
        globals.Error("Error making new req", err)
        return err
    }

    // Add header to req
    req.Header.Set("metadata", "1")

    // Make default client
    client := http.Client{}

    // Client DOES the req
    _, err = client.Do(req)
    if err != nil {
        globals.Error("Error doing req", err)
        return err
    }

    return nil
}
