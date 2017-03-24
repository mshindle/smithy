package crypt

import (
	"encoding/gob"
	"os"
)

// Saves a key using gob serialization.
func saveKey(filename string, key interface{}) error {
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()

	encoder := gob.NewEncoder(outfile)
	return encoder.Encode(key)
}

// Loads a key from filename into key
func loadKey(filename string, key interface{}) error {
	infile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer infile.Close()

	decoder := gob.NewDecoder(infile)
	return decoder.Decode(key)
}

//func extractPemDataBlock(file string) (*pem.Block, error) {
//	fileContents, err := ioutil.ReadFile(file)
//	if err != nil {
//		log.WithField("file", file).Error("could not read file")
//		return nil, err
//	}
//	block, _ := pem.Decode(fileContents)
//	if block == nil {
//		log.WithField("file", file).Error("could not decode file contents - not PEM encoded")
//		return nil, errors.New("data not PEM encoded")
//	}
//	return block, nil
//}
