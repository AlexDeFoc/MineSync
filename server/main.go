package server

import (
	"encoding/json"
	"io"
	"mineSync/core"
	"mineSync/globals"
	"net/http"
	"strings"

	tunnel "github.com/jonasfj/go-localtunnel"
)

var metadataStore map[string][]byte

var changed, removed []string

var changedDataStore map[string][]byte = make(map[string][]byte, 0)

func Run() error {
    // Set for errors
	globals.SetFileName("server.go")
    globals.SetFuncName("Run")

    // Create listener
    ln, err := tunnel.Listen(tunnel.Options{})
    if err != nil {
        globals.Error("Error creating ln", err)
        return err
    }

    // Log URL
    url := ln.URL()
    url = strings.ReplaceAll(url, `https://`, ``)
    url = strings.ReplaceAll(url, `.loca.lt`, ``)
    globals.Info("URL:", url)

    // Route
    mx := http.NewServeMux()
    route(mx, "/", exchangeHandler)

    // Create server
    server := http.Server{Handler: mx}

    // Start server
    err = server.Serve(ln)
    if err != nil {
        globals.Error("Error starting server", err)
        return err
    }

    // End
    return nil
}

// FLOW
// CLIENT STORES METADATA ON SERVER
// SERVER COMPARES LOCAL vs CLIENT
// SERVER STORES for each CHANGED the file data on a URL at /files/
// CLIENT REQUESTS in the mean time LIST OF CHANGED
// CLIENT REQUESTS right after LIST OF REMOVED
// CLIENT REQUESTS for each CHANGED - ITS DATA

func exchangeHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPut:
        // Check if client storing metadata
        if r.Header.Get("metadata") == "1" {
            globals.Info("Client is storing METADATA on server")

            // Read & Decode metadata
            err := json.NewDecoder(r.Body).Decode(&metadataStore)
            if err != nil {
                globals.Error("Error reading & decoding metadata body", err)
                return
            }

            // Get server metadata
            metadata, err := core.GetMetadata()
            if err != nil {
                globals.Error("Error getting server metadata")
                return
            }

            // Compare metadata lists
            err = core.CompareMetadata(&changed, &removed, metadataStore, metadata)
            if err != nil {
                globals.Error("Error comparing metadata")
                return
            }

            // Load onto server changed data
            err = core.LoadChanged(&changedDataStore, changed)
            if err != nil {
                globals.Error("Error loading changed data")
                return
            }

            // Log changed & removed lists
            for i := range changed {
                globals.Info("C/N:", changed[i])
            }
            for i := range removed {
                globals.Info("R:", removed[i])
            }

            // Log if empty changed || removed lists
            if len(changed) == 0 {
                globals.Info("Nothing new")
            }
            if len(removed) == 0 {
                globals.Info("Nothing removed")
            }
        }

        if r.Header.Get("finish") == "1" {
            globals.Info("FINISH SYNC")
        }
    case http.MethodGet:
        // Check if client requesting list of changed
        if r.Header.Get("changed") == "1" {
            globals.Info("Client is requesting CHANGES")

            // Encode changed list
            jsonData, err := json.Marshal(changed)
            if err != nil {
                globals.Error("Error encoding changed list", err)
                return
            }

            // Respond with data to req
            _, err = w.Write(jsonData)
            if err != nil {
                globals.Error("Error sending changed list to client", err)
            }
        }

        // Check if client requesting list of removed
        if r.Header.Get("removed") == "1" {
            globals.Info("Client is requesting REMOVED")

            // Encode removed list
            jsonData, err := json.Marshal(removed)
            if err != nil {
                globals.Error("Error encoding removed list", err)
                return
            }

            // Respond with data to req
            _, err = w.Write(jsonData)
            if err != nil {
                globals.Error("Error sending removed list to client", err)
            }
        }

        // Check if client requesting changed data
        if r.Header.Get("data") == "1" {
            // Get file path
            path, err := io.ReadAll(r.Body)
            if err != nil {
                globals.Error("Error reading file path", err)
                return
            }

            // Log request info
            globals.Info("DATA requested for FILE:", string(path))

            // Respond with data
            _, err = w.Write(changedDataStore[string(path)])
            if err != nil {
                globals.Error("Error sending data to client", err)
                return
            }

            // Log finish sending
            globals.Info("Data sent")
        }
    }
}

func route(mux *http.ServeMux, dir string, fun http.HandlerFunc) {
    mux.HandleFunc(dir, fun)
}
