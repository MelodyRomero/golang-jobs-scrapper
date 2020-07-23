# Project Title

## Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Contributing](../CONTRIBUTING.md)

## About <a name = "about"></a>

Write about 1-2 paragraphs describing the purpose of your project.

## Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them.

```
Give examples
```

### Installing

A step by step series of examples that tell you how to get a development env running.

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo.

## Usage <a name = "usage"></a>

Add notes about how to use the system.

Call order of callbacks
1. OnRequest
Called before a request

2. OnError
Called if error occured during the request

3. OnResponseHeaders
Called after response headers received

4. OnResponse
Called after response received

5. OnHTML
Called right after OnResponse if the received content is HTML

6. OnXML
Called right after OnHTML if the received content is HTML or XML

7. OnScraped
Called after OnXML callbacks

c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
})

c.OnError(func(_ *colly.Response, err error) {
    log.Println("Something went wrong:", err)
})

c.OnResponseHeaders(func(r *colly.Response) {
    fmt.Println("Visited", r.Request.URL)
})

c.OnResponse(func(r *colly.Response) {
    fmt.Println("Visited", r.Request.URL)
})

c.OnHTML("a[href]", func(e *colly.HTMLElement) {
    e.Request.Visit(e.Attr("href"))
})

c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
    fmt.Println("First column of a table row:", e.Text)
})

c.OnXML("//h1", func(e *colly.XMLElement) {
    fmt.Println(e.Text)
})

c.OnScraped(func(r *colly.Response) {
    fmt.Println("Finished", r.Request.URL)
})