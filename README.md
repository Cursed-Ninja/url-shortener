<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->

<a name="readme-top"></a>

<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/cursed-ninja/url-shortener">
    <img src="images/logo.png" alt="Logo" width="80" height="80">
  </a>

<h3 align="center">URL Shortener</h3>

  <p align="center">
    This is a simple URL shortener project built using Go.
    <br />
    <a href="https://github.com/cursed-ninja/url-shortener"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/cursed-ninja/url-shortener/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/cursed-ninja/url-shortener/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#workflow">Workflow</a></li>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#architecture">Architecture</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#testing">Testing</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->

## About The Project

The URL Shortener is a powerful and efficient system designed to handle high traffic loads seamlessly. Built with Go, this project transforms long URLs into concise, manageable short links. These short URLs redirect users to the original long URLs.

The shortened URLs, along with their metadata, are stored as documents in MongoDB. To ensure efficient logging of all requests, the system leverages Kafka in conjunction with MongoDB. Additionally, Redis is utilized for caching the shortened URLs, providing faster access and improved performance.

The project has separate directories for each service, thereby, allowing easy deployment to cloud platforms. The services are containerized using Docker, ensuring seamless deployment and scaling.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Workflow

- Shorten Request Workflow

  1. The user sends a POST request to the main service with the long URL.
  2. The main service calls the database service to create and store the shortened url.
  3. The database service returns the shortened URL to the main service.
  4. The main service prepends the base URL to the shortened URL and returns it to the user.

- Redirect Request Workflow
  1. The user sends a GET request to the shortened URL.
  2. The main service calls the cache service to get the original URL.
  3. If the URL is not found in the cache, the cache service calls the database service to get the original URL.
  4. The database service returns the original URL to the cache service.
  5. The cache service stores the original URL in the cache and returns it to the main service.
  6. The main service redirects the user to the original URL.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Built With

- [![Go][GoLang]][GoLang-url]
- [![MongoDB][MongoDB]][MongoDB-url]
- [![Kafka][Kafka]][Kafka-url]
- [![Redis][Redis]][Redis-url]
- [![Docker][Docker]][Docker-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Architecture

![Architecture][architecture-image]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->

## Getting Started

To get a local copy up and running follow these simple example steps.

### Prerequisites

Make sure you have the following installed on your system:

- Go
- Docker
- Kafka
- Redis
- Docker Compose

Additionally, you will need to create an account on the following platforms:

- [MongoDB](https://www.mongodb.com/)

### Installation

1. Create a MongoDB cluster and get the connection url.
2. Run the redis server on your local machine.
3. Clone the repo
   ```sh
   git clone https://github.com/cursed-ninja/url-shortener.git
   ```
4. For each service, create a folder `config` and a file `app.env` inside it. The below is an example of the `app.env` file for each of the service.

   - Main Service

   ```env
   DATABASE_SERVICE_BASE_URL=http://localhost:8081
   BASE_URL=http://localhost:8080
   CACHE_SERVICE_BASE_URL=http://localhost:8082
   KAFKA_SERVICE_BASE_URL=localhost:29092
   ```

   - Database Service

   ```env
   MONGO_URI="ENTER YOUR MONGO URI"
   DB_NAME=url-shortener
   COLLECTION_NAME=urls
   KAFKA_SERVICE_BASE_URL=localhost:29092
   ```

   - Cache Service

   ```env
   DATABASE_SERVICE_BASE_URL=http://localhost:8081
    REDIS_ADDR=localhost:6379
    REDIS_PASSWORD=12345678
    REDIS_DB=0
    KAFKA_SERVICE_BASE_URL=localhost:29092
   ```

   - Kafka Service

   ```env
   MONGO_URI="ENTER YOUR MONGO URI"
   DB_NAME=url-shortener
   MAIN_SERVER_COLLECTION_NAME=main-server-logs
   CACHE_SERVER_COLLECTION_NAME=cache-server-logs
   DATABASE_SERVER_COLLECTION_NAME=database-server-logs
   KAFKA_SERVICE_BASE_URL=localhost:29092
   ```

5. Run the following command to start the docker containers for kafka. Make sure docker engine is running in the background.

   ```sh
   docker-compose up -d
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->

## Usage

After ensuring all the services and dependencies are up and running, you can use the service by following the below steps.

- Make a POST request to the main service with the long URL in the body. The response will contain the shortened URL. The body should be in the below format.

  ```json
  {
    "url": "https://www.google.com"
  }
  ```

- Make a GET request to the shortened URL to be redirected to the original long URL.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- TESTING -->

## Testing

Each service has some unit tests written to ensure the functionality of the service. To run the tests, follow the below steps.

<strong>Some tests would require adding urls (like mongo_uri, redis, etc.) to the tests. Make sure to fill all the placeholders correctly otherwise the tests would fail.</strong>

1. Navigate to the service directory.
2. Run the below command to run the tests.

   ```sh
   go test ./... -v
   ```

This will run a verbose test for all the packages in the service.

<!-- CONTRIBUTING -->

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->

## Contact

Shivam Mahajan - [@Cursed-Ninja](https://linkedin.com/in/cursed-ninja) - shivam.sm2002@gmail.com

Project Link: [https://github.com/cursed-ninja/url-shortener](https://github.com/cursed-ninja/url-shortener)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[contributors-shield]: https://img.shields.io/github/contributors/cursed-ninja/url-shortener.svg?style=for-the-badge
[contributors-url]: https://github.com/cursed-ninja/url-shortener/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/cursed-ninja/url-shortener.svg?style=for-the-badge
[forks-url]: https://github.com/cursed-ninja/url-shortener/network/members
[stars-shield]: https://img.shields.io/github/stars/cursed-ninja/url-shortener.svg?style=for-the-badge
[stars-url]: https://github.com/cursed-ninja/url-shortener/stargazers
[issues-shield]: https://img.shields.io/github/issues/cursed-ninja/url-shortener.svg?style=for-the-badge
[issues-url]: https://github.com/cursed-ninja/url-shortener/issues
[license-shield]: https://img.shields.io/github/license/cursed-ninja/url-shortener.svg?style=for-the-badge
[license-url]: https://github.com/Cursed-Ninja/url-Shortener/blob/main/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/cursed-ninja
[GoLang]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[GoLang-url]: https://golang.org/
[MongoDB]: https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white
[MongoDB-url]: https://www.mongodb.com/
[Kafka]: https://img.shields.io/badge/Apache%20Kafka-231F20?style=for-the-badge&logo=apachekafka&logoColor=white
[Kafka-url]: https://kafka.apache.org/
[Redis]: https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white
[Redis-url]: https://redis.io/
[Docker]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[architecture-image]: images/architecture.png
