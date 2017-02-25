# Charcoal [![Build Status](https://travis-ci.org/dadleyy/charcoal.api.svg?branch=master)](https://travis-ci.org/dadleyy/charcoal.api)

[golang](https://golang.org) backend responsible for persisting information to the database and providing an HTTP/JSON api for [charcoal.ui](https://github.com/dadleyy/charcoal.ui) client consumption.

For detailed information about the api, please see the [API Docs](https://documenter.getpostman.com/view/1070956/charcoal-api/6YsXcyj).

### Installing

This application relies on mysql for it's persistence layer; the migrations for the database can be found in the [charcoal.db](https://github.com/dadleyy/charcoal.db) repository. For more detailed instructions on preparing your database, please refer to that project's [readme](https://github.com/dadleyy/charcoal.db).

After cloning this repository, it's dependencies can be installed using [glide](https://github.com/Masterminds/glide):

```
$ go get -u github.com/Masterminds/glide
$ glide install
```

You will also need to prepare a `.env` file at the root of the repository - an example can be seen [here](https://github.com/dadleyy/charcoal.api/blob/master/.env.example). The third party services required to run the application are:

1. [aws](http://aws.amazon.com/) - s3 is used to store files uploaded by clients
2. [google](https://developers.google.com/apis-explorer/#p/) - the google apis are used by the application for user authentication
3. [mailgun](https://www.mailgun.com/) - the application provides an api endpoint for use w/ mailgun webhooks to process data
