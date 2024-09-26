package index

import (
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Handler struct{}

const pageContent = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>PetPet-Go</title>
    <meta content="PetPet GIF Generator" property="og:title"/>
    <meta content="To use send 'https://pet.wavycat.ru/ds/user_id.gif'" property="og:description"/>
    <meta content="https://github.com/wavy-cat/petpet-go" property="og:url"/>
    <meta content="#AF47DC" data-react-helmet="true" name="theme-color"/>
    <link href="https://fonts.googleapis.com/css2?family=Pixelify+Sans&display=swap" rel="stylesheet">
    <style>
        body {
            font-family: "Pixelify Sans", sans-serif;
            font-optical-sizing: auto;
            font-weight: bold;
            font-style: normal;
            color: white;
            background-color: black;
            height: 100%;
            text-align: center;
        }

        .outer {
            display: table;
            position: absolute;
            top: 0;
            left: 0;
            height: 100%;
            width: 100%;
        }

        .middle {
            display: table-cell;
            vertical-align: middle;
        }
    </style>
</head>
<body>
<div class="outer">
    <div class="middle">
        <p>Waiting for something to happen?</p>
        <a href="https://github.com/wavy-cat/petpet-go">Click</a>
    </div>
</div>
</body>
</html>`

func (Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*zap.Logger)

	w.Header().Set("Content-Type", "text/html")

	_, err := w.Write([]byte(strings.ReplaceAll(pageContent, "'", "`")))
	if err != nil {
		logger.Error("Error sending response", zap.Error(err))
	}
}
