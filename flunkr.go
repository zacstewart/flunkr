package main

import (
    "flag"
    "fmt"
    "github.com/zacstewart/flunkr/flickr"
    "net/http"
    "io"
    "os"
    "strings"
)

func usage() {
    fmt.Fprintf(os.Stderr, "Usage: flunkr <ARGUMENTS>\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func main() {
    var user_id string
    var api_key string
    var response *flickr.Response
    var err error

    flag.Usage = usage
    flag.StringVar(&user_id, "user", "", "The id of the user you want to download")
    flag.StringVar(&api_key, "key", "", "Your Flickr API key")
    flag.Parse()

    if user_id == "" || api_key == "" {
        usage()
    }

    f := flickr.Flickr{api_key}

    response, err = f.Request(
        "flickr.photosets.getList",
        map[string]string{
            "user_id": user_id,
            "per_page": "500",
        },
    )

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    photosets := response.Message.Photosets.Photoset
    for ps := range photosets {
        response, err = f.Request(
            "flickr.photosets.getPhotos",
            map[string]string{
                "photoset_id": photosets[ps].Id,
            },
        )

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        err = os.MkdirAll(response.Message.Photoset.Title, os.ModeDir | 0777)

        if err != nil && !os.IsExist(err) {
            fmt.Println(err)
            os.Exit(1)
        }

        photos := response.Message.Photoset.Photo

        for p := range photos {
            response, err = f.Request(
                "flickr.photos.getSizes",
                map[string]string{
                    "photo_id": photos[p].Id,
                },
            )

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            original := getOriginalSize(response.Message.Sizes.Size)

            filename := fmt.Sprintf(
                "%s/%s",
                photosets[ps].Title.Content,
                photos[p].Id,
            )
            if photos[p].Title != "" {
                filename += fmt.Sprintf("-%s", photos[p].Title)
            }

            // Extract extension
            parts := strings.Split(original.Source, ".")
            filename += fmt.Sprintf(".%s", parts[len(parts) - 1])

            file, err := os.Create(filename)
            defer file.Close()

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            image, err := http.Get(original.Source)
            defer image.Body.Close()

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            _, err = io.Copy(file, image.Body)

            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            fmt.Println(filename)
        }
    }
}

func getOriginalSize(sizes []flickr.Size) (size flickr.Size) {
    for s := range sizes {
        if sizes[s].Label == "Original" {
            size = sizes[s]
        }
    }
    return
}
