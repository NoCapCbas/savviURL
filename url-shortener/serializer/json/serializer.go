package json 

import (
  "encoding/json"
  "github.com/pkg/errors"
  "github.com/NoCapCbas/url-shortener/urlshortener"
)

type Redirect struct {}

func (*r Redirect) Decode(input []byte) (*urlshortener.Redirect, error){
  redirect := &urlshortener.Redirect{}
  if err := json.Unmarshal(input, redirect); err != nil {
    return nil, errors.Wrap(err, "serializer.Redirect.Decode")
  }
  return redirect, nil
}

func (*r Redirect) Encode(input *urlshortener.Redirect) ([]byte, error){
  rawMsg, err := json.Marshal(input)
  if err != nil {
    return nil, errors.Wrap(err, "serializer.Redirect.Encode")
  }
  return rawMsg, nil
}
