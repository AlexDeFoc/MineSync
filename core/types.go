package core

type File struct{
    Path string `json:"path"`
    Hash []byte `json:"hash"`
}

type ServerFile struct{
    Relative_Path string `json:"path"`
    Data []byte `json:"hash"`
}
